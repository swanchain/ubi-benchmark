package main

import (
	"bytes"
	"encoding/json"
	"github.com/swanchain/ubi-benchmark/utils"
	"net/http"
)

const (
	CPU512 = 1
	CPU32G = 2
	GPU512 = 3
	GPU32G = 4
)

type Task struct {
	Name        string `json:"name"`
	Type        int    `json:"type"`
	ZkType      string `json:"zk_type"` // fil-c2-512M
	InputParam  string `json:"input_param"`
	VerifyParam string `json:"verify_param"`
	ResourceID  int    `json:"resource_id"`
}

func DoSend(task Task) {
	jsonData, err := json.Marshal(task)
	if err != nil {
		log.Errorf("JSON encoding failed: %v", err)
		return
	}

	url := utils.GetConfig().MCS.HubUrl
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Errorf("POST request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode == http.StatusOK {
		log.Infof("Request successful, status code: %d", resp.StatusCode)
		// Process the successful response
	} else {
		log.Infof("Request failed, status code: %d", resp.StatusCode)

		// Read the error response
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		errBody := buf.String()
		log.Infof("Error message: %s", errBody)
	}
}
