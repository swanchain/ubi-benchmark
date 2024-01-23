package bucket

import (
	"testing"

	"github.com/filswan/go-mcs-sdk/mcs/api/user"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/logs"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/constants"

	"github.com/stretchr/testify/assert"
)

var buketClient *BucketClient
var network = constants.MCS_NETWORK_VERSION_TESTNET
var apikey = "MCS_28ede6fe0e753a331584d3a0"
var file2Upload = "/Users/daniel/code/go-code/go-mcs-sdk/log_mcs.png"
var folder2Upload = "/Users/dorachen/work/test3"

func init() {
	if buketClient != nil {
		return
	}

	mcsClient, err := user.LoginByApikeyV2(apikey, network)
	if err != nil {
		logs.GetLogger().Fatal(err)
	}

	buketClient = GetBucketClient(*mcsClient)
}

func TestGetGateway(t *testing.T) {
	gateway, err := buketClient.GetGateway()
	assert.Nil(t, err)
	assert.NotEmpty(t, gateway)

	logs.GetLogger().Info(*gateway)
}
