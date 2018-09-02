package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/djhworld/gomeboycolor/utils"
)

type Config struct {
	//mandatory settings
	Title      string
	ScreenSize int
	SkipBoot   bool
	DisplayFPS bool
	ColorMode  bool

	//optional
	Headless  bool
	Debug     bool
	BreakOn   string
	DumpState bool
}

func (c *Config) String() string {
	return fmt.Sprintln("Configuration settings") +
		fmt.Sprintln(strings.Repeat("-", 50)) +
		fmt.Sprintln(utils.PadRight("Title: ", 19, " "), c.Title) +
		fmt.Sprintln(utils.PadRight("Skip Boot: ", 19, " "), c.SkipBoot) +
		fmt.Sprintln(utils.PadRight("GB Color Mode: ", 19, " "), c.ColorMode) +
		fmt.Sprintln(utils.PadRight("Display FPS: ", 19, " "), c.DisplayFPS) +
		fmt.Sprintln(utils.PadRight("Screen Size: ", 19, " "), c.ScreenSize) +
		fmt.Sprintln(utils.PadRight("Debug mode?: ", 19, " "), c.Debug) +
		fmt.Sprintln(utils.PadRight("Breakpoint: ", 19, " "), c.BreakOn) +
		fmt.Sprintln(utils.PadRight("CPU Dump?: ", 19, " "), c.DumpState) +
		fmt.Sprintln(utils.PadRight("Headless: ", 19, " "), c.Headless) +
		fmt.Sprint(strings.Repeat("-", 50))
}

func (c *Config) Validate() error {
	if c.Title == "" {
		return ConfigValidationError("\"Title\" attribute cannot be blank")
	}

	if c.ScreenSize <= 0 || c.ScreenSize > 6 {
		return ConfigValidationError("\"ScreenSize\" attribute must be between 1 and 6")
	}

	return nil
}

func ConfigValidationError(message string) error {
	return errors.New(fmt.Sprintf("Config validation error: %s", message))
}
