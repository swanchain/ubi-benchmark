package web

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/constants"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/utils"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/logs"
)

const HTTP_CONTENT_TYPE_FORM = "application/x-www-form-urlencoded"
const HTTP_CONTENT_TYPE_JSON = "application/json; charset=UTF-8"

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HttpPostTimeout(uri, tokenString string, params interface{}, timeoutSecond int, result interface{}) error {
	err := HttpRequest(http.MethodPost, uri, &tokenString, params, &timeoutSecond, result)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}

func HttpPost(uri, tokenString string, params, result interface{}) error {
	err := HttpRequest(http.MethodPost, uri, &tokenString, params, nil, result)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}

func HttpPut(uri, tokenString string, params, result interface{}) error {
	err := HttpRequest(http.MethodPut, uri, &tokenString, params, nil, result)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}

func HttpGet(uri, tokenString string, params, result interface{}) error {
	err := HttpRequest(http.MethodGet, uri, &tokenString, params, nil, result)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}

type McsResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HttpRequest(httpMethod, uri string, tokenString *string, params interface{}, timeoutSecond *int, result interface{}) error {
	var request *http.Request
	var err error

	switch params := params.(type) {
	case io.Reader:
		request, err = http.NewRequest(httpMethod, uri, params)
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		request.Header.Set("Content-Type", HTTP_CONTENT_TYPE_FORM)
	default:
		jsonReq, errJson := json.Marshal(params)
		if errJson != nil {
			logs.GetLogger().Error(errJson)
			return err
		}

		request, err = http.NewRequest(httpMethod, uri, bytes.NewBuffer(jsonReq))
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		request.Header.Set("Content-Type", HTTP_CONTENT_TYPE_JSON)
	}

	if !utils.IsStrEmpty(tokenString) {
		request.Header.Set("Authorization", "Bearer "+*tokenString)
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Transport: customTransport}
	if timeoutSecond != nil {
		client.Timeout = time.Duration(*timeoutSecond) * time.Second
	}

	response, err := client.Do(request)

	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("http status: %s, code:%d, url:%s", response.Status, response.StatusCode, uri)
		logs.GetLogger().Error(err)
		switch response.StatusCode {
		case http.StatusNotFound:
			logs.GetLogger().Error("please check your url:", uri)
		case http.StatusUnauthorized:
			logs.GetLogger().Error("Please check your token:", tokenString)
		}
	}

	var responseBody []byte
	var mcsResponse McsResponse
	if response.Body != nil && response.StatusCode != http.StatusNotFound {
		responseBody, err = io.ReadAll(response.Body)
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}

		err = json.Unmarshal(responseBody, &mcsResponse)
		if err != nil {
			logs.GetLogger().Error(err)
			logs.GetLogger().Error(string(responseBody))
			return err
		}

		if !strings.EqualFold(mcsResponse.Status, constants.HTTP_STATUS_SUCCESS) {
			err := fmt.Errorf("%s failed, status:%s, message:%s", uri, mcsResponse.Status, mcsResponse.Message)
			logs.GetLogger().Error(err)
			return err
		}

		if mcsResponse.Data != nil {
			mcsResponseDataJson, err := json.Marshal(mcsResponse.Data)
			if err != nil {
				logs.GetLogger().Error(err)
				logs.GetLogger().Error(string(responseBody))
				return err
			}

			err = json.Unmarshal(mcsResponseDataJson, result)
			if err != nil {
				logs.GetLogger().Error(err)
				logs.GetLogger().Error(string(responseBody))
				return err
			}
		}
	}

	return nil
}

func HttpUploadFileByStream(uri, filefullpath string) ([]byte, error) {
	fileReader, err := os.Open(filefullpath)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	filename := filepath.Base(filefullpath)

	boundary := "MyMultiPartBoundary12345"
	token := "DEPLOY_GATE_TOKEN"
	message := "Uploaded by Nebula"
	releaseNote := "Built by Nebula"
	fieldFormat := "--%s\r\nContent-Disposition: form-data; name=\"%s\"\r\n\r\n%s\r\n"
	tokenPart := fmt.Sprintf(fieldFormat, boundary, "token", token)
	messagePart := fmt.Sprintf(fieldFormat, boundary, "message", message)
	releaseNotePart := fmt.Sprintf(fieldFormat, boundary, "release_note", releaseNote)
	fileName := filename
	fileHeader := "Content-type: application/octet-stream"
	fileFormat := "--%s\r\nContent-Disposition: form-data; name=\"file\"; filename=\"%s\"\r\n%s\r\n\r\n"
	filePart := fmt.Sprintf(fileFormat, boundary, fileName, fileHeader)
	bodyTop := fmt.Sprintf("%s%s%s%s", tokenPart, messagePart, releaseNotePart, filePart)
	bodyBottom := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	body := io.MultiReader(strings.NewReader(bodyTop), fileReader, strings.NewReader(bodyBottom))

	contentType := fmt.Sprintf("multipart/form-data; boundary=%s", boundary)

	response, err := http.Post(uri, contentType, body)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, nil
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("http status:%s, code:%d, url:%s", response.Status, response.StatusCode, uri)
		logs.GetLogger().Error(err)
		switch response.StatusCode {
		case http.StatusNotFound:
			logs.GetLogger().Error("please check your url:", uri)
		}
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	responseStr := string(responseBody)
	//logs.GetLogger().Info(responseStr)
	filesInfo := strings.Split(responseStr, "\n")
	if len(filesInfo) < 4 {
		err := fmt.Errorf("not enough files info returned, ipfs response:%s", responseStr)
		logs.GetLogger().Error(err)
		return nil, err
	}
	responseStr = filesInfo[3]
	return []byte(responseStr), nil
}
