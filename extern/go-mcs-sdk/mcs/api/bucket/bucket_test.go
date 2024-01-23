package bucket

import (
	"testing"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/logs"

	"github.com/stretchr/testify/assert"
)

func TestListBuckets(t *testing.T) {
	buckets, err := buketClient.ListBuckets()
	assert.Nil(t, err)
	assert.NotEmpty(t, buckets)

	for _, bucket := range buckets {
		logs.GetLogger().Info(*bucket)
	}
}

func TestCreateBucket(t *testing.T) {
	bucketUid, err := buketClient.CreateBucket("test23")
	assert.Nil(t, err)
	assert.NotEmpty(t, bucketUid)

	logs.GetLogger().Info(*bucketUid)
}

func TestDeleteBucket(t *testing.T) {
	err := buketClient.DeleteBucket("abc")
	assert.Nil(t, err)
}

func TestGetBucket(t *testing.T) {
	bucket, err := buketClient.GetBucket("test23", "")
	assert.Nil(t, err)
	assert.NotNil(t, bucket)
	logs.GetLogger().Info(*bucket)
}

func TestGetBucketUid(t *testing.T) {
	bucketUid, err := buketClient.GetBucketUid("test23")
	assert.Nil(t, err)
	assert.NotNil(t, bucketUid)
	assert.NotEmpty(t, bucketUid)
	logs.GetLogger().Info(*bucketUid)
}

func TestRenameBucket(t *testing.T) {
	err := buketClient.RenameBucket("aaa", "31a04949-1e3f-494b-8285-516d2a048322")
	assert.Nil(t, err)
}

func TestGetTotalStorageSize(t *testing.T) {
	totalStorageSize, err := buketClient.GetTotalStorageSize()
	assert.Nil(t, err)
	assert.NotNil(t, totalStorageSize)

	logs.GetLogger().Info(*totalStorageSize)
}
