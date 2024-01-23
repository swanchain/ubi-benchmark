package bucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/constants"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/utils"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/web"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/logs"

	"github.com/codingsince1985/checksum"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/jinzhu/gorm"
)

type OssFile struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	Prefix     string `json:"prefix"`
	BucketUid  string `json:"bucket_uid"`
	FileHash   string `json:"file_hash"`
	Size       int64  `json:"size"`
	PayloadCid string `json:"payload_cid"`
	PinStatus  string `json:"pin_status"`
	IsDeleted  bool   `json:"is_deleted"`
	IsFolder   bool   `json:"is_folder"`
	ObjectName string `json:"object_name"`
	Type       int    `json:"type"`
	gorm.Model
}

func (bucketClient *BucketClient) GetFile(bucketName, objectName string) (*OssFile, error) {
	bucketUid, err := bucketClient.GetBucketUid(bucketName)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_FILE_GET_FILE_INFO_BY_OBJECT_NAME)
	apiUrl = apiUrl + "?bucket_uid=" + *bucketUid + "&object_name=" + objectName

	var fileInfo OssFile
	err = web.HttpGet(apiUrl, bucketClient.JwtToken, nil, &fileInfo)
	if err != nil {
		//logs.GetLogger().Error(err)
		return nil, err
	}

	return &fileInfo, nil
}

func (bucketClient *BucketClient) CreateFolder(bucketName, folderName, prefix string) (*string, error) {
	bucketUid, err := bucketClient.GetBucketUid(bucketName)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_FILE_CREATE_FOLDER)

	var params struct {
		FileName  string `json:"file_name"`
		Prefix    string `json:"prefix"`
		BucketUid string `json:"bucket_uid"`
	}

	params.FileName = folderName
	params.Prefix = prefix
	params.BucketUid = *bucketUid

	var folderNameRetuned string
	err = web.HttpPost(apiUrl, bucketClient.JwtToken, &params, &folderNameRetuned)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &folderName, nil
}

func (bucketClient *BucketClient) DeleteFile(bucketName, objectName string) error {
	ossFile, err := bucketClient.GetFile(bucketName, objectName)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_FILE_DELETE_FILE)
	apiUrl = apiUrl + "?file_id=" + strconv.Itoa(int(ossFile.ID))

	err = web.HttpGet(apiUrl, bucketClient.JwtToken, nil, nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}

func (bucketClient *BucketClient) ListFiles(bucketName, prefix string, limit, offset int) ([]*OssFile, *int, error) {
	bucketUid, err := bucketClient.GetBucketUid(bucketName)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, nil, err
	}

	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_FILE_GET_FILE_LIST) +
		"?bucket_uid=" + *bucketUid + "&prefix=" + prefix +
		"&limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)

	var files struct {
		OssFiles []*OssFile `json:"file_list"`
		Count    int        `json:"count"`
	}

	err = web.HttpGet(apiUrl, bucketClient.JwtToken, nil, &files)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, nil, err
	}

	return files.OssFiles, &files.Count, nil
}

func (bucketClient *BucketClient) UploadFile(bucketName, objectName, filePath string, replace bool) error {
	prefix, fileName := getPrefixFileName(objectName)
	//todo check prefix exist,if not exist, create it
	// for example, objectName is a/b/c/d.txt, then prefix is a/b/c, fileName is d.txt
	// first, check a exist, if not exist, create it
	// second, check a/b exist, if not exist, create it
	// third, check a/b/c exist, if not exist, create it
	// Split the prefix by '/' to get all directories
	// If the prefix is empty, then there are no directories to create
	dirs := strings.Split(prefix, "/")
	if len(dirs) == 1 && dirs[0] == "" {
		dirs = []string{}
	}

	// Iterate over the directories, checking and creating each one
	currentPath := ""
	for _, dir := range dirs {
		if currentPath != "" {
			currentPath += "/"
		}

		// Create the directory if it does not exist
		exists, err := bucketClient.GetFile(bucketName, currentPath+dir)
		if exists != nil {
			currentPath += dir
			continue
		}
		_, err = bucketClient.CreateFolder(bucketName, dir, strings.TrimSuffix(currentPath, "/"))
		currentPath += dir
		if err != nil {
			//logs.GetLogger().Error(err)
			return err
		}
	}
	bucketUid, err := bucketClient.GetBucketUid(bucketName)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	if bucketUid == nil {
		err := fmt.Errorf("bucket:%s not exists", bucketName)
		logs.GetLogger().Error(err)
		return err
	}

	osFileInfo, err := os.Stat(filePath)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	fileSize := osFileInfo.Size()

	fileHashMd5, err := checksum.MD5sum(filePath)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	ossFileInfo, err := bucketClient.checkFile(*bucketUid, prefix, fileHashMd5, fileName)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	if ossFileInfo.FileIsExist && replace {
		err = bucketClient.DeleteFile(bucketName, objectName)
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		ossFileInfo, err = bucketClient.checkFile(*bucketUid, prefix, fileHashMd5, fileName)
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
	}

	if !ossFileInfo.FileIsExist {
		if !ossFileInfo.IpfsIsExist {
			file, err := os.Open(filePath)
			if err != nil {
				logs.GetLogger().Error(err)
				return err
			}
			bytesReadTotal := int64(0)
			chunNo := 0

			var wg sync.WaitGroup

			for bytesReadTotal < fileSize {
				var chunkSize int64
				bytesLeft := fileSize - bytesReadTotal
				if bytesLeft >= constants.FILE_CHUNK_SIZE_MAX {
					chunkSize = constants.FILE_CHUNK_SIZE_MAX
				} else {
					chunkSize = bytesLeft
				}
				chunk := make([]byte, chunkSize)
				_, err := file.ReadAt(chunk, bytesReadTotal)
				if err != nil {
					logs.GetLogger().Error(err)
					return err
				}
				bytesReadTotal = bytesReadTotal + chunkSize
				chunNo = chunNo + 1

				partFileName := strconv.Itoa(chunNo) + "_" + fileName

				wg.Add(1)
				go func() {
					logs.GetLogger().Info("file name:", partFileName, ", chunk size:", chunkSize)
					_, err = bucketClient.uploadFileChunk(fileHashMd5, partFileName, chunk)
					if err != nil {
						logs.GetLogger().Error(err)
					}
					wg.Done()
				}()
			}

			wg.Wait()
			_, err = bucketClient.mergeFile(*bucketUid, fileHashMd5, fileName, prefix)
			if err != nil {
				logs.GetLogger().Error(err)
				return err
			}
		}
	}

	return nil
}

func (bucketClient *BucketClient) UploadFolder(bucketName, folderPath, prefix string) error {
	folderName := filepath.Base(folderPath)
	_, err := bucketClient.CreateFolder(bucketName, folderName, "")
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(folderPath, file.Name())
		objectName := folderName + "/" + file.Name()
		err := bucketClient.UploadFile(bucketName, objectName, filePath, false)
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
	}

	return nil
}

func (bucketClient *BucketClient) GetFileInfo(fileId int) (*OssFile, error) {
	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_FILE_GET_FILE_INFO)
	apiUrl = apiUrl + "?file_id=" + strconv.Itoa(fileId)

	var fileInfo OssFile
	err := web.HttpGet(apiUrl, bucketClient.JwtToken, nil, &fileInfo)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &fileInfo, nil
}

func getPrefixFileName(objectName string) (string, string) {
	lastIndex := strings.LastIndex(objectName, "/")

	if lastIndex == -1 {
		return "", objectName
	}

	prefix := objectName[0:lastIndex]
	fileName := objectName[lastIndex+1:]

	return prefix, fileName
}

type OssFileInfo struct {
	FileId      uint   `form:"file_id" json:"file_id"`
	FileHash    string `form:"file_hash" json:"file_hash"`
	FileIsExist bool   `form:"file_is_exist" json:"file_is_exist"`
	IpfsIsExist bool   `form:"ipfs_is_exist" json:"ipfs_is_exist"`
	Size        int64  `form:"size" json:"size"`
	PayloadCid  string `form:"payload_cid" json:"payload_cid"`
	//IpfsUrl     string `form:"ipfs_url" json:"ipfs_url"`
}

func (bucketClient *BucketClient) checkFile(bucketUid, prefix, fileHash, fileName string) (*OssFileInfo, error) {
	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_FILE_CHECK_UPLOAD)

	var params struct {
		FileName  string `json:"file_name"`
		FileHash  string `json:"file_hash"`
		Prefix    string `json:"prefix"`
		BucketUid string `json:"bucket_uid"`
	}

	params.FileName = fileName
	params.FileHash = fileHash
	params.Prefix = prefix
	params.BucketUid = bucketUid

	var ossFileInfo OssFileInfo
	err := web.HttpPost(apiUrl, bucketClient.JwtToken, &params, &ossFileInfo)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &ossFileInfo, nil
}

func (bucketClient *BucketClient) uploadFileChunk(fileHash, fileName string, chunk []byte) ([]string, error) {
	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_FILE_UPLOAD_CHUNK)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	chunkReader := bytes.NewReader(chunk)

	//chunkReader.WriteTo(part)

	_, err = io.Copy(part, chunkReader)
	//n, err := part.Write(chunk)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	err = writer.WriteField("hash", fileHash)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	request, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bucketClient.JwtToken))
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("http status: %s, code:%d, url:%s", response.Status, response.StatusCode, apiUrl)
		logs.GetLogger().Error(err)
		switch response.StatusCode {
		case http.StatusNotFound:
			logs.GetLogger().Error("please check your url:", apiUrl)
		case http.StatusUnauthorized:
			logs.GetLogger().Error("Please check your token:", bucketClient.JwtToken)
		}
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	var mcsResponse struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Data    []string `json:"data"`
	}

	err = json.Unmarshal(responseBytes, &mcsResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if !strings.EqualFold(mcsResponse.Status, constants.HTTP_STATUS_SUCCESS) {
		err := fmt.Errorf("%s failed, status:%s, message:%s", apiUrl, mcsResponse.Status, mcsResponse.Message)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return mcsResponse.Data, nil
}

func (bucketClient *BucketClient) mergeFile(bucketUid, fileHash, fileName, prefix string) (*OssFileInfo, error) {
	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_FILE_MERGE_FILE)

	var params struct {
		FileName  string `json:"file_name"`
		FileHash  string `json:"file_hash"`
		Prefix    string `json:"prefix"`
		BucketUid string `json:"bucket_uid"`
	}

	params.FileName = fileName
	params.FileHash = fileHash
	params.Prefix = prefix
	params.BucketUid = bucketUid

	var ossFileInfo OssFileInfo
	err := web.HttpPostTimeout(apiUrl, bucketClient.JwtToken, &params, 600, &ossFileInfo)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &ossFileInfo, nil
}

type PinFiles2IpfsResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Data    OssFile `json:"data"`
}

func (bucketClient *BucketClient) UploadIpfsFolder(bucketName, objectName, folderPath string) (*OssFile, error) {
	folderName := filepath.Base(objectName)
	prefix := strings.TrimRight(objectName, folderName)

	if strings.Trim(folderName, " ") == "" {
		folderName = filepath.Base(folderPath)
	}

	bucketUid, err := bucketClient.GetBucketUid(bucketName)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	err = writer.WriteField("folder_name", folderName)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	err = writer.WriteField("prefix", prefix)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	err = writer.WriteField("bucket_uid", *bucketUid)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	fsFiles, err := os.ReadDir(folderPath)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	for _, fsFile := range fsFiles {
		file, err := os.Open(filepath.Join(folderPath, fsFile.Name()))
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}
		defer file.Close()

		part1, err := writer.CreateFormFile("files", folderName+"/"+fsFile.Name())
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}

		_, err = io.Copy(part1, file)
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_FILE_PIN_FILES_2_IPFS)
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, payload)

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bucketClient.JwtToken))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := httpClient.Do(req)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	var pinFiles2IpfsResponse PinFiles2IpfsResponse
	err = json.Unmarshal(body, &pinFiles2IpfsResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if !strings.EqualFold(pinFiles2IpfsResponse.Status, constants.HTTP_STATUS_SUCCESS) {
		err := fmt.Errorf("get parameters failed, status:%s,message:%s", pinFiles2IpfsResponse.Status, pinFiles2IpfsResponse.Message)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &pinFiles2IpfsResponse.Data, nil
}

func (bucketClient *BucketClient) DownloadFile(bucketName, objectName, destFileDir string) error {
	ossFile, err := bucketClient.GetFile(bucketName, objectName)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	gateway, err := bucketClient.GetGateway()
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	sh := shell.NewShell(*gateway)
	logs.GetLogger().Info(*gateway)
	err = sh.Get(ossFile.PayloadCid, filepath.Base(destFileDir))
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}
