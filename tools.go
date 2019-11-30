package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func postReq(ep string, payload string) error {
	body := strings.NewReader(payload)
	req, err := http.NewRequest("PUT", os.Getenv(EnvBackendUrl)+"/galaxy/"+ep, body)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		fmt.Println("Non-OK HTTP status:", response.StatusCode)
	}

	return nil
}

func delReq(ep string) error {
	req, err := http.NewRequest("DELETE", os.Getenv(EnvBackendUrl)+"/galaxy/"+ep, nil)
	client := &http.Client{}
	response, err := client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		fmt.Println("Non-OK HTTP status:", response.StatusCode)
	}

	return nil
}

func writeToLog(path string, line string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file: ", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(line)

	return nil
}
