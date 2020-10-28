package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const uAgent = "golang:beebot:0.1 (by /u/phlyingpenguin)"
const pushoverURL = "https://api.pushover.net/1/messages.json"

const envToken = "PUSHOVER_TOKEN"
const envUser = "PUSHOVER_USER"
const envURL = "MODQ_JSON"

var (
	pushoverToken = os.Getenv(envToken)
	pushoverUser  = os.Getenv(envUser)
	modQURL       = os.Getenv(envURL)

	dataFile = flag.String("data", "data.json", "data file")
)

func readData() (map[string]bool, error) {
	seen := map[string]bool{}
	_, err := os.Stat(*dataFile)
	if err == nil {
		seenData, err := ioutil.ReadFile(*dataFile)
		if err != nil {
			return seen, err
		}
		err = json.Unmarshal(seenData, &seen)
		if err != nil {
			return seen, err
		}
	} else {
		return seen, err
	}

	log.Debug().Interface("seen", seen).Msgf("data")

	return seen, nil
}

func getModQ() (ModResp, error) {
	req, err := http.NewRequest(http.MethodGet, modQURL, nil)
	modQ := ModResp{}

	if err != nil {
		return modQ, err
	}
	req.Header.Set("User-Agent", uAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return modQ, err
	}
	res, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(res, &modQ)
	if err != nil {
		return modQ, err
	}
	return modQ, nil
}

func saveData(seen map[string]bool) error {
	seenData, _ := json.Marshal(seen)
	return ioutil.WriteFile(*dataFile, seenData, 0666)
}

func sendPush(seen map[string]bool, modQ ModResp) (map[string]bool, error) {
	for _, item := range modQ.Data.Children {
		if seen[item.Data.ID] {
			log.Debug().Msgf("Already pushed %s, skipping.", item.Data.ID)
			continue
		}
		reports := ""
		for _, r := range item.Data.ModReports {
			reports += fmt.Sprintf("%v\n", r)
		}
		msg := fmt.Sprintf("%s from %s reports: %v",
			item.Data.Title,
			item.Data.Subreddit,
			reports)
		log.Debug().Msgf(msg)

		values := url.Values{}
		values.Set("token", pushoverToken)
		values.Set("user", pushoverUser)
		values.Set("message", msg)
		values.Set("url", item.Data.URL)
		values.Set("url_title", item.Data.Title)
		resp, err := http.PostForm(pushoverURL, values)
		if err != nil {
			return seen, err
		}
		r, _ := ioutil.ReadAll(resp.Body)
		log.Debug().Str("resp", string(r)).Msgf("post to pushover")
		seen[item.Data.ID] = true
	}
	return seen, nil
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal().Err(err).Msgf("error")
	}
}

func checkEnv() error {
	if pushoverUser == "" {
		return fmt.Errorf("must provide a pushover user ID in %s", envUser)
	}
	if pushoverToken == "" {
		return fmt.Errorf("must provide a pushover app token in %s", envToken)
	}
	if modQURL == "" {
		return fmt.Errorf("must provide a reddit Mod Queue URL in %s", envURL)
	}
	return nil
}

func main() {
	flag.Parse()

	err := checkEnv()
	fatalErr(err)
	seen, err := readData()
	fatalErr(err)
	modQ, err := getModQ()
	fatalErr(err)
	seen, err = sendPush(seen, modQ)
	fatalErr(err)
	err = saveData(seen)
	fatalErr(err)
}