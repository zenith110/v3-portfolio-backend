package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"os"

	"github.com/zenith110/portfilo/graph/generated"
	"github.com/zenith110/portfilo/graph/model"
	"github.com/zenith110/portfilo/graph/routes"
)

func (r *queryResolver) GithubProjects(ctx context.Context) (*model.GithubProjects, error) {
	githubUser := os.Getenv("GITHUBUSER")
	github, err := routes.FetchProjects(githubUser)
	return github, err
}

func (r *queryResolver) NotionGoals(ctx context.Context) (*model.NotionGoals, error) {
	intergrationToken := os.Getenv("INTERGRATIONTOKEN")
	notionGoals, err := routes.FetchNotionPage(intergrationToken)
	return notionGoals, err
}

func (r *queryResolver) Profile(ctx context.Context) (*model.GithubBio, error) {
	githubUser := os.Getenv("GITHUBUSER")
	githubBio, err := routes.GrabProfileData(githubUser)
	return githubBio, err
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
