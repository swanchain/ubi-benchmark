package utils

import (
	"context"
	"fmt"
	titan_storage "github.com/utopiosphe/titan-storage-sdk"
	"strings"
)

const (
	titanStorageURL = "https://api-test1.container1.titannet.io"
)

type TiTanClient struct {
	titanStorage titan_storage.Storage
}

func NewTiTanClient(apiKey string) (*TiTanClient, error) {
	client := &TiTanClient{}

	var err error
	client.titanStorage, err = titan_storage.Initialize(&titan_storage.Config{
		TitanURL: titanStorageURL,
		APIKey:   apiKey,
	})

	if err != nil {
		return nil, err
	}
	return client, nil
}

func (client *TiTanClient) UploadFile(filePath string, folderId int) (string, error) {
	progress := func(doneSize int64, totalSize int64) {
		if doneSize == totalSize {
			fmt.Printf("%s upload success\n", filePath)
		}
	}

	root, err := client.titanStorage.UploadFilesWithPath(context.Background(), filePath, progress, false, titan_storage.WithGroupID(folderId))
	if err != nil {
		fmt.Println("UploadFile error ", err.Error())
		return "", err
	}

	assetResult, err := client.titanStorage.GetURL(context.Background(), root.String())

	var url string
	for _, l := range assetResult.URLs {
		fmt.Println("UploadFile url: ", l)
		if strings.HasPrefix(l, "https://"+root.String()) {
			filenameIndex := strings.Index(l, "filename")
			url = l[:strings.LastIndex(l, "?")] + "?" + l[filenameIndex:]
			break
		}
	}
	return url, nil
}

func (client *TiTanClient) CreateFolder(rootId int, taskDir string) (int, error) {
	return client.titanStorage.CreateFolderV2(context.Background(), taskDir, rootId)
}
