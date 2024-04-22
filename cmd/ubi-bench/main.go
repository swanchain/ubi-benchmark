package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/filecoin-project/go-paramfetch"
	"github.com/filecoin-project/go-state-types/abi"
	prooftypes "github.com/filecoin-project/go-state-types/proof"
	lapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/storage/sealer/ffiwrapper"
	"github.com/filecoin-project/lotus/storage/sealer/storiface"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("ubi-bench")

var tmpDir string
var accessToken string

const UbiProofDir = "zk-proof"

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
	r := gin.Default()
	r.Use(AuthMiddleware(), cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		ValidateHeaders: false,
	}))

	router := r.Group("/api")
	router.GET("/ubi/:miner_id/:sector_id", getC2Proof)
	router.POST("/ubi", doC2Req)

	token, ok := os.LookupEnv("access_token")
	if !ok {
		log.Fatalf("must be set access_token env")
	}
	accessToken = token

	var err error
	tmpDir, err = os.MkdirTemp("", UbiProofDir)
	if err != nil {
		log.Errorf("create ubi proof dir failed, error: %v", err)
		return
	}

	srv := &http.Server{
		Addr:              ":9000",
		Handler:           r,
		ReadHeaderTimeout: 60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %v\n", err)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("access_token")
		if token != accessToken {
			c.JSON(http.StatusUnauthorized, createResponse(AuthFailedCode, ""))
			c.Abort()
			return
		}
		c.Next()
	}
}

func doC2Req(c *gin.Context) {
	var c2in Commit2In
	if err := c.ShouldBindJSON(&c2in); err != nil {
		c.JSON(http.StatusInternalServerError, createResponse(JsonError, ""))
		return
	}

	if err := paramfetch.GetParams(context.TODO(), build.ParametersJSON(), build.SrsJSON(), c2in.SectorSize); err != nil {
		log.Errorf("getting params: %v", err)
		c.JSON(http.StatusInternalServerError, createResponse(ServerError, ""))
		return
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("catch painc error: %v", err)
			}
		}()

		sb, err := ffiwrapper.New(nil)
		if err != nil {
			log.Errorf("create ffiwrapper failed, error: %v", err)
			return
		}

		start := time.Now()
		proof, err := sb.SealCommit2(context.TODO(), c2in.Sid, c2in.Phase1Out)
		if err != nil {
			log.Errorf("miner_id: %s, sector_id: %d, do SealCommit2 failed, error: %v", c2in.Sid.ID.Miner.String(), c2in.SectorNum, err)
			return
		}
		totalTime := time.Since(start)

		var c2proof C2Proof
		c2proof.Proof = base64.StdEncoding.EncodeToString(proof)
		if c2in.Ticket != nil {
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

			ok, err := ffiwrapper.ProofVerifier.VerifySeal(svi)
			if err != nil {
				log.Errorf("verify miner_id: %s, sector_id: %d, c2 proof failed, error: %v", svi.Miner.String(), c2in.SectorNum, err)
				return
			}
			c2proof.Verify = ok
			log.Infof("miner_id: %s, sector_id: %d, ", svi.Miner.String(), svi.SectorID.Number)
		}

		if c2proof.Verify {
			log.Infof("seal: commit phase 2 finished, total time: %f, miner_id: %s, sector_id: %d, proof was valid", totalTime.Seconds(), c2in.Sid.ID.Miner.String(), c2in.SectorNum)
		} else {
			log.Infof("seal: commit phase 2 finished, total time: %f, miner_id: %s, sector_id: %d", totalTime.Seconds(), c2in.Sid.ID.Miner.String(), c2in.SectorNum)
		}

		result, _ := json.Marshal(&c2proof)
		c2JsonFile := filepath.Join(tmpDir, fmt.Sprintf("c2-%d-%d.json", c2in.Sid.ID.Miner, c2in.Sid.ID.Number))
		if err = os.WriteFile(c2JsonFile, result, 0666); err != nil {
			log.Errorf("save miner_id: %s, sector_id: %s, proof file failed, error: %v", c2in.Sid.ID.Miner, c2in.Sid.ID.Number, err)
			return
		}
	}()

	c.JSON(http.StatusOK, createResponse(SubmitProof, ""))
}

func getC2Proof(c *gin.Context) {
	minerId := c.Param("miner_id")
	sectorId := c.Param("sector_id")

	if strings.TrimSpace(minerId) == "" {
		c.JSON(http.StatusBadRequest, createResponse(JsonError, "the miner_id field is required"))
		return
	}
	if strings.TrimSpace(sectorId) == "" {
		c.JSON(http.StatusBadRequest, createResponse(JsonError, "the sector_id field is required"))
		return
	}

	c2File := filepath.Join(tmpDir, fmt.Sprintf("c2-%s-%s.json", minerId, sectorId))

	if _, err := os.Stat(c2File); err != nil {
		log.Errorf("get miner_id: %s, sector_id: %s proof file failed, error: %v", minerId, sectorId, err)
		c.JSON(http.StatusBadRequest, createResponse(ServerError, "not found the proof"))
		return
	}

	data, err := os.ReadFile(c2File)
	if err != nil {
		log.Errorf("read miner: %s, sector_id: %s proof file failed, error: %v", minerId, sectorId, err)
		c.JSON(http.StatusBadRequest, createResponse(ServerError, "get the proof failed"))
		return
	}

	c.JSON(http.StatusOK, createDataResponse(SuccessCode, string(data)))
}

type C2Proof struct {
	Proof  string
	Verify bool
}

type Response struct {
	Code    int         `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func createResponse(code int, msg string) Response {
	var message string
	if msg != "" {
		message = msg
	} else {
		message = codeMsg[code]
	}

	return Response{
		Code:    code,
		Message: message,
	}
}

func createDataResponse(code int, data interface{}) Response {
	return Response{
		Code: code,
		Data: data,
	}
}

const (
	SuccessCode    = 200
	AuthFailedCode = 401
	JsonError      = 400
	ServerError    = 500
	BadParamError  = 5001

	SubmitProof = 2000
	ProofError  = 7003
)

var codeMsg = map[int]string{
	BadParamError:  "The request parameter is not valid",
	JsonError:      "An error occurred while converting to json",
	ServerError:    "server failed",
	SubmitProof:    "The proof task has been submitted",
	ProofError:     "An error occurred while executing the calculation task",
	AuthFailedCode: "the request header missing access_token or access_token is incorrect",
}
