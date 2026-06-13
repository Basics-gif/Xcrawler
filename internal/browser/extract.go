package browser

import (
	"encoding/json"
	"fmt"
)

type initialPageJson struct {
	LayoutPage struct {
		VideoListProps struct {
			VideoThumbProps []Video `json:"videoThumbProps"`
		} `json:"videoListProps"`
	} `json:"layoutPage"`
}

type videoPageJSON struct {
	RalatedVideosComponent struct {
		VideoTabInitialData struct {
			VideoListProps struct {
				VideoThumbProps []Video `json:"videoThumbProps"`
			} `json:"videoListProps"`
		} `json:"videoTabInitialData"`
	} `json:"relatedVideosComponent"`
}

func (m *Manager) ExtractInitials() ([]byte, error) {
	result, err := m.Page.Evaluate(`() => JSON.stringify(window.initials)`)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate window.initials: %w", err)
	}

	raw, ok := result.(string)
	if !ok {
		return nil, fmt.Errorf("window.initials não é string, tipo: %T", result)
	}

	return []byte(raw), nil
}

func ParseInitialPage(raw []byte) (*VideoList, error) {
	var data initialPageJson
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("could not parse initial page: %w", err)
	}

	return &VideoList{Videos: data.LayoutPage.VideoListProps.VideoThumbProps}, nil
}

func ParseVideoPage(raw []byte) (*VideoList, error) {
	var data videoPageJSON
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("could not parse video page: %w", err)
	}
	return &VideoList{Videos: data.RalatedVideosComponent.VideoTabInitialData.VideoListProps.VideoThumbProps}, nil
}
