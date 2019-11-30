package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func (a *App) handleProtocol(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)

	// Parse JSON
	var data map[string]interface{}
	json.Unmarshal([]byte(b), &data)

	text := data["text"].(string)
	var msg map[string]interface{}
	json.Unmarshal([]byte(text), &msg)

	user := msg["user"].(map[string]interface{})
	u, _ := json.Marshal(user)

	// Filter event
	switch t := msg["type"].(string); t {
	case "question", "camera":
		role := user["role"].(string)
		id := user["id"].(string)
		ep := "groups/" + id
		if role == "user" {
			ep = "users/" + id
		}
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

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) handelEvent(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)

	// Parse JSON
	var data map[string]interface{}
	json.Unmarshal([]byte(b), &data)
	msg := data["event"].(map[string]interface{})
	plugin := msg["plugin"].(string)

	// Filter plugin events
	if plugin == "janus.plugin.videoroom" {

		event_data := msg["data"].(map[string]interface{})

		// Check if property exist
		if _, ok := event_data["event"].(string); ok {

			user_data := event_data["display"].(string)
			var user map[string]interface{}
			json.Unmarshal([]byte(user_data), &user)
			u, _ := json.Marshal(user)

			// User role parsing
			role := user["role"].(string)
			id := user["id"].(string)
			ep := "groups/" + id
			if role == "user" {
				ep = "users/" + id
			}

			// Filter data event
			switch t := event_data["event"].(string); t {
			case "joined":
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

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
