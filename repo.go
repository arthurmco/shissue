package main

/**
 * Basic things for a generic repository
 *
 * Contains an interface for a generic repository host, a repository description
 * and an issue description
 *
 * Copyright (C) 2018 Arthur M
 */

import (
	"time"
)

/* Authentication data, for when you need your users permission
 * Typically when you're modifying or adding something, like an issue
 * Or when you're been rate-limited.
 */
type TAuthentication struct {
	username string
	password string
	token string
}

type TRepository struct {
	name   string // Repository name
	desc   string // Repository description
	author string // Repository author

	url string // Repository external URL
	base_url string // Base URL
	api_url string     // Repository 'api' URL
	host    *TRepoHost // Pointer to the repository host
}

type TIssueLabel struct {
	name string // Issue label name

	colorR, colorG, colorB uint8 // Color data, only for decoration
}

/* The issue itself */
type TIssue struct {
	// Issue ID, in the repository host API
	//	(only meaningful to the repohost interface)
	id uint

	number    uint          // Issue number
	name      string        // Issue name
	url       string        // Issue URL, to view it online
	author    string        // Issue author
	assignees []string      // People assigned with the issue
	labels    []TIssueLabel // Issue labels

	creation time.Time // Issue creation date

	content string // Issue content

	is_closed bool // Is the issue closed?
}

/* Error type that happened when you couldn't connect to your repository */
type RepoConnectError struct {
	err string
	ErrorCode int
}

func (e *RepoConnectError) Error() string {
	return e.err
}


/* Issue comment.
 * They might be as important as the issue itself, because they contain additional info and
 * decisions made that can change the issue meaning
 *
 * So it's good to include them
 */
type TIssueComment struct {
	id       uint      // Comment ID, for the repo host
	url      string    // Comment URL
	author   string    // Comment author
	creation time.Time // Comment creation date

	content string // Comment content
}

/* Filter for the issue query
 * If any of the issue query filters are 'null', it means that it shouldn't be
 * considered
 */
type TIssueFilter struct {
	labels   *[]TIssueLabel // Label list, 'nil' if no label filter
	assignee *string // Only get issues assigned to 'assignee'
	getOpen, getClosed bool // True if you want to get open or closed issues
	creator *string // Only get issues made by 'creator'
}

type TRepoHost interface {

	/* "Initialize" the host, with info from the repository
	 * This is used to setup the URLs related to that repo
	 *
	 * Return nil on success, together with the 'api_url' string.
	 * It needs to fill all fields of the 'repo' structure
	 * Returns an error object on error
	 */
	Initialize(auth *TAuthentication, repo *TRepository) (string, error)

	/* Download all issues from the repository
	 * You can use the TAuthentication struct to pass authentication info
	 * Send it nil for no authentication, but take note that the host
	 * might not send everything to unauthenticated users
	 *
	 * Return nil on the issue list and on the error if no issues exist.
	 * Return an issue list on success, or nil on issue list and an error
	 * on error
	 */
	DownloadAllIssues(auth *TAuthentication, filter TIssueFilter) ([]TIssue, error)

	/* Download an specific issue by ID,
	 */
	DownloadIssue(auth *TAuthentication, id uint) (*TIssue, error)

	/* Download all comments from that issue */
	DownloadIssueComments(auth *TAuthentication, issue_id uint) ([]TIssueComment, error)
}
