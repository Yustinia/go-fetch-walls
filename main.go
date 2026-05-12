package main

import (
	"encoding/json"
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

func main() {
	baseURL := "https://wallhaven.cc/api/v1/search"
	configData, err := os.ReadFile("configs/settings.json")
	if err != nil {
		panic(err)
	}

	var settings Settings
	err = json.Unmarshal(configData, &settings)
	if err != nil {
		panic(err)
	}
}
