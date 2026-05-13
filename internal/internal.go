package internal

import (
	"encoding/json"
	"fmt"
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

func LoadSettings(configPath string, settings *Settings) error {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(configData, settings)
}

func ValidateSettings(settings *Settings) error {
	if len(settings.API) < 20 {
		return fmt.Errorf("'%v' is an invalid api key", settings.API)
	}

	return nil
}
