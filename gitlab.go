package main

import (
	"github.com/xanzy/go-gitlab"
	"strconv"
)

/**
 * Gitlab implementation for the generic repository
 *
 * Copyright (C) 2018 Arthur M
 */

// Handler to a gitlab repo
type TGitLabRepo struct {
	client  *gitlab.Client
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
	if auth != nil && auth.token != "" {
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

/* Get a map with all labels and the hex colors used in this project */
func (gl *TGitLabRepo) getLabels() (map[string]string, error) {
	// Get the label colors. For the visuals!
	var labelColors = map[string]string{}
	labels, _, err := gl.client.Labels.ListLabels(gl.project.ID)

	if err != nil {
		return nil, err
	}

	for _, l := range labels {
		labelColors[l.Name] = l.Color
	}

	return labelColors, nil
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

	// Create filter
	if filter.labels != nil {
		glabels := make([]string, 0, len(*filter.labels))
		for _, label := range *filter.labels {
			glabels = append(glabels, label.name)
		}

		goptions.Labels = glabels
	}

	gstate := ""
	if filter.getOpen && !filter.getClosed {
		gstate = "opened"
	} else if !filter.getOpen && filter.getClosed {
		gstate = "closed"
	}

	if gstate != "" {
		goptions.State = &gstate
	}

	labelColors, err := gl.getLabels()
	if err != nil {
		return nil, err
	}

	// Do the request

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
			// Gitlab returned label colors have a "#" prefix,
			// like in '#ff0000'. We need to take it out
			lcolor := labelColors[label][1:]
			cR, _ := strconv.ParseUint(lcolor[0:2], 16, 8)
			cG, _ := strconv.ParseUint(lcolor[2:4], 16, 8)
			cB, _ := strconv.ParseUint(lcolor[4:6], 16, 8)

			labels = append(labels, TIssueLabel{name: label,
				colorR: uint8(cR),
				colorG: uint8(cG),
				colorB: uint8(cB)})
		}

		issues = append(issues, TIssue{
			id:        uint(issue.ID),
			number:    uint(issue.IID),
			name:      issue.Title,
			url:       issue.WebURL,
			author:    issue.Author.Name,
			assignees: assignees,
			labels:    labels,
			creation:  *issue.CreatedAt,
			content:   issue.Description,
			is_closed: (issue.State == "closed"),
		})
	}

	return issues, nil
}

/* Download an specific issue by ID,
 *
* In Gitlab, the ID used to return the issue is the databse ID.
* The ID in this parameter is the issue number, what gitlab calls 'iid'
*/
func (gl *TGitLabRepo) DownloadIssue(auth *TAuthentication, id uint) (*TIssue, error) {

	var goptions gitlab.ListProjectIssuesOptions
	goptions.Page = 1
	goptions.PerPage = 1000
	goptions.IIDs = append(make([]int, 0, 1), int(id))

	glissues, _, err := gl.client.Issues.ListProjectIssues(gl.project.ID,
		&goptions)

	if err != nil {
		return nil, err
	}

	if len(glissues) == 0 {
		return nil, nil // No issues do not mean error
	}

	labelColors, err := gl.getLabels()
	if err != nil {
		return nil, err
	}

	glissue := glissues[0]

	assignees := make([]string, 0, len(glissue.Assignees))
	for _, assignee := range glissue.Assignees {
		assignees = append(assignees, assignee.Name)
	}

	labels := make([]TIssueLabel, 0, len(glissue.Labels))
	for _, label := range glissue.Labels {
		// Gitlab returned label colors have a "#" prefix,
		// like in '#ff0000'. We need to take it out
		lcolor := labelColors[label][1:]
		cR, _ := strconv.ParseUint(lcolor[0:2], 16, 8)
		cG, _ := strconv.ParseUint(lcolor[2:4], 16, 8)
		cB, _ := strconv.ParseUint(lcolor[4:6], 16, 8)

		labels = append(labels, TIssueLabel{name: label,
			colorR: uint8(cR),
			colorG: uint8(cG),
			colorB: uint8(cB)})
	}

	issue := &TIssue{
		id:        uint(glissue.ID),
		number:    uint(glissue.IID),
		name:      glissue.Title,
		url:       glissue.WebURL,
		author:    glissue.Author.Name,
		assignees: assignees,
		labels:    labels,
		creation:  *glissue.CreatedAt,
		content:   glissue.Description,
		is_closed: (glissue.State == "closed"),
	}

	return issue, nil
}

/* Download all comments from that issue */
func (gl *TGitLabRepo) DownloadIssueComments(auth *TAuthentication, issue_id uint) ([]TIssueComment, error) {

	return nil, nil
}
