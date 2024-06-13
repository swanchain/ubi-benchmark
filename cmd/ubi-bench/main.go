package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-paramfetch"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/types"
	lcli "github.com/filecoin-project/lotus/cli"
	"github.com/mitchellh/go-homedir"
	"github.com/swanchain/ubi-benchmark/utils"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/blake2b"

	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	prooftypes "github.com/filecoin-project/go-state-types/proof"
	lapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/storage/sealer/ffiwrapper"
	"github.com/filecoin-project/lotus/storage/sealer/ffiwrapper/basicfs"
	"github.com/filecoin-project/lotus/storage/sealer/storiface"
)

var log = logging.Logger("ubi-bench")
var latestHeight int64

type BenchResults struct {
	EnvVar map[string]string

	SectorSize   abi.SectorSize
	SectorNumber int

	SealingSum     SealingResult
	SealingResults []SealingResult
}

func (bo *BenchResults) SumSealingTime() error {
	if len(bo.SealingResults) <= 0 {
		return xerrors.Errorf("BenchResults SealingResults len <= 0")
	}
	if len(bo.SealingResults) != bo.SectorNumber {
		return xerrors.Errorf("BenchResults SealingResults len(%d) != bo.SectorNumber(%d)", len(bo.SealingResults), bo.SectorNumber)
	}

	for _, sealing := range bo.SealingResults {
		bo.SealingSum.AddPiece += sealing.AddPiece
		bo.SealingSum.PreCommit1 += sealing.PreCommit1
		bo.SealingSum.PreCommit2 += sealing.PreCommit2
	}
	return nil
}

type SealingResult struct {
	AddPiece   time.Duration
	PreCommit1 time.Duration
	PreCommit2 time.Duration
}

type Commit2In struct {
	SectorNum  int64
	Phase1Out  []byte `json:"Phase1Out,omitempty"`
	SectorSize uint64
	Commit1In
}

type Commit1In struct {
	Sid        storiface.SectorRef
	Ticket     abi.SealRandomness
	Piece      []abi.PieceInfo `json:"piece,omitempty"`
	Cids       storiface.SectorCids
	Seed       lapi.SealSeed  `json:"seed,omitempty"`
	SectorSize abi.SectorSize `json:"sector_size,omitempty"`
}

func main() {
	logging.SetLogLevel("*", "INFO")

	log.Info("Starting ubi-bench")

	app := &cli.App{
		Name:                      "ubi-bench",
		Usage:                     "Benchmark performance of ubi on your hardware",
		Version:                   "v0.0.1",
		DisableSliceFlagSeparator: true,
		Commands: []*cli.Command{
			sealCmd,
			seedCmd,
			c1Cmd,
			c2Cmd,
			verifyCmd,
			batchC1Cmd,
			uploadC1Cmd,
			daemonCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Warnf("%+v", err)
		return
	}
}

var sealCmd = &cli.Command{
	Name:   "sealing",
	Usage:  "Benchmark seal",
	Hidden: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "storage-dir",
			Value: "~/.ubi-bench",
			Usage: "path to the storage directory that will store sectors long term",
		},
		&cli.StringFlag{
			Name:  "sector-size",
			Value: "512MiB",
			Usage: "size of the sectors in bytes, i.e. 32GiB",
		},
		&cli.BoolFlag{
			Name:  "no-gpu",
			Usage: "disable gpu usage for the benchmark run",
		},
		&cli.StringFlag{
			Name:  "miner-addr",
			Usage: "pass miner address (only necessary if using existing sectorbuilder)",
			Value: "t01000",
		},
		&cli.StringFlag{
			Name:  "ticket-preimage",
			Usage: "ticket random",
		},

		&cli.IntFlag{
			Name:  "num-sectors",
			Usage: "select number of sectors to seal",
			Value: 1,
		},
		&cli.IntFlag{
			Name:  "parallel",
			Usage: "num run in parallel",
			Value: 1,
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("no-gpu") {
			err := os.Setenv("BELLMAN_NO_GPU", "1")
			if err != nil {
				return xerrors.Errorf("setting no-gpu flag: %w", err)
			}
		}

		var sbdir string
		sdir, err := homedir.Expand(c.String("storage-dir"))
		if err != nil {
			return err
		}

		err = os.MkdirAll(sdir, 0775) //nolint:gosec
		if err != nil {
			return xerrors.Errorf("creating sectorbuilder dir: %w", err)
		}

		tsdir, err := os.MkdirTemp(sdir, "bench")
		if err != nil {
			return err
		}

		if err := os.MkdirAll(tsdir, 0775); err != nil {
			return err
		}
		sbdir = tsdir

		// miner address
		maddr, err := address.NewFromString(c.String("miner-addr"))
		if err != nil {
			return err
		}
		amid, err := address.IDFromAddress(maddr)
		if err != nil {
			return err
		}
		mid := abi.ActorID(amid)

		// sector size
		sectorSizeInt, err := units.RAMInBytes(c.String("sector-size"))
		if err != nil {
			return err
		}
		sectorSize := abi.SectorSize(sectorSizeInt)

		sbfs := &basicfs.Provider{
			Root: sbdir,
		}

		sb, err := ffiwrapper.New(sbfs)
		if err != nil {
			return err
		}

		sectorNumber := c.Int("num-sectors")

		var sealTimings []SealingResult
		var extendedSealedSectors []prooftypes.ExtendedSectorInfo
		var sealedSectors []prooftypes.SectorInfo

		parCfg := ParCfg{
			PreCommit1: c.Int("parallel"),
			PreCommit2: 1,
		}
		sealTimings, extendedSealedSectors, err = runSeals(sb, sectorNumber, parCfg, mid, sectorSize, []byte(c.String("ticket-preimage")), sbdir)
		if err != nil {
			return xerrors.Errorf("failed to run seals: %w", err)
		}
		for _, s := range extendedSealedSectors {
			sealedSectors = append(sealedSectors, prooftypes.SectorInfo{
				SealedCID:    s.SealedCID,
				SectorNumber: s.SectorNumber,
				SealProof:    s.SealProof,
			})
		}

		bo := BenchResults{
			SectorSize:     sectorSize,
			SectorNumber:   sectorNumber,
			SealingResults: sealTimings,
		}
		if err := bo.SumSealingTime(); err != nil {
			return err
		}

		bo.EnvVar = make(map[string]string)
		for _, envKey := range []string{"BELLMAN_NO_GPU", "FIL_PROOFS_USE_GPU_COLUMN_BUILDER",
			"FIL_PROOFS_USE_GPU_TREE_BUILDER", "FIL_PROOFS_USE_MULTICORE_SDR", "BELLMAN_CUSTOM_GPU"} {
			envValue, found := os.LookupEnv(envKey)
			if found {
				bo.EnvVar[envKey] = envValue
			}
		}

		if c.Bool("json-out") {
			data, err := json.MarshalIndent(bo, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(data))
		} else {
			fmt.Println("environment variable list:")
			for envKey, envValue := range bo.EnvVar {
				fmt.Printf("%s=%s\n", envKey, envValue)
			}
			fmt.Printf("----\nresults (v28) SectorSize:(%d), SectorNumber:(%d)\n", sectorSize, sectorNumber)
			fmt.Printf("seal: addPiece: %s (%s)\n", bo.SealingSum.AddPiece, bps(bo.SectorSize, bo.SectorNumber, bo.SealingSum.AddPiece))
			fmt.Printf("seal: preCommit phase 1: %s (%s)\n", bo.SealingSum.PreCommit1, bps(bo.SectorSize, bo.SectorNumber, bo.SealingSum.PreCommit1))
			fmt.Printf("seal: preCommit phase 2: %s (%s)\n", bo.SealingSum.PreCommit2, bps(bo.SectorSize, bo.SectorNumber, bo.SealingSum.PreCommit2))
			fmt.Println("")
		}
		return nil
	},
}

type ParCfg struct {
	PreCommit1 int
	PreCommit2 int
}

func runSeals(sb *ffiwrapper.Sealer, numSectors int, par ParCfg, mid abi.ActorID, sectorSize abi.SectorSize, ticketPreimage []byte, sbdir string) ([]SealingResult, []prooftypes.ExtendedSectorInfo, error) {
	var pieces []abi.PieceInfo
	sealTimings := make([]SealingResult, numSectors)
	sealedSectors := make([]prooftypes.ExtendedSectorInfo, numSectors)

	preCommit2Sema := make(chan struct{}, par.PreCommit2)

	if numSectors%par.PreCommit1 != 0 {
		return nil, nil, fmt.Errorf("parallelism factor must cleanly divide numSectors")
	}
	for i := abi.SectorNumber(0); i < abi.SectorNumber(numSectors); i++ {
		sid := storiface.SectorRef{
			ID: abi.SectorID{
				Miner:  mid,
				Number: i,
			},
			ProofType: spt(sectorSize, false),
		}

		start := time.Now()
		log.Infof("[%d] Writing piece into sector...", i)

		r := rand.New(rand.NewSource(100 + int64(i)))

		pi, err := sb.AddPiece(context.TODO(), sid, nil, abi.PaddedPieceSize(sectorSize).Unpadded(), r)
		if err != nil {
			return nil, nil, err
		}

		pieces = append(pieces, pi)

		sealTimings[i].AddPiece = time.Since(start)
	}

	sectorsPerWorker := numSectors / par.PreCommit1

	errs := make(chan error, par.PreCommit1)
	for wid := 0; wid < par.PreCommit1; wid++ {
		go func(worker int) {
			sealerr := func() error {
				start := worker * sectorsPerWorker
				end := start + sectorsPerWorker
				for i := abi.SectorNumber(start); i < abi.SectorNumber(end); i++ {
					sid := storiface.SectorRef{
						ID: abi.SectorID{
							Miner:  mid,
							Number: i,
						},
						ProofType: spt(sectorSize, false),
					}

					start := time.Now()

					trand := blake2b.Sum256(ticketPreimage)
					ticket := abi.SealRandomness(trand[:])

					log.Infof("[%d] Running replication(1)...", i)
					piece := []abi.PieceInfo{pieces[i]}
					pc1o, err := sb.SealPreCommit1(context.TODO(), sid, ticket, piece)
					if err != nil {
						return xerrors.Errorf("commit: %w", err)
					}

					precommit1 := time.Now()

					preCommit2Sema <- struct{}{}
					pc2Start := time.Now()
					log.Infof("[%d] Running replication(2)...", i)
					cids, err := sb.SealPreCommit2(context.TODO(), sid, pc1o)
					if err != nil {
						return xerrors.Errorf("commit: %w", err)
					}

					precommit2 := time.Now()
					<-preCommit2Sema

					sealedSectors[i] = prooftypes.ExtendedSectorInfo{
						SealProof:    sid.ProofType,
						SectorNumber: i,
						SealedCID:    cids.Sealed,
						SectorKey:    nil,
					}

					log.Infof("[%d] Generating Commit1 for sector:", i)

					var c1in = new(Commit1In)
					c1in.Sid = sid
					c1in.Ticket = ticket
					c1in.Piece = piece
					c1in.Cids = cids
					c1in.SectorSize = sectorSize
					bytes, err := json.Marshal(c1in)
					if err != nil {
						return err
					}

					fileName := filepath.Join(filepath.Dir(sbdir), fmt.Sprintf("c1in-%d-%s.json", mid, i.String()))
					if err = os.WriteFile(fileName, bytes, 0644); err != nil {
						return err
					}

					sealTimings[i].PreCommit1 = precommit1.Sub(start)
					sealTimings[i].PreCommit2 = precommit2.Sub(pc2Start)
				}
				return nil
			}()
			if sealerr != nil {
				errs <- sealerr
				return
			}
			errs <- nil
		}(wid)
	}

	for i := 0; i < par.PreCommit1; i++ {
		err := <-errs
		if err != nil {
			return nil, nil, err
		}
	}

	return sealTimings, sealedSectors, nil
}

var seedCmd = &cli.Command{
	Name:   "seed",
	Usage:  "Generate random numbers",
	Hidden: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "miner-addr",
			Usage: "miner address",
			Value: "t01000",
		},
		&cli.Int64Flag{
			Name:  "height",
			Usage: "specify a height",
		},
	},
	Action: func(c *cli.Context) error {
		height := c.Int64("height")
		if height == 0 {
			return fmt.Errorf("must be specify a height")
		}

		maddr, err := address.NewFromString(c.String("miner-addr"))
		if err != nil {
			return err
		}

		randomness, err := utils.GetRandomness(maddr, crypto.DomainSeparationTag_InteractiveSealChallengeSeed, height)
		if err != nil {
			return err
		}
		fmt.Printf("randomness: %v \n", randomness)
		return nil
	},
}

var c1Cmd = &cli.Command{
	Name:      "c1",
	Usage:     "execute Commit1 task",
	ArgsUsage: "[input.json]",
	Hidden:    true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "storage-dir",
			Usage: "path to the storage directory that will store sectors long term",
		},
		&cli.Int64Flag{
			Name:  "height",
			Usage: "specify a height",
		},
	},
	Action: func(c *cli.Context) error {
		if !c.Args().Present() {
			return xerrors.Errorf("Usage: ubi-bench c1 [input.json]")
		}

		height := c.Int64("height")
		if height == 0 {
			return fmt.Errorf("must be specify a height")
		}

		sdir := c.String("storage-dir")
		if _, err := os.Stat(sdir); err != nil && os.IsNotExist(err) {
			return err
		}
		sbfs := &basicfs.Provider{
			Root: sdir,
		}
		sb, err := ffiwrapper.New(sbfs)
		if err != nil {
			return err
		}

		inb, err := os.ReadFile(c.Args().First())
		if err != nil {
			return xerrors.Errorf("reading input file: %w", err)
		}
		var c1in Commit1In
		if err := json.Unmarshal(inb, &c1in); err != nil {
			return xerrors.Errorf("unmarshalling input file: %w", err)
		}

		maddr, err := address.NewFromString("t0" + c1in.Sid.ID.Miner.String())
		if err != nil {
			return err
		}

		randomness, err := utils.GetRandomness(maddr, crypto.DomainSeparationTag_InteractiveSealChallengeSeed, height)
		if err != nil {
			return err
		}

		seed := lapi.SealSeed{
			Epoch: abi.ChainEpoch(height),
			Value: randomness,
		}

		c1o, err := sb.SealCommit1(context.TODO(), c1in.Sid, c1in.Ticket, seed.Value, c1in.Piece, c1in.Cids)
		if err != nil {
			return err
		}

		var c2in = new(Commit2In)
		c2in.SectorNum = int64(c1in.Sid.ID.Number)
		c2in.Phase1Out = c1o
		c2in.SectorSize = uint64(c1in.SectorSize)
		c2in.Cids = c1in.Cids
		c2in.Sid = c1in.Sid
		c2in.Ticket = c1in.Ticket
		c2in.Seed = seed

		c2inBytes, err := json.Marshal(c2in)
		if err != nil {
			return err
		}

		c1JsonFile := filepath.Join(filepath.Dir(sdir), fmt.Sprintf("c1-%d-%d-%d.json", c1in.Sid.ID.Miner, c1in.Sid.ID.Number, seed.Epoch))
		if err = os.WriteFile(c1JsonFile, c2inBytes, 0666); err != nil {
			return err
		}

		fmt.Printf("seal: commit phase 1 finished, sector_id: %d \n", c1in.Sid.ID.Number)
		return nil
	},
}

var c2Cmd = &cli.Command{
	Name:      "c2",
	Usage:     "execute c2 task for a proof computation",
	ArgsUsage: "[input.json]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "no-gpu",
			Usage: "disable gpu usage for the benchmark run",
		},
		&cli.StringFlag{
			Name:  "storage-dir",
			Usage: "path to the storage directory that will store sectors long term",
			Value: "/var/tmp",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("no-gpu") {
			err := os.Setenv("BELLMAN_NO_GPU", "1")
			if err != nil {
				return xerrors.Errorf("setting no-gpu flag: %w", err)
			}
		}

		paramsFile := c.Args().First()
		if strings.TrimSpace(paramsFile) == "" {
			return xerrors.Errorf("input param is empty")
		}

		sdir := c.String("storage-dir")
		if _, err := os.Stat(sdir); err != nil && os.IsNotExist(err) {
			return err
		}

		inb, err := os.ReadFile(paramsFile)
		if err != nil {
			return xerrors.Errorf("reading input file: %w", err)
		}

		var c2in Commit2In
		if err := json.Unmarshal(inb, &c2in); err != nil {
			return xerrors.Errorf("unmarshalling input file: %w", err)
		}

		if err := paramfetch.GetParams(lcli.ReqContext(c), build.ParametersJSON(), build.SrsJSON(), c2in.SectorSize); err != nil {
			return xerrors.Errorf("getting params: %w", err)
		}

		sb, err := ffiwrapper.New(nil)
		if err != nil {
			return err
		}

		start := time.Now()
		proof, err := sb.SealCommit2(context.TODO(), c2in.Sid, c2in.Phase1Out)
		log.Info("proof: %v", proof)
		if err != nil {
			return err
		}
		totalTime := time.Since(start)
		svi := prooftypes.SealVerifyInfo{
			SectorID:              c2in.Sid.ID,
			SealedCID:             c2in.Cids.Sealed,
			SealProof:             c2in.Sid.ProofType,
			Proof:                 proof,
			DealIDs:               nil,
			Randomness:            c2in.Ticket,
			InteractiveRandomness: c2in.Seed.Value,
			UnsealedCID:           c2in.Cids.Unsealed,
		}
		c2OutBytes, err := json.Marshal(svi)
		if err != nil {
			return err
		}
		log.Info("c2OutBytes: %v", string(c2OutBytes))
		c2JsonFile := filepath.Join(filepath.Dir(sdir), fmt.Sprintf("c2-%d-%d-%d.json", c2in.Sid.ID.Miner, c2in.Sid.ID.Number, c2in.Seed.Epoch))
		if err = os.WriteFile(c2JsonFile, c2OutBytes, 0666); err != nil {
			return err
		}

		log.Info("seal: commit phase 2 finished, total time: %f, sector_id: %d \n", totalTime.Seconds(), c2in.SectorNum)
		return nil
	},
}

var batchC1Cmd = &cli.Command{
	Name:      "batch",
	Usage:     "execute batch Commit1 task",
	ArgsUsage: "[input.json]",
	Hidden:    true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "storage-dir",
			Usage: "path to the storage directory that will store sectors long term",
		},
		&cli.Int64Flag{
			Name:  "start",
			Usage: "specify a height",
		},
		&cli.IntFlag{
			Name:  "num",
			Usage: "number of batches generated",
			Value: 1,
		},
	},
	Action: func(c *cli.Context) error {
		if !c.Args().Present() {
			return xerrors.Errorf("Usage: ubi-bench batch [input.json]")
		}
		num := c.Int("num")

		height := c.Int64("start")
		if height == 0 {
			return fmt.Errorf("must be specify a height")
		}

		sdir := c.String("storage-dir")
		if _, err := os.Stat(sdir); err != nil && os.IsNotExist(err) {
			return err
		}
		sbfs := &basicfs.Provider{
			Root: sdir,
		}
		sb, err := ffiwrapper.New(sbfs)
		if err != nil {
			return err
		}

		inb, err := os.ReadFile(c.Args().First())
		if err != nil {
			return xerrors.Errorf("reading input file: %w", err)
		}
		var c1in Commit1In
		if err := json.Unmarshal(inb, &c1in); err != nil {
			return xerrors.Errorf("unmarshalling input file: %w", err)
		}

		maddr, err := address.NewFromString("t0" + c1in.Sid.ID.Miner.String())
		if err != nil {
			return err
		}

		for i := 0; i < num; i++ {
			_, _, err := generaC1Out(maddr, sb, sdir, c1in, height+int64(i))
			if err != nil {
				return err
			}
		}
		return nil
	},
}

var uploadC1Cmd = &cli.Command{
	Name:   "upload",
	Usage:  "Batch upload the results of c1 to mcs",
	Hidden: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "c1-dir",
			Usage: "path to the c1 out directory",
		},
		&cli.StringFlag{
			Name:  "type",
			Usage: "example: 512 or 32",
			Value: "512",
		},
	},
	Action: func(c *cli.Context) error {
		c1Dir := c.String("c1-dir")
		if _, err := os.Stat(c1Dir); err != nil && os.IsNotExist(err) {
			return err
		}

		if err := utils.InitConfig(); err != nil {
			return err
		}

		numType := c.String("type")
		var dirName string
		if numType == "512" {
			dirName = "fil-c2/512M"
		} else if numType == "32" {
			dirName = "fil-c2/32G"
		}

		var count int
		storageService := utils.NewStorageService()
		err := filepath.WalkDir(c1Dir, func(path string, d fs.DirEntry, err error) error {
			split := strings.Split(d.Name(), "-")
			if d.IsDir() && len(split) == 4 {
				storageService.CreateFolder(dirName, d.Name())
				files, err := os.ReadDir(path)
				if err != nil {
					return err
				}
				var inputParam, verifyParam string
				for _, f := range files {
					mcsOssFile, err := storageService.UploadFileToBucket(filepath.Join(dirName, d.Name(), f.Name()), filepath.Join(path, f.Name()), true)
					if err != nil {
						log.Errorf("Failed upload file to bucket, error: %v", err)
						return err
					}
					gatewayUrl, err := storageService.GetGatewayUrl()
					if err != nil {
						log.Errorf("Failed get mcs ipfs gatewayUrl, error: %v", err)
						return err
					}

					fileUrl := *gatewayUrl + "/ipfs/" + mcsOssFile.PayloadCid
					fmt.Printf("file name: %s, url: %s \n", f.Name(), fileUrl)
					if strings.Contains(f.Name(), "verify") {
						verifyParam = fileUrl
					} else {
						inputParam = fileUrl
					}
				}

				var task Task
				if count%2 == 0 {
					var resourceId int
					if numType == "512" {
						resourceId = CPU512
					} else if numType == "32" {
						resourceId = CPU32G
					}

					task = Task{
						Name:        d.Name(),
						Type:        0,
						InputParam:  inputParam,
						VerifyParam: verifyParam,
						ResourceID:  resourceId,
					}
				} else {
					var resourceId int
					if numType == "512" {
						resourceId = GPU512
					} else if numType == "32" {
						resourceId = GPU32G
					}

					task = Task{
						Name:        d.Name(),
						Type:        1,
						InputParam:  inputParam,
						VerifyParam: verifyParam,
						ResourceID:  resourceId,
					}
				}
				DoSend(task)
				count++
				fmt.Println("==============")
			}
			if err != nil {
				return err
			}

			return nil
		})

		return err
	},
}

var daemonCmd = &cli.Command{
	Name:      "daemon",
	Usage:     "Auto generate c1 out and upload the results of c1 to mcs",
	ArgsUsage: "[c1in-input.json]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "storage-dir",
			Usage: "path to the storage directory that will store sectors long term",
		},
		&cli.Int64Flag{
			Name:  "last-height",
			Usage: "specify a height",
		},
		&cli.Int64Flag{
			Name:  "sector-type",
			Usage: "512 or 32",
		},
	},
	Action: func(c *cli.Context) error {
		if !c.Args().Present() {
			return xerrors.Errorf("Usage: ubi-bench batch [input.json]")
		}

		height := c.Int64("last-height")
		if height == 0 {
			return fmt.Errorf("must be specify a last-height")
		}

		sectorType := c.Int64("sector-type")
		if sectorType == 0 {
			return fmt.Errorf("must be specify a sector-type")
		}

		if sectorType != 512 && sectorType != 32 {
			return fmt.Errorf("sector-type value is wrong")
		}

		sdir := c.String("storage-dir")
		if _, err := os.Stat(sdir); err != nil && os.IsNotExist(err) {
			return err
		}
		sbfs := &basicfs.Provider{
			Root: sdir,
		}
		sb, err := ffiwrapper.New(sbfs)
		if err != nil {
			return err
		}

		inb, err := os.ReadFile(c.Args().First())
		if err != nil {
			return xerrors.Errorf("reading input file: %w", err)
		}
		var c1in Commit1In
		if err := json.Unmarshal(inb, &c1in); err != nil {
			return xerrors.Errorf("unmarshalling input file: %w", err)
		}

		maddr, err := address.NewFromString("t0" + c1in.Sid.ID.Miner.String())
		if err != nil {
			return err
		}

		if err := utils.InitConfig(); err != nil {
			return err
		}

		ticker := time.NewTicker(time.Duration(utils.GetConfig().HUB.CheckInterval) * time.Minute)
		defer ticker.Stop()

		latestHeight = height
		for {
			select {
			case <-ticker.C:
				checkTaskCount(maddr, sb, sdir, c1in, latestHeight, sectorType)
			}
		}
	},
}

var verifyCmd = &cli.Command{
	Name:      "verify",
	Usage:     "Verify a proof computation",
	ArgsUsage: "[input.json]",
	Flags: []cli.Flag{
		&cli.Int64Flag{
			Name:  "height",
			Usage: "specify a height",
		},
	},
	Action: func(c *cli.Context) error {
		if !c.Args().Present() && !c.IsSet("s") {
			return xerrors.Errorf("Usage: ubi verify [input.json]")
		}

		height := c.Int64("height")
		if height == 0 {
			return fmt.Errorf("must be specify a height")
		}
		var inb []byte
		if c.IsSet("s") {
			inb = []byte(c.String("s"))
		} else {
			var err error
			inb, err = os.ReadFile(c.Args().First())
			if err != nil {
				return xerrors.Errorf("reading input file: %w", err)
			}
		}

		var svi prooftypes.SealVerifyInfo
		if err := json.Unmarshal(inb, &svi); err != nil {
			return xerrors.Errorf("unmarshalling input file: %w", err)
		}

		maddr, err := address.NewFromString("t0" + svi.Miner.String())
		if err != nil {
			return err
		}

		randomness, err := utils.GetRandomness(maddr, crypto.DomainSeparationTag_InteractiveSealChallengeSeed, height)
		if err != nil {
			return err
		}
		svi.InteractiveRandomness = randomness

		ok, err := ffiwrapper.ProofVerifier.VerifySeal(svi)
		if err != nil {
			return err
		}
		if !ok {
			return xerrors.Errorf("proof for sector %d was invalid", svi.SectorID.Number)
		}

		fmt.Printf("seal: proof for sector %d was valid. \n", svi.SectorID.Number)
		return nil
	},
}

func checkTaskCount(mAddr address.Address, sealer *ffiwrapper.Sealer, sdir string, c1in Commit1In, height, sectorType int64) {
	response, err := http.Get(utils.GetConfig().HUB.TaskUrl)
	if err != nil {
		log.Errorf("Error making GET request: %v", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error reading response body: %v", err)
		return
	}

	var taskStats TaskStats
	err = json.Unmarshal(body, &taskStats)
	if err != nil {
		log.Errorf("Error response convet to json: %v", err)
		return
	}

	log.Infof("current task stats: %+v", taskStats.Data)
	var rc512, rc32 ResourceCountList
	for _, t := range taskStats.Data {
		if sectorType == 512 && (t.ResourceId == CPU512 || t.ResourceId == GPU512) {
			rc512 = append(rc512, t)
		}
		if sectorType == 32 && (t.ResourceId == CPU32G || t.ResourceId == GPU32G) {
			rc32 = append(rc32, t)
		}
	}

	var dirName, zkType string
	var taskType int     //  1:Fil-C2-512M, 2:Aleo, 3:AI, 4:Fil-C2-32G
	var resourceType int // 0: cpu 1: gpu
	var needTask ResourceCount
	if sectorType == 512 {
		taskType = 1
		sort.Sort(rc512)
		if len(rc512) > 0 {
			needTask = rc512[0]
			if needTask.Count < 20000 {
				dirName = "fil-c2/512M"
				zkType = "fil-c2-512M"
				if needTask.ResourceId == CPU512 {
					resourceType = 0
				} else {
					resourceType = 1
				}
			} else {
				return
			}
		} else {
			return
		}
	} else {
		taskType = 4
		sort.Sort(rc32)
		if len(rc32) > 0 {
			needTask = rc32[0]
			if needTask.Count < 20000 {
				dirName = "fil-c2/32G"
				zkType = "fil-c2-32G"
				if needTask.ResourceId == CPU32G {
					resourceType = 0
				} else {
					resourceType = 1
				}
			} else {
				return
			}
		} else {
			return
		}
	}

	for i := 0; i < utils.GetConfig().HUB.BatchNum; i++ {
		rootDir, taskDir, err := generaC1Out(mAddr, sealer, sdir, c1in, height+int64(i))
		if err != nil {
			log.Errorf("Error response convet to json: %v", err)
			return
		}

		storageService := utils.NewStorageService()
		storageService.CreateFolder(dirName, taskDir)

		var inputParam, verifyParam string
		err = filepath.Walk(rootDir, func(path string, f fs.FileInfo, err error) error {
			if f.IsDir() {
				return nil
			}
			mcsOssFile, err := storageService.UploadFileToBucket(filepath.Join(dirName, taskDir, f.Name()), path, true)
			if err != nil {
				log.Errorf("Failed upload file to bucket, error: %v", err)
				return err
			}
			gatewayUrl, err := storageService.GetGatewayUrl()
			if err != nil {
				log.Errorf("Failed get mcs ipfs gatewayUrl, error: %v", err)
				return err
			}

			fileUrl := *gatewayUrl + "/ipfs/" + mcsOssFile.PayloadCid
			fmt.Printf("file name: %s, url: %s \n", f.Name(), fileUrl)
			if strings.Contains(f.Name(), "verify") {
				verifyParam = fileUrl
			} else {
				inputParam = fileUrl
			}
			if err != nil {
				return err
			}
			return nil
		})

		if strings.HasSuffix(inputParam, "ipfs/") || strings.HasSuffix(verifyParam, "ipfs/") {
			continue
		}

		var task = Task{
			Name:         taskDir,
			Type:         taskType,
			InputParam:   inputParam,
			VerifyParam:  verifyParam,
			ResourceID:   needTask.ResourceId,
			ResourceType: resourceType,
		}
		DoSend(task)
		fmt.Println("==============")

		latestHeight++
		if err != nil {
			log.Errorf("upload file to mcs failed, task: %s, zkType: %s", taskDir, zkType)
			continue
		}
		os.RemoveAll(rootDir)
	}

}

func generaC1Out(mAddr address.Address, sealer *ffiwrapper.Sealer, sdir string, c1in Commit1In, height int64) (string, string, error) {
	randomness, err := utils.GetRandomness(mAddr, crypto.DomainSeparationTag_InteractiveSealChallengeSeed, height)
	if err != nil {
		return "", "", err
	}

	seed := lapi.SealSeed{
		Epoch: abi.ChainEpoch(height),
		Value: randomness,
	}

	c1o, err := sealer.SealCommit1(context.TODO(), c1in.Sid, c1in.Ticket, seed.Value, c1in.Piece, c1in.Cids)
	if err != nil {
		return "", "", err
	}

	var c2in = new(Commit2In)
	c2in.SectorNum = int64(c1in.Sid.ID.Number)
	c2in.SectorSize = uint64(c1in.SectorSize)
	c2in.Cids = c1in.Cids
	c2in.Sid = c1in.Sid
	c2in.Ticket = c1in.Ticket
	c2in.Seed = seed
	c2inBytes, err := json.Marshal(c2in)
	if err != nil {
		return "", "", err
	}

	taskDir := fmt.Sprintf("%d-%d-%d-%d", c2in.Sid.ID.Miner, c2in.Sid.ID.Number, c2in.Sid.ProofType, c2in.Seed.Epoch)
	rootDir := filepath.Join(filepath.Dir(sdir), taskDir)

	err = os.MkdirAll(rootDir, 0775) //nolint:gosec
	if err != nil {
		return "", "", xerrors.Errorf("creating task dir: %w", err)
	}
	log.Infof("create dir: %s", rootDir)

	c2JsonFile := filepath.Join(rootDir, fmt.Sprintf("c1out-%d-%d-%d-verify.json", c2in.Sid.ID.Miner, c2in.Sid.ID.Number, c2in.Seed.Epoch))
	if err = os.WriteFile(c2JsonFile, c2inBytes, 0666); err != nil {
		return "", "", err
	}

	c2in.Phase1Out = c1o
	c2inBytesWithC1, err := json.Marshal(c2in)
	if err != nil {
		return "", "", err
	}
	c1JsonFile := filepath.Join(rootDir, fmt.Sprintf("c1out-%d-%d-%d.zst", c1in.Sid.ID.Miner, c1in.Sid.ID.Number, seed.Epoch))
	if err = utils.CompressDataToFile(c1JsonFile, c2inBytesWithC1); err != nil {
		return "", "", err
	}

	log.Infof("seal: commit phase1 finished, sector_id: %d, num: %d\n", c1in.Sid.ID.Number, height)
	return rootDir, taskDir, nil
}

func bps(sectorSize abi.SectorSize, sectorNum int, d time.Duration) string {
	bdata := new(big.Int).SetUint64(uint64(sectorSize))
	bdata = bdata.Mul(bdata, big.NewInt(int64(sectorNum)))
	bdata = bdata.Mul(bdata, big.NewInt(time.Second.Nanoseconds()))
	bps := bdata.Div(bdata, big.NewInt(d.Nanoseconds()))
	return types.SizeStr(types.BigInt{Int: bps}) + "/s"
}

func spt(ssize abi.SectorSize, synth bool) abi.RegisteredSealProof {
	spt, err := miner.SealProofTypeFromSectorSize(ssize, build.TestNetworkVersion, synth)
	if err != nil {
		panic(err)
	}

	return spt
}
