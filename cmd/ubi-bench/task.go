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
	log.Infof("send req: %s", string(jsonData))

	url := utils.GetConfig().HUB.HubUrl
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Errorf("POST request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Infof("Request successful, status code: %d", resp.StatusCode)
	} else {
		log.Errorf("Request failed, status code: %d", resp.StatusCode)

		// Read the error response
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		errBody := buf.String()
		log.Infof("Error message: %s", errBody)
	}
}

type TaskStats struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data ResourceCountList `json:"data"`
}

type ResourceCount struct {
	ResourceId int `json:"resource_id"`
	Count      int `json:"count"`
}

type ResourceCountList []ResourceCount

func (t ResourceCountList) Len() int {
	return len(t)
}

func (t ResourceCountList) Less(i, j int) bool {
	return t[i].Count < t[j].Count
}

func (t ResourceCountList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
