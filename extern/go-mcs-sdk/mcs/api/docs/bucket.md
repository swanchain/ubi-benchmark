* [GetBucketClient](#GetBucketClient)
* [ListBuckets](#ListBuckets)
* [CreateBucket](#CreateBucket)
* [DeleteBucket](#DeleteBucket)
* [GetBucket](#GetBucket)
* [GetBucketUid](#GetBucketUid)
* [RenameBucket](#RenameBucket)
* [GetTotalStorageSize](#GetTotalStorageSize)
* [GetGateway](#GetGateway)
* [GetFile](#GetFile)
* [CreateFolder](#CreateFolder)
* [DeleteFile](#DeleteFile)
* [ListFiles](#ListFiles)
* [UploadFile](#UploadFile)
* [UploadFolder](#UploadFolder)
* [GetFileInfo](#GetFileInfo)
* [UploadIpfsFolder](#UploadIpfsFolder)
* [DownloadFile](#DownloadFile)

## GetBucketClient

Definition:

```shell
func GetBucketClient(mcsClient user.McsClient) *BucketClient
```

Outputs:

```shell
*BucketClient  # includes jwt token and other information for use when call the other apis
```

## ListBuckets

Definition:

```shell
func (bucketClient *BucketClient) ListBuckets() ([]*Bucket, error)
```

Outputs:

```shell
[]*Bucket  # the bucket list belong to current user
error      # error or nil
```

## CreateBucket

Definition:

```shell
func (bucketClient *BucketClient) CreateBucket(bucketName string) (*string, error)
```

Outputs:

```shell
*string  # bucket uid
error    # error or nil
```

## DeleteBucket

Definition:

```shell
func (bucketClient *BucketClient) DeleteBucket(bucketName string) error
```

Outputs:

```shell
error    # error or nil
```

## GetBucket

Definition:

```shell
func (bucketClient *BucketClient) GetBucket(bucketName, bucketUid string) (*Bucket, error)
```

Outputs:

```shell
*Bucket  # bucket info whose bucket name or bucket uid is the same as the parameter
error    # error or nil
```

## GetBucketUid

Definition:

```shell
func (bucketClient *BucketClient) GetBucketUid(bucketName string) (*string, error)
```

Outputs:

```shell
*string  # bucket uid whose bucket name is the same as the parameter
error    # error or nil
```

## RenameBucket

Definition:

```shell
func (bucketClient *BucketClient) RenameBucket(newBucketName string, bucketUid string) error
```

Outputs:

```shell
error    # error or nil
```

## GetTotalStorageSize

Definition:

```shell
func (bucketClient *BucketClient) GetTotalStorageSize() (*int64, error)
```

Outputs:

```shell
*int64   # total storage size
error    # error or nil
```

## GetGateway

Definition:

```shell
func (bucketClient *BucketClient) GetGateway() (*string, error)
```

Outputs:

```shell
*string   # ipfs gateway
error     # error or nil
```

## GetFile

Definition:

```shell
func (bucketClient *BucketClient) GetFile(bucketName, objectName string) (*OssFile, error)
```

Outputs:

```shell
*OssFile   # file with the object name in the bucket whose name is the same as parameter
error      # error or nil
```

## CreateFolder

Definition:

```shell
func (bucketClient *BucketClient) CreateFolder(bucketName, folderName, prefix string) (*string, error)
```

Outputs:

```shell
*string   # folder name
error      # error or nil
```

## DeleteFile

Definition:

```shell
func (bucketClient *BucketClient) DeleteFile(bucketName, objectName string) error
```

Outputs:

```shell
error      # error or nil
```

## ListFiles

Definition:

```shell
func (bucketClient *BucketClient) ListFiles(bucketName, prefix string, limit, offset int) ([]*OssFile, *int, error)
```

Outputs:

```shell
[]*OssFile # file list
*int       # total file count
error      # error or nil
```

## UploadFile

Definition:

```shell
func (bucketClient *BucketClient) UploadFile(bucketName, objectName, filePath string, replace bool) error
```

Outputs:

```shell
error      # error or nil
```

## UploadFolder

Definition:

```shell
func (bucketClient *BucketClient) UploadFolder(bucketName, folderPath, prefix string) error
```

Outputs:

```shell
error      # error or nil
```

## GetFileInfo

Definition:

```shell
func (bucketClient *BucketClient) GetFileInfo(fileId int) (*OssFile, error)
```

Outputs:

```shell
*OssFile   # file info
error      # error or nil
```

## UploadIpfsFolder

Definition:

```shell
func (bucketClient *BucketClient) UploadIpfsFolder(bucketName, objectName, folderPath string) (*OssFile, error)
```

Outputs:

```shell
*OssFile   # file info
error      # error or nil
```

## DownloadFile

Definition:

```shell
func (bucketClient *BucketClient) DownloadFile(bucketName, objectName, destFileDir string) error
```

Outputs:

```shell
*OssFile   # file info
error      # error or nil
```
