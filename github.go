package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
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

type TGitHubIssueLabel struct {
	ID    uint
	Name  string
	Color string
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
	Labels     []TGitHubIssueLabel
}

type TGitHubIssueComment struct {
	ID         uint
	Html_url   string
	Body       string
	User       TGitHubUser
	Created_at time.Time
}

type errGithubRepo struct {
	err string
}

func (e *errGithubRepo) Error() string {
	return e.err
}

func (gh *TGitHubRepo) Initialize(auth *TAuthentication, repo *TRepository) (string, error) {

	// We need to get the api URL from the repository URL
	// URL is https://github.com/arthurmco/clinancial
	// API url is https://api.github.com/repos/arthurmco/clinancial

	username := repo.author
	reponame := repo.name

	api_url := "https://api.github.com/repos/" + username + "/" + reponame

	// Now we need to download this.
	resp, err := gh.buildGetRequest(api_url, auth, "")

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

/* Build and send a common GET request to the API.
 * 'params' is an "HTTP GET"-like parameter string, without the '?'
 * Return the response object on success, or an error.
 */
func (gh *TGitHubRepo) buildGetRequest(url string, auth *TAuthentication, params string) (*http.Response, error) {

	client := http.Client{}

	// Replace spaces with HTTP-allowed spaces and + with HTTP-blessed ones
	params = strings.Replace(params, " ", "%20", -1)
	params = strings.Replace(params, "%", "%2B", -1)

	// Build the request, and then do it
	req, err := http.NewRequest("GET", url+"?per_page=100&"+params, nil)
	if err != nil {
		return nil, err
	}

	// Only send authorization data when we have an username
	if auth != nil && auth.username != "" {
		req.SetBasicAuth(auth.username, auth.password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 {
		return nil, &errGithubRepo{"Authentication failed: wrong username and/or password"}
	}

	return resp, nil
}

/* Download a range of issues based on a certain filter, return the amount of what are really issues
 * Because Github API doesn't differentiate issues from pull requests, this is needed
 *
 * It fills the 'issuelist' with the found issues
 *
 * The maximum allowed by the API is 100, and it downloads 100 by 199, so count+start needs to be less
 * than 100.
 */
func (gh *TGitHubRepo) downloadIssueRange(start, count, page int,
	auth *TAuthentication, 	url, params string,
	issuelist *[]TGitHubIssue) (int, error) {

	params = params + "&page="+strconv.Itoa(page)

	resp, err := gh.buildGetRequest(url, auth, params)

	if err != nil {
		return 0, err
	}
	
	// Read the result and build the JSON
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var ghissues []TGitHubIssue
	err = json.Unmarshal(body, &ghissues)

	if len(ghissues) < start {
		// No issues to return
		return 0, nil
	}

	// TODO: Differentiate issues from pull requests

	count = len(ghissues)
	icount := 0
	for _, iss := range ghissues[start:(count-start)] {
		icount += 1
		*issuelist = append(*issuelist, iss)
	}

	return icount, nil
}

/*
 * Download all issues from this github repository
 * You can use the TAuthentication struct to pass authentication info
 * Send it nil for no authentication, but take note that the host
 * might not send everything to unauthenticated users
 *
 * Return nil on the issue list and on the error if no issues exist.
 * Return an issue list on success, or nil on issue list and an error
 * on error
 */
func (gh *TGitHubRepo) DownloadAllIssues(auth *TAuthentication, filter TIssueFilter) ([]TIssue, error) {
	if !gh.Has_issues {
		// Return a nil list, since this repository doesn't has issues
		// Return no errors too, since no error has been found
		return nil, nil
	}

	issue_url := strings.Replace(gh.Issues_url, "{/number}", "", 1)

	// Build filters
	paramstr := make([]string, 0, 4)
	if filter.labels != nil {
		// Build label parameter filter
		// They need to be comma-separated
		labelarr := make([]string, 0)
		for _, l := range *filter.labels {
			labelarr = append(labelarr, l.name)
		}

		paramstr = append(paramstr, "labels="+strings.Join(
			labelarr, ","))

	}

	if filter.assignee != nil {
		paramstr = append(paramstr, "assignee="+*filter.assignee)
	}

	if filter.getOpen && filter.getClosed {
		paramstr = append(paramstr, "state=all")
	} else if !filter.getOpen && filter.getClosed {
		paramstr = append(paramstr, "state=closed")
	} else {
		paramstr = append(paramstr, "state=open")
	}

	if filter.creator != nil {
		paramstr = append(paramstr, "creator="+*filter.creator)
	}


	var ghissues []TGitHubIssue
	imin, imax := 0, 1000
	icurr := imin
	ioff := icurr
	pagen := 1
	const pagecount = 100

	params := strings.Join(paramstr, "&")

	for icurr < imax {
		var pageissues []TGitHubIssue
		iret, err := gh.downloadIssueRange(ioff, pagecount, pagen, auth,
			issue_url, params, &pageissues)

		if err != nil {
			return nil, err
		}

		// No pages
		if iret == 0 {
			break
		}

		ghissues = append(ghissues, pageissues...)

		// Didn't reached the page count
		if iret < pagecount {
			break
		}

		pagen += 1
		icurr += pagecount
		ioff = 0
	}
	
	
	

	issues := make([]TIssue, len(ghissues))

	for idx, ghissue := range ghissues {
		issues[idx].id = ghissue.ID
		issues[idx].number = ghissue.Number
		issues[idx].name = ghissue.Title
		issues[idx].url = ghissue.Html_url
		issues[idx].author = ghissue.User.Login

		assignees := make([]string, 0)
		for _, assignee := range ghissue.Assignees {
			assignees = append(assignees, assignee.Login)
		}
		issues[idx].assignees = assignees

		labels := make([]TIssueLabel, 0)
		for _, ghlabel := range ghissue.Labels {
			colorR, _ := strconv.ParseUint(ghlabel.Color[0:2], 16, 8)
			colorG, _ := strconv.ParseUint(ghlabel.Color[2:4], 16, 8)
			colorB, _ := strconv.ParseUint(ghlabel.Color[4:6], 16, 8)

			labels = append(labels, TIssueLabel{
				name:   ghlabel.Name,
				colorR: uint8(colorR),
				colorG: uint8(colorG),
				colorB: uint8(colorB),
			})
		}
		issues[idx].labels = labels

		issues[idx].creation = ghissue.Created_at
		issues[idx].content = ghissue.Body
		issues[idx].is_closed = (ghissue.State == "closed")
	}

	return issues, nil
}

/*
 * Download an specific issue from this github repository
 */
func (gh *TGitHubRepo) DownloadIssue(auth *TAuthentication, id uint) (*TIssue, error) {
	if !gh.Has_issues {
		// Return a nil list, since this repository doesn't has issues
		// Return no errors too, since no error has been found
		return nil, nil
	}

	issue_url := strings.Replace(gh.Issues_url, "{/number}",
		"/"+strconv.Itoa(int(id)), 1)

	resp, err := gh.buildGetRequest(issue_url, auth, "")
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

	assignees := make([]string, 0)
	for _, assignee := range ghissue.Assignees {
		assignees = append(assignees, assignee.Login)
	}
	issue.assignees = assignees

	labels := make([]TIssueLabel, 0)
	for _, ghlabel := range ghissue.Labels {
		colorR, _ := strconv.ParseUint(ghlabel.Color[0:2], 16, 8)
		colorG, _ := strconv.ParseUint(ghlabel.Color[2:4], 16, 8)
		colorB, _ := strconv.ParseUint(ghlabel.Color[4:6], 16, 8)

		labels = append(labels, TIssueLabel{
			name:   ghlabel.Name,
			colorR: uint8(colorR),
			colorG: uint8(colorG),
			colorB: uint8(colorB),
		})
	}
	issue.labels = labels

	issue.creation = ghissue.Created_at
	issue.content = ghissue.Body
	issue.is_closed = (ghissue.State == "closed")

	return issue, nil
}

/* Download all comments from that issue */
func (gh *TGitHubRepo) DownloadIssueComments(auth *TAuthentication, issue_id uint) ([]TIssueComment, error) {

	if !gh.Has_issues {
		// Return a nil list, since this repository doesn't has issues
		// Return no errors too, since no error has been found
		return nil, nil
	}

	comment_url := strings.Replace(gh.Issues_url, "{/number}",
		"/"+strconv.Itoa(int(issue_id))+"/comments", 1)

	resp, err := gh.buildGetRequest(comment_url, auth, "")
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
