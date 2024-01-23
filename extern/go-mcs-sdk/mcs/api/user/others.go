package user

import (
	"strconv"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/constants"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/utils"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/web"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/logs"
)

func (mcsClient *McsClient) CheckLogin() (*string, *string, error) {
	apiUrl := utils.UrlJoin(mcsClient.BaseUrl, constants.API_URL_USER_CHECK_LOGIN)

	var response struct {
		NetworkName   string `json:"network_name"`
		WalletAddress string `json:"wallet_address"`
	}

	err := web.HttpPost(apiUrl, mcsClient.JwtToken, nil, &response)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, nil, err
	}

	return &response.NetworkName, &response.WalletAddress, nil
}

func (mcsClient *McsClient) GenerateApikey(validDys int) (*string, *string, error) {
	apiUrl := utils.UrlJoin(mcsClient.BaseUrl, constants.API_URL_USER_GENERATE_APIKEY)
	apiUrl = apiUrl + "?valid_days=" + strconv.Itoa(validDys)

	var response struct {
		Apikey      string `json:"apikey"`
		AccessToken string `json:"access_token"`
	}

	err := web.HttpGet(apiUrl, mcsClient.JwtToken, nil, &response)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, nil, err
	}

	return &response.Apikey, &response.AccessToken, nil
}

func (mcsClient *McsClient) DeleteApikey(apikey string) error {
	apiUrl := utils.UrlJoin(mcsClient.BaseUrl, constants.API_URL_USER_DELETE_APIKEY)
	apiUrl = apiUrl + "?apikey=" + apikey

	err := web.HttpPut(apiUrl, mcsClient.JwtToken, nil, nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}

type Apikey struct {
	ID          int64  `json:"id"`
	WalletId    int64  `json:"wallet_id"`
	ApiKey      string `json:"api_key"`
	AccessToken string `json:"access_token"`
	ValidDays   int32  `json:"valid_days"`
	CreateAt    int64  `json:"create_at"`
	UpdateAt    int64  `json:"update_at"`
}

func (mcsClient *McsClient) GetApikeys() ([]*Apikey, error) {
	apiUrl := utils.UrlJoin(mcsClient.BaseUrl, constants.API_URL_USER_GET_APIKEYS)

	var response struct {
		Apikey []*Apikey `json:"apikey"`
	}

	err := web.HttpGet(apiUrl, mcsClient.JwtToken, nil, &response)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return response.Apikey, nil
}

func (mcsClient *McsClient) RegisterEmail(email string) (*string, error) {
	apiUrl := utils.UrlJoin(mcsClient.BaseUrl, constants.API_URL_USER_REGISTER_EMAIL)
	var params struct {
		Email string `json:"email"`
	}
	params.Email = email

	var response string
	err := web.HttpPost(apiUrl, mcsClient.JwtToken, &params, &response)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &response, nil
}

type Wallet struct {
	ID           int64   `json:"id"`
	Address      string  `json:"address"`
	Email        *string `json:"email"`
	EmailStatus  *int    `json:"email_status"`
	EmailPopupAt *int64  `json:"email_popup_at"`
	CreateAt     int64   `json:"create_at"`
	UpdateAt     int64   `json:"update_at"`
}

func (mcsClient *McsClient) GetWallet() (*Wallet, error) {
	apiUrl := utils.UrlJoin(mcsClient.BaseUrl, constants.API_URL_USER_GET_WALLET)

	var response struct {
		Wallet *Wallet `json:"wallet"`
	}

	err := web.HttpGet(apiUrl, mcsClient.JwtToken, nil, &response)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return response.Wallet, nil
}

func (mcsClient *McsClient) SetPopupTime() error {
	apiUrl := utils.UrlJoin(mcsClient.BaseUrl, constants.API_URL_USER_SET_POPUP_TIME)

	err := web.HttpPut(apiUrl, mcsClient.JwtToken, nil, nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}

func (mcsClient *McsClient) DeleteEmail() error {
	apiUrl := utils.UrlJoin(mcsClient.BaseUrl, constants.API_URL_USER_DELETE_EMAIL)

	err := web.HttpPut(apiUrl, mcsClient.JwtToken, nil, nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}
