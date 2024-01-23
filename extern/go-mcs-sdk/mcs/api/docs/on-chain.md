* [GetOnChainClient](#GetOnChainClient)
* [GetSystemParam](#GetSystemParam)
* [GetHistoricalAveragePriceVerified](#GetHistoricalAveragePriceVerified)
* [GetAmount](#GetAmount)
* [Upload](#Upload)
* [GetMintInfo](#GetMintInfo)
* [GetUserTaskDeals](#GetUserTaskDeals)
* [GetDealDetail](#GetDealDetail)
* [GetDealLogs](#GetDealLogs)
* [GetSourceFileUpload](#GetSourceFileUpload)
* [UnpinSourceFile](#UnpinSourceFile)
* [WriteNftCollection](#WriteNftCollection)
* [GetNftCollections](#GetNftCollections)
* [RecordMintInfo](#RecordMintInfo)
* [Pay](#Pay)
* [GetPaymentInfo](#GetPaymentInfo)
* [GetFileCoinPrice](#GetFileCoinPrice)
* [GetBillingHistory](#GetBillingHistory)
* [GetDeals2PreSign](#GetDeals2PreSign)
* [GetDeals2Sign](#GetDeals2Sign)
* [GetDeals2SignHash](#GetDeals2SignHash)

## GetOnChainClient

Definition:

```shell
func GetOnChainClient(mcsClient user.McsClient) *OnChainClient
```

Outputs:

```shell
*OnChainClient  # includes jwt token and other information for use when call the other apis
```
## GetSystemParam

Definition:

```shell
func (onChainClient *OnChainClient) GetSystemParam() (*SystemParam, error)
```

Outputs:

```shell
*SystemParam  # system parameters
error         # error or nil
```

## GetHistoricalAveragePriceVerified

Definition:

```shell
func GetHistoricalAveragePriceVerified() (float64, error)
```

Outputs:

```shell
float64    # historical average verified price
error      # error or nil
```

## GetAmount

Definition:

```shell
func GetAmount(fizeSizeByte int64, historicalAveragePriceVerified, fileCoinPrice float64, copyNumber int) (int64, error)
```

Outputs:

```shell
int64    # amount to pay
error    # error or nil
```

## Upload

Definition:

```shell
func (onChainClient *OnChainClient) Upload(filePath string, fileType int) (*UploadFile, error)
```

Outputs:

```shell
*UploadFile # upload file information
error       # error or nil
```

## GetMintInfo

Definition:

```shell
func (onChainClient *OnChainClient) GetMintInfo(sourceFileUploadId int64) ([]*SourceFileMintOut, error)
```

Outputs:

```shell
*SourceFileMintOut # file mint info
error              # error or nil
```

## GetUserTaskDeals

Definition:

```shell
func (onChainClient *OnChainClient) GetUserTaskDeals(dealsParams DealsParams) ([]*Deal, *int64, error)
```

Outputs:

```shell
[]*Deal # deal list
*int64  # total count
error   # error or nil
```

## GetDealDetail

Definition:

```shell
func (onChainClient *OnChainClient) GetDealDetail(sourceFileUploadId, dealId int64) (*SourceFileUploadDeal, []*DaoSignature, *int, error)
```

Outputs:

```shell
*SourceFileUploadDeal # deal list
[]*DaoSignature       # dao signature list
*int                  # dao threshold
error                 # error or nil
```

## GetDealLogs

Definition:

```shell
func (onChainClient *OnChainClient) GetDealLogs(offlineDealId int64) ([]*OfflineDealLog, error)
```

Outputs:

```shell
[]*OfflineDealLog  # deal logs
error              # error or nil
```

## GetSourceFileUpload

Definition:

```shell
func (onChainClient *OnChainClient) GetSourceFileUpload(sourceFileUploadId int64) (*SourceFileUpload, error)
```

Outputs:

```shell
*SourceFileUpload  # source file upload information
error              # error or nil
```

## UnpinSourceFile

Definition:

```shell
func (onChainClient *OnChainClient) UnpinSourceFile(sourceFileUploadId int64) error
```

Outputs:

```shell
error              # error or nil
```

## WriteNftCollection

Definition:

```shell
func (onChainClient *OnChainClient) WriteNftCollection(nftCollectionParams NftCollectionParams) error
```

Outputs:

```shell
error              # error or nil
```

## GetNftCollections

Definition:

```shell
func (onChainClient *OnChainClient) GetNftCollections() ([]*NftCollection, error)
```

Outputs:

```shell
[]*NftCollection   # NFT collections
error              # error or nil
```

## RecordMintInfo

Definition:

```shell
func (onChainClient *OnChainClient) RecordMintInfo(recordMintInfoParams *RecordMintInfoParams) (*SourceFileMint, error)
```

Outputs:

```shell
*SourceFileMint   # Mint info
error             # error or nil
```

## Pay

Definition:

```shell
func (onChainClient *OnChainClient) Pay(sourceFileUploadId int64, privateKeyStr string, rpcUrl string) (*string, error)
```

Outputs:

```shell
*string   # payment transaction hash
error     # error or nil
```

## GetPaymentInfo

Definition:

```shell
func (client *OnChainClient) GetPaymentInfo(fileUploadId int64) (*LockPaymentInfo, error)
```

Outputs:

```shell
*LockPaymentInfo # payment information
error            # error or nil
```

## GetFileCoinPrice

Definition:

```shell
func (onChainClient *OnChainClient) GetFileCoinPrice() (*float64, error)
```

Outputs:

```shell
*float64 # filecoin price
error    # error or nil
```

## GetBillingHistory

Definition:

```shell
func (onChainClient *OnChainClient) GetBillingHistory(billingHistoryParams BillingHistoryParams) ([]*BillingHistory, *int64, error)
```

Outputs:

```shell
[]*BillingHistory # billing list
*int64            # total record number
error             # error or nil
```

## GetDeals2PreSign

Definition:

```shell
func (onChainClient *OnChainClient) GetDeals2PreSign() ([]*Deal2PreSign, error)
```

Outputs:

```shell
[]*Deal2PreSign # deals to pre sign
error           # error or nil
```

## GetDeals2Sign

Definition:

```shell
func (onChainClient *OnChainClient) GetDeals2Sign() ([]*Deal2Sign, error)
```

Outputs:

```shell
[]*Deal2Sign # deals to sign
error           # error or nil
```

## GetDeals2SignHash

Definition:

```shell
func (onChainClient *OnChainClient) GetDeals2SignHash() ([]*Deal2Sign, error)
```

Outputs:

```shell
[]*Deal2Sign # deals to sign hash
error        # error or nil
```
