package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bans1mp/mcp-server/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	spotifyTopSongs = "https://api.spotify.com/v1/me/tracks"
)

type SpotifyInput struct {
	NumberOfSongs int `json:"number_of_songs" jsonschema:"the number of top songs to retrieve"`
}

type SpotifyOutput struct {
	TopSongs []string `json:"top_songs" jsonschema:"the list of top songs"`
}	

type ApiResponse struct {
	Items []*Item `json:"items"`
}

type Item struct {
	Track *Track `json:"track"`
}

type Track struct {
	Name string `json:"name"`
}

var GetTopSongsTool = &mcp.Tool{
	Name:        "get_top_songs",
	Description: "Retrieve the top songs from Spotify",
}

func GetTopSongs(ctx context.Context, req *mcp.CallToolRequest, input *SpotifyInput) (*mcp.CallToolResult, *SpotifyOutput, error) {
	songsStr := strconv.Itoa(input.NumberOfSongs)
	url := spotifyTopSongs + "?limit=" + songsStr
	
	httpClient := &http.Client{}

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Authorization header
	httpReq.Header.Set("Authorization", "Bearer "+ getSpotifyAccessToken())

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Non 200 means token invalid or missing scopes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("spotify error: %s", string(body))
	}

	apiResp := ApiResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Build output
	output := &SpotifyOutput{
		TopSongs: []string{},
	}

	for _, t := range apiResp.Items {
		if t.Track == nil {
			continue
		}
		output.TopSongs = append(output.TopSongs, t.Track.Name)
	}

	return &mcp.CallToolResult{}, output, nil
}

func getSpotifyAccessToken() string {
	return auth.SpotifyAccessToken
}
