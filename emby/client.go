package emby

import (
"net/http"
"fmt"
"time"
"io/ioutil"
	"encoding/json"
	url2 "net/url"
)

func NewClient(url string, apikey string, userId string) Client {
	client := http.Client{Transport:http.DefaultTransport.(*http.Transport), Timeout:time.Second * 10}
	return Client{client:client, baseURL:url, apiKey:apikey, userId:userId}
}

type Client struct {
	client http.Client
	baseURL string
	apiKey string
	userId string
}

type SearchHints struct {
	Results []Result `json:"SearchHints"`
}

type Show struct {
	Episodes []Episode `json:"Items"`
}

type Episode struct {
	Id string `json:"Id"`
	Type string `json:"Type"`
	Name string `json:"Name"`
	IndexNumber int `json:"IndexNumber"`
	SeasonName string `json:"SeasonName"`
	UserData UserData `json:"UserData"`
}

type UserData struct {
	PlaybackPositionTicks int `json:"Id"`
	PlayCount int `json:"PlayCount"`
	IsFavorite string `json:"IsFavorite"`
	Played bool `json:"Played"`
	Key string  `json:"Key"`
}

type Result struct {
	Id string `json:"Id"`
	Name string `json:"Name"`
	Type string `json:"Type"`
}


type Emby interface {
	Search(searchTerm string, itemType string) ([]Result, error)
	GetItem(guid string) ([]Episode, error)
	MarkItemAsPlayed(guid string, datePlayed string) error

}

func(c *Client) Search(searchTerm string, itemType string) ([]Result, error) {
	url := fmt.Sprintf("%s/Search/Hints?SearchTerm=%s&IncludeItemTypes=%s&api_key=%s", c.baseURL, url2.PathEscape(searchTerm), itemType, c.apiKey)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SearchHints
	json.Unmarshal(byteValue, &result)
	return result.Results, nil
}

func(c *Client) GetItem(guid string) ([]Episode, error) {
	url := fmt.Sprintf("%s/Shows/%s/Episodes?UserId=%s&api_key=%s", c.baseURL, guid, c.userId, c.apiKey)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Show
	json.Unmarshal(byteValue, &result)
	return result.Episodes, nil
}

func(c *Client) MarkItemAsPlayed(guid string, datePlayed string) error {
	url := fmt.Sprintf("%s/Users/%s/PlayedItems/%s?DatePlayed=%s&api_key=%s", c.baseURL, c.userId, guid, datePlayed, c.apiKey)
	resp, err := c.client.Post(url, "application/json", nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%d status code", resp.StatusCode)
	}

	return nil
}