# go-mcs-sdk

[![Made by FilSwan](https://img.shields.io/badge/made%20by-FilSwan-green.svg)](https://www.filswan.com/)
[![Chat on discord](https://img.shields.io/badge/join%20-discord-brightgreen.svg)](https://discord.com/invite/KKGhy8ZqzK)

# Table of Contents <!-- omit in toc -->

- [Introduction](#introduction)
    - [Functions](#Functions)
    - [Data Structures](#Data-Structures)
    - [Constants](#Constants)
- [Prerequisites](#Prerequisites)
- [Usage](#usage)
    - [Download SDK](#Download-SDK)
    - [Call SDK](#Call-SDK)
    - [Documentation](#documentation)
- [MCS API](#mcs-api)
- [Contributing](#contributing)
- [Sponsors](#Sponsors)

## Introduction

A Golang software development kit for the [Multi-Chain Storage (MCS) Service](https://multichain.storage/) . It provides a
convenient interface for working with the MCS API. 

### Functions:

- [User Functions](https://github.com/filswan/go-mcs-sdk/blob/dev/mcs/api/docs/user.md)
- [Bucket Functions](https://github.com/filswan/go-mcs-sdk/blob/dev/mcs/api/docs/bucket.md)
- [On-chain Functions](https://github.com/filswan/go-mcs-sdk/blob/dev/mcs/api/docs/on-chain.md)

### Data Structures:
- [Struct](https://github.com/filswan/go-mcs-sdk/blob/dev/mcs/api/docs/struct.md)

### Constants:
- [Constants](https://github.com/filswan/go-mcs-sdk/blob/dev/mcs/api/common/constants/constants.go)

## Prerequisites
- [Metamask Wallet](https://docs.filswan.com/getting-started/beginner-walkthrough/public-testnet/setup-metamask)
- [Polygon RPC](https://www.alchemy.com/)
- [USDC and MATIC balance]
- [Optional: apikey](https://multichain.storage/) -> Setting -> Create API Key

## Usage

### Download SDK
```
go get github.com/filswan/go-mcs-sdk
```


### Call SDK
1. Login using apikey/accessToken
```go
import (
	"github.com/filswan/go-mcs-sdk/mcs/api/user"
)

mcsClient, err := user.LoginByApikeyV2(apikey, network)
// apikey: your apikey
// network: (MCS mainnet: mainnet, MCS testnet: testnet . Default is the main network.)

// mcsClient: result including the information to access the other API(s)
// err: when err generated while accessing this api, the error info will store in err
```
- See [Constants](#Constants) to show optional network

2. Call `bucket` related api(s)
- Step :one: Change `McsClient` to `BucketClient`
```go
import (
	"github.com/filswan/go-mcs-sdk/mcs/api/bucket"
)
bucketClient := bucket.GetBucketClient(*mcsClient)
```
- Step :two: Create a bucket
```go
bucketUid, err := buketClient.CreateBucket([BUCKET_NAME])
// bucketUid: the new created bucket UID
// err: when err generated while accessing this api, the error info will store in err
```
- Step :three: Upload a file to the bucket
```go
// err := buketClient.UploadFile([BUCKET_NAME], [YOUR_FILE_NAME], [YOUR_FILE_PATH], true)
// err: when err generated while accessing this api, the error info will store in err
```


### Documentation

For more examples please see the [SDK documentation](https://docs.filswan.com/multi-chain-storage/developer-quickstart/sdk)

## MCS API

For more information about the API usage, check out the MCS API
documentation (https://docs.filswan.com/development-resource/mcp-api).

## Contributing

Feel free to join in and discuss. Suggestions are welcome! [Open an issue](https://github.com/filswan/go-mcs-sdk/issues) or [Join the Discord](https://discord.com/invite/KKGhy8ZqzK)!

## Sponsors

This project is sponsored by Filecoin Foundation

[Flink SDK - A data provider offers Chainlink Oracle service for Filecoin Network ](https://github.com/filecoin-project/devgrants/issues/463)

<img src="https://github.com/filswan/flink/blob/main/filecoin.png" width="200">
