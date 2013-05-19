package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"utils"
)

const (
	SKIP_BOOT_FLAG     string = "skipboot"
	SCREEN_SIZE_FLAG          = "size"
	SAVE_DIR_FLAG             = "savedir"
	SHOW_FPS_FLAG             = "showfps"
	SETTINGS_FILE_FLAG        = "settings"
	TITLE_FLAG                = "title"
	DUMP_FLAG                 = "dump"
	DEBUGGER_ON_FLAG          = "debug"
	BREAK_WHEN_FLAG           = "b"
)

var settingsFile *string = flag.String(SETTINGS_FILE_FLAG, "settings.json", "Location of settings file")
var title *string = flag.String(TITLE_FLAG, TITLE, "Title to use")
var fps *bool = flag.Bool(SHOW_FPS_FLAG, false, "Calculate and display frames per second")
var screenSizeMultiplier *int = flag.Int(SCREEN_SIZE_FLAG, 1, "Screen size multiplier")
var skipBoot *bool = flag.Bool(SKIP_BOOT_FLAG, false, "Skip boot sequence")
var savesDir *string = flag.String(SAVE_DIR_FLAG, "saves", "Location where game saves are stored")

//debug stuff...
var dumpState *bool = flag.Bool(DUMP_FLAG, false, "Print state of machine after each cycle (WARNING - WILL RUN SLOW)")
var debug *bool = flag.Bool(DEBUGGER_ON_FLAG, false, "Enable debugger")
var breakOn *string = flag.String(BREAK_WHEN_FLAG, "0x0000", "Break into debugger when PC equals a given value between 0x0000 and 0xFFFF")

type Config struct {
	//mandatory settings
	Title      string
	ScreenSize int
	SkipBoot   bool
	SavesDir   string
	DisplayFPS bool

	//optional
	Debug          bool
	BreakOn        string
	DumpState      bool
	OnCloseHandler func()
}

func NewConfig() *Config {
	var c *Config = new(Config)
	return c
}

func NewConfigFromFile(settingsFile string) (*Config, error) {
	file, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		return nil, err
	}

	//Make sure all settings keys are present
	var initialMap map[string]interface{}
	err = json.Unmarshal(file, &initialMap)

	if err != nil {
		return nil, ConfigSettingsParseError(fmt.Sprintf("%v", err))
	}

	for _, v := range []string{"Title", "ScreenSize", "SkipBoot", "SavesDir", "DisplayFPS"} {
		if _, ok := initialMap[v]; !ok {
			return nil, ConfigValidationError("Could not find settings key: \"" + v + "\" in settings file")
		}
	}

	//Now parse into struct
	var config Config
	err = json.Unmarshal(file, &config)

	if err != nil {
		return nil, err
	}

	//Perform validations
	err = config.Validate()
	if err != nil {
		return nil, err
	}

	//these are defaults
	config.Debug = *debug
	config.BreakOn = *breakOn
	config.DumpState = *dumpState
	return &config, nil
}

func DefaultConfig() *Config {
	var c *Config = new(Config)
	c.ScreenSize = *screenSizeMultiplier
	c.SkipBoot = *skipBoot
	c.SavesDir = *savesDir
	c.DisplayFPS = *fps
	c.Title = *title
	c.Debug = *debug
	c.BreakOn = *breakOn
	c.DumpState = *dumpState
	return c
}

func (c *Config) String() string {
	return fmt.Sprintln("Configuration settings") +
		fmt.Sprintln(strings.Repeat("-", 50)) +
		fmt.Sprintln(utils.PadRight("Title: ", 19, " "), c.Title) +
		fmt.Sprintln(utils.PadRight("Skip Boot: ", 19, " "), c.SkipBoot) +
		fmt.Sprintln(utils.PadRight("Saves Directory: ", 19, " "), c.SavesDir) +
		fmt.Sprintln(utils.PadRight("Display FPS: ", 19, " "), c.DisplayFPS) +
		fmt.Sprintln(utils.PadRight("Screen Size: ", 19, " "), c.ScreenSize) +
		fmt.Sprintln(utils.PadRight("Debug mode?: ", 19, " "), c.Debug) +
		fmt.Sprintln(utils.PadRight("Breakpoint: ", 19, " "), c.BreakOn) +
		fmt.Sprintln(utils.PadRight("CPU Dump?: ", 19, " "), c.DumpState) +
		fmt.Sprint(strings.Repeat("-", 50))
}

func (c *Config) Validate() error {
	if c.Title == "" {
		return ConfigValidationError("\"Title\" attribute cannot be blank")
	}

	if c.SavesDir == "" {
		return ConfigValidationError("\"SavesDir\" attribute cannot be blank")
	}

	if c.ScreenSize <= 0 || c.ScreenSize > 6 {
		return ConfigValidationError("\"ScreenSize\" attribute must be between 1 and 6")
	}

	return nil
}

func (currentConfig *Config) OverrideConfigWithAnySetFlags() {
	overrideFn := func(f *flag.Flag) {
		log.Println("Overriding configuration in settings file for flag: -" + f.Name)
		switch f.Name {
		case SKIP_BOOT_FLAG:
			currentConfig.SkipBoot = *skipBoot
		case SCREEN_SIZE_FLAG:
			currentConfig.ScreenSize = *screenSizeMultiplier
		case SHOW_FPS_FLAG:
			currentConfig.DisplayFPS = *fps
		case SAVE_DIR_FLAG:
			currentConfig.SavesDir = *savesDir
		case TITLE_FLAG:
			currentConfig.Title = *title
		}
	}
	flag.Visit(overrideFn)
}

func ConfigValidationError(message string) error {
	return errors.New(fmt.Sprintf("Config validation error: %s", message))
}

func ConfigSettingsParseError(message string) error {
	return errors.New(fmt.Sprintf("Config parse error: %s", message))
}
