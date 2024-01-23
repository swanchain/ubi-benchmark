package user

import (
	"github.com/filswan/go-mcs-sdk/mcs/api/common/constants"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/logs"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/utils"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/web"
)

type McsClient struct {
	BaseUrl  string `json:"base_url"`
	JwtToken string `json:"jwt_token"`
}

func getBaseApiUrl(network string) (string, string) {
	apiUrlBase := ""
	switch network {
	case constants.PAYMENT_CHAIN_NAME_POLYGON_MAINNET:
		apiUrlBase = constants.API_URL_MCS_POLYGON_MAINNET
	case constants.PAYMENT_CHAIN_NAME_POLYGON_MUMBAI:
		apiUrlBase = constants.API_URL_MCS_POLYGON_MUMBAI
	case constants.PAYMENT_CHAIN_NAME_BSC_TESTNET:
		apiUrlBase = constants.API_URL_MCS_BSC_TESTNET
	default:
		apiUrlBase = constants.API_URL_MCS_POLYGON_MAINNET
		network = constants.PAYMENT_CHAIN_NAME_POLYGON_MAINNET
	}

	return apiUrlBase, network
}

func LoginByApikey(apikey, accessToken, network string) (*McsClient, error) {
	apiUrlBase, network := getBaseApiUrl(network)

	var params struct {
		Apikey      string `json:"apikey" binding:"required,min=1,max=100"`
		AccessToken string `json:"access_token" binding:"required,min=1,max=100"`
		Network     string `json:"network" binding:"required,min=1,max=65535"`
	}

	params.Apikey = apikey
	params.AccessToken = accessToken
	params.Network = network

	apiUrl := utils.UrlJoin(apiUrlBase, constants.API_URL_USER_LOGIN_BY_APIKEY)

	var loginByApikeyResponse struct {
		JwtToken string `json:"jwt_token"`
	}

	err := web.HttpPost(apiUrl, "", params, &loginByApikeyResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	mcsClient := McsClient{
		BaseUrl:  apiUrlBase,
		JwtToken: loginByApikeyResponse.JwtToken,
	}

	return &mcsClient, nil
}

func Register(publicKeyAddress, network string) (*string, error) {
	apiUrlBase, _ := getBaseApiUrl(network)

	var params struct {
		PublicKeyAddress string `json:"public_key_address" binding:"required,min=1,max=100"`
	}

	params.PublicKeyAddress = publicKeyAddress

	apiUrl := utils.UrlJoin(apiUrlBase, constants.API_URL_USER_REGISTER)

	var response struct {
		Nonce string `json:"nonce"`
	}

	err := web.HttpPost(apiUrl, "", params, &response)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &response.Nonce, nil
}

func LoginByPublicKeySignature(nonce, publicKeyAddress, signature, network string) (*McsClient, error) {
	apiUrlBase, network := getBaseApiUrl(network)

	var params struct {
		Nonce            string `json:"nonce" binding:"required,min=1,max=100"`
		PublicKeyAddress string `json:"public_key_address" binding:"required,min=1,max=100"`
		Signature        string `json:"signature" binding:"required,min=1,max=100"`
		Network          string `json:"network" binding:"required,min=1,max=65535"`
	}

	params.Nonce = nonce
	params.PublicKeyAddress = publicKeyAddress
	params.Signature = signature
	params.Network = network

	apiUrl := utils.UrlJoin(apiUrlBase, constants.API_URL_USER_LOGIN)

	var response struct {
		JwtToken string `json:"jwt_token"`
	}

	err := web.HttpPost(apiUrl, "", params, &response)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	mcsClient := McsClient{
		BaseUrl:  apiUrlBase,
		JwtToken: response.JwtToken,
	}

	return &mcsClient, nil
}

func GetBaseApiUrlV2(network string) (apiUrlBase string) {
	switch network {
	case constants.MCS_NETWORK_VERSION_TESTNET:
		apiUrlBase = constants.API_URL_MCS_POLYGON_MUMBAI
	case constants.MCS_NETWORK_VERSION_MAINNET:
		apiUrlBase = constants.API_URL_MCS_POLYGON_MAINNET
	default:
		apiUrlBase = constants.API_URL_MCS_POLYGON_MAINNET
	}

	return apiUrlBase
}

func LoginByApikeyV2(apikey string, network string) (*McsClient, error) {
	apiUrlBase := GetBaseApiUrlV2(network)
	var params struct {
		Apikey string `json:"apikey" binding:"required,min=1,max=100"`
	}

	params.Apikey = apikey

	apiUrl := utils.UrlJoin(apiUrlBase, constants.API_URL_USER_LOGIN_BY_APIKEY_V2)

	var result string

	err := web.HttpPost(apiUrl, "", params, &result)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	//if loginByApikeyResponse.Status != "success" {
	//	return nil, errors.Errorf("Login failed")
	//}
	//
	mcsClient := McsClient{
		BaseUrl:  apiUrlBase,
		JwtToken: result,
	}
	return &mcsClient, nil
}
