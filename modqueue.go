package main

type ModResp struct {
	Data struct {
		After    interface{} `json:"after"`
		Before   interface{} `json:"before"`
		Children []struct {
			Data struct {
				ID                         string        `json:"id"`
				ModNote                    interface{}   `json:"mod_note"`
				ModReasonBy                interface{}   `json:"mod_reason_by"`
				ModReasonTitle             interface{}   `json:"mod_reason_title"`
				ModReports                 [][]string    `json:"mod_reports"`
				Subreddit                  string        `json:"subreddit"`
				Title                      string        `json:"title"`
				URL                        string        `json:"url"`
			} `json:"data"`
			Kind string `json:"kind"`
		} `json:"children"`
		Dist    int64       `json:"dist"`
		Modhash interface{} `json:"modhash"`
	} `json:"data"`
	Kind string `json:"kind"`
}
