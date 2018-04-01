package main

/*
 *  Main file for shissue
 *  Copyright (C) 2018 Arthur M
 */

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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

type ArgumentData struct {
	auth                *TAuthentication
	allowUntrustedCerts bool
}

type CCommandFunc func(ArgumentData, []string)
type CCommand struct {
	name     string
	desc     string
	function CCommandFunc
}

var commands = make([]CCommand, 0)

func printHelp() {
	fmt.Println(" shissue - view github/gitlab issues in command line")
	fmt.Println()
	fmt.Printf(" Usage: %s [options] command [commandargs...]\n", os.Args[0])
	fmt.Println()
	fmt.Println(" Commands: ")

	for _, c := range commands {
		fmt.Printf("\t%-20s %s\n", c.name, c.desc)
	}

	fmt.Println()
	fmt.Println(" Options: ")
	fmt.Println(" [-U|--username] <<username>>\n\tspecify the username used in your repo account")
	fmt.Println(" [-P|--password] <<password>>\n\tspecify the password used in your repo account")
	fmt.Println(" --allow-untrusted-certs\n\tAllow connecting to certificates not trusted by the system")
	fmt.Println()
}

/* Parse the arguments
 * Return the argument index of the subcommand
 */
func parseArgs(ad *ArgumentData) uint {

	commandstart := uint(1)
	for idx, par := range os.Args {
		if par == "--allow-untrusted-certs" {
			ad.allowUntrustedCerts = true
			commandstart = uint(idx + 1)
		}

		if par == "-U" || par == "--username" {
			if len(os.Args) < int(idx+1) {
				panic("Username not specified")
			}

			if ad.auth == nil {
				ad.auth = new(TAuthentication)
			}
			ad.auth.username = os.Args[idx+1]
			commandstart = uint(idx + 2)
		}

		if par == "-P" || par == "--password" {
			if len(os.Args) < int(idx+1) {
				panic("Password not specified")
			}

			if ad.auth == nil {
				panic("Specify username before password")
			}

			ad.auth.password = os.Args[idx+1]
			commandstart = uint(idx + 2)
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

	// Process general parameters
	var ad ArgumentData
	ad.auth = nil

	// Get username and token from git configuration
	username, _ := getGitProperty("shissue.username")
	token, _ := getGitProperty("shissue.token")

	if username != "" || token != "" {
		ad.auth = new(TAuthentication)
		ad.auth.username = username
		ad.auth.token = token
	}

	commandstart := parseArgs(&ad)

	if ad.allowUntrustedCerts == true {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if len(os.Args) <= int(commandstart) {
		// No subcommand called. Print help
		printHelp()
		fmt.Println("\nPlease specify a subcommand")
		return
	}

	if ad.auth != nil && ad.auth.username != "" && ad.auth.password == "" {
		// gets() is made in a java-like way. Congrats, Go!
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Password for %s: ", ad.auth.username)

		/* Disable echo for you to type password, then enable it */
		raw := exec.Command("stty", "-echo")
		raw.Stdin = os.Stdin
		_ = raw.Run()

		pwd, _ := reader.ReadString('\n')
		raw = exec.Command("stty", "echo")
		raw.Stdin = os.Stdin
		_ = raw.Run()

		pwd = strings.Trim(pwd, "\n\r")
		ad.auth.password = pwd
		fmt.Println()
	}

	// Check what command you want
	for _, c := range commands {
		if c.name == os.Args[commandstart] {
			c.function(ad, os.Args[commandstart:])
			return
		}
	}

	fmt.Println("No command named " + os.Args[1])
}

func _printHelp(ad ArgumentData, args []string) {
	printHelp()
}

func _printIssues(ad ArgumentData, args []string) {
	printMode := "long"
	if len(args) > 1 {
		if args[1] == "long" || args[1] == "full" || args[1] == "short" || args[1] == "oneline" {
			printMode = args[1]
		}
	}

	if len(args) > 1 && args[1] == "help" {
		fmt.Println(args[0] + " [full|short|<issue_num>] [filters] ")
		fmt.Println(" Get an issue list ")
		fmt.Println()
		fmt.Println(" filters can be one or more of:")
		fmt.Println(" \tlabels <label1,[label2...]> - Filter by labels")
		fmt.Println(" \tassignee <assignee> - Filter by users that have an issue assigned to them")
		fmt.Println(" \tcreator <creator>  - Filter by issue creators,")
		fmt.Println(" \t[open|closed|all] - Get only open, only closed or all issues")
		fmt.Println()
		return
	}

	r := getRepositoryHost(ad.auth)

	fnBold := func(s string) string {
		return "\033[37;1m" + s + "\033[0m"
	}

	fnBoldYellow := func(s string) string {
		return "\033[33;1m" + s + "\033[0m"
	}

	fnBoldRed := func(s string) string {
		return "\033[31;1m" + s + "\033[0m"
	}

	fnYellow := func(s string) string {
		return "\033[33m" + s + "\033[0m"
	}

	fnBoldBlue := func(s string) string {
		return "\033[36;1m" + s + "\033[0m"
	}

	// Write the string 's' with a background color 'r','g','b'
	// It will convert the color to a 256-color compatible one for printing
	// to the terminal
	// TODO: Check if 256 color is supported
	// TODO: Print directly to the 24-bit color if supported

	fnPrintBackColor := func(s string, r, g, b uint8) string { return s }
	if os.Getenv("COLORTERM") == "truecolor" || os.Getenv("COLORTERM") == "24bit" {
		fnPrintBackColor = func(s string, r, g, b uint8) string {
			// (255 / 51 = 5, the number we have to limit it to convert the
			cR, cG, cB := float32(r/51.0), float32(g/51.0), float32(b/51.0)

			if cR+cG*2.5+cB > 9.0 {
				s = "\033[30m" + s
			}
			return fmt.Sprintf("\033[48;2;%d;%d;%dm%s\033[0m",
				r, g, b, s)
		}
	} else {
		fnPrintBackColor = func(s string, r, g, b uint8) string {
			// (255 / 51 = 5, the number we have to limit it to convert the
			// number to a 256-color compatible one
			cR, cG, cB := r/51, g/51, b/51

			if cR+uint8(float32(cG)*2.5)+cB > 9 {
				s = "\033[30m" + s
			}

			// taken from https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit
			cColorNum := 16 + 36*cR + 6*cG + cB

			return "\033[48;5;" + strconv.Itoa(int(cColorNum)) +
				"m" + s + "\033[0m"
		}
	}

	// If arg is a number, it might be the issue number
	if len(args) > 1 {
		if issuen, err := strconv.ParseUint(args[1], 10, 64); err == nil {
			issue, err := r.DownloadIssue(ad.auth, uint(issuen))
			if err != nil {
				panic(err)
			}

			if issue == nil {
				panic("No issue found with that number")
			}

			slabels := ""
			for _, label := range issue.labels {
				slabels = slabels + " " + fnPrintBackColor(
					" "+label.name+" ",
					label.colorR, label.colorG, label.colorB)
			}

			printIssue := fnBoldYellow
			if issue.is_closed {
				printIssue = fnBoldRed
			}

			fmt.Printf("\t#"+fnBold("%d")+" - "+printIssue("%s")+" %s\n",
				issue.number, issue.name, slabels)
			fmt.Printf("\tCreated by "+fnBoldBlue("%s")+" in %v\n",
				issue.author, issue.creation)

			strissue := "no one"
			if len(issue.assignees) > 0 {
				strissue = fnYellow(strings.Join(issue.assignees,
					", "))
			}
			fmt.Printf("\tAssigned to %s\n", strissue)
			if issue.is_closed {
				fmt.Println("\tThis issue has been closed")
			}

			fmt.Println("\tView it online: " + issue.url)
			fmt.Println()
			fmt.Println(issue.content)
			fmt.Println()

			icomments, err := r.DownloadIssueComments(ad.auth,
				uint(issuen))
			if err != nil {
				panic(err)
			}

			for _, comment := range icomments {
				fmt.Printf("\t\t comment by "+
					fnYellow("%s")+" in %v\n",
					comment.author, comment.creation)

				contentlines := strings.Split(comment.content, "\n")

				for _, cline := range contentlines {
					fmt.Println("\t\t\t" + cline)
				}

				fmt.Println()
			}

			return
		}
	}

	filter := TIssueFilter{
		labels:    nil,
		assignee:  nil,
		getOpen:   true,
		getClosed: false,
		creator:   nil,
	}

	// Create the filter structure
	// Do not need to be done if you want to get a specific issue
	if len(args) > 1 {
		for idx, param := range args[1:] {
			if param == "labels" || param == "label" {
				// Get the labels
				// They are comma-separated values
				if len(args) < idx+1 {
					panic("Label list not specified!")
				}

				labelarr := strings.Split(args[1+idx+1], ",")
				labellist := make([]TIssueLabel, 0, len(labelarr))

				for _, l := range labelarr {
					labellist = append(labellist, TIssueLabel{
						name: strings.Trim(l, " "),
					})
				}

				filter.labels = &labellist
				continue
			}

			if param == "assignee" {
				// Get the assignee
				if len(args) < idx+1 {
					panic("Assignee not specified!")
				}
				filter.assignee = &args[1+idx+1]
				continue
			}

			if param == "creator" {
				// Get the assignee
				if len(args) < idx+1 {
					panic("Creator not specified!")
				}
				filter.creator = &args[1+idx+1]
				continue
			}

			if param == "closed" {
				filter.getClosed = true
				filter.getOpen = false
			}

			if param == "all" {
				filter.getClosed = true
				filter.getOpen = true
			}

		}
	}

	// If not, it might be the type. Download everybody, then!
	issues, err := r.DownloadAllIssues(ad.auth, filter)
	if err != nil {
		panic(err)
	}

	for _, issue := range issues {

		slabels := ""
		for _, label := range issue.labels {
			slabels = slabels + " " + fnPrintBackColor(
				" "+label.name+" ",
				label.colorR, label.colorG, label.colorB)
		}

		printIssue := fnBoldYellow
		if issue.is_closed {
			printIssue = fnBoldRed
		}

		printIssueShort := func(s string) string { return s }
		if issue.is_closed {
			printIssueShort = fnBoldRed
		}

		if printMode == "long" || printMode == "full" {

			fmt.Printf("\t#"+fnBold("%d")+" - "+printIssue("%s")+" %s\n",
				issue.number, issue.name, slabels)
			fmt.Printf("\tCreated by "+fnBoldBlue("%s")+" in %v\n",
				issue.author, issue.creation)

			strissue := "no one"
			if len(issue.assignees) > 0 {
				strissue = fnYellow(strings.Join(issue.assignees,
					", "))
			}
			fmt.Printf("\tAssigned to %s\n", strissue)
			if issue.is_closed {
				fmt.Println("\tThis issue has been closed")
			}
			fmt.Println("\tView it online: " + issue.url)
			fmt.Println()
			fmt.Println(issue.content)
			fmt.Println("\n ")

		} else if printMode == "oneline" || printMode == "short" {
			fmt.Printf(" #"+fnBold("%d")+" "+printIssueShort("%s")+" (by "+fnYellow("%s")+")  %s\n",
				issue.number, issue.name, issue.author, slabels)
		} else {
			panic("Mode " + printMode + " is unknown. \n" +
				"Try 'long' or 'full' for a complete detail of issues\n" +
				"or 'oneline' or 'short' for a simple listing, with only name and number\n" +
				" or try issues <num> to see the issue of number <num>")
		}
	}

}
