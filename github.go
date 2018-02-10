package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"strings"
)

/**
 * Github implementation for the generic repository
 *
 * Copyright (C) 2018 Arthur M
 */

/* Structures that represent the used fields from the JSON returned by the
 * Github API.
 *
 * TGitHubOwner is the repository owner information
 * TGitHubRepo is the repository information
 */
type TGitHubUser struct {
	Login string
	ID    uint
}

type TGitHubRepo struct {
	ID          uint
	Name        string
	Full_name   string
	Owner       TGitHubUser
	Description string

	Issues_url        string
	Issue_comment_url string

	Has_issues bool
}

type TGitHubIssue struct {
	ID         uint
	Number     uint
	Title      string
	User       TGitHubUser
	Assignees  []TGitHubUser
	Html_url   string
	State      string
	Created_at time.Time
	Body       string
}

type errGithubRepo struct {
	err string
}

func (e *errGithubRepo) Error() string {
	return e.err
}

func (gh *TGitHubRepo) Initialize(repo *TRepository) (string, error) {

	// We need to get the api URL from the repository URL
	// URL is https://github.com/arthurmco/clinancial
	// API url is https://api.github.com/repos/arthurmco/clinancial

	username := repo.author
	reponame := repo.name

	api_url := "https://api.github.com/repos/" + username + "/" + reponame

	// Now we need to download this.
	resp, err := http.Get(api_url)

	if err != nil {
		return "", err
	}

	if resp.StatusCode == 404 {
		return "", &errGithubRepo{"Repository not found!"}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, gh)
	if err != nil {
		return "", err
	}

	repo.name = gh.Name
	repo.desc = gh.Description

	return api_url, nil

}

/*
 * Download all issues from this github repository
 */
func (gh *TGitHubRepo) DownloadAllIssues() ([]TIssue, error) {
	if !gh.Has_issues {
		// Return a nil list, since this repository doesn't has issues
		// Return no errors too, since no error has been found
		return nil, nil
	}
	
	issue_url := strings.Replace(gh.Issues_url, "{/number}", "", 1)

	resp, err := http.Get(issue_url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var ghissues []TGitHubIssue
	err = json.Unmarshal(body, &ghissues)
	if err != nil {
		return nil, err
	}

	issues := make([]TIssue, len(ghissues))

	for idx, ghissue := range ghissues {
		issues[idx].id = ghissue.ID
		issues[idx].number = ghissue.Number
		issues[idx].name = ghissue.Title
		issues[idx].url = ghissue.Html_url
		issues[idx].author = ghissue.User.Login
		issues[idx].creation = ghissue.Created_at
		issues[idx].content = ghissue.Body
	}

	return issues, nil
}
