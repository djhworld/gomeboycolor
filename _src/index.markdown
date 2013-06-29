<div id="nav">
	<ul>
		<li><a href="/#features">Features</a></li>
		<li><a href="/#screenies">Screenshots</a></li>
		<li><a href="/#about">About</a></li>
		<li><a href="/#roadmap">Roadmap</a></li>
		<li><a href="/#documentation">Documentation</a></li>
	</ul>
</div>

*"Yet another emulator?"*

There is no point beating around the bush, people have written Gameboy emulators before with great success, however I've always been curious as to how they are done. This project was set up as a learning exercise to help me understand how the basic principles of how computers function at the lowest level and how you can simulate them in software.

Since then it's become my side project and I feel I'm at a point where I can share it with the world. 

# Features

* Cross platform support for 64-bit Windows, Linux and OSX operating systems, binaries available [here]()
* Supports most<sup>&#8224;</sup> of the traditional black and white Gameboy ROMs
* Supports most<sup>&#8224;</sup> Gameboy Color ROMs with colour support<sup>&#8225;</sup> (and black and white mode for supported GBC titles if desired)
* Battery saves facility, with compressed save files for lightweight storage
* Six screen resolutions
* Passes all blargg CPU instruction and instruction timing tests
* Delightful emulation of the scrolling Nintendo "boot screen" when you load the emulator up (can be disabled)

<div class="footnotes">
	<span class="small"><sup>&#8224;</sup> *This project is still under active development so there will inevitably be ROMs out there that don't work (yet) or have bugs.</span>
	<br/>
	<span class="small"><sup>&#8225;</sup> *Gameboy Color features are mostly done but there are a few [outstanding tasks remaining]()*</span>
</div>

# Screenshots 

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

# About 

The emulator is written entirely in [go](http://golang.org/). Why? It's fast, it's fun to code in, has great tooling and is really simple to pick up. It also has a history of [people writing emulators in](http://dave.cheney.net/2013/01/09/go-the-language-for-emulators), maybe it's the next "Hello World"?

I started working on this after reading [Code: The Hidden Language of Computer Hardware and Software](http://www.amazon.co.uk/gp/product/0735611319/ref=as_li_tf_tl?ie=UTF8&camp=1634&creative=6738&creativeASIN=0735611319&linkCode=as2&tag=djhworld-21) by Charles Petzold, a delightful read that takes you through the history of how computers came to be and how they function, right from the days of telegraph relays to the modern transistor. This book inspired me to investigate further and what better way to do it than writing your own computer? Well, *technically* it's Nintendo's computer, implemented in software. 

People have often asked me how challenging it is to actually write an emulator that works, and as it turns out, it really isn't that difficult. However, it's immensely frustrating, baffling, tedious, anger inducing, soul crushing, boring, exciting, challenging and thankfully *hugely rewarding*. A journey of highs and lows with sometimes very little to show for it, but as soon as all the pieces start to fall into place and you see your software boot into its own environment all by itself, a small bubble of pride makes it all worthwhile.

The project is open sourced under the MIT license, details of which you can [view here.](https://raw.github.com/djhworld/gomeboycolor/master/LICENSE.txt)

# Roadmap

* Sound. I've been putting this off for quite some time but as I get nearer to feature complete status, this is the next priority (see [#10](https://github.com/djhworld/gomeboycolor/issues/10))
* More game support. Memory bank controller MBC2 is currently unsupported, along with a few other cartridge types
* GUI based launcher with ROM and RAM saves administration (see [#35](https://github.com/djhworld/gomeboycolor/issues/35))

# Documentation 

## Installing

Installation involves simply unzipping the downloaded ZIP file to a location of your choice. \*nix based distributions can add the ````bin```` directory to their PATH if desired

## Running

To launch the emulator, simply invoke the executable with the location of a ROM file on your machine, e.g.

	./gomeboycolor ~/location/to/my/romfile.gbc

## Usage 

You can pass some optional flags to the emulator (before the ROM file argument) that change some parameters of how the emulator runs, details of these are as follows 

``````
Usage: -

To launch the emulator, simply run and pass it the location of your ROM file, e.g. 

	gomeboycolor location/of/romfile.gbc

Flags: -

	-help			->	Show this help message
	-skipboot		->	Disables the boot sequence and will boot you straight into the ROM you have provided. Defaults to false
	-color			->	Turns color GB features on. Defaults to true
	-showfps		->	Prints average frames per second to the console. Defaults to false
	-dump			->	Dump CPU state after every cycle. Will be very SLOW and resource intensive. Defaults to false
	-size=(1-6)		->	Set screen size. Defaults to 1.
	-title=(title)		->	Change window title. Defaults to 'gomeboycolor'.

You can pass an option argument to the boolean flags if you want to enable that particular option. e.g. to disable the boot screen you would do the following

	gomeboycolor -skipboot=false location/of/romfile.gbc
``````

## Storage

The emulator will create a directory under your home directory for storing ROM saves and other settings. So for Windows this would be something like

	C:\Users\joebloggs\.gomeboycolor

Or on OSX/Linux

	~/.gomeboycolor

You can optionally create the file ````config.json```` and store it under there. This will allow you to set the defaults for the flags so you don't need to keep providing them when running the emulator. An example file is detailed below

``````
{
	"Title":"gomeboycolor",
	"ScreenSize":2,
	"ColorMode":true,
	"SkipBoot":true,
	"DisplayFPS":false
}
``````

Future plans will negate the need for storing lots of items in this folder as I plan to move it to a small database such as sqlite. 
