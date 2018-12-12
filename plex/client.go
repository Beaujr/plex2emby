package plex

import (
	"net/http"
	"fmt"
	"time"
	"io/ioutil"
	"encoding/xml"
)

func NewClient(url string, apikey string) Client {
	client := http.Client{Transport:http.DefaultTransport.(*http.Transport), Timeout:time.Second * 10}
	return Client{client:client, baseURL:url, apiKey:apikey}
}

type Client struct {
	client http.Client
	baseURL string
	apiKey string
}

type DirectoryContainer struct {
	Container
	Sections   []Section   `xml:"Directory"`
}

type Container struct {
	XMLName xml.Name `xml:"MediaContainer"`
}

type Section struct {
	Data
	XMLName xml.Name `xml:"Directory"`
}


type VideoContainer struct {
	Container
	Video   []Video   `xml:"Video"`
}

type Video struct {
	Data
	XMLName xml.Name `xml:"Video"`
}

type Data struct {
	Key string  `xml:"key,attr"`
	Type string `xml:"type,attr"`
	Title string  `xml:"title,attr"`
	ViewCount string `xml:"viewCount,attr"`
	EpisodeNumber string `xml:"index,attr"`
	LastViewedAt string` xml:"lastViewedAt,attr"`
}

type Plex interface {
	GetSections() ([]Section, error)
	GetFilmSection(section string) ([]Video, error)
	GetTVSection(section string) ([]Section, error)
	GetShow(key string) ([]Section, error)
	GetSeason(key string) ([]Video, error)
}

func(c *Client) GetSections() ([]Section, error) {
	url := fmt.Sprintf("%s/library/sections?X-Plex-Token=%s&X-Plex-Language=en&X-Plex-Text-Format=plain", c.baseURL, c.apiKey)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var directoryContainer DirectoryContainer
	xml.Unmarshal(byteValue, &directoryContainer)

	return directoryContainer.Sections, nil
}

func(c *Client) GetFilmSection(section string) ([]Video, error) {
	url := fmt.Sprintf("%s/library/sections/%s/all?X-Plex-Token=%s&X-Plex-Language=en&X-Plex-Text-Format=plain", c.baseURL, section, c.apiKey)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var videoContainer VideoContainer
	xml.Unmarshal(byteValue, &videoContainer)

	return videoContainer.Video, nil
}

func(c *Client) GetTVSection(section string) ([]Section, error) {
	url := fmt.Sprintf("%s/library/sections/%s/all?X-Plex-Token=%s&X-Plex-Language=en&X-Plex-Text-Format=plain", c.baseURL, section, c.apiKey)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var directoryContainer DirectoryContainer
	xml.Unmarshal(byteValue, &directoryContainer)

	return directoryContainer.Sections, nil
}

func(c *Client) GetShow(key string) ([]Section, error) {
	url := fmt.Sprintf("%s%s?X-Plex-Token=%s&X-Plex-Language=en&X-Plex-Text-Format=plain", c.baseURL, key, c.apiKey)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var directoryContainer DirectoryContainer
	xml.Unmarshal(byteValue, &directoryContainer)

	return directoryContainer.Sections, nil
}

func(c *Client) GetSeason(key string) ([]Video, error) {
	url := fmt.Sprintf("%s%s?X-Plex-Token=%s&X-Plex-Language=en&X-Plex-Text-Format=plain", c.baseURL, key, c.apiKey)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var mediacontainer VideoContainer
	xml.Unmarshal(byteValue, &mediacontainer)
	return mediacontainer.Video, nil
}