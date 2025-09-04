package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ElrohirGT/Redes_RemoteMCPServer/lib"
)

type GetCardsResponse struct {
	Cards []MTGCard `json:"cards"`
}

var GET_CARDS_CACHE = lib.NewCache(&GetCardsResponse{}, time.Now())

func GetCardsCore(ctx context.Context) (GetCardsResponse, error, bool) {
	l.Println("Checking cache validity...")
	if data, valid := GET_CARDS_CACHE.GetData(); valid {
		l.Println("Cache is valid, returning cache...")
		return *data, nil, true
	}

	var response GetCardsResponse
	req, err := http.NewRequestWithContext(ctx, "GET", API_BASE_URL+"/v1/cards", nil)
	if err != nil {
		l.Println("Failed to create request:", err)
		return response, err, true
	}
	l.Println("Calling:", req.URL.String())

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
	l.Println("Body bytes:\n", string(bodyBytes))

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		l.Println("Failed to unmarshal body:", err)
		return response, err, true
	}

	l.Println("Updating cache...")
	GET_CARDS_CACHE.Update(&response, time.Now().Add(CACHE_DURATION))
	return response, nil, false
}
