package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type githubUserData struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
}

func fetchGithubUserData(username string) ([]githubUserData, error) {
	url := fmt.Sprintf(
		"https://api.github.com/users/%s/events",
		username,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []githubUserData{}, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return []githubUserData{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > 399 {
		return []githubUserData{}, errors.New("couldn't get response from github")
	}

	var dat []githubUserData
	if err := json.NewDecoder(res.Body).Decode(&dat); err != nil {
		return []githubUserData{}, err
	}

	for _, item := range dat {
		fmt.Printf("  - %s to %s\n", item.Type, item.Repo.Name)
	}

	return dat, nil
}
