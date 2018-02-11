package main

/*
 *  Main file for shissue
 *  Copyright (C) 2018 Arthur M
 */

import (
	"fmt"
	"os"
	"strings"
	"strconv"
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
	fmt.Println()
	fmt.Printf(" Usage: %s [options] command [commandargs...]\n", os.Args[0])
	fmt.Println()
	fmt.Println(" Commands: ")

	for _, c := range commands {
		fmt.Printf("\t%-20s %s\n", c.name, c.desc)
	}

	fmt.Println();
	fmt.Println(" Options: ")
	fmt.Println(" -U, --username: specify the username used in your account")
	fmt.Println(" -P, --password: specify the password used in your account")
	fmt.Println()
}


type ArgumentData struct {
	auth *TAuthentication
}

/* Parse the arguments
 * Return the argument index of the subcommand
 */
func parseArgs(ad *ArgumentData) uint {

	commandstart := uint(1)
	for idx, par := range os.Args {
		if (par == "-p" || par == "--username") {
			if (len(os.Args) < int(idx+1)) {
				panic("Username not specified")
			}
			
			ad.auth = new(TAuthentication)
			ad.auth.username = os.Args[idx+1]
			commandstart = uint(idx+2)
		}

		if (par == "-U" || par == "--password") {
			if (len(os.Args) < int(idx+1)) {
				panic("Password not specified")
			}

			if ad.auth == nil {
				panic("Specify username before password")
			}
			
			ad.auth.password = os.Args[idx+1]
			commandstart = uint(idx+2)
		}
	}

	return commandstart
}

func main() {

	commands = append(commands,
		CCommand{name: "help", desc: "Print this help text",
			function: _printHelp},
		CCommand{name: "issues", desc: "List repository issues",
			function: _printIssues},
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


func _printIssues(args []string) {
	printMode := "long"
	if len(args) > 1 {
		printMode = args[1]
	}

	
	r := getRepositoryHost()

	fnBold := func(s string) string {
		return "\033[37;1m" + s + "\033[0m"
	}

	fnBoldYellow := func(s string) string {
		return "\033[33;1m" + s + "\033[0m"
	}

	fnYellow := func(s string) string {
		return "\033[33m" + s + "\033[0m"
	}
	
	fnBoldBlue := func(s string) string {
		return "\033[36;1m" + s + "\033[0m"
	}

	// If arg is a number, it might be the issue number
	if len(args) > 1 {
		if issuen, err := strconv.ParseUint(args[1], 10, 64); err == nil {
			issue, err := r.DownloadIssue(uint(issuen))
			if err != nil {
				panic(err)
			}

			if issue == nil {
				panic("No issue found with that number")
			}

			fmt.Printf("\t#"+fnBold("%d")+" - "+fnBoldYellow("%s")+"\n",
				issue.number, issue.name)
			fmt.Printf("\tCreated by "+fnBoldBlue("%s")+" in %v\n",
				issue.author, issue.creation)
			fmt.Println("\tView it online: " + issue.url)
			fmt.Println()
			fmt.Println(issue.content)
			fmt.Println()

			icomments, err := r.DownloadIssueComments(uint(issuen))
			if err != nil {
				panic(err)
			}

			for _, comment := range icomments {
				fmt.Printf("\t\t comment by "+
					fnYellow("%s")+" in %v\n",
					comment.author, comment.creation)

				contentlines := strings.Split(comment.content, "\n")

				for _, cline := range contentlines {
					fmt.Println("\t\t\t"+cline)
				}
								
				fmt.Println()				
			}
			
			return
		}
	}

	// If not, it might be the type. Download everybody, then!	
	issues, err := r.DownloadAllIssues()
	if err != nil {
		panic(err)
	}


	for _, issue := range issues {
		if printMode == "long" || printMode == "full" {		
			fmt.Printf("\t#"+fnBold("%d")+" - "+fnBoldYellow("%s")+"\n",
				issue.number, issue.name)
			fmt.Printf("\tCreated by "+fnBoldBlue("%s")+" in %v\n",
				issue.author, issue.creation)
			fmt.Println("\tView it online: " + issue.url)
			fmt.Println()
			fmt.Println(issue.content)
			fmt.Println()
			fmt.Println("________________________________________________")
			fmt.Println()
		} else if printMode == "oneline" || printMode == "short" {
			fmt.Printf(" #"+fnBold("%d")+" %s (by "+fnYellow("%s")+")\n",
				issue.number, issue.name, issue.author)
		} else {
			panic("Mode "+printMode+" is unknown. \n"+
				"Try 'long' or 'full' for a complete detail of issues\n"+
				"or try 'oneline' or 'short' for a simple listing, with only name and number")
		}
	}
	
}
