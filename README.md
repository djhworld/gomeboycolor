gomeboycolor
============================
**This project is a work in progress and is no way near complete!**

Nintendo Gameboy Color emulator, this is my first emulator so I'm learning as I go along...

You are welcome to visit the github page for this project by [clicking here](http://djhworld.github.io/gomeboycolor), this includes links to [executables](http://djhworld.github.io/gomeboycolor/#downloads), [documentation](http://djhworld.github.io/gomeboycolor/#documentation) and some [background about the project](http://djhworld.github.io/gomeboycolor/#about).

FAQ
============================

#### How do I build it?

You will need an installation of [Go](http://golang.org) (version >= 1.1.1) available on your PATH, as well as a version of Ruby with rake. Personally I've had great success using [jruby](http://jruby.org/).

To build the project you need to run ````rake```` from the root of this repository, this will

1. Fetch dependencies
2. Setup folders
3. Build the project

On Windows and OSX you shouldn't have to do anything more then that as I've included the shared libraries this project needs to build with (glew and GLFW)

However on Linux you will need to have libGLEW and libglfw installed.

Builds have been tested on the following systems, note that these are all 64-bit

* OSX Mountain Lion (10.5.8)
* Ubuntu Linux 13.04
* Windows 7

License
-----------------------------
Copyright (c) 2013. Daniel James Harper

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

Progress
---------------------------
![Emulator boots to boot screen](https://github.com/djhworld/gomeboycolor/raw/master/images/boot_sequence.png)&nbsp;
![lionking1](https://github.com/djhworld/gomeboycolor/raw/master/images/lionking1.png)&nbsp;
![lionking2](https://github.com/djhworld/gomeboycolor/raw/master/images/lionking2.png)&nbsp;
![Mario Tennis 1](https://github.com/djhworld/gomeboycolor/raw/master/images/mariotennis1.png)&nbsp;
![Mario Tennis 2](https://github.com/djhworld/gomeboycolor/raw/master/images/mariotennis2.png)&nbsp;
![Warioland](https://github.com/djhworld/gomeboycolor/raw/master/images/warioland1.png)&nbsp;
![Warioland 2](https://github.com/djhworld/gomeboycolor/raw/master/images/warioland2.png)&nbsp;
![Tetris DX](https://f.cloud.github.com/assets/529730/619306/8d6f4d6a-ceca-11e2-9789-f11a0545e643.png)
![Tetris DX](https://f.cloud.github.com/assets/529730/619308/96ecdae2-ceca-11e2-8941-c5e6ba79c5c8.png)
![TetrisDX](https://f.cloud.github.com/assets/529730/619303/86a964c0-ceca-11e2-8c04-ace874c45957.png)
![Pokemon Red](https://github.com/djhworld/gomeboycolor/raw/master/images/pokemonred1.png)&nbsp;
![Pokemon Red](https://github.com/djhworld/gomeboycolor/raw/master/images/pokemonred2.png)&nbsp;
![Zelda](https://github.com/djhworld/gomeboycolor/raw/master/images/zelda.gb.png)&nbsp;
![Super Mario Land](https://github.com/djhworld/gomeboycolor/raw/master/images/sml.gb.png)&nbsp;
![Super Mario Land game](https://github.com/djhworld/gomeboycolor/raw/master/images/sml_game.gb.png)&nbsp;
![Super Mario Land 2](https://github.com/djhworld/gomeboycolor/raw/master/images/sml2.gb.png)&nbsp;
![Tetris](https://github.com/djhworld/gomeboycolor/raw/master/images/tetris.gb.png)&nbsp;
![Metroid](https://github.com/djhworld/gomeboycolor/raw/master/images/metroid1.png)&nbsp;
![Metroid](https://github.com/djhworld/gomeboycolor/raw/master/images/metroid2.png)&nbsp;
![Passes all blargg CPU tests](https://github.com/djhworld/gomeboycolor/raw/master/images/cpu_instrs.gb.png)&nbsp;
![Passes blargg instruction timing test](https://github.com/djhworld/gomeboycolor/raw/master/images/instr_timing.gb.png)&nbsp;
![test](https://github.com/djhworld/gomeboycolor/raw/master/images/test.gb.png)&nbsp;
