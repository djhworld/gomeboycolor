*"Yet another emulator?"*

There is no point beating around the bush, people have written Gameboy emulators before with great success, however I've always been curious as to how they are done. This project was set up as a learning exercise to help me understand how the basic principles of how computers function at the lowest level and how you can simulate them in software.

Since then it's become my side project and I feel I'm at a point where I can share it with the world. 

# Features

* Cross platform support for 64-bit Windows, Linux and OSX operating systems
* Supports most<sup>&#8224;</sup> of the traditional black and white Gameboy ROMs
* Supports most<sup>&#8224;</sup> Gameboy Color ROMs with colour support<sup>&#8225;</sup> (and black and white mode for supported GBC titles if desired)
* Battery saves facility, with compressed save files for lightweight storage
* Six screen resolutions
* Passes all blargg CPU instruction and instruction timing tests
* Delightful emulation of the scrolling Nintendo "boot screen" when you load the emulator up (can be disabled)

<span class="small"><sup>&#8224;</sup> *This project is still under active development so there will inevitably be ROMs out there that don't work (yet) or have bugs. For a full list of tested games, [click here](tested-games.html)*</span>
<br/>
<span class="small"><sup>&#8225;</sup> *Gameboy Color features are mostly done but there are a few [outstanding tasks remaining]()*</span>

# Screenies

<div id="screenshots">
	<ul>
		<li><img src="images/screenshots/tetrisdx.png" alt="Tetris DX" title="Tetris DX" /></li>
		<li><img src="images/screenshots/lionking1.png" alt="The Lion King" title="The Lion King"  /></li>
		<li><img src="images/screenshots/mariotennis1.png"  alt="Mario Tennis" title="Mario Tennis"/></li>
		<li><img src="images/screenshots/warioland2-1.png"  alt="Warioland 2" title="Warioland 2" /></li>
		<li><img src="images/screenshots/pokemonyellow1.png" alt="Pokemon Yellow" title="Pokemon Yellow" /></li>
		<li><img src="images/screenshots/zelda1.png" alt="Zelda" title="Zelda" /></li>
	</ul>
</div>

# Coming soon

* Sound. I've been putting this off for quite some time but as I get nearer to feature complete status, this is the next priority (see [#10](https://github.com/djhworld/gomeboycolor/issues/10))
* More game support. Memory bank controller MBC2 is currently unsupported, along with a few other cartridge types
* GUI based launcher with ROM and RAM saves administration (see [#35](https://github.com/djhworld/gomeboycolor/issues/35))


# Project

The emulator is written entirely in [go](http://golang.org/). Why? It's fast, it's fun to write code in, has great tooling and is really simple to pick up. It also has a history of [people writing emulators in](http://dave.cheney.net/2013/01/09/go-the-language-for-emulators), maybe it's the next "Hello World"?

I started the project after reading 

