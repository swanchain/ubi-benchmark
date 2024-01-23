package user

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//
//func TestLoginByApikey(t *testing.T) {
//	apikey := ""
//	accessToken := ""
//
//	mcsClient, err := LoginByApikey(apikey, accessToken, network)
//	assert.Nil(t, err)
//	assert.NotNil(t, mcsClient)
//	assert.NotEmpty(t, mcsClient.BaseUrl)
//	assert.NotEmpty(t, mcsClient.JwtToken)
//}
//
//func TestRegister(t *testing.T) {
//	nonce, err := Register("0xbE14Eb1ffcA54861D3081560110a45F4A1A9e9c5", network)
//	assert.Nil(t, err)
//	assert.NotEmpty(t, nonce)
//	logs.GetLogger().Info(*nonce)
//}
//
//func TestLoginByPublicKeySignature(t *testing.T) {
//	mcsClient, err := LoginByPublicKeySignature(
//		"1067049846399020981103631740110767813482",
//		"0xbE14Eb1ffcA54861D3081560110a45F4A1A9e9c5",
//		"0xff93680ae74eaccc9858ef12a83592038d6b4bf6e2ef166f792cd14f8247bb1d22c01bdfb496f798c7574342ea3d919c15c4af137932e46c5bca7873e7d4ae121c",
//		network)
//
//	assert.Nil(t, err)
//	assert.NotNil(t, mcsClient)
//	assert.NotEmpty(t, mcsClient.BaseUrl)
//	assert.NotEmpty(t, mcsClient.JwtToken)
//}

func TestLoginByApikeyV2(t *testing.T) {
	fmt.Print("TestLoginByApikeyV2")
	mcsClient, err := LoginByApikeyV2("MCS_28ede6fe0e753a331584d3a0", "testnet")
	assert.Nil(t, err)
	assert.NotNil(t, mcsClient)
	assert.NotEmpty(t, mcsClient.BaseUrl)
	assert.NotEmpty(t, mcsClient.JwtToken)
}
