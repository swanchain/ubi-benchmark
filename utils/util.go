package utils

import (
	"bytes"
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/client"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/beacon"
	lrand "github.com/filecoin-project/lotus/chain/rand"
	"github.com/filecoin-project/lotus/chain/types"
	cliutil "github.com/filecoin-project/lotus/cli/util"
	logging "github.com/ipfs/go-log/v2"
	"github.com/valyala/gozstd"
	"golang.org/x/xerrors"
	"io"
	"os"
	"time"
)

var log = logging.Logger("utils")

type NodeService struct {
	api    api.FullNode
	closer jsonrpc.ClientCloser
}

func NewNodeService(ctx context.Context, fullNodeApi string) (*NodeService, error) {
	ainfo := cliutil.ParseApiInfo(fullNodeApi)
	addr, err := ainfo.DialArgs("v1")
	if err != nil {
		return nil, err
	}

	nodeApi, closer, err := client.NewFullNodeRPCV1(ctx, addr, ainfo.AuthHeader())
	if err != nil {
		return nil, err
	}

	return &NodeService{
		api:    nodeApi,
		closer: closer,
	}, nil
}

func (s *NodeService) FullNodeAPI() api.FullNode {
	return s.api
}

func (s *NodeService) Close() error {
	if s.closer == nil {
		return xerrors.Errorf("Node services already closed")
	}
	s.closer()
	s.closer = nil
	return nil
}

func GetRandomness(minerId address.Address, tag crypto.DomainSeparationTag, height int64) ([]byte, error) {
	beaconEntries, err := generateBeaconEntry(abi.ChainEpoch(height))
	if err != nil {
		return nil, err
	}

	var rbase types.BeaconEntry
	if len(beaconEntries) > 0 {
		rbase = beaconEntries[len(beaconEntries)-1]
	}

	buf := new(bytes.Buffer)
	if err := minerId.MarshalCBOR(buf); err != nil {
		err = xerrors.Errorf("failed to marshal miner address: %w", err)
		return nil, err
	}

	round := abi.ChainEpoch(height)

	rand, err := lrand.DrawRandomnessFromBase(rbase.Data, tag, round, buf.Bytes())
	if err != nil {
		err = xerrors.Errorf("failed to get randomness for computing seal proof: %w", err)
		return nil, err
	}
	return rand, nil
}

func generateBeaconEntry(epoch abi.ChainEpoch) ([]types.BeaconEntry, error) {
	bSchedule := beacon.Schedule{{Start: 0, Beacon: beacon.NewMockBeacon(time.Second)}}
	randomBeacon := bSchedule.BeaconForEpoch(epoch)

	start := build.Clock.Now()
	maxRound := randomBeacon.MaxBeaconRoundForEpoch(network.Version16, epoch)

	cur := maxRound
	var out []types.BeaconEntry
	ctx := context.Background()
	rch := randomBeacon.Entry(ctx, cur)
	select {
	case resp := <-rch:
		if resp.Err != nil {
			return nil, xerrors.Errorf("beacon entry request returned error: %w", resp.Err)
		}

		out = append(out, resp.Entry)
		cur = resp.Entry.Round - 1
	case <-ctx.Done():
		return nil, xerrors.Errorf("context timed out waiting on beacon entry to come back for epoch %d: %w", epoch, ctx.Err())
	}
	log.Debugw("fetching beacon entries", "took", build.Clock.Since(start), "numEntries", len(out))
	reverse(out)
	return out, nil
}

func reverse(arr []types.BeaconEntry) {
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[len(arr)-(1+i)] = arr[len(arr)-(1+i)], arr[i]
	}
}

func CompressDataToFile(fileName string, in []byte) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	result := gozstd.Compress(nil, in)
	_, err = io.Copy(f, bytes.NewBuffer(result))
	return err
}

func DecompressFileToData(fileName string) ([]byte, error) {
	_, err := os.Stat(fileName)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return gozstd.Decompress(nil, data)
}
