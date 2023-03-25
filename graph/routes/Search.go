package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/zenith110/Portfolio-Backend/graph/model"
)

func ZincSearch(input *model.SearchInput) (*model.Articles, error) {
	var articles *model.Articles
	var err error
	projectUUID := os.Getenv("PROJECTUUID")
	keyword := input.Term
	username := os.Getenv("ZINCUSERNAME")
	password := os.Getenv("ZINCPASSWORD")
	searchQuery := map[string]string{
		"query": fmt.Sprintf(`
			{
				articlesPublic(input: {project_uuid: "%s", keyword: "%s", username: "%s", password: "%s"})
			}
			articles{
				title
				titleCard
				author
				contentData
				dateWritten
				url
				description
				uuid
				tags{
					tag
				}
			}
			total
		`, projectUUID, keyword, username, password),
	}
	jsonValue, _ := json.Marshal(searchQuery)
	request, err := http.NewRequest("POST", "https://cms-backend.abrahannevarez.dev/query", bytes.NewBuffer(jsonValue))
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		panic(fmt.Sprintf("Error had occured while grabbing the articles! Error code is %v", err))
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(data))
	return articles, err
}
