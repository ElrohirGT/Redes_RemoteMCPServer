package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ElrohirGT/Redes_MCPServer/lib"
)

type GetGameFormatsResponse struct {
	Formats []string `json:"formats"`
}

var GET_GAME_FORMATS_CACHE = lib.NewCache(&GetGameFormatsResponse{}, time.Now())

func GetGameFormats(ctx context.Context) (GetGameFormatsResponse, error, bool) {
	l.Println("Checking cache validity...")
	if data, valid := GET_GAME_FORMATS_CACHE.GetData(); valid {
		l.Println("Cache is valid, returning cache...")
		return *data, nil, true
	}

	var response GetGameFormatsResponse
	req, err := http.NewRequestWithContext(ctx, "GET", API_BASE_URL+"/formats", nil)
	if err != nil {
		l.Println("Failed to create request:", err)
		return response, err, true
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Println("Failed to make request:", err)
		return response, err, false
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Println("Failed to read body:", err)
		return response, err, false
	}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		l.Println("Failed to unmarshal body:", err)
		return response, err, true
	}

	l.Println("Updating cache...")
	GET_GAME_FORMATS_CACHE.Update(&response, time.Now().Add(CACHE_DURATION))
	return response, nil, false
}
