// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type GithubProjects struct {
	Projects []*Project `json:"projects"`
}

type Project struct {
	Name           string   `json:"name"`
	Githublink     string   `json:"githublink"`
	Description    string   `json:"description"`
	Createdon      string   `json:"createdon"`
	Languages      []Tag    `json:"languages"`
	Topics         []string `json:"topics"`
	Deploymentlink string   `json:"deploymentlink"`
}

type Tag struct {
	Language string `json:"language"`
}
