package matsurihime

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"mltd-linebot-border-go/common"
)

// events struct
type (
	Event struct {
		ID         int      `json:"id"`
		Type       int      `json:"type"`
		AppealType int      `json:"appealType"`
		Name       string   `json:"name"`
		Schedule   Schedule `json:"schedule"`
		Item       Item     `json:"item"`
	}

	Schedule struct {
		BeginAt      time.Time  `json:"beginAt"`
		EndAt        time.Time  `json:"endAt"`
		PageOpenedAt time.Time  `json:"pageOpenedAt"`
		PageClosedAt time.Time  `json:"pageClosedAt"`
		BoostBeginAt *time.Time `json:"boostBeginAt"`
		BoostEndAt   *time.Time `json:"boostEndAt"`
	}

	Item struct {
		Name      *string `json:"name"`
		ShortName *string `json:"shortName"`
	}
)

// rankings struct
type (
	ScoreData struct {
		Score        int       `json:"score"`
		AggregatedAt time.Time `json:"aggregatedAt"`
	}

	RankEntry struct {
		Rank int         `json:"rank"`
		Data []ScoreData `json:"data"`
	}
)

var (
	apiUrl = "https://api.matsurihi.me/api/mltd/v2/events/"
)

func GetEvents() ([]Event, error) {
	var events []Event

	req, err := sendGetRequest(apiUrl)
	if err != nil {
		return events, err
	}

	err = json.Unmarshal(req, &events)
	if err != nil {
		return events, err
	}
	return events, nil
}

func GetRankings(eventID int, logType string, rankFormat string) ([]RankEntry, error) {
	var rankings []RankEntry
	req, err := sendGetRequest(fmt.Sprintf("%s%d/rankings/%s/logs/%s", apiUrl, eventID, logType, rankFormat))
	if err != nil {
		return rankings, err
	}

	err = json.Unmarshal(req, &rankings)
	if err != nil {
		return rankings, err
	}
	return rankings, nil
}

func sendGetRequest(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%w %d", common.ErrStatusCodeAbnormal, resp.StatusCode)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return r, nil
}
