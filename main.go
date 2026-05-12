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

func main() {
	baseURL := "https://wallhaven.cc/api/v1/search?"
	configData, err := os.ReadFile("configs/settings.json")
	if err != nil {
		panic(err)
	}

	var settings Settings
	err = json.Unmarshal(configData, &settings)
	if err != nil {
		panic(err)
	}

	params := url.Values{}
	params.Set("apikey", settings.API)
	params.Set("categories", settings.Categories)
	params.Set("purity", settings.Purity)
	params.Set("sorting", settings.Sorting)
	params.Set("order", settings.Order)
	params.Set("atleast", settings.AtLeast)
	params.Set("ratios", settings.Ratios)
	params.Set("page", fmt.Sprintf("%d", settings.Page))

	fullURL := baseURL + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic(err)
	}

	for index, wall := range result.Data {
		fmt.Printf("[%d]: %v\n", index, wall.Path)
		fmt.Printf("	Category: %v\n", wall.Category)
		fmt.Printf("   	Purity: %v\n", wall.Purity)
		fmt.Printf("   	Resolution: %v\n", wall.Resolution)
		fmt.Println()
	}
}
