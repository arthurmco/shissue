package main

/*
 *  Main file for shissue
 *  Copyright (C) 2018 Arthur Mendes
 */

import (
	"fmt"
	"os"
)

/**
 * shissue, just like Git, has a series of subcommands that
 * can be called after the executable name
 *
 * This structure lists them
 * 'name' is the command name, 'desc' is the help text that will show, and 'function' is a
 * function pointer that the software will call when you specificate that command
 */

type CCommandFunc func([]string)
type CCommand struct {
	name     string
	desc     string
	function CCommandFunc
}

var commands = make([]CCommand, 0)

func printHelp() {
	fmt.Println(" shissue - view github issues in command line")
	fmt.Println("")
	fmt.Printf(" Usage: %s [command] [commandargs...]\n", os.Args[0])
	fmt.Println("")
	fmt.Println(" Commands: ")

	for _, c := range commands {
		fmt.Printf("\t%-20s %s\n", c.name, c.desc)
	}
}

func main() {

	commands = append(commands,
		CCommand{name: "help", desc: "Print this help text",
			function: _printHelp},
	)

	if len(os.Args) <= 1 {
		// No subcommand called. Print help
		printHelp()
		return
	}

	// Check what command you want
	for _, c := range commands {
		if c.name == os.Args[1] {
			c.function(os.Args[1:])
			return
		}
	}

	fmt.Println("No command named " + os.Args[1])
}

func _printHelp(args []string) {
	printHelp()
}
