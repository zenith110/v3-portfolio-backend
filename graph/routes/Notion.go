package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/zenith110/Portfolio-Backend/graph/model"
)

type NotionResults struct {
	Object  string `json:"object"`
	Results []struct {
		Object         string    `json:"object"`
		ID             string    `json:"id"`
		CreatedTime    time.Time `json:"created_time"`
		LastEditedTime time.Time `json:"last_edited_time"`
		CreatedBy      struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"created_by"`
		LastEditedBy struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"last_edited_by"`
		Cover  interface{} `json:"cover"`
		Icon   interface{} `json:"icon"`
		Parent struct {
			Type      string `json:"type"`
			Workspace bool   `json:"workspace"`
		} `json:"parent"`
		Archived   bool `json:"archived"`
		Properties struct {
			Title struct {
				ID    string `json:"id"`
				Type  string `json:"type"`
				Title []struct {
					Type string `json:"type"`
					Text struct {
						Content string      `json:"content"`
						Link    interface{} `json:"link"`
					} `json:"text"`
					Annotations struct {
						Bold          bool   `json:"bold"`
						Italic        bool   `json:"italic"`
						Strikethrough bool   `json:"strikethrough"`
						Underline     bool   `json:"underline"`
						Code          bool   `json:"code"`
						Color         string `json:"color"`
					} `json:"annotations"`
					PlainText string      `json:"plain_text"`
					Href      interface{} `json:"href"`
				} `json:"title"`
			} `json:"title"`
		} `json:"properties"`
		URL string `json:"url"`
	} `json:"results"`
	NextCursor     interface{} `json:"next_cursor"`
	HasMore        bool        `json:"has_more"`
	Type           string      `json:"type"`
	PageOrDatabase struct {
	} `json:"page_or_database"`
}

type BlockData struct {
	Object  string `json:"object"`
	Results []struct {
		Object string `json:"object"`
		ID     string `json:"id"`
		Parent struct {
			Type   string `json:"type"`
			PageID string `json:"page_id"`
		} `json:"parent"`
		CreatedTime    time.Time `json:"created_time"`
		LastEditedTime time.Time `json:"last_edited_time"`
		CreatedBy      struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"created_by"`
		LastEditedBy struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"last_edited_by"`
		HasChildren bool   `json:"has_children"`
		Archived    bool   `json:"archived"`
		Type        string `json:"type"`
		Paragraph   struct {
			RichText []struct {
				Type string `json:"type"`
				Text struct {
					Content string      `json:"content"`
					Link    interface{} `json:"link"`
				} `json:"text"`
				Annotations struct {
					Bold          bool   `json:"bold"`
					Italic        bool   `json:"italic"`
					Strikethrough bool   `json:"strikethrough"`
					Underline     bool   `json:"underline"`
					Code          bool   `json:"code"`
					Color         string `json:"color"`
				} `json:"annotations"`
				PlainText string      `json:"plain_text"`
				Href      interface{} `json:"href"`
			} `json:"rich_text"`
			Color string `json:"color"`
		} `json:"paragraph"`
	} `json:"results"`
	NextCursor interface{} `json:"next_cursor"`
	HasMore    bool        `json:"has_more"`
	Type       string      `json:"type"`
	Block      struct {
	} `json:"block"`
}

func GetGoalsPageId(IntergrationToken string) string {
	year, _, _ := time.Now().Date()
	url := "https://api.notion.com/v1/search"
	query := fmt.Sprintf(`{
		"query": "List of things to learn in %d",
		"filter": {
			"value": "page",
			"property": "object"
			}
		}`, year)
	payload := strings.NewReader(query)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", IntergrationToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error is %v", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var notionData NotionResults
	unmarshalErr := json.Unmarshal(body, &notionData)
	if unmarshalErr != nil {
		fmt.Printf("Error while unmarshaling project! %v", unmarshalErr)
	}

	return notionData.Results[0].ID
}
func FetchNotionPage(IntergrationToken string) (*model.NotionGoals, error) {
	notionId := GetGoalsPageId(IntergrationToken)
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=100", notionId)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", IntergrationToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error found in fetch notion page! %v", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var blockData BlockData
	unmarshalErr := json.Unmarshal(body, &blockData)
	if unmarshalErr != nil {
		fmt.Printf("Error while unmarshaling project! %v", unmarshalErr)
	}
	var goals []string
	goalsList := blockData.Results

	for goalsIndex := 0; goalsIndex < len(goalsList); goalsIndex++ {
		if len(goalsList[goalsIndex].Paragraph.RichText) >= 1 {
			goals = append(goals, goalsList[goalsIndex].Paragraph.RichText[0].PlainText)
		}
	}
	notionResults := model.NotionGoals{Goals: goals}
	return &notionResults, unmarshalErr
}
