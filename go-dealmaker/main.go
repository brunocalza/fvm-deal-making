package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ipfs/go-cid"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "dealmaker",
		Usage: "dealmaker lets you make deals using a FVM smart contract",
		Commands: []*cli.Command{
			createCmd,
			statusCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var createCmd = &cli.Command{
	Name:  "create",
	Usage: "Create deal",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "rpc-endpoint",
			Usage: "Gateway RPC endpoint",
		},
		&cli.StringFlag{
			Name:  "contract",
			Usage: "The Smart Contract address",
		},
		&cli.StringFlag{
			Name:  "piece-cid",
			Usage: "The piece CID",
		},
		&cli.Int64Flag{
			Name:  "piece-size",
			Usage: "The piece size in bytes",
		},
		&cli.BoolFlag{
			Name:  "verified",
			Usage: "If it's a verified deal or not",
		},
		&cli.StringFlag{
			Name:  "payload-cid",
			Usage: "The payload CID",
		},
		&cli.Int64Flag{
			Name:  "start-epoch",
			Usage: "When the deal starts",
		},
		&cli.Int64Flag{
			Name:  "end-epoch",
			Usage: "When the deal ends",
		},
		&cli.StringFlag{
			Name:  "location-ref",
			Usage: "Where the CAR file can be downloaded",
		},
		&cli.Int64Flag{
			Name:  "car-size",
			Usage: "The size of the CAR file",
		},
		&cli.StringFlag{
			Name:  "private-key",
			Usage: "The private key",
		},
		&cli.Int64Flag{
			Name:  "chain-id",
			Usage: "The network id",
		},
	},
	Action: func(cCtx *cli.Context) error {
		flags, err := ValidateCreateFlags(cCtx)
		if err != nil {
			return err
		}

		conn, err := ethclient.Dial(flags.RPCEndpoint.String())
		if err != nil {
			log.Fatal(err)
		}

		contract, err := NewContract(flags.Contract, conn)
		if err != nil {
			log.Fatal(err)
		}

		dealRequest := DealRequest{
			PieceCid:             flags.PieceCID.Bytes(),
			PieceSize:            uint64(flags.PieceSize),
			VerifiedDeal:         flags.Verified,
			Label:                flags.PayloadCID.String(),
			StartEpoch:           flags.StartEpoch,
			EndEpoch:             flags.EndEpoch,
			StoragePricePerEpoch: big.NewInt(0),
			ProviderCollateral:   big.NewInt(0),
			ClientCollateral:     big.NewInt(0),
			ExtraParamsVersion:   1,
			ExtraParams: ExtraParamsV1{
				LocationRef:        flags.LocationRef.String(),
				CarSize:            uint64(flags.CarSize),
				SkipIpniAnnounce:   false,
				RemoveUnsealedCopy: false,
			},
		}

		auth, err := bind.NewKeyedTransactorWithChainID(flags.PrivateKey, big.NewInt(flags.ChainID))
		if err != nil {
			return err
		}

		tx, err := contract.MakeDealProposal(auth, dealRequest)
		if err != nil {
			return err
		}

		fmt.Println(tx.Hash())
		return nil
	},
}

var statusCmd = &cli.Command{
	Name:  "status",
	Usage: "Check the status of a deal",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "rpc-endpoint",
			Usage: "Gateway RPC endpoint",
		},
		&cli.StringFlag{
			Name:  "contract",
			Usage: "The Smart Contract address",
		},
		&cli.StringFlag{
			Name:  "piece-cid",
			Usage: "The piece CID",
		},
	},
	Action: func(cCtx *cli.Context) error {
		flags, err := ValidateStatusFlags(cCtx)
		if err != nil {
			return err
		}
		conn, err := ethclient.Dial(flags.RPCEndpoint.String())
		if err != nil {
			log.Fatal(err)
		}

		contract, err := NewContract(flags.Contract, conn)
		if err != nil {
			log.Fatal(err)
		}

		status, err := contract.PieceStatus(nil, flags.PieceCID.Bytes())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(status)

		return nil
	},
}

type CreateFlags struct {
	RPCEndpoint *url.URL
	Contract    common.Address
	PieceCID    cid.Cid
	PieceSize   int64
	Verified    bool
	PayloadCID  cid.Cid
	StartEpoch  int64
	EndEpoch    int64
	LocationRef *url.URL
	CarSize     int64
	PrivateKey  *ecdsa.PrivateKey
	ChainID     int64
}

func ValidateCreateFlags(cCtx *cli.Context) (CreateFlags, error) {
	rpcEndpoint, err := url.Parse(cCtx.String("rpc-endpoint"))
	if err != nil {
		return CreateFlags{}, err
	}

	if !common.IsHexAddress(cCtx.String("contract")) {
		return CreateFlags{}, fmt.Errorf("contract is not an ETH address")
	}

	pieceCID, err := cid.Decode(cCtx.String("piece-cid"))
	if err != nil {
		return CreateFlags{}, err
	}

	payloadCID, err := cid.Decode(cCtx.String("payload-cid"))
	if err != nil {
		return CreateFlags{}, err
	}

	locationRef, err := url.Parse(cCtx.String("location-ref"))
	if err != nil {
		return CreateFlags{}, err
	}

	pk, err := crypto.HexToECDSA(cCtx.String("private-key"))
	if err != nil {
		log.Fatal(err)
	}

	return CreateFlags{
		RPCEndpoint: rpcEndpoint,
		Contract:    common.HexToAddress(cCtx.String("contract")),
		PieceCID:    pieceCID,
		PieceSize:   cCtx.Int64("piece-size"),
		Verified:    cCtx.Bool("verified"),
		PayloadCID:  payloadCID,
		StartEpoch:  cCtx.Int64("start-epoch"),
		EndEpoch:    cCtx.Int64("end-epoch"),
		LocationRef: locationRef,
		CarSize:     cCtx.Int64("car-size"),
		PrivateKey:  pk,
		ChainID:     cCtx.Int64("chain-id"),
	}, nil
}

type StatusFlags struct {
	RPCEndpoint *url.URL
	Contract    common.Address
	PieceCID    cid.Cid
}

func ValidateStatusFlags(cCtx *cli.Context) (CreateFlags, error) {
	rpcEndpoint, err := url.Parse(cCtx.String("rpc-endpoint"))
	if err != nil {
		return CreateFlags{}, err
	}

	if !common.IsHexAddress(cCtx.String("contract")) {
		return CreateFlags{}, fmt.Errorf("contract is not an ETH address")
	}

	pieceCID, err := cid.Decode(cCtx.String("piece-cid"))
	if err != nil {
		return CreateFlags{}, err
	}

	return CreateFlags{
		RPCEndpoint: rpcEndpoint,
		Contract:    common.HexToAddress(cCtx.String("contract")),
		PieceCID:    pieceCID,
	}, nil
}
