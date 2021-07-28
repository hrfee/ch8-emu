package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	FILE  string
	DEBUG bool
	STEP  = 0
	SCALE = 20
	PAD   = 2
)

func Println(a ...interface{}) {
	if !DEBUG {
		return
	}
	fmt.Println(a...)
}

func main() {
	flag.StringVar(&FILE, "file", FILE, "file to run")
	flag.BoolVar(&DEBUG, "debug", DEBUG, "debug info")
	flag.IntVar(&STEP, "step", STEP, "wait n ms between FDE cycles")
	flag.IntVar(&SCALE, "scale", SCALE, "scale pixels by n")
	flag.IntVar(&PAD, "pad", PAD, "pad (scaled) pixels with n pixels on each side")
	flag.Parse()

	m := newMachine(STEP, SCALE, PAD)

	f, err := os.Open(FILE)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	off := int64(0)
	ops := make([]byte, 1)
	for {
		_, err := f.ReadAt(ops, off)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		m.Mem[off+0x200] = ops[0]
		off++

	}
	m.Run()
}
