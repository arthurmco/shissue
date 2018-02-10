package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

/**
 * Github implementation for the generic repository
 *
 * Copyright (C) 2018 Arthur M
 */

type TGitHubOwner struct {
	Login string
	ID    uint
}

type TGitHubRepo struct {
	ID          uint
	Name        string
	Full_name   string
	Owner       TGitHubOwner
	Description string

	Issues_url        string
	Issue_comment_url string
}

func (gh *TGitHubRepo) Initialize(repo *TRepository) (string, error) {

	// We need to get the api URL from the repository URL
	// URL is https://github.com/arthurmco/clinancial
	// API url is https://api.github.com/repos/arthurmco/clinancial

	username := repo.author
	reponame := repo.name

	api_url := "https://api.github.com/repos/" + username + "/" + reponame

	// Now we need to download this.
	resp, err := http.Get("https://api.github.com/repos/arthurmco/clinancial")

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var ghrepo TGitHubRepo
	err = json.Unmarshal(body, &ghrepo)
	if err != nil {
		return "", err
	}

	repo.name = ghrepo.Name
	repo.desc = ghrepo.Description

	return api_url, nil

}
