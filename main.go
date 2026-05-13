package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Settings struct {
	API        string `json:"api_key"`
	Categories string `json:"categories"`
	Purity     string `json:"purity"`
	Sorting    string `json:"sorting"`
	Order      string `json:"order"`
	AtLeast    string `json:"atleast"`
	Ratios     string `json:"ratios"`
	Page       int    `json:"page"`
}

type Wallpaper struct {
	Path       string `json:"path"`
	Purity     string `json:"purity"`
	Category   string `json:"category"`
	Resolution string `json:"resolution"`
}

type Response struct {
	Data []Wallpaper `json:"data"`
}

func loadSettings(configPath string, settings *Settings) error {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(configData, settings)
}

func validateSettings(settings *Settings) error {
	if len(settings.API) < 20 {
		return fmt.Errorf("'%v' is an invalid api key", settings.API)
	}

	return nil
}

func buildParams(settings *Settings) url.Values {
	params := url.Values{}
	params.Set("apikey", settings.API)
	params.Set("categories", settings.Categories)
	params.Set("purity", settings.Purity)
	params.Set("sorting", settings.Sorting)
	params.Set("order", settings.Order)
	params.Set("atleast", settings.AtLeast)
	params.Set("ratios", settings.Ratios)
	params.Set("page", fmt.Sprintf("%d", settings.Page))

	return params
}

func getResponse(fullURL string) (Response, error) {
	resp, err := http.Get(fullURL)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Response{}, err
	}

	return result, nil
}

func printWallData(result *Response) {
	for index, wall := range result.Data {
		fmt.Printf("[%d]: %v\n", index, wall.Path)
		fmt.Printf("	Category: %v\n", wall.Category)
		fmt.Printf("   	Purity: %v\n", wall.Purity)
		fmt.Printf("   	Resolution: %v\n", wall.Resolution)
		fmt.Println()
	}
}

func main() {
	baseURL := "https://wallhaven.cc/api/v1/search?"
	configPath := "configs/settings.json"
	var settings Settings

	err := loadSettings(configPath, &settings)
	if err != nil {
		panic(err)
	}
	err = validateSettings(&settings)
	if err != nil {
		panic(err)
	}

	params := buildParams(&settings)
	fullURL := baseURL + params.Encode()

	result, err := getResponse(fullURL)
	if err != nil {
		panic(err)
	}

	printWallData(&result)
}
