#### chip-8 emulator

bare minimum chip-8 emulator. Should support all instructions. Some programs don't work. No sound yet. Keypad uses normal layout:
```
1234
qwer
asdf
zxcv
```

```
Usage of ./ch8-emu:
  -debug
    	debug info
  -file string
    	file to run
  -pad int
    	pad (scaled) pixels with n pixels on each side (default 2)
  -scale int
    	scale pixels by n (default 20)
  -speed int
    	rough CPU speed in Hz. (default 800)
```

#### references

[1](https://tobiasvl.github.io/blog/write-a-chip-8-emulator), [2](https://github.com/mattmikolay/chip-8/wiki/CHIP%E2%80%908-Instruction-Set)
