gomeboycolor
============================

**This project is a work in progress and is no way near complete!**

Nintendo Gameboy Color emulator

This is a backend library that provides the core emulator runtime, it is designed to be used in conjuction with a frontend. 

Available frontends: 

* [gomeboycolor-glfw](https://github.com/djhworld/gomeboycolor-glfw) - Uses `libglfw` to provide you a native, windowed version of the emulator.
* [\_examples](https://github.com/djhworld/gomeboycolor/tree/master/_examples) - Uses [tcell](https://github.com/gdamore/tcell) to render the emulator in your terminal.


You are welcome to visit the github page for this project by [clicking here](http://djhworld.github.io/gomeboycolor)

FAQ
============================


### Backend?

This module will emulate the hardware of the Gameboy Color.

You can write a 'frontend' to receive the screen data and render it to a medium of your choosing, along with handling keyboard inputs, and something to handle saving battery saves.

See the [\_examples](https://github.com/djhworld/gomeboycolor/tree/master/_examples)  directory for a simple example of how to write a frontend. Alternatively, look at the 'Available frontends' above.


### Features?

* ⚠️ Mostly works. It is not a perfect emulator by any means and some games might not function correctly.
  * ✅ blargg CPU tests pass
  * ❌ Memory timing tests don't pass
* ✅ Supports battery saves for ROMS that allow you to save state
* ❌ Audio is NOT implemented right now
* ❌ Does not support games that require the Gameboy Color HDMA extensions


### How do I build it?

This is a 'library' module, no build required. 


License
-----------------------------

MIT License

Copyright (c) 2013-2018 Daniel James Harper

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
