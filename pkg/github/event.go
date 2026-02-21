package github

import (
	"encoding/json"
	"fmt"
	"os"
)

// Event は GitHub Actions のイベント情報を保持します。
type Event struct {
	Action      string `json:"action"`
	PullRequest struct {
		Number int    `json:"number"`
		Merged bool   `json:"merged"`
		Body   string `json:"body"`
		User   struct {
			Login string `json:"login"`
		} `json:"user"`
		Base struct {
			Repo struct {
				FullName string `json:"full_name"`
			} `json:"repo"`
		} `json:"base"`
	} `json:"pull_request"`
}

// LoadEvent は GITHUB_EVENT_PATH からイベント JSON を読み込みます。
func LoadEvent(path string) (*Event, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read event file: %w", err)
	}

	var event Event
	if err := json.Unmarshal(file, &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event JSON: %w", err)
	}

	return &event, nil
}
