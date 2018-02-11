package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"strconv"
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

type TGitHubIssueComment struct {
	ID uint
	Html_url string
	Body string
	User TGitHubUser
	Created_at time.Time
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

	if resp.StatusCode == 403 {
		if resp.Header.Get("X-RateLimit-Remaining") == "0" {
			return "", &errGithubRepo{"Github API rate limit exceeded"}
		} else {
			return "", &errGithubRepo{"Permission error!"}
		}
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

/*
 * Download an specific issue from this github repository
 */
func (gh *TGitHubRepo) DownloadIssue(id uint) (*TIssue, error) {
	if !gh.Has_issues {
		// Return a nil list, since this repository doesn't has issues
		// Return no errors too, since no error has been found
		return nil, nil
	}

	issue_url := strings.Replace(gh.Issues_url, "{/number}",
		"/"+strconv.Itoa(int(id)), 1)

	resp, err := http.Get(issue_url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var ghissue TGitHubIssue
	err = json.Unmarshal(body, &ghissue)
	if err != nil {
		return nil, err
	}

	issue := new(TIssue)

	issue.id = ghissue.ID
	issue.number = ghissue.Number
	issue.name = ghissue.Title
	issue.url = ghissue.Html_url
	issue.author = ghissue.User.Login
	issue.creation = ghissue.Created_at
	issue.content = ghissue.Body

	return issue, nil
}

/* Download all comments from that issue */
func (gh *TGitHubRepo) DownloadIssueComments(issue_id uint) ([]TIssueComment, error) {
	
	if !gh.Has_issues {
		// Return a nil list, since this repository doesn't has issues
		// Return no errors too, since no error has been found
		return nil, nil
	}

	comment_url := strings.Replace(gh.Issues_url, "{/number}",
		"/"+strconv.Itoa(int(issue_id))+"/comments", 1)

	resp, err := http.Get(comment_url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var ghcomments []TGitHubIssueComment
	err = json.Unmarshal(body, &ghcomments)
	if err != nil {
		return nil, err
	}

	comments := make([]TIssueComment, len(ghcomments))

	for idx, ghissue := range ghcomments {
		comments[idx].id = ghissue.ID
		comments[idx].url = ghissue.Html_url
		comments[idx].author = ghissue.User.Login
		comments[idx].creation = ghissue.Created_at
		comments[idx].content = ghissue.Body
	}

	return comments, nil
}

