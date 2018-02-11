package main

/*
 *  Main file for shissue
 *  Copyright (C) 2018 Arthur Mendes
 */

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
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

type errRepoLimit struct {
	err string
}

func (e *errRepoLimit) Error() string {
	return e.err
}

/* Get the repository from the directory 'dir' */
func getRepository(dir string) (*TRepository, error) {

	// Parse the 'git remote -v' output to get the remote
	// The remote is the remote URL of the repo, almost always the web repo
	bout, err := exec.Command("/bin/sh", "-c", "git -C '"+dir+"' remote -v").Output()
	if err != nil {
		// Git could not find anything
		if err.Error() == "exit status 128" {
			return nil, &errRepoLimit{"No git repository found even in root directory"}
		}

		return nil, err
	}

	// Regexes to discover if we have a git or an ssh repository
	// Doing this here so we don't be bothered to recompile again on
	// every loop iteration
	sshregex := regexp.MustCompile(`git@(.*):(.*)/(.*)`)
	httpregex := regexp.MustCompile(`https://(.*)/(.*)/(.*)`)

	out := string(bout[:len(bout)])
	lines := strings.Split(strings.Trim(out, "\n\t "), "\n")

	for _, line := range lines {
		remote := strings.Split(line, "\t")

		/* Prefer remotes named 'origin'
		 * TODO: Allow specifying other remote names
		 */

		if remote[0] != "origin" {
			continue
		}

		remoteurl := strings.Split(remote[1], " ")

		sshregexres := sshregex.FindAllStringSubmatch(remoteurl[0], -1)
		if sshregexres != nil {
			fmt.Println(sshregexres[0])
			return &TRepository{
				name:    strings.Replace(sshregexres[0][3], ".git", "", -1),
				desc:    "",
				author:  sshregexres[0][2],
				url:     sshregexres[0][0],
				api_url: ""}, nil

		}

		httpregexres := httpregex.FindAllStringSubmatch(remoteurl[0], -1)
		if httpregexres != nil {
			// fuck regexes
			return &TRepository{
				name:    strings.Replace(httpregexres[0][3], ".git", "", -1),
				desc:    "",
				author:  httpregexres[0][2],
				url:     httpregexres[0][0],
				api_url: ""}, nil
		}

	}

	return nil, &errRepoLimit{"This git repository doesn't have a remote"}
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

	/* Get an repository */
	cwd, err := os.Getwd()
	if err != nil {
		os.Stderr.WriteString("Error while getcwd()ing: " + err.Error() + "\n")
		return
	}

	repo, err := getRepository(cwd)
	if err != nil {
		os.Stderr.WriteString("Error while getting the repository: " + err.Error() + "\n")
		return
	}

	/* Get the github information.
	 * TODO: support gitlab, bitbucket...
	 */
	var r TGitHubRepo
	_, err = r.Initialize(repo)
	if err != nil {
		panic(err)
	}

	issues, err := r.DownloadAllIssues()
	if err != nil {
		panic(err)
	}

	fnBold := func(s string) string {
		return "\033[37;1m" + s + "\033[0m"
	}

	fnBoldYellow := func(s string) string {
		return "\033[33;1m" + s + "\033[0m"
	}

	fnBoldBlue := func(s string) string {
		return "\033[36;1m" + s + "\033[0m"
	}

	for _, issue := range issues {
		fmt.Printf("\t#"+fnBold("%d")+" - "+fnBoldYellow("%s")+"\n",
			issue.number, issue.name)
		fmt.Printf("\tCreated by "+fnBoldBlue("%s")+" on %v\n",
			issue.author, issue.creation)
		fmt.Println("\tView it online: " + issue.url)
		fmt.Println()
		fmt.Println(issue.content)
		fmt.Println()
		fmt.Println("________________________________________________")
		fmt.Println()
	}
}

func _printHelp(args []string) {
	printHelp()
}
