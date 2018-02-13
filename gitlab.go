package main

import (
	"github.com/xanzy/go-gitlab"
)

/**
 * Gitlab implementation for the generic repository
 *
 * Copyright (C) 2018 Arthur M
 */

// Handler to a gitlab repo
type TGitLabRepo struct {
	client *gitlab.Client
	project *gitlab.Project
}

/* "Initialize" the host, with info from the repository
 * This is used to setup the URLs related to that repo
 *
 * Return nil on success, together with the 'api_url' string.
 * It needs to fill all fields of the 'repo' structure
 * Returns an error object on error
 */
func (gl *TGitLabRepo) Initialize(auth *TAuthentication, repo *TRepository) (string, error) {

	token := ""
	if (auth != nil && auth.token != "") {
		token = auth.token
	}

	git := gitlab.NewClient(nil, token)

	project, _, err := git.Projects.GetProject(repo.author + "/" + repo.name)
	if err != nil {
		return "", err
	}

	gl.client = git
	gl.project = project

	author := ""
	if project.Owner != nil {
		author = project.Owner.Name
	}
	
	repo.name = project.Name
	repo.desc = project.Description
	repo.author = author
	repo.url = project.WebURL
	repo.api_url = project.SSHURLToRepo // sort of an api url

	return project.SSHURLToRepo, nil
}

/* Download all issues from the repository
 * You can use the TAuthentication struct to pass authentication info
 * Send it nil for no authentication, but take note that the host
 * might not send everything to unauthenticated users
 *
 * Return nil on the issue list and on the error if no issues exist.
 * Return an issue list on success, or nil on issue list and an error
 * on error
 */
func (gl *TGitLabRepo) DownloadAllIssues(auth *TAuthentication, filter TIssueFilter) ([]TIssue, error) {

	var goptions gitlab.ListProjectIssuesOptions
	goptions.Page = 1
	goptions.PerPage = 1000
	
	glissues, _, err := gl.client.Issues.ListProjectIssues(gl.project.ID,
		&goptions)

	if err != nil {
		return nil, err
	}

	issues := make([]TIssue, 0, len(glissues))

	if issues == nil {
		return nil, nil
	}
	
	for _, issue := range glissues {

		assignees := make([]string, 0, len(issue.Assignees))
		for _, assignee := range issue.Assignees {
			assignees = append(assignees, assignee.Name)
		}

		labels := make([]TIssueLabel, 0, len(issue.Labels))
		for _, label := range issue.Labels {
			labels = append(labels, TIssueLabel{name: label,
				colorR: uint8(255),
				colorG: uint8(255),
				colorB: uint8(255)})
		}
			
		
		issues = append(issues, TIssue{
			id: uint(issue.ID),
			number: uint(issue.IID),
			name: issue.Title,
			url: issue.WebURL,
			author: issue.Author.Name,
			assignees: assignees,
			labels: labels,
			creation: *issue.CreatedAt,
			content: issue.Description,
			is_closed: (issue.State == "closed"),
		})
	}

	return issues, nil
}

/* Download an specific issue by ID,
 */
func (gl *TGitLabRepo) DownloadIssue(auth *TAuthentication, id uint) (*TIssue, error) {

	return nil, nil
}

/* Download all comments from that issue */
func (gl *TGitLabRepo) DownloadIssueComments(auth *TAuthentication, issue_id uint) ([]TIssueComment, error) {

	return nil, nil
}
