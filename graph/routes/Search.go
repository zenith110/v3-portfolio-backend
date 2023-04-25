package routes

import (
	"context"
	"fmt"
	"os"

	"github.com/machinebox/graphql"
	"github.com/zenith110/Portfolio-Backend/graph/model"
)

type ArticlesResults struct {
	ArticlesPublic struct {
		Article []struct {
			Author struct {
				Name     string `json:"name"`
				Picture  string `json:"picture"`
				Username string `json:"username"`
			} `json:"author"`
			ContentData string `json:"contentData"`
			DateWritten string `json:"dateWritten"`
			Description string `json:"description"`
			Tags        []any  `json:"tags"`
			Title       string `json:"title"`
			TitleCard   string `json:"titleCard"`
			URL         string `json:"url"`
			UUID        string `json:"uuid"`
		} `json:"article"`
		Total int `json:"total"`
	} `json:"articlesPublic"`
}

type ArticleResults struct {
	ArticlePublic struct {
		Author struct {
			Name     string `json:"name"`
			Picture  string `json:"picture"`
			Username string `json:"username"`
		} `json:"author"`
		ContentData string `json:"contentData"`
		DateWritten string `json:"dateWritten"`
		Description string `json:"description"`
		Tags        []struct {
			Tag string `json:"tag"`
		} `json:"tags"`
		Title     string `json:"title"`
		TitleCard string `json:"titleCard"`
		URL       string `json:"url"`
		UUID      string `json:"uuid"`
	} `json:"articlePublic"`
}

func ZincSearch(input *model.SearchInput) (*model.Articles, error) {
	var err error
	projectUUID := os.Getenv("PROJECTUUID")
	keyword := input.Term
	username := os.Getenv("ZINCUSERNAME")
	password := os.Getenv("ZINCPASSWORD")
	graphqlURL := os.Getenv("BLOGGRAPHQLURL")
	graphqlClient := graphql.NewClient(graphqlURL)

	graphqlRequest := graphql.NewRequest(`
			query($projectUuid: String!, $keyword: String!, $username: String!, $password: String!){
				articlesPublic(input: {project_uuid: $projectUuid, keyword: $keyword, username: $username, password: $password} ){
				article{
					title
					titleCard
					uuid
					author{
						name
						picture
						username
					}
					contentData
					dateWritten
					url
					description
					tags{
					tag
					}
				}
				total
			}
		}
	`)
	graphqlRequest.Var("projectUuid", projectUUID)
	graphqlRequest.Var("keyword", keyword)
	graphqlRequest.Var("username", username)
	graphqlRequest.Var("password", password)
	var articlesResults ArticlesResults
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &articlesResults); err != nil {
		fmt.Printf("error is %v", err)
	}

	var articlesStorage []model.Article
	var tags []model.ArticleTag
	for hit := range articlesResults.ArticlesPublic.Article {
		fmt.Print(articlesResults.ArticlesPublic.Article[hit])
		author := model.Author{Name: articlesResults.ArticlesPublic.Article[hit].Author.Name, Profile: "", Picture: articlesResults.ArticlesPublic.Article[hit].Author.Picture, Username: articlesResults.ArticlesPublic.Article[hit].Author.Username}
		article := model.Article{Author: &author, ContentData: articlesResults.ArticlesPublic.Article[hit].ContentData, DateWritten: articlesResults.ArticlesPublic.Article[hit].DateWritten, Description: articlesResults.ArticlesPublic.Article[hit].Description, Tags: tags, Title: articlesResults.ArticlesPublic.Article[hit].Title, TitleCard: articlesResults.ArticlesPublic.Article[hit].TitleCard, UUID: articlesResults.ArticlesPublic.Article[hit].UUID, URL: articlesResults.ArticlesPublic.Article[hit].URL}
		articlesStorage = append(articlesStorage, article)
	}
	articlesCollection := model.Articles{ArticleCollection: articlesStorage, Total: articlesResults.ArticlesPublic.Total}

	return &articlesCollection, err
}

func PublicArticle(articleuuid string) (*model.Article, error) {
	var err error

	projectUUID := os.Getenv("PROJECTUUID")
	graphqlURL := os.Getenv("BLOGGRAPHQLURL")
	graphqlClient := graphql.NewClient(graphqlURL)
	graphqlRequest := graphql.NewRequest(`
			query($projectUuid: String!, $articleUuid: String!){
				articlePublic(input: {project_uuid: $projectUuid, article_uuid: $articleUuid} ){
					title
					titleCard
					uuid
					author{
						name
						picture
						username
					}
					contentData
					dateWritten
					url
					description
					tags{
					tag
					}
			}
		}
	`)
	graphqlRequest.Var("projectUuid", projectUUID)
	graphqlRequest.Var("articleUuid", articleuuid)
	var articleResult ArticleResults
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &articleResult); err != nil {
		fmt.Printf("error is %v", err)
	}
	author := model.Author{Name: articleResult.ArticlePublic.Author.Name, Profile: "", Picture: articleResult.ArticlePublic.Author.Picture, Username: articleResult.ArticlePublic.Author.Username}
	article := model.Article{Title: articleResult.ArticlePublic.Title, TitleCard: articleResult.ArticlePublic.TitleCard, UUID: articleResult.ArticlePublic.UUID, Author: &author, ContentData: articleResult.ArticlePublic.ContentData, DateWritten: articleResult.ArticlePublic.DateWritten, URL: articleResult.ArticlePublic.URL, Description: articleResult.ArticlePublic.Description}
	return &article, err
}
