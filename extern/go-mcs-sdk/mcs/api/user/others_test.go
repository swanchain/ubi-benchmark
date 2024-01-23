package user

//
//import (
//	"strings"
//	"testing"
//
//	"github.com/filswan/go-mcs-sdk/mcs/api/common/constants"
//	"github.com/filswan/go-mcs-sdk/mcs/api/common/logs"
//
//	"github.com/stretchr/testify/assert"
//)
//
//var mcsClient *McsClient
//var network = constants.PAYMENT_CHAIN_NAME_POLYGON_MUMBAI
//
//func init() {
//	if mcsClient != nil {
//		return
//	}
//
//	apikey := ""
//	accessToken := ""
//
//	var err error
//	mcsClient, err = LoginByApikey(apikey, accessToken, network)
//	if err != nil {
//		logs.GetLogger().Fatal(err)
//	}
//}
//
//func TestCheckLogin(t *testing.T) {
//	networkName, walletAddress, err := mcsClient.CheckLogin()
//	assert.Nil(t, err)
//	assert.NotEmpty(t, networkName)
//	assert.NotEmpty(t, walletAddress)
//	assert.Equal(t, network, *networkName)
//	assert.Contains(t, strings.ToUpper(*walletAddress), "0X")
//}
//
//func TestGenerateApikey(t *testing.T) {
//	apikey, accessToken, err := mcsClient.GenerateApikey(30)
//	assert.Nil(t, err)
//	assert.NotEmpty(t, apikey)
//	assert.NotEmpty(t, accessToken)
//}
//
//func TestDeleteApikey(t *testing.T) {
//	err := mcsClient.DeleteApikey("2dkFLDsWNYDTkZkz6qB6PG")
//	assert.Nil(t, err)
//}
//
//func TestRegisterEmail(t *testing.T) {
//	response, err := mcsClient.RegisterEmail("fchen@nbai.io")
//	assert.Nil(t, err)
//	assert.NotEmpty(t, response)
//}
//
//func TestGetApikeys(t *testing.T) {
//	apikeys, err := mcsClient.GetApikeys()
//	assert.Nil(t, err)
//	assert.NotEmpty(t, apikeys)
//
//	for _, apikey := range apikeys {
//		logs.GetLogger().Info(*apikey)
//	}
//}
//
//func TestGetWallet(t *testing.T) {
//	wallet, err := mcsClient.GetWallet()
//	assert.Nil(t, err)
//	assert.NotEmpty(t, wallet)
//
//	logs.GetLogger().Info(*wallet)
//}
//
//func TestSetPopupTime(t *testing.T) {
//	err := mcsClient.SetPopupTime()
//	assert.Nil(t, err)
//}
//
//func TestDeleteEmail(t *testing.T) {
//	err := mcsClient.DeleteEmail()
//	assert.Nil(t, err)
//}
