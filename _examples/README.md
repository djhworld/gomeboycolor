This is an example of how to write a frontend for the emulator. 

It uses [tcell](https://github.com/gdamore/tcell) to render the screen to your terminal and handle keyboard input. It requires a terminal with 256 colour support (haven't tested it on others).

You will probably want to adjust your terminals font settings and vertical character spacing as the screen will probably scroll off the end of your terminal window. 


![screenshot](https://github.com/djhworld/gomeboycolor/raw/master/_examples/terminal.png)


## How to run

```
go run . <path-to-rom-file>
```

Note: Pressing `Esc` will quit the application

## Overview of files

### terminal\_io.go

Sets up the IO related stuff to handle screen drawing and keyboard updates. 

In this particular case we are rendering to the terminal. 

### noop\_store.go

Sets up a no-op battery save store. You can change this to write to a filesystem or other storage medium, but for this example it just does nothing.

### main.go

Glues everything together and runs the application

