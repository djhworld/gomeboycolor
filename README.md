gomeboycolor
============================
**This project is a work in progress and is no way near complete!**

Nintendo Gameboy Color emulator, this is my first emulator so I'm learning as I go along...

TODO
---------------------------
1. **Windowing** -> currently no windowing is rendered
2. **8x16 sprites** -> I'm not sure how these work yet
3. ~~ **Investigate scrolling bug** -> there are some issues with the drawing code that I think are connected to scrolling, might need to rewrite to fix ~~
4. **MBC1 RAM saving** -> allow users to save games!
5. **MBC2** -> to implement
6. **MBC3** -> partially implemented but needs work
7. **Colour** -> no gameboy color features are supported yet
8. **MBC4** -> to implement
9. **MBC5** -> to implement
10. **Sound** -> looks tricky, but lower priority to the above
11. Investigate why the clock is ticking so fast in games like Super Mario Land

Progress
---------------------------
![Emulator boots to boot screen](https://github.com/djhworld/gomeboycolor/raw/master/images/boot_sequence.png)&nbsp;
![Super Mario Land](https://github.com/djhworld/gomeboycolor/raw/master/images/sml.gb.png)&nbsp;
![Super Mario Land game](https://github.com/djhworld/gomeboycolor/raw/master/images/sml_game.gb.png)&nbsp;
![Super Mario Land 2](https://github.com/djhworld/gomeboycolor/raw/master/images/sml2.gb.png)&nbsp;
![Tetris](https://github.com/djhworld/gomeboycolor/raw/master/images/tetris.gb.png)&nbsp;
![Passes all blargg CPU tests](https://github.com/djhworld/gomeboycolor/raw/master/images/cpu_instrs.gb.png)&nbsp;
![Passes blargg instruction timing test](https://github.com/djhworld/gomeboycolor/raw/master/images/instr_timing.gb.png)&nbsp;
![test](https://github.com/djhworld/gomeboycolor/raw/master/images/test.gb.png)&nbsp;

Nice to haves
---------------------------
1. Ability to pause and serialize the state of the emulator to disk (and can reload state from disk)
2. Ability to enter the debugger while the emulator is running
3. Frame rate counter on screen
4. Would be nice if there was a way of packaging GLFW/GL libs statically into the binary, I believe this is a limitation right now that makes distributing this emulator a bit of a pain
5. Provide an interface to the emulator so you can maybe hook in AI or whatever?

