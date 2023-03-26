// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Article struct {
	Title       string       `json:"title"`
	TitleCard   string       `json:"titleCard"`
	Author      *Author      `json:"author"`
	ContentData string       `json:"contentData"`
	DateWritten string       `json:"dateWritten"`
	URL         string       `json:"url"`
	Description string       `json:"description"`
	UUID        string       `json:"uuid"`
	Tags        []ArticleTag `json:"tags"`
}

type ArticleTag struct {
	Tag string `json:"tag"`
}

type Articles struct {
	ArticleCollection []Article `json:"articleCollection"`
	Total             int       `json:"total"`
}

type Author struct {
	Name    string `json:"name"`
	Profile string `json:"profile"`
	Picture string `json:"picture"`
}

type GithubBio struct {
	Position string `json:"position"`
	Company  string `json:"company"`
	Readme   string `json:"readme"`
}

type GithubProjects struct {
	Projects         []*Project `json:"projects"`
	ContributorCount int        `json:"contributorCount"`
}

type NotionGoals struct {
	Goals []string `json:"goals"`
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

type SearchInput struct {
	Term string `json:"term"`
}
