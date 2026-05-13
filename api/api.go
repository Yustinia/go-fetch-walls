package api

import (
	"encoding/json"
	"fmt"
	"go-fetch-walls/internal"
	"net/http"
	"net/url"
)

type Response struct {
	Data []internal.Wallpaper `json:"data"`
}

func BuildParams(settings *internal.Settings) url.Values {
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

func GetResponse(fullURL string) (Response, error) {
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
