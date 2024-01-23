package bucket

import (
	"fmt"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/constants"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/utils"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/web"
	"github.com/filswan/go-mcs-sdk/mcs/api/user"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/logs"
)

type BucketClient struct {
	user.McsClient
}

func GetBucketClient(mcsClient user.McsClient) *BucketClient {
	var bucketClient = &BucketClient{}

	bucketClient.BaseUrl = mcsClient.BaseUrl
	bucketClient.JwtToken = mcsClient.JwtToken

	return bucketClient
}

type Bucket struct {
	BucketUid  string `json:"bucket_uid"`
	Address    string `json:"address"`
	MaxSize    int64  `json:"max_size"`
	Size       int64  `json:"size"`
	IsFree     bool   `json:"is_free"`
	PaymentTx  string `json:"payment_tx"`
	IsActive   bool   `json:"is_active"`
	IsDeleted  bool   `json:"is_deleted"`
	BucketName string `json:"bucket_name"`
	FileNumber int64  `json:"file_number"`
}

func (bucketClient *BucketClient) ListBuckets() ([]*Bucket, error) {
	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_GET_BUCKET_LIST)

	var buckets []*Bucket
	err := web.HttpGet(apiUrl, bucketClient.JwtToken, nil, &buckets)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return buckets, nil
}
func (bucketClient *BucketClient) CreateBucket(bucketName string) (*string, error) {
	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_CREATE_BUCKET)

	var bucket struct {
		BucketName string `json:"bucket_name"`
	}
	bucket.BucketName = bucketName

	var bucketUid string
	err := web.HttpPost(apiUrl, bucketClient.JwtToken, &bucket, &bucketUid)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &bucketUid, nil
}

func (bucketClient *BucketClient) DeleteBucket(bucketName string) error {
	bucketUid, err := bucketClient.GetBucketUid(bucketName)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_DELETE_BUCKET)
	apiUrl = apiUrl + "?bucket_uid=" + *bucketUid

	err = web.HttpGet(apiUrl, bucketClient.JwtToken, nil, nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}

func (bucketClient *BucketClient) GetBucket(bucketName, bucketUid string) (*Bucket, error) {
	buckets, err := bucketClient.ListBuckets()
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	for _, bucket := range buckets {
		if (bucketName == "" || bucketName != "" && bucket.BucketName == bucketName) &&
			(bucketUid == "" || bucketUid != "" && bucket.BucketUid == bucketUid) {
			return bucket, nil
		}
	}

	err = fmt.Errorf("bucket not exists")

	return nil, err
}

func (bucketClient *BucketClient) GetBucketUid(bucketName string) (*string, error) {
	bucket, err := bucketClient.GetBucket(bucketName, "")
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &bucket.BucketUid, nil
}

func (bucketClient *BucketClient) RenameBucket(newBucketName string, bucketUid string) error {
	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_RENAME_BUCKET)

	var params struct {
		BucketName string `json:"bucket_name"`
		BucketUid  string `json:"bucket_uid"`
	}
	params.BucketName = newBucketName
	params.BucketUid = bucketUid

	var result string
	err := web.HttpPost(apiUrl, bucketClient.JwtToken, &params, &result)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}

func (bucketClient *BucketClient) GetTotalStorageSize() (*int64, error) {
	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_GET_TOTAL_STORAGE_SIZE)

	var totalStorageSize int64

	err := web.HttpGet(apiUrl, bucketClient.JwtToken, nil, &totalStorageSize)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &totalStorageSize, nil
}
