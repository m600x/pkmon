/*
 * Send Slack message to a channel in random interval for a random time at
 * specific hours avoiding spamming at 3AM...
 *
 * Default settings:
 * - Between 9AM to 7PM
 * - Rate between 15 and 30 minutes
 * - TTL between 30 and 60 seconds
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type slackResponse struct {
	Ok      bool   `json:"ok"`
	Channel string `json:"channel"`
	Ts      string `json:"ts"`
	Message struct {
		Type     string `json:"type"`
		Subtype  string `json:"subtype"`
		Text     string `json:"text"`
		Ts       string `json:"ts"`
		Username string `json:"username"`
		BotID    string `json:"bot_id"`
	} `json:"message"`
}

const (
	botName = "Pkmon"                                                                                       // Name displayed
	botPayload = "A wild PKCHU appeared! @pk @pk ? pkchuuuuuuu!   <--- Type pk multiple time to catch him!" // Message content
	token = "---Your token goes here---"                                                                    // Auth token
	channel = "C0G07RYJD"                                                                                   // This is 42paris_global_random
	logfile = "/var/log/pkchu.log"                                                                          // Path to log file
	timeFrameMin = 9                                                                                        // Start hour authorized
	timeFrameMax = 19                                                                                       // End hour authorized
	timeRandMin  = 15                                                                                       // Rate minimal value
	timeRandMax  = 30                                                                                       // Rate maximal value
)

func sendMessage() slackResponse {
	req, _ := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", nil)
	q := req.URL.Query()
	q.Add("token", token)
	q.Add("channel", channel)
	q.Add("text", botPayload)
	q.Add("username", botName)
	q.Add("link_names", "1")
	req.URL.RawQuery = q.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error with sendMessage: %s", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var rt slackResponse
	json.Unmarshal(body, &rt)
	return rt
}

func deleteMessage(ts string) {
	req, _ := http.NewRequest("POST", "https://slack.com/api/chat.delete", nil)
	q := req.URL.Query()
	q.Add("token", token)
	q.Add("channel", channel)
	q.Add("ts", ts)
	q.Add("username", botName)
	req.URL.RawQuery = q.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error with deleteMessage: %s", err)
	}
	defer res.Body.Close()
}

func sender() {
	launch := sendMessage()
	pkchuTTL := (rand.Intn(timeRandMax-timeRandMin+1) + timeRandMin) * 2
	log.Printf("Pkchu is alive for %d seconds\n", pkchuTTL)
	time.Sleep(time.Duration(pkchuTTL) * time.Second)
	deleteMessage(launch.Ts)
}

func main() {
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	for {
		rand.Seed(time.Now().UnixNano())
		sleeper := rand.Intn(timeRandMax-timeRandMin+1) + timeRandMin
		log.Printf("Entering sleep for %d minutes\n", sleeper)
		for sleeper > 0 {
			if sleeper%5 == 0 {
				log.Printf("Sleep time left :%d minutes\n", sleeper)
			}
			time.Sleep(time.Minute)
			sleeper--
		}
		currentTime := time.Now()
		timeStampString := currentTime.Format("2020-01-02 13:37:42")
		timeStamp, _ := time.Parse("2020-01-02 13:37:42", timeStampString)
		hr, _, _ := timeStamp.Clock()
		if hr >= timeFrameMin && hr <= timeFrameMax {
			log.Printf("Within frame, launching command")
			sender()
		} else {
			log.Printf("It's %d h. Let's wait a bit...", hr)
		}
	}
}
