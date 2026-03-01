package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	listPath = "list.json"
)

type Station struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type FetchRule struct {
	Name      string
	PageURL   string
	MustAll   []string // 同一行必須全部包含
	MustAny   []string // 同一行只要包含其中一個（可空）
	PickFirst bool     // 多個符合時取第一個
}

func main() {
	// 1) 抓兩個電台 URL
	rules := []FetchRule{
		{
			Name:      "港都電台",
			PageURL:   "http://www.bestradio.com.tw/",
			MustAll:   []string{"m3u8", "RA000012"},
			PickFirst: true,
		},
		{
			Name:      "台北愛樂電台",
			PageURL:   "https://www.e-classical.com.tw/index.html",
			MustAll:   []string{"m3u8"},
			PickFirst: true,
		},
	}

	fetched := make([]Station, 0, len(rules))
	for _, r := range rules {
		u, err := fetchStreamURLByYTDLP(r.PageURL, r.MustAll, r.MustAny, r.PickFirst)
		if err != nil {
			fatal(fmt.Errorf("%s 抓取失敗: %w", r.Name, err))
		}
		fetched = append(fetched, Station{Name: r.Name, URL: u})
		fmt.Println("Fetched:", r.Name, "=>", u)
	}

	// 2) 讀 JSON
	stations, err := readStations(listPath)
	if err != nil {
		fatal(err)
	}

	// 3) 依序 upsert（不存在就插到開頭）
	//    注意：連續插到開頭時，後插入的會在最前面
	//    我們希望最後順序是：港都電台 在最上面，其次 台北愛樂電台
	//    所以倒序處理插入，讓港都最後插入到最前。
	for i := len(fetched) - 1; i >= 0; i-- {
		stations = upsertAtFront(stations, fetched[i])
	}

	// 4) 寫回 JSON
	if err := writeStations(listPath, stations); err != nil {
		fatal(err)
	}

	fmt.Println("OK: updated", listPath)
}

// fetchStreamURLByYTDLP runs: yt-dlp -g <url>
// then finds the first line that matches:
// - contains ALL strings in mustAll
// - and (if mustAny not empty) contains at least one in mustAny
func fetchStreamURLByYTDLP(pageURL string, mustAll, mustAny []string, pickFirst bool) (string, error) {
	cmd := exec.Command("yt-dlp", "-g", pageURL)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("yt-dlp failed: %s", msg)
	}

	lines := strings.Split(out.String(), "\n")
	var matches []string

LineLoop:
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// mustAll: 全部都要出現
		for _, s := range mustAll {
			if s == "" {
				continue
			}
			if !strings.Contains(line, s) {
				continue LineLoop
			}
		}

		// mustAny: 至少要命中一個（若有設定）
		if len(mustAny) > 0 {
			ok := false
			for _, s := range mustAny {
				if s == "" {
					continue
				}
				if strings.Contains(line, s) {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}

		matches = append(matches, line)
		if pickFirst {
			return line, nil
		}
	}

	if len(matches) > 0 {
		return matches[0], nil
	}

	preview := strings.TrimSpace(out.String())
	if len(preview) > 800 {
		preview = preview[:800] + "..."
	}
	return "", fmt.Errorf("no matching url found. yt-dlp output preview:\n%s", preview)
}

func readStations(path string) ([]Station, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return []Station{}, nil
	}

	var stations []Station
	if err := json.Unmarshal(data, &stations); err != nil {
		return nil, err
	}
	return stations, nil
}

func writeStations(path string, stations []Station) error {
	data, err := json.MarshalIndent(stations, "", "    ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0644)
}

// upsertAtFront updates matching name's URL, else inserts at front.
// (更新時不會移到最前面；若你想更新也移到最前，跟我說我再改。)
func upsertAtFront(stations []Station, s Station) []Station {
	for i := range stations {
		if stations[i].Name == s.Name {
			stations[i].URL = s.URL
			return stations
		}
	}
	return append([]Station{s}, stations...)
}

func fatal(err error) {
	if err == nil {
		err = errors.New("unknown error")
	}
	fmt.Fprintln(os.Stderr, "ERROR:", err)
	os.Exit(1)
}
