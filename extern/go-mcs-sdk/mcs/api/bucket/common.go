package bucket

import (
	"fmt"
	"strings"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/constants"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/utils"
	"github.com/filswan/go-mcs-sdk/mcs/api/common/web"

	"github.com/filswan/go-mcs-sdk/mcs/api/common/logs"
)

func (bucketClient *BucketClient) GetGateway() (*string, error) {
	apiUrl := utils.UrlJoin(bucketClient.BaseUrl, constants.API_URL_BUCKET_GATEWAY_GET_GATEWAY)

	var subDomains []string

	err := web.HttpGet(apiUrl, bucketClient.JwtToken, nil, &subDomains)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if len(subDomains) <= 0 {
		err := fmt.Errorf("no gateway returned")
		logs.GetLogger().Error(err)
		return nil, err
	}

	gateway := subDomains[0]
	if !strings.HasPrefix(gateway, "http") {
		gateway = "https://" + gateway
	}

	return &gateway, nil
}
