package bucket

import (
	"os"
	"testing"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/logs"

	"github.com/stretchr/testify/assert"
)

func TestGetFile(t *testing.T) {
	ossFile, err := buketClient.GetFile("aaa", "duration6")
	assert.Nil(t, err)
	assert.NotEmpty(t, ossFile)

	logs.GetLogger().Info(*ossFile)
}

func TestCreateFolder(t *testing.T) {
	folderName, err := buketClient.CreateFolder("aaa", "test2", "test1")
	assert.Nil(t, err)
	assert.NotEmpty(t, folderName)

	logs.GetLogger().Info(*folderName)
}

func TestDeleteFile(t *testing.T) {
	err := buketClient.DeleteFile("aaa", "duration6")
	assert.Nil(t, err)
}

func TestListFiles(t *testing.T) {
	ossFiles, count, err := buketClient.ListFiles("aaa", "", 10, 0)
	assert.Nil(t, err)
	assert.NotNil(t, ossFiles)
	assert.NotNil(t, count)

	for _, ossFile := range ossFiles {
		logs.GetLogger().Info(*ossFile)
	}

	logs.GetLogger().Info(*count)
}

func TestUploadFile(t *testing.T) {
	err := buketClient.UploadFile("aaa", "test1/test2/test3/duration23", file2Upload, true)
	assert.Nil(t, err)
}

func TestUploadFolder(t *testing.T) {
	err := buketClient.UploadFolder("aaa", folder2Upload, "")
	assert.Nil(t, err)
}

func TestGetFileInfo(t *testing.T) {
	fileInfo, err := buketClient.GetFileInfo(6674)
	assert.Nil(t, err)
	assert.NotEmpty(t, fileInfo)

	logs.GetLogger().Info(*fileInfo)
}

func TestUploadIpfsFolder(t *testing.T) {
	ossFile, err := buketClient.UploadIpfsFolder("aaa", "aaa", folder2Upload)
	assert.Nil(t, err)
	assert.NotEmpty(t, ossFile)

	logs.GetLogger().Info(*ossFile)
}

func TestDownloadFile(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		logs.GetLogger().Fatal(err)
	}
	err = buketClient.DownloadFile("aaa", "aaa", path)
	assert.Nil(t, err)
}
