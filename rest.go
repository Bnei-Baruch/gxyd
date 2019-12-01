package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func (a *App) handleProtocol(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)

	// Parse JSON
	var data map[string]interface{}
	json.Unmarshal([]byte(b), &data)

	text := data["text"].(string)
	var msg map[string]interface{}
	json.Unmarshal([]byte(text), &msg)

	// User role parsing
	user := msg["user"].(map[string]interface{})
	role := user["role"].(string)
	chk, err := regexp.MatchString(`^(user|group)$`, role)
	if err != nil {
		fmt.Println("Regexp Failed:", err)
	}

	// Allow only user and group role
	if chk == true {
		id := user["id"].(string)
		ep := role + "s/" + id

		// Filter event
		switch t := msg["type"].(string); t {
		case "question", "camera":
			user[t] = msg["status"]
			u, _ := json.Marshal(user)
			err := postReq(ep, string(u))
			if err != nil {
				fmt.Println("Post Request Failed:", err)
			}
		default:
			fmt.Println("no cases in switch")
		}

		// Write to log
		err := writeToLog(os.Getenv(EnvLogPath)+"/protocol.log", string(b))
		if err != nil {
			fmt.Println("Log Failed:", err)
		}
		defer r.Body.Close()
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) handelEvent(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)

	// Parse JSON
	var data map[string]interface{}
	json.Unmarshal([]byte(b), &data)
	msg := data["event"].(map[string]interface{})
	plugin := msg["plugin"].(string)

	// Filter videroom plugin events
	if plugin == "janus.plugin.videoroom" {

		// Check if property exist
		event_data := msg["data"].(map[string]interface{})
		if _, ok := event_data["event"].(string); ok {

			// User role parsing
			user_data := event_data["display"].(string)
			var user map[string]interface{}
			json.Unmarshal([]byte(user_data), &user)
			role := user["role"].(string)
			chk, err := regexp.MatchString(`^(user|group)$`, role)
			if err != nil {
				fmt.Println("Regexp Failed:", err)
			}

			// Allow only user and group role
			if chk == true {
				id := user["id"].(string)
				ep := role + "s/" + id

				// Filter data event
				switch t := event_data["event"].(string); t {
				case "joined":
					u, _ := json.Marshal(user)
					err := postReq(ep, string(u))
					if err != nil {
						fmt.Println("Post Request Failed:", err)
					}

					err = writeToLog(os.Getenv(EnvLogPath)+"/events.log", string(b))
					if err != nil {
						fmt.Println("Log Failed:", err)
					}
					defer r.Body.Close()
				case "leaving":
					err := delReq(ep)
					if err != nil {
						fmt.Println("Del Request Failed:", err)
					}

					err = writeToLog(os.Getenv(EnvLogPath)+"/events.log", string(b))
					if err != nil {
						fmt.Println("Log Failed:", err)
					}
					defer r.Body.Close()
				}
			}
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
