package main

type githubUserData struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Ref     string `json:"ref"`
		RefType string `json:"ref_type"`
		Commits []struct {
			Message string `json:"message"`
		} `json:"commits"`
		Action string `json:"action"`
		Forkee struct {
			FullName string `json:"full_name"`
			Owner    struct {
				Login string `json:"login"`
			} `json:"owner"`
		} `json:"forkee"`
		Issue struct {
			Title string `json:"title"`
		} `json:"issue"`
		PullReq struct {
			Title string `json:"title"`
		} `json:"pull_request"`
		Member struct {
			Login string `json:"login"`
		} `json:"member"`
		Release struct {
			Name string `json:"name"`
		} `json:"release"`
	} `json:"payload"`
}
