package routes

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"github.com/zenith110/portfilo/graph/model"
	"golang.org/x/oauth2"
)

/*
Given a language url, loop through the json contents and append it into a tag struct
*/
func ParseLanguages(languages []struct{ Name string }) []model.Tag {
	var tags []model.Tag
	for language := range languages {
		var tag model.Tag
		tag.Language = languages[language].Name
		tags = append(tags, tag)
	}
	return tags
}

func ParseTopics(topicsList []struct{ Topic struct{ Name string } }, contributorCount int) []string {
	var topics []string
	for topicIndex := 0; topicIndex < len(topicsList); topicIndex++ {
		if topicsList[topicIndex].Topic.Name == "contributor"{
			contributorCount += 1
		}
		topics = append(topics, topicsList[topicIndex].Topic.Name)
	}
	return topics
}

/*
@Params - nil
@Description - Fetches github projects and returns the data in an array for the graphql endpoint to grab
*/
func FetchProjects(githubUser string) (*model.GithubProjects, error) {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUBACCESSTOKEN")},
	)
	httpClient := oauth2.NewClient(ctx, src)

	client := githubv4.NewClient(httpClient)
	var query struct {
		RepositoryOwner struct {
			Login        string
			Repositories struct {
				Nodes []struct {
					Name        string
					Description string
					Url         string
					HomepageUrl string
					Languages   struct {
						Nodes []struct {
							Name string
						}
					}`graphql:"languages(first: 100)"`
					RepositoryTopic struct {
						Nodes []struct {
							Topic struct {
								Name string
							}
						}
					} `graphql:"repositoryTopics(first: 100)"`
				}
			} `graphql:"repositories(privacy: PUBLIC, first: 100, ownerAffiliations: [OWNER])"`
		} `graphql:"repositoryOwner(login: $user)"`
	}
	variables := map[string]interface{}{
		"user": githubv4.String(githubUser),
	}
	queryErr := client.Query(ctx, &query, variables)
	if queryErr != nil {
		fmt.Printf("Error is %v", queryErr)
	}
	contributorCount := 0
	repos := query.RepositoryOwner.Repositories.Nodes
	var projects []*model.Project
	for repo := 0; repo < len(repos); repo++ {
		topics := ParseTopics(repos[repo].RepositoryTopic.Nodes, contributorCount)
		projects = append(projects, &model.Project{
			Name:           repos[repo].Name,
			Githublink:     repos[repo].Url,
			Description:    repos[repo].Description,
			Languages:      ParseLanguages(repos[repo].Languages.Nodes),
			Topics:         topics,
			Deploymentlink: repos[repo].HomepageUrl,
		})
	}

	var gitProjects = model.GithubProjects{Projects: projects, ContributorCount: contributorCount}
	return &gitProjects, queryErr
}

func GrabProfileData(githubUser string) (*model.GithubBio, error){
	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUBACCESSTOKEN")},
	)
	httpClient := oauth2.NewClient(ctx, src)

	client := githubv4.NewClient(httpClient)
	var query struct {
		Viewer struct {
			Bio 		   string
			Company 	   string
		}
	}
	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		fmt.Printf("Error is %v", err)
	}
	
	var readmeQuery struct{
	Repository struct {
			Object struct{
				Blob struct{
					Text string
				}`graphql:"... on Blob"`
			}`graphql:"object(expression: \"main:README.md\")"`
		}`graphql:"repository(owner: $user, name: $repoName)"`
	}

	variables := map[string]interface{}{
			"user": githubv4.String(githubUser),
			"repoName": githubv4.String(githubUser),
	}
	readmeQueryErr := client.Query(ctx, &readmeQuery, variables)
	if readmeQueryErr != nil {
		fmt.Printf("Error is %v", readmeQueryErr)
	}
	
	var GithubBioInfo = model.GithubBio{
		Position: query.Viewer.Bio,
		Company: query.Viewer.Company,
		Readme: readmeQuery.Repository.Object.Blob.Text,
	}
	return &GithubBioInfo, err
}
