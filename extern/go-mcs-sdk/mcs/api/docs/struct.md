* [McsClient](#McsClient)
* [Apikey](#Apikey)
* [Wallet](#Wallet)
* [BucketClient](#BucketClient)
* [Bucket](#Bucket)
* [OssFile](#OssFile)
* [OnChainClient](#OnChainClient)
* [SystemParam](#SystemParam)
* [UploadFile](#UploadFile)
* [OfflineDeal](#OfflineDeal)
* [Deal](#Deal)
* [DealsParams](#DealsParams)
* [SourceFileUploadDeal](#SourceFileUploadDeal)
* [DaoSignature](#DaoSignature)
* [OfflineDealLog](#OfflineDealLog)
* [SourceFileUpload](#SourceFileUpload)
* [NftCollectionParams](#NftCollectionParams)
* [NftCollection](#NftCollection)
* [RecordMintInfoParams](#RecordMintInfoParams)
* [SourceFileMint](#SourceFileMint)
* [SourceFileMintOut](#SourceFileMintOut)
* [LockPaymentInfo](#LockPaymentInfo)
* [BillingHistory](#BillingHistory)
* [BillingHistoryParams](#BillingHistoryParams)
* [Deal2PreSign](#Deal2PreSign)
* [Deal2SignBatchInfo](#Deal2SignBatchInfo)
* [Deal2Sign](#Deal2Sign)

## McsClient
```
type McsClient struct {
	BaseUrl  string `json:"base_url"`
	JwtToken string `json:"jwt_token"`
}
```

## Apikey
```
type Apikey struct {
	ID          int64  `json:"id"`
	WalletId    int64  `json:"wallet_id"`
	ApiKey      string `json:"api_key"`
	AccessToken string `json:"access_token"`
	ValidDays   int32  `json:"valid_days"`
	CreateAt    int64  `json:"create_at"`
	UpdateAt    int64  `json:"update_at"`
}
```

## Wallet
```
type Wallet struct {
	ID           int64   `json:"id"`
	Address      string  `json:"address"`
	Email        *string `json:"email"`
	EmailStatus  *int    `json:"email_status"`
	EmailPopupAt *int64  `json:"email_popup_at"`
	CreateAt     int64   `json:"create_at"`
	UpdateAt     int64   `json:"update_at"`
}
```

## BucketClient
```
type BucketClient struct {
	user.McsClient
}
```

## Bucket
```
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
```

## OssFile
```
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
```

## OnChainClient
```
type OnChainClient struct {
	BaseUrl  string `json:"base_url"`
	JwtToken string `json:"jwt_token"`
}
```

## SystemParam
```
type SystemParam struct {
	ChainName                   string  `json:"chain_name"`
	PaymentContractAddress      string  `json:"payment_contract_address"`
	PaymentRecipientAddress     string  `json:"payment_recipient_address"`
	DaoContractAddress          string  `json:"dao_contract_address"`
	DexAddress                  string  `json:"dex_address"`
	UsdcWFilPoolContract        string  `json:"usdc_wFil_pool_contract"`
	DefaultNftCollectionAddress string  `json:"default_nft_collection_address"`
	NftCollectionFactoryAddress string  `json:"nft_collection_factory_address"`
	UsdcAddress                 string  `json:"usdc_address"`
	GasLimit                    uint64  `json:"gas_limit"`
	LockTime                    int     `json:"lock_time"`
	PayMultiplyFactor           float32 `json:"pay_multiply_factor"`
	DaoThreshold                int     `json:"dao_threshold"`
	FilecoinPrice               float64 `json:"filecoin_price"`
}
```

## UploadFile
```
type UploadFile struct {
	SourceFileUploadId int64  `json:"source_file_upload_id"`
	PayloadCid         string `json:"payload_cid"`
	IpfsUrl            string `json:"ipfs_url"`
	FileSize           int64  `json:"file_size"`
	WCid               string `json:"w_cid"`
	Status             string `json:"status"`
}
```

## OfflineDeal
```
type OfflineDeal struct {
	Id             int64   `json:"id"`
	CarFileId      int64   `json:"car_file_id"`
	DealCid        string  `json:"deal_cid"`
	MinerId        int64   `json:"miner_id"`
	Verified       bool    `json:"verified"`
	StartEpoch     int     `json:"start_epoch"`
	SenderWalletId int64   `json:"sender_wallet_id"`
	Status         string  `json:"status"`
	DealId         *int64  `json:"deal_id"`
	OnChainStatus  *string `json:"on_chain_status"`
	UnlockTxHash   *string `json:"unlock_tx_hash"`
	UnlockAt       *int64  `json:"unlock_at"`
	Note           *string `json:"note"`
	NetworkId      int64   `json:"network_id"`
	MinerFid       string  `json:"miner_fid"`
	CreateAt       int64   `json:"create_at"`
	UpdateAt       int64   `json:"update_at"`
}
```

## Deal
```
type Deal struct {
	SourceFileUploadId int64          `json:"source_file_upload_id"`
	FileName           string         `json:"file_name"`
	FileSize           int64          `json:"file_size"`
	UploadAt           int64          `json:"upload_at"`
	Duration           int            `json:"duration"`
	IpfsUrl            string         `json:"ipfs_url"`
	PinStatus          string         `json:"pin_status"`
	PayAmount          string         `json:"pay_amount"`
	Status             string         `json:"status"`
	Note               string         `json:"note"`
	IsFree             bool           `json:"is_free"`
	IsMinted           bool           `json:"is_minted"`
	RefundedBySelf     bool           `json:"refunded_by_self"`
	OfflineDeals       []*OfflineDeal `json:"offline_deal"`
}
```

## DealsParams
```
type DealsParams struct {
	PageNumber *int    `json:"page_number"`
	PageSize   *int    `json:"page_size"`
	FileName   *string `json:"file_name"`
	Status     *string `json:"status"`
	IsMinted   *string `json:"is_minted"`
	OrderBy    *string `json:"order_by"`
	IsAscend   *string `json:"is_ascend"`
}
```

## SourceFileUploadDeal
```
type SourceFileUploadDeal struct {
	DealID                   *int    `json:"deal_id"`
	DealCid                  *string `json:"deal_cid"`
	MessageCid               *string `json:"message_cid"`
	Height                   *int    `json:"height"`
	PieceCid                 *string `json:"piece_cid"`
	VerifiedDeal             *bool   `json:"verified_deal"`
	StoragePricePerEpoch     *int    `json:"storage_price_per_epoch"`
	Signature                *string `json:"signature"`
	SignatureType            *string `json:"signature_type"`
	CreatedAt                *int    `json:"created_at"`
	PieceSizeFormat          *string `json:"piece_size_format"`
	StartHeight              *int    `json:"start_height"`
	EndHeight                *int    `json:"end_height"`
	Client                   *string `json:"client"`
	ClientCollateralFormat   *string `json:"client_collateral_format"`
	Provider                 *string `json:"provider"`
	ProviderTag              *string `json:"provider_tag"`
	VerifiedProvider         *int    `json:"verified_provider"`
	ProviderCollateralFormat *string `json:"provider_collateral_format"`
	Status                   *int    `json:"status"`
	NetworkName              *string `json:"network_name"`
	StoragePrice             *int    `json:"storage_price"`
	IpfsUrl                  string  `json:"ipfs_url"`
	FileName                 string  `json:"file_name"`
	WCid                     string  `json:"w_cid"`
	CarFilePayloadCid        string  `json:"car_file_payload_cid"`
	LockedAt                 int64   `json:"locked_at"`
	LockedFee                string  `json:"locked_fee"`
	Unlocked                 bool    `json:"unlocked"`
}
```

## DaoSignature
```
type DaoSignature struct {
	WalletSigner string  `json:"wallet_signer"`
	TxHash       *string `json:"tx_hash"`
	Status       *string `json:"status"`
	CreateAt     *int64  `json:"create_at"`
}
```

## OfflineDealLog
```
type OfflineDealLog struct {
	Id             int64  `json:"id"`
	OfflineDealId  int64  `json:"offline_deal_id"`
	OnChainStatus  string `json:"on_chain_status"`
	OnChainMessage string `json:"on_chain_message"`
	CreateAt       int64  `json:"create_at"`
}
```

## SourceFileUpload
```
type SourceFileUpload struct {
	WCid     string `json:"w_cid"`
	Status   string `json:"status"`
	IsFree   bool   `json:"is_free"`
	FileSize int64  `json:"file_size"`
}
```

## NftCollectionParams
```
type NftCollectionParams struct {
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	ImageUrl        *string `json:"image_url"`
	ExternalLink    *string `json:"external_link"`
	SellerFee       *int    `json:"seller_fee"`
	WalletRecipient *string `json:"wallet_recipient"`
	TxHash          string  `json:"tx_hash"`
}
```

## NftCollection
```
type NftCollection struct {
	ID                int64   `json:"id"`
	Address           *string `json:"address"`
	WalletId          int64   `json:"wallet_id"`
	Name              string  `json:"name"`
	Description       *string `json:"description"`
	ImageUrl          *string `json:"image_url"`
	ExternalLink      *string `json:"external_link"`
	SellerFee         *int    `json:"seller_fee"`
	WalletIdRecipient *int64  `json:"wallet_id_recipient"`
	TxHash            string  `json:"tx_hash"`
	CreateAt          int64   `json:"create_at"`
	UpdateAt          int64   `json:"update_at"`
	WalletRecipient   string  `json:"wallet_recipient"`
	IsDefault         bool    `json:"is_default"`
}
```

## RecordMintInfoParams
```
type RecordMintInfoParams struct {
	SourceFileUploadId int64   `json:"source_file_upload_id"`
	NftCollectionId    int64   `json:"nft_collection_id"`
	TxHash             string  `json:"tx_hash"`
	TokenId            int64   `json:"token_id"`
	Name               *string `json:"name"`
	Description        *string `json:"description"`
}
```

## SourceFileMint
```
type SourceFileMint struct {
	ID                 int64   `json:"id"`
	SourceFileUploadId int64   `json:"source_file_upload_id"`
	NftTxHash          string  `json:"nft_tx_hash"`
	MintAddress        string  `json:"mint_address"`
	NftCollectionId    int64   `json:"nft_collection_id"`
	TokenId            int64   `json:"token_id"`
	Name               *string `json:"name"`
	Description        *string `json:"description"`
	CreateAt           int64   `json:"create_at"`
	UpdateAt           int64   `json:"update_at"`
}
```

## SourceFileMintOut
```
type SourceFileMintOut struct {
	SourceFileMint
	NftCollectionAddress  string  `json:"nft_collection_address"`
	NftCollectionName     *string `json:"nft_collection_name"`
	NftCollectionImageUrl *string `json:"nft_collection_image_url"`
}
```

## LockPaymentInfo
```
type LockPaymentInfo struct {
	WCid         string `json:"w_cid"`
	PayAmount    string `json:"pay_amount"`
	PayTxHash    string `json:"pay_tx_hash"`
	TokenAddress string `json:"token_address"`
}
```

## BillingHistory
```
type BillingHistory struct {
	PayId        int64  `json:"pay_id"`
	PayTxHash    string `json:"pay_tx_hash"`
	PayAmount    string `json:"pay_amount"`
	UnlockAmount string `json:"unlock_amount"`
	FileName     string `json:"file_name"`
	PayloadCid   string `json:"payload_cid"`
	PayAt        int64  `json:"pay_at"`
	UnlockAt     int64  `json:"unlock_at"`
	Deadline     int64  `json:"deadline"`
	NetworkName  string `json:"network_name"`
	TokenName    string `json:"token_name"`
}
```

## BillingHistoryParams
```
type BillingHistoryParams struct {
	PageNumber *int    `json:"page_number"`
	PageSize   *int    `json:"page_size"`
	FileName   *string `json:"file_name"`
	TxHash     *string `json:"tx_hash"`
	OrderBy    *string `json:"order_by"`
	IsAscend   *string `json:"is_ascend"`
}
```

## Deal2PreSign
```
type Deal2PreSign struct {
	DealId              int64 `json:"deal_id"`
	SourceFileUploadCnt int   `json:"source_file_upload_cnt"`
	BatchCount          int   `json:"batch_count"`
}
```

## Deal2SignBatchInfo
type Deal2SignBatchInfo struct {
	BatchNo int      `json:"batch_no"`
	WCid    []string `json:"w_cid"`
}

## Deal2Sign
type Deal2Sign struct {
	OfflineDealId int64                 `json:"offline_deal_id"`
	CarFileId     int64                 `json:"car_file_id"`
	DealId        int64                 `json:"deal_id"`
	BatchCount    int                   `json:"batch_count"`
	BatchSizeMax  int                   `json:"batch_size_max"`
	BatchInfo     []*Deal2SignBatchInfo `json:"batch_info"`
}
