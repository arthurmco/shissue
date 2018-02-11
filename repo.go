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

type TRepository struct {
	name   string // Repository name
	desc   string // Repository description
	author string // Repository author

	url string // Repository external URL

	api_url string     // Repository 'api' URL
	host    *TRepoHost // Pointer to the repository host
}

/* The issue itself */
type TIssue struct {
	// Issue ID, in the repository host API
	//	(only meaningful to the repohost interface)
	id uint

	number   uint      // Issue number
	name     string    // Issue name
	url      string    // Issue URL, to view it online
	author   string    // Issue author
	creation time.Time // Issue creation date

	content string // Issue content
}

/* Issue comment.
 * They might be as important as the issue itself, because they contain additional info and made
 * decisions
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

type TRepoHost interface {

	/* "Initialize" the host, with info from the repository
	 * This is used to setup the URLs related to that repo
	 *
	 * Return nil on success, together with the 'api_url' string.
	 * It needs to fill all fields of the 'repo' structure
	 * Returns an error object on error
	 */
	Initialize(repo *TRepository) (string, error)

	/* Download all issues from the repository
	 * Return nil on the issue list and on the error if no issues exist.
	 * Return an issue list on success, or nil on issue list and an error
	 * on error
	 */
	DownloadAllIssues() ([]TIssue, error)

	/* Download an specific issue by ID,
	 */
	DownloadIssue(id uint) (*TIssue, error)

	/* Download all comments from that issue */
	DownloadIssueComments(issue_id uint) ([]TIssueComment, error)
}
