package main

import (
	"fmt"
	"go-fetch-walls/api"
	"go-fetch-walls/cmd"
	"go-fetch-walls/internal"
)

func printWallData(result *api.Response) {
	for index, wall := range result.Data {
		fmt.Printf("[%d]: %v\n", index, wall.Path)
		fmt.Printf("	Category: %v\n", wall.Category)
		fmt.Printf("   	Purity: %v\n", wall.Purity)
		fmt.Printf("   	Resolution: %v\n", wall.Resolution)
		fmt.Println()
	}
}

func devModeSettings(mode bool) string {
	if mode {
		return "configs/dev_settings.json"
	}

	return "configs/settings.json"
}

func main() {
	configPath := devModeSettings(true)
	baseURL := "https://wallhaven.cc/api/v1/search?"
	var settings internal.Settings

	err := internal.LoadSettings(configPath, &settings)
	if err != nil {
		panic(err)
	}
	err = internal.ValidateSettings(&settings)
	if err != nil {
		panic(err)
	}

	params := api.BuildParams(&settings)
	fullURL := baseURL + params.Encode()

	result, err := api.GetResponse(fullURL)
	if err != nil {
		panic(err)
	}

	printWallData(&result)

	err = cmd.Downloader(&result)
	if err != nil {
		panic(err)
	}
}
