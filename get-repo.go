/**
 * Functions used for repository and host identification
 * 
 * They identify if the directory is a Git repository, and what host they are
 * 
 * Currently only github is supported
 *
 * Copyright (C) 2018 Arthur M
 */

package main

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

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

/* Gets the correct repository host, based in the remote data
 * 'auth' is an authentication object, for the cases we need to authenticate
 * to even see the repository (e.g private repos)
 *
 * Panics if you can't get it, but it doesn't matter. You wouldn't be able to do
 * nothing if it didn't fail...
 */
func getRepositoryHost(auth *TAuthentication) *TGitHubRepo {
	/* Get an repository */
	cwd, err := os.Getwd()
	if err != nil {
		panic("Error while getcwd()ing: " + err.Error() + "\n")
	}

	repo, err := getRepository(cwd)
	if err != nil {
		panic("Error while getting the repository: " + err.Error() + "\n")
	}

	/* Get the github information.
	 * TODO: support gitlab, bitbucket...
	 */
	r := new(TGitHubRepo)
	_, err = r.Initialize(auth, repo)
	if err != nil {
		panic(err)
	}

	return r

}
