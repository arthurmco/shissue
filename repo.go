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

	api_url string // Repository 'api' URL
}

type TIssue struct {
	id       uint      // Issue ID
	name     string    // Issue name
	author   string    // Issue author
	creation time.Time // Issue creation date

	content string // Issue content
}

type TRepoHost interface {

	/* "Initialize" the host, with info from the repository
	 * This is used to setup the URLs related to that repo
	 *
	 * Return nil on success, together with the 'api_url' string.
	 * Returns an error object on error
	 */
	Initialize(repo TRepository) (string, error)

	/* Download all issues from the repository */
	DownloadAllIssues() ([]TIssue, error)

	/* Download an specific issue range by ID, including the
	 * first up to 'last-1'
	 */
	DownloadIssueRange(first uint, last uint) ([]TIssue, error)
}
