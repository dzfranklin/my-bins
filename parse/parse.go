package parse

import (
	"encoding/json"
	"fmt"
)

type Collection struct {
	Bins []string
	Date string
}

func Parse(body []byte) ([]Collection, error) {
	var data struct {
		Data struct {
			TabCollections []struct {
				Colour string
				Date   string
				Type   string
			} `json:"tab_collections"`
		}
	}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	var out []Collection
	for _, c := range data.Data.TabCollections {
		var last *Collection
		if len(out) > 0 {
			last = &out[len(out)-1]
		}

		if last == nil || last.Date != c.Date {
			out = append(out, Collection{
				Bins: []string{
					c.Type,
				},
				Date: c.Date,
			})
		} else {
			last.Bins = append(last.Bins, c.Type)
		}
	}

	return out, nil
}
