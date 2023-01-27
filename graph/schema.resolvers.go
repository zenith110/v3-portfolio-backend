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

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
