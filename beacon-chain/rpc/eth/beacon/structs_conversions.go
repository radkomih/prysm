package beacon

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	bytesutil2 "github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/wealdtech/go-bytesutil"
)

func (b *SignedBeaconBlock) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	sig, err := hexutil.Decode(b.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Signature")
	}
	slot, err := strconv.ParseUint(b.Message.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Slot")
	}
	proposerIndex, err := strconv.ParseUint(b.Message.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(b.Message.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(b.Message.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(b.Message.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(b.Message.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(b.Message.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(b.Message.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(b.Message.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(b.Message.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(b.Message.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(b.Message.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(b.Message.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(b.Message.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}

	block := &eth.SignedBeaconBlock{
		Block: &eth.BeaconBlock{
			Slot:          primitives.Slot(slot),
			ProposerIndex: primitives.ValidatorIndex(proposerIndex),
			ParentRoot:    parentRoot,
			StateRoot:     stateRoot,
			Body: &eth.BeaconBlockBody{
				RandaoReveal: randaoReveal,
				Eth1Data: &eth.Eth1Data{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				},
				Graffiti:          graffiti,
				ProposerSlashings: proposerSlashings,
				AttesterSlashings: attesterSlashings,
				Attestations:      atts,
				Deposits:          deposits,
				VoluntaryExits:    exits,
			},
		},
		Signature: sig,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_Phase0{Phase0: block}}, nil
}

func convertInternalSignedBeaconBlock(b *eth.SignedBeaconBlock) (*SignedBeaconBlock, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Block.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Block.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Block.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Block.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Block.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	return &SignedBeaconBlock{
		Message: &BeaconBlock{
			Slot:          fmt.Sprintf("%d", b.Block.Slot),
			ProposerIndex: fmt.Sprintf("%d", b.Block.ProposerIndex),
			ParentRoot:    hexutil.Encode(b.Block.ParentRoot),
			StateRoot:     hexutil.Encode(b.Block.StateRoot),
			Body: &BeaconBlockBody{
				RandaoReveal: hexutil.Encode(b.Block.Body.RandaoReveal),
				Eth1Data: &Eth1Data{
					DepositRoot:  hexutil.Encode(b.Block.Body.Eth1Data.DepositRoot),
					DepositCount: fmt.Sprintf("%d", b.Block.Body.Eth1Data.DepositCount),
					BlockHash:    hexutil.Encode(b.Block.Body.Eth1Data.BlockHash),
				},
				Graffiti:          hexutil.Encode(b.Block.Body.Graffiti),
				ProposerSlashings: proposerSlashings,
				AttesterSlashings: attesterSlashings,
				Attestations:      atts,
				Deposits:          deposits,
				VoluntaryExits:    exits,
			},
		},
		Signature: hexutil.Encode(b.Signature),
	}, nil
}

func convertInternalBeaconBlock(b *eth.BeaconBlock) (*BeaconBlock, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	return &BeaconBlock{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BeaconBlockBody{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
		},
	}, nil
}

func convertInternalBeaconBlockAltair(b *eth.BeaconBlockAltair) (*BeaconBlockAltair, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}

	return &BeaconBlockAltair{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BeaconBlockBodyAltair{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
		},
	}, nil
}

func convertInternalBlindedBeaconBlockBellatrix(b *eth.BlindedBeaconBlockBellatrix) (*BlindedBeaconBlockBellatrix, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}

	return &BlindedBeaconBlockBellatrix{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BlindedBeaconBlockBodyBellatrix{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
			ExecutionPayloadHeader: &ExecutionPayloadHeader{
				ParentHash:       hexutil.Encode(b.Body.ExecutionPayloadHeader.ParentHash),
				FeeRecipient:     hexutil.Encode(b.Body.ExecutionPayloadHeader.FeeRecipient),
				StateRoot:        hexutil.Encode(b.Body.ExecutionPayloadHeader.StateRoot),
				ReceiptsRoot:     hexutil.Encode(b.Body.ExecutionPayloadHeader.ReceiptsRoot),
				LogsBloom:        hexutil.Encode(b.Body.ExecutionPayloadHeader.LogsBloom),
				PrevRandao:       hexutil.Encode(b.Body.ExecutionPayloadHeader.PrevRandao),
				BlockNumber:      fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.BlockNumber),
				GasLimit:         fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.GasLimit),
				GasUsed:          fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.GasUsed),
				Timestamp:        fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.Timestamp),
				ExtraData:        hexutil.Encode(b.Body.ExecutionPayloadHeader.ExtraData),
				BaseFeePerGas:    hexutil.Encode(b.Body.ExecutionPayloadHeader.BaseFeePerGas),
				BlockHash:        hexutil.Encode(b.Body.ExecutionPayloadHeader.BlockHash),
				TransactionsRoot: hexutil.Encode(b.Body.ExecutionPayloadHeader.TransactionsRoot),
			},
		},
	}, nil
}

func convertInternalBeaconBlockBellatrix(b *eth.BeaconBlockBellatrix) (*BeaconBlockBellatrix, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	transactions := make([]string, len(b.Body.ExecutionPayload.Transactions))
	for i, tx := range b.Body.ExecutionPayload.Transactions {
		transactions[i] = hexutil.Encode(tx)
	}
	return &BeaconBlockBellatrix{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BeaconBlockBodyBellatrix{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
			ExecutionPayload: &ExecutionPayload{
				ParentHash:    hexutil.Encode(b.Body.ExecutionPayload.ParentHash),
				FeeRecipient:  hexutil.Encode(b.Body.ExecutionPayload.FeeRecipient),
				StateRoot:     hexutil.Encode(b.Body.ExecutionPayload.StateRoot),
				ReceiptsRoot:  hexutil.Encode(b.Body.ExecutionPayload.ReceiptsRoot),
				LogsBloom:     hexutil.Encode(b.Body.ExecutionPayload.LogsBloom),
				PrevRandao:    hexutil.Encode(b.Body.ExecutionPayload.PrevRandao),
				BlockNumber:   fmt.Sprintf("%d", b.Body.ExecutionPayload.BlockNumber),
				GasLimit:      fmt.Sprintf("%d", b.Body.ExecutionPayload.GasLimit),
				GasUsed:       fmt.Sprintf("%d", b.Body.ExecutionPayload.GasUsed),
				Timestamp:     fmt.Sprintf("%d", b.Body.ExecutionPayload.Timestamp),
				ExtraData:     hexutil.Encode(b.Body.ExecutionPayload.ExtraData),
				BaseFeePerGas: hexutil.Encode(b.Body.ExecutionPayload.BaseFeePerGas),
				BlockHash:     hexutil.Encode(b.Body.ExecutionPayload.BlockHash),
				Transactions:  transactions,
			},
		},
	}, nil
}

func convertInternalBlindedBeaconBlockCapella(b *eth.BlindedBeaconBlockCapella) (*BlindedBeaconBlockCapella, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}

	blsChanges, err := convertInternalBlsChanges(b.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}

	return &BlindedBeaconBlockCapella{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BlindedBeaconBlockBodyCapella{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
			ExecutionPayloadHeader: &ExecutionPayloadHeaderCapella{
				ParentHash:       hexutil.Encode(b.Body.ExecutionPayloadHeader.ParentHash),
				FeeRecipient:     hexutil.Encode(b.Body.ExecutionPayloadHeader.FeeRecipient),
				StateRoot:        hexutil.Encode(b.Body.ExecutionPayloadHeader.StateRoot),
				ReceiptsRoot:     hexutil.Encode(b.Body.ExecutionPayloadHeader.ReceiptsRoot),
				LogsBloom:        hexutil.Encode(b.Body.ExecutionPayloadHeader.LogsBloom),
				PrevRandao:       hexutil.Encode(b.Body.ExecutionPayloadHeader.PrevRandao),
				BlockNumber:      fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.BlockNumber),
				GasLimit:         fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.GasLimit),
				GasUsed:          fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.GasUsed),
				Timestamp:        fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.Timestamp),
				ExtraData:        hexutil.Encode(b.Body.ExecutionPayloadHeader.ExtraData),
				BaseFeePerGas:    hexutil.Encode(b.Body.ExecutionPayloadHeader.BaseFeePerGas),
				BlockHash:        hexutil.Encode(b.Body.ExecutionPayloadHeader.BlockHash),
				TransactionsRoot: hexutil.Encode(b.Body.ExecutionPayloadHeader.TransactionsRoot),
				WithdrawalsRoot:  hexutil.Encode(b.Body.ExecutionPayloadHeader.WithdrawalsRoot), // new in capella
			},
			BlsToExecutionChanges: blsChanges, // new in capella
		},
	}, nil
}

func convertInternalBeaconBlockCapella(b *eth.BeaconBlockCapella) (*BeaconBlockCapella, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	transactions := make([]string, len(b.Body.ExecutionPayload.Transactions))
	for i, tx := range b.Body.ExecutionPayload.Transactions {
		transactions[i] = hexutil.Encode(tx)
	}
	withdrawals := make([]*Withdrawal, len(b.Body.ExecutionPayload.Withdrawals))
	for i, w := range b.Body.ExecutionPayload.Withdrawals {
		withdrawals[i] = &Withdrawal{
			WithdrawalIndex:  fmt.Sprintf("%d", w.Index),
			ValidatorIndex:   fmt.Sprintf("%d", w.ValidatorIndex),
			ExecutionAddress: hexutil.Encode(w.Address),
			Amount:           fmt.Sprintf("%d", w.Amount),
		}
	}
	blsChanges, err := convertInternalBlsChanges(b.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}
	return &BeaconBlockCapella{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BeaconBlockBodyCapella{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
			ExecutionPayload: &ExecutionPayloadCapella{
				ParentHash:    hexutil.Encode(b.Body.ExecutionPayload.ParentHash),
				FeeRecipient:  hexutil.Encode(b.Body.ExecutionPayload.FeeRecipient),
				StateRoot:     hexutil.Encode(b.Body.ExecutionPayload.StateRoot),
				ReceiptsRoot:  hexutil.Encode(b.Body.ExecutionPayload.ReceiptsRoot),
				LogsBloom:     hexutil.Encode(b.Body.ExecutionPayload.LogsBloom),
				PrevRandao:    hexutil.Encode(b.Body.ExecutionPayload.PrevRandao),
				BlockNumber:   fmt.Sprintf("%d", b.Body.ExecutionPayload.BlockNumber),
				GasLimit:      fmt.Sprintf("%d", b.Body.ExecutionPayload.GasLimit),
				GasUsed:       fmt.Sprintf("%d", b.Body.ExecutionPayload.GasUsed),
				Timestamp:     fmt.Sprintf("%d", b.Body.ExecutionPayload.Timestamp),
				ExtraData:     hexutil.Encode(b.Body.ExecutionPayload.ExtraData),
				BaseFeePerGas: hexutil.Encode(b.Body.ExecutionPayload.BaseFeePerGas),
				BlockHash:     hexutil.Encode(b.Body.ExecutionPayload.BlockHash),
				Transactions:  transactions,
				Withdrawals:   withdrawals, // new in capella
			},
			BlsToExecutionChanges: blsChanges, // new in capella
		},
	}, nil
}

func convertInternalBlindedBeaconBlockContentsDeneb(b *eth.BlindedBeaconBlockAndBlobsDeneb) (*BlindedBeaconBlockContentsDeneb, error) {
	if b == nil || b.Block == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	var blindedBlobSidecars []*BlindedBlobSidecar
	if len(b.Blobs) != 0 {
		blindedBlobSidecars = make([]*BlindedBlobSidecar, len(b.Blobs))
		for i, s := range b.Blobs {
			signedBlob, err := convertInternalToBlindedBlobSidecar(s)
			if err != nil {
				return nil, err
			}
			blindedBlobSidecars[i] = signedBlob
		}
	}
	blindedBlock, err := convertInternalToBlindedDenebBlock(b.Block)
	if err != nil {
		return nil, err
	}
	return &BlindedBeaconBlockContentsDeneb{
		BlindedBlock:        blindedBlock,
		BlindedBlobSidecars: blindedBlobSidecars,
	}, nil
}

func convertInternalBeaconBlockContentsDeneb(b *eth.BeaconBlockAndBlobsDeneb) (*BeaconBlockContentsDeneb, error) {
	if b == nil || b.Block == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	var blobSidecars []*BlobSidecar
	if len(b.Blobs) != 0 {
		blobSidecars = make([]*BlobSidecar, len(b.Blobs))
		for i, s := range b.Blobs {
			blob, err := convertInternalToBlobSidecar(s)
			if err != nil {
				return nil, err
			}
			blobSidecars[i] = blob
		}
	}
	block, err := convertInternalToDenebBlock(b.Block)
	if err != nil {
		return nil, err
	}
	return &BeaconBlockContentsDeneb{
		Block:        block,
		BlobSidecars: blobSidecars,
	}, nil
}

func (b *SignedBeaconBlockAltair) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	sig, err := hexutil.Decode(b.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Signature")
	}
	slot, err := strconv.ParseUint(b.Message.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Slot")
	}
	proposerIndex, err := strconv.ParseUint(b.Message.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(b.Message.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(b.Message.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(b.Message.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(b.Message.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(b.Message.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(b.Message.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(b.Message.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(b.Message.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(b.Message.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(b.Message.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(b.Message.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(b.Message.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := bytesutil.FromHexString(b.Message.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(b.Message.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.SyncAggregate.SyncCommitteeSignature")
	}

	block := &eth.SignedBeaconBlockAltair{
		Block: &eth.BeaconBlockAltair{
			Slot:          primitives.Slot(slot),
			ProposerIndex: primitives.ValidatorIndex(proposerIndex),
			ParentRoot:    parentRoot,
			StateRoot:     stateRoot,
			Body: &eth.BeaconBlockBodyAltair{
				RandaoReveal: randaoReveal,
				Eth1Data: &eth.Eth1Data{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				},
				Graffiti:          graffiti,
				ProposerSlashings: proposerSlashings,
				AttesterSlashings: attesterSlashings,
				Attestations:      atts,
				Deposits:          deposits,
				VoluntaryExits:    exits,
				SyncAggregate: &eth.SyncAggregate{
					SyncCommitteeBits:      syncCommitteeBits,
					SyncCommitteeSignature: syncCommitteeSig,
				},
			},
		},
		Signature: sig,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_Altair{Altair: block}}, nil
}

func (b *SignedBeaconBlockBellatrix) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	sig, err := hexutil.Decode(b.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Signature")
	}
	slot, err := strconv.ParseUint(b.Message.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Slot")
	}
	proposerIndex, err := strconv.ParseUint(b.Message.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(b.Message.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(b.Message.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(b.Message.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(b.Message.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(b.Message.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(b.Message.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(b.Message.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(b.Message.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(b.Message.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(b.Message.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(b.Message.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(b.Message.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := bytesutil.FromHexString(b.Message.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(b.Message.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.SyncAggregate.SyncCommitteeSignature")
	}
	payloadParentHash, err := hexutil.Decode(b.Message.Body.ExecutionPayload.ParentHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.ParentHash")
	}
	payloadFeeRecipient, err := hexutil.Decode(b.Message.Body.ExecutionPayload.FeeRecipient)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.FeeRecipient")
	}
	payloadStateRoot, err := hexutil.Decode(b.Message.Body.ExecutionPayload.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.StateRoot")
	}
	payloadReceiptsRoot, err := hexutil.Decode(b.Message.Body.ExecutionPayload.ReceiptsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.ReceiptsRoot")
	}
	payloadLogsBloom, err := hexutil.Decode(b.Message.Body.ExecutionPayload.LogsBloom)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.LogsBloom")
	}
	payloadPrevRandao, err := hexutil.Decode(b.Message.Body.ExecutionPayload.PrevRandao)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.PrevRandao")
	}
	payloadBlockNumber, err := strconv.ParseUint(b.Message.Body.ExecutionPayload.BlockNumber, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.BlockNumber")
	}
	payloadGasLimit, err := strconv.ParseUint(b.Message.Body.ExecutionPayload.GasLimit, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.GasLimit")
	}
	payloadGasUsed, err := strconv.ParseUint(b.Message.Body.ExecutionPayload.GasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.GasUsed")
	}
	payloadTimestamp, err := strconv.ParseUint(b.Message.Body.ExecutionPayload.Timestamp, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.Timestamp")
	}
	payloadExtraData, err := hexutil.Decode(b.Message.Body.ExecutionPayload.ExtraData)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.ExtraData")
	}
	payloadBaseFeePerGas, err := uint256ToHex(b.Message.Body.ExecutionPayload.BaseFeePerGas)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.BaseFeePerGas")
	}
	payloadBlockHash, err := hexutil.Decode(b.Message.Body.ExecutionPayload.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.BlockHash")
	}
	payloadTxs := make([][]byte, len(b.Message.Body.ExecutionPayload.Transactions))
	for i, tx := range b.Message.Body.ExecutionPayload.Transactions {
		payloadTxs[i], err = hexutil.Decode(tx)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ExecutionPayload.Transactions[%d]", i)
		}
	}

	block := &eth.SignedBeaconBlockBellatrix{
		Block: &eth.BeaconBlockBellatrix{
			Slot:          primitives.Slot(slot),
			ProposerIndex: primitives.ValidatorIndex(proposerIndex),
			ParentRoot:    parentRoot,
			StateRoot:     stateRoot,
			Body: &eth.BeaconBlockBodyBellatrix{
				RandaoReveal: randaoReveal,
				Eth1Data: &eth.Eth1Data{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				},
				Graffiti:          graffiti,
				ProposerSlashings: proposerSlashings,
				AttesterSlashings: attesterSlashings,
				Attestations:      atts,
				Deposits:          deposits,
				VoluntaryExits:    exits,
				SyncAggregate: &eth.SyncAggregate{
					SyncCommitteeBits:      syncCommitteeBits,
					SyncCommitteeSignature: syncCommitteeSig,
				},
				ExecutionPayload: &enginev1.ExecutionPayload{
					ParentHash:    payloadParentHash,
					FeeRecipient:  payloadFeeRecipient,
					StateRoot:     payloadStateRoot,
					ReceiptsRoot:  payloadReceiptsRoot,
					LogsBloom:     payloadLogsBloom,
					PrevRandao:    payloadPrevRandao,
					BlockNumber:   payloadBlockNumber,
					GasLimit:      payloadGasLimit,
					GasUsed:       payloadGasUsed,
					Timestamp:     payloadTimestamp,
					ExtraData:     payloadExtraData,
					BaseFeePerGas: payloadBaseFeePerGas,
					BlockHash:     payloadBlockHash,
					Transactions:  payloadTxs,
				},
			},
		},
		Signature: sig,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_Bellatrix{Bellatrix: block}}, nil
}

func (b *SignedBlindedBeaconBlockBellatrix) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	sig, err := hexutil.Decode(b.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Signature")
	}
	slot, err := strconv.ParseUint(b.Message.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Slot")
	}
	proposerIndex, err := strconv.ParseUint(b.Message.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(b.Message.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(b.Message.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(b.Message.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(b.Message.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(b.Message.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(b.Message.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(b.Message.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(b.Message.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(b.Message.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(b.Message.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(b.Message.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(b.Message.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := bytesutil.FromHexString(b.Message.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(b.Message.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.SyncAggregate.SyncCommitteeSignature")
	}
	payloadParentHash, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.ParentHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.ParentHash")
	}
	payloadFeeRecipient, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.FeeRecipient)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.FeeRecipient")
	}
	payloadStateRoot, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.StateRoot")
	}
	payloadReceiptsRoot, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.ReceiptsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.ReceiptsRoot")
	}
	payloadLogsBloom, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.LogsBloom)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.LogsBloom")
	}
	payloadPrevRandao, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.PrevRandao)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.PrevRandao")
	}
	payloadBlockNumber, err := strconv.ParseUint(b.Message.Body.ExecutionPayloadHeader.BlockNumber, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.BlockNumber")
	}
	payloadGasLimit, err := strconv.ParseUint(b.Message.Body.ExecutionPayloadHeader.GasLimit, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.GasLimit")
	}
	payloadGasUsed, err := strconv.ParseUint(b.Message.Body.ExecutionPayloadHeader.GasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.GasUsed")
	}
	payloadTimestamp, err := strconv.ParseUint(b.Message.Body.ExecutionPayloadHeader.Timestamp, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.Timestamp")
	}
	payloadExtraData, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.ExtraData)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.ExtraData")
	}
	payloadBaseFeePerGas, err := uint256ToHex(b.Message.Body.ExecutionPayloadHeader.BaseFeePerGas)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.BaseFeePerGas")
	}
	payloadBlockHash, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.BlockHash")
	}
	payloadTxsRoot, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.TransactionsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.TransactionsRoot")
	}

	block := &eth.SignedBlindedBeaconBlockBellatrix{
		Block: &eth.BlindedBeaconBlockBellatrix{
			Slot:          primitives.Slot(slot),
			ProposerIndex: primitives.ValidatorIndex(proposerIndex),
			ParentRoot:    parentRoot,
			StateRoot:     stateRoot,
			Body: &eth.BlindedBeaconBlockBodyBellatrix{
				RandaoReveal: randaoReveal,
				Eth1Data: &eth.Eth1Data{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				},
				Graffiti:          graffiti,
				ProposerSlashings: proposerSlashings,
				AttesterSlashings: attesterSlashings,
				Attestations:      atts,
				Deposits:          deposits,
				VoluntaryExits:    exits,
				SyncAggregate: &eth.SyncAggregate{
					SyncCommitteeBits:      syncCommitteeBits,
					SyncCommitteeSignature: syncCommitteeSig,
				},
				ExecutionPayloadHeader: &enginev1.ExecutionPayloadHeader{
					ParentHash:       payloadParentHash,
					FeeRecipient:     payloadFeeRecipient,
					StateRoot:        payloadStateRoot,
					ReceiptsRoot:     payloadReceiptsRoot,
					LogsBloom:        payloadLogsBloom,
					PrevRandao:       payloadPrevRandao,
					BlockNumber:      payloadBlockNumber,
					GasLimit:         payloadGasLimit,
					GasUsed:          payloadGasUsed,
					Timestamp:        payloadTimestamp,
					ExtraData:        payloadExtraData,
					BaseFeePerGas:    payloadBaseFeePerGas,
					BlockHash:        payloadBlockHash,
					TransactionsRoot: payloadTxsRoot,
				},
			},
		},
		Signature: sig,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_BlindedBellatrix{BlindedBellatrix: block}}, nil
}

func (b *SignedBeaconBlockCapella) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	sig, err := hexutil.Decode(b.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Signature")
	}
	slot, err := strconv.ParseUint(b.Message.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Slot")
	}
	proposerIndex, err := strconv.ParseUint(b.Message.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(b.Message.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(b.Message.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(b.Message.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(b.Message.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(b.Message.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(b.Message.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(b.Message.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(b.Message.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(b.Message.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(b.Message.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(b.Message.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(b.Message.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := bytesutil.FromHexString(b.Message.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(b.Message.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.SyncAggregate.SyncCommitteeSignature")
	}
	payloadParentHash, err := hexutil.Decode(b.Message.Body.ExecutionPayload.ParentHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.ParentHash")
	}
	payloadFeeRecipient, err := hexutil.Decode(b.Message.Body.ExecutionPayload.FeeRecipient)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.FeeRecipient")
	}
	payloadStateRoot, err := hexutil.Decode(b.Message.Body.ExecutionPayload.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.StateRoot")
	}
	payloadReceiptsRoot, err := hexutil.Decode(b.Message.Body.ExecutionPayload.ReceiptsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.ReceiptsRoot")
	}
	payloadLogsBloom, err := hexutil.Decode(b.Message.Body.ExecutionPayload.LogsBloom)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.LogsBloom")
	}
	payloadPrevRandao, err := hexutil.Decode(b.Message.Body.ExecutionPayload.PrevRandao)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.PrevRandao")
	}
	payloadBlockNumber, err := strconv.ParseUint(b.Message.Body.ExecutionPayload.BlockNumber, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.BlockNumber")
	}
	payloadGasLimit, err := strconv.ParseUint(b.Message.Body.ExecutionPayload.GasLimit, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.GasLimit")
	}
	payloadGasUsed, err := strconv.ParseUint(b.Message.Body.ExecutionPayload.GasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.GasUsed")
	}
	payloadTimestamp, err := strconv.ParseUint(b.Message.Body.ExecutionPayload.Timestamp, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.Timestamp")
	}
	payloadExtraData, err := hexutil.Decode(b.Message.Body.ExecutionPayload.ExtraData)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.ExtraData")
	}
	payloadBaseFeePerGas, err := uint256ToHex(b.Message.Body.ExecutionPayload.BaseFeePerGas)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.BaseFeePerGas")
	}
	payloadBlockHash, err := hexutil.Decode(b.Message.Body.ExecutionPayload.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayload.BlockHash")
	}
	txs := make([][]byte, len(b.Message.Body.ExecutionPayload.Transactions))
	for i, tx := range b.Message.Body.ExecutionPayload.Transactions {
		txs[i], err = hexutil.Decode(tx)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ExecutionPayload.Transactions[%d]", i)
		}
	}
	withdrawals := make([]*enginev1.Withdrawal, len(b.Message.Body.ExecutionPayload.Withdrawals))
	for i, w := range b.Message.Body.ExecutionPayload.Withdrawals {
		withdrawalIndex, err := strconv.ParseUint(w.WithdrawalIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ExecutionPayload.Withdrawals[%d].WithdrawalIndex", i)
		}
		validatorIndex, err := strconv.ParseUint(w.ValidatorIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ExecutionPayload.Withdrawals[%d].ValidatorIndex", i)
		}
		address, err := hexutil.Decode(w.ExecutionAddress)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ExecutionPayload.Withdrawals[%d].ExecutionAddress", i)
		}
		amount, err := strconv.ParseUint(w.Amount, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ExecutionPayload.Withdrawals[%d].Amount", i)
		}
		withdrawals[i] = &enginev1.Withdrawal{
			Index:          withdrawalIndex,
			ValidatorIndex: primitives.ValidatorIndex(validatorIndex),
			Address:        address,
			Amount:         amount,
		}
	}
	blsChanges, err := convertBlsChanges(b.Message.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}

	block := &eth.SignedBeaconBlockCapella{
		Block: &eth.BeaconBlockCapella{
			Slot:          primitives.Slot(slot),
			ProposerIndex: primitives.ValidatorIndex(proposerIndex),
			ParentRoot:    parentRoot,
			StateRoot:     stateRoot,
			Body: &eth.BeaconBlockBodyCapella{
				RandaoReveal: randaoReveal,
				Eth1Data: &eth.Eth1Data{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				},
				Graffiti:          graffiti,
				ProposerSlashings: proposerSlashings,
				AttesterSlashings: attesterSlashings,
				Attestations:      atts,
				Deposits:          deposits,
				VoluntaryExits:    exits,
				SyncAggregate: &eth.SyncAggregate{
					SyncCommitteeBits:      syncCommitteeBits,
					SyncCommitteeSignature: syncCommitteeSig,
				},
				ExecutionPayload: &enginev1.ExecutionPayloadCapella{
					ParentHash:    payloadParentHash,
					FeeRecipient:  payloadFeeRecipient,
					StateRoot:     payloadStateRoot,
					ReceiptsRoot:  payloadReceiptsRoot,
					LogsBloom:     payloadLogsBloom,
					PrevRandao:    payloadPrevRandao,
					BlockNumber:   payloadBlockNumber,
					GasLimit:      payloadGasLimit,
					GasUsed:       payloadGasUsed,
					Timestamp:     payloadTimestamp,
					ExtraData:     payloadExtraData,
					BaseFeePerGas: payloadBaseFeePerGas,
					BlockHash:     payloadBlockHash,
					Transactions:  txs,
					Withdrawals:   withdrawals,
				},
				BlsToExecutionChanges: blsChanges,
			},
		},
		Signature: sig,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_Capella{Capella: block}}, nil
}

func (b *SignedBlindedBeaconBlockCapella) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	sig, err := hexutil.Decode(b.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Signature")
	}
	slot, err := strconv.ParseUint(b.Message.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Slot")
	}
	proposerIndex, err := strconv.ParseUint(b.Message.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(b.Message.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(b.Message.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(b.Message.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(b.Message.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(b.Message.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(b.Message.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(b.Message.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(b.Message.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(b.Message.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(b.Message.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(b.Message.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(b.Message.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := bytesutil.FromHexString(b.Message.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(b.Message.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.SyncAggregate.SyncCommitteeSignature")
	}
	payloadParentHash, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.ParentHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.ParentHash")
	}
	payloadFeeRecipient, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.FeeRecipient)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.FeeRecipient")
	}
	payloadStateRoot, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.StateRoot")
	}
	payloadReceiptsRoot, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.ReceiptsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.ReceiptsRoot")
	}
	payloadLogsBloom, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.LogsBloom)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.LogsBloom")
	}
	payloadPrevRandao, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.PrevRandao)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.PrevRandao")
	}
	payloadBlockNumber, err := strconv.ParseUint(b.Message.Body.ExecutionPayloadHeader.BlockNumber, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.BlockNumber")
	}
	payloadGasLimit, err := strconv.ParseUint(b.Message.Body.ExecutionPayloadHeader.GasLimit, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.GasLimit")
	}
	payloadGasUsed, err := strconv.ParseUint(b.Message.Body.ExecutionPayloadHeader.GasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.GasUsed")
	}
	payloadTimestamp, err := strconv.ParseUint(b.Message.Body.ExecutionPayloadHeader.Timestamp, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.Timestamp")
	}
	payloadExtraData, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.ExtraData)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.ExtraData")
	}
	payloadBaseFeePerGas, err := uint256ToHex(b.Message.Body.ExecutionPayloadHeader.BaseFeePerGas)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.BaseFeePerGas")
	}
	payloadBlockHash, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.BlockHash")
	}
	payloadTxsRoot, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.TransactionsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.TransactionsRoot")
	}
	payloadWithdrawalsRoot, err := hexutil.Decode(b.Message.Body.ExecutionPayloadHeader.WithdrawalsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode b.Message.Body.ExecutionPayloadHeader.WithdrawalsRoot")
	}
	blsChanges, err := convertBlsChanges(b.Message.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}

	block := &eth.SignedBlindedBeaconBlockCapella{
		Block: &eth.BlindedBeaconBlockCapella{
			Slot:          primitives.Slot(slot),
			ProposerIndex: primitives.ValidatorIndex(proposerIndex),
			ParentRoot:    parentRoot,
			StateRoot:     stateRoot,
			Body: &eth.BlindedBeaconBlockBodyCapella{
				RandaoReveal: randaoReveal,
				Eth1Data: &eth.Eth1Data{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				},
				Graffiti:          graffiti,
				ProposerSlashings: proposerSlashings,
				AttesterSlashings: attesterSlashings,
				Attestations:      atts,
				Deposits:          deposits,
				VoluntaryExits:    exits,
				SyncAggregate: &eth.SyncAggregate{
					SyncCommitteeBits:      syncCommitteeBits,
					SyncCommitteeSignature: syncCommitteeSig,
				},
				ExecutionPayloadHeader: &enginev1.ExecutionPayloadHeaderCapella{
					ParentHash:       payloadParentHash,
					FeeRecipient:     payloadFeeRecipient,
					StateRoot:        payloadStateRoot,
					ReceiptsRoot:     payloadReceiptsRoot,
					LogsBloom:        payloadLogsBloom,
					PrevRandao:       payloadPrevRandao,
					BlockNumber:      payloadBlockNumber,
					GasLimit:         payloadGasLimit,
					GasUsed:          payloadGasUsed,
					Timestamp:        payloadTimestamp,
					ExtraData:        payloadExtraData,
					BaseFeePerGas:    payloadBaseFeePerGas,
					BlockHash:        payloadBlockHash,
					TransactionsRoot: payloadTxsRoot,
					WithdrawalsRoot:  payloadWithdrawalsRoot,
				},
				BlsToExecutionChanges: blsChanges,
			},
		},
		Signature: sig,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_BlindedCapella{BlindedCapella: block}}, nil
}

func (b *SignedBeaconBlockContentsDeneb) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	var signedBlobSidecars []*eth.SignedBlobSidecar
	if len(b.SignedBlobSidecars) != 0 {
		signedBlobSidecars = make([]*eth.SignedBlobSidecar, len(b.SignedBlobSidecars))
		for i, s := range b.SignedBlobSidecars {
			signedBlob, err := convertToSignedBlobSidecar(i, s)
			if err != nil {
				return nil, err
			}
			signedBlobSidecars[i] = signedBlob
		}
	}
	signedDenebBlock, err := convertToSignedDenebBlock(b.SignedBlock)
	if err != nil {
		return nil, err
	}
	block := &eth.SignedBeaconBlockAndBlobsDeneb{
		Block: signedDenebBlock,
		Blobs: signedBlobSidecars,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_Deneb{Deneb: block}}, nil
}

func convertToSignedDenebBlock(signedBlock *SignedBeaconBlockDeneb) (*eth.SignedBeaconBlockDeneb, error) {
	sig, err := hexutil.Decode(signedBlock.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock .Signature")
	}
	slot, err := strconv.ParseUint(signedBlock.Message.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Slot")
	}
	proposerIndex, err := strconv.ParseUint(signedBlock.Message.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(signedBlock.Message.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(signedBlock.Message.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(signedBlock.Message.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(signedBlock.Message.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(signedBlock.Message.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(signedBlock.Message.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(signedBlock.Message.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(signedBlock.Message.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(signedBlock.Message.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(signedBlock.Message.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(signedBlock.Message.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(signedBlock.Message.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := bytesutil.FromHexString(signedBlock.Message.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(signedBlock.Message.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.SyncAggregate.SyncCommitteeSignature")
	}
	payloadParentHash, err := hexutil.Decode(signedBlock.Message.Body.ExecutionPayload.ParentHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.ParentHash")
	}
	payloadFeeRecipient, err := hexutil.Decode(signedBlock.Message.Body.ExecutionPayload.FeeRecipient)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.FeeRecipient")
	}
	payloadStateRoot, err := hexutil.Decode(signedBlock.Message.Body.ExecutionPayload.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.StateRoot")
	}
	payloadReceiptsRoot, err := hexutil.Decode(signedBlock.Message.Body.ExecutionPayload.ReceiptsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.ReceiptsRoot")
	}
	payloadLogsBloom, err := hexutil.Decode(signedBlock.Message.Body.ExecutionPayload.LogsBloom)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.LogsBloom")
	}
	payloadPrevRandao, err := hexutil.Decode(signedBlock.Message.Body.ExecutionPayload.PrevRandao)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.PrevRandao")
	}
	payloadBlockNumber, err := strconv.ParseUint(signedBlock.Message.Body.ExecutionPayload.BlockNumber, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.BlockNumber")
	}
	payloadGasLimit, err := strconv.ParseUint(signedBlock.Message.Body.ExecutionPayload.GasLimit, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.GasLimit")
	}
	payloadGasUsed, err := strconv.ParseUint(signedBlock.Message.Body.ExecutionPayload.GasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.GasUsed")
	}
	payloadTimestamp, err := strconv.ParseUint(signedBlock.Message.Body.ExecutionPayload.Timestamp, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayloadHeader.Timestamp")
	}
	payloadExtraData, err := hexutil.Decode(signedBlock.Message.Body.ExecutionPayload.ExtraData)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.ExtraData")
	}
	payloadBaseFeePerGas, err := uint256ToHex(signedBlock.Message.Body.ExecutionPayload.BaseFeePerGas)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.BaseFeePerGas")
	}
	payloadBlockHash, err := hexutil.Decode(signedBlock.Message.Body.ExecutionPayload.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.BlockHash")
	}
	txs := make([][]byte, len(signedBlock.Message.Body.ExecutionPayload.Transactions))
	for i, tx := range signedBlock.Message.Body.ExecutionPayload.Transactions {
		txs[i], err = hexutil.Decode(tx)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode signedBlock.Message.Body.ExecutionPayload.Transactions[%d]", i)
		}
	}
	withdrawals := make([]*enginev1.Withdrawal, len(signedBlock.Message.Body.ExecutionPayload.Withdrawals))
	for i, w := range signedBlock.Message.Body.ExecutionPayload.Withdrawals {
		withdrawalIndex, err := strconv.ParseUint(w.WithdrawalIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode signedBlock.Message.Body.ExecutionPayload.Withdrawals[%d].WithdrawalIndex", i)
		}
		validatorIndex, err := strconv.ParseUint(w.ValidatorIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode signedBlock.Message.Body.ExecutionPayload.Withdrawals[%d].ValidatorIndex", i)
		}
		address, err := hexutil.Decode(w.ExecutionAddress)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ExecutionPayload.Withdrawals[%d].ExecutionAddress", i)
		}
		amount, err := strconv.ParseUint(w.Amount, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ExecutionPayload.Withdrawals[%d].Amount", i)
		}
		withdrawals[i] = &enginev1.Withdrawal{
			Index:          withdrawalIndex,
			ValidatorIndex: primitives.ValidatorIndex(validatorIndex),
			Address:        address,
			Amount:         amount,
		}
	}
	blsChanges, err := convertBlsChanges(signedBlock.Message.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}
	payloadDataGasUsed, err := strconv.ParseUint(signedBlock.Message.Body.ExecutionPayload.DataGasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.DataGasUsed")
	}
	payloadExcessDataGas, err := strconv.ParseUint(signedBlock.Message.Body.ExecutionPayload.ExcessDataGas, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlock.Message.Body.ExecutionPayload.ExcessDataGas")
	}
	return &eth.SignedBeaconBlockDeneb{
		Block: &eth.BeaconBlockDeneb{
			Slot:          primitives.Slot(slot),
			ProposerIndex: primitives.ValidatorIndex(proposerIndex),
			ParentRoot:    parentRoot,
			StateRoot:     stateRoot,
			Body: &eth.BeaconBlockBodyDeneb{
				RandaoReveal: randaoReveal,
				Eth1Data: &eth.Eth1Data{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				},
				Graffiti:          graffiti,
				ProposerSlashings: proposerSlashings,
				AttesterSlashings: attesterSlashings,
				Attestations:      atts,
				Deposits:          deposits,
				VoluntaryExits:    exits,
				SyncAggregate: &eth.SyncAggregate{
					SyncCommitteeBits:      syncCommitteeBits,
					SyncCommitteeSignature: syncCommitteeSig,
				},
				ExecutionPayload: &enginev1.ExecutionPayloadDeneb{
					ParentHash:    payloadParentHash,
					FeeRecipient:  payloadFeeRecipient,
					StateRoot:     payloadStateRoot,
					ReceiptsRoot:  payloadReceiptsRoot,
					LogsBloom:     payloadLogsBloom,
					PrevRandao:    payloadPrevRandao,
					BlockNumber:   payloadBlockNumber,
					GasLimit:      payloadGasLimit,
					GasUsed:       payloadGasUsed,
					Timestamp:     payloadTimestamp,
					ExtraData:     payloadExtraData,
					BaseFeePerGas: payloadBaseFeePerGas,
					BlockHash:     payloadBlockHash,
					Transactions:  txs,
					Withdrawals:   withdrawals,
					DataGasUsed:   payloadDataGasUsed,
					ExcessDataGas: payloadExcessDataGas,
				},
				BlsToExecutionChanges: blsChanges,
			},
		},
		Signature: sig,
	}, nil
}

func convertToSignedBlobSidecar(i int, signedBlob *SignedBlobSidecar) (*eth.SignedBlobSidecar, error) {
	blobSig, err := hexutil.Decode(signedBlob.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlob.Signature")
	}
	if signedBlob.Message == nil {
		return nil, fmt.Errorf("blobsidecar message was empty at index %d", i)
	}
	blockRoot, err := hexutil.Decode(signedBlob.Message.BlockRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.BlockRoot at index %d", i))
	}
	index, err := strconv.ParseUint(signedBlob.Message.Index, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.Index at index %d", i))
	}
	slot, err := strconv.ParseUint(signedBlob.Message.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.Index at index %d", i))
	}
	blockParentRoot, err := hexutil.Decode(signedBlob.Message.BlockParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.BlockParentRoot at index %d", i))
	}
	proposerIndex, err := strconv.ParseUint(signedBlob.Message.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.ProposerIndex at index %d", i))
	}
	blob, err := hexutil.Decode(signedBlob.Message.Blob)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.Blob at index %d", i))
	}
	kzgCommitment, err := hexutil.Decode(signedBlob.Message.KzgCommitment)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.KzgCommitment at index %d", i))
	}
	kzgProof, err := hexutil.Decode(signedBlob.Message.KzgProof)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.KzgProof at index %d", i))
	}
	bsc := &eth.BlobSidecar{
		BlockRoot:       blockRoot,
		Index:           index,
		Slot:            primitives.Slot(slot),
		BlockParentRoot: blockParentRoot,
		ProposerIndex:   primitives.ValidatorIndex(proposerIndex),
		Blob:            blob,
		KzgCommitment:   kzgCommitment,
		KzgProof:        kzgProof,
	}
	return &eth.SignedBlobSidecar{
		Message:   bsc,
		Signature: blobSig,
	}, nil
}

func (b *SignedBlindedBeaconBlockContentsDeneb) ToGeneric() (*eth.GenericSignedBeaconBlock, error) {
	var signedBlindedBlobSidecars []*eth.SignedBlindedBlobSidecar
	if len(b.SignedBlindedBlobSidecars) != 0 {
		signedBlindedBlobSidecars = make([]*eth.SignedBlindedBlobSidecar, len(b.SignedBlindedBlobSidecars))
		for i, s := range b.SignedBlindedBlobSidecars {
			signedBlob, err := convertToSignedBlindedBlobSidecar(i, s)
			if err != nil {
				return nil, err
			}
			signedBlindedBlobSidecars[i] = signedBlob
		}
	}
	signedBlindedBlock, err := convertToSignedBlindedDenebBlock(b.SignedBlindedBlock)
	if err != nil {
		return nil, err
	}
	block := &eth.SignedBlindedBeaconBlockAndBlobsDeneb{
		Block: signedBlindedBlock,
		Blobs: signedBlindedBlobSidecars,
	}
	return &eth.GenericSignedBeaconBlock{Block: &eth.GenericSignedBeaconBlock_BlindedDeneb{BlindedDeneb: block}}, nil
}

func convertToSignedBlindedDenebBlock(signedBlindedBlock *SignedBlindedBeaconBlockDeneb) (*eth.SignedBlindedBeaconBlockDeneb, error) {
	if signedBlindedBlock == nil {
		return nil, errors.New("signed blinded block is empty")
	}
	sig, err := hexutil.Decode(signedBlindedBlock.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Signature")
	}
	slot, err := strconv.ParseUint(signedBlindedBlock.Message.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Slot")
	}
	proposerIndex, err := strconv.ParseUint(signedBlindedBlock.Message.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.ProposerIndex")
	}
	parentRoot, err := hexutil.Decode(signedBlindedBlock.Message.ParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.ParentRoot")
	}
	stateRoot, err := hexutil.Decode(signedBlindedBlock.Message.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.StateRoot")
	}
	randaoReveal, err := hexutil.Decode(signedBlindedBlock.Message.Body.RandaoReveal)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.RandaoReveal")
	}
	depositRoot, err := hexutil.Decode(signedBlindedBlock.Message.Body.Eth1Data.DepositRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.Eth1Data.DepositRoot")
	}
	depositCount, err := strconv.ParseUint(signedBlindedBlock.Message.Body.Eth1Data.DepositCount, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.Eth1Data.DepositCount")
	}
	blockHash, err := hexutil.Decode(signedBlindedBlock.Message.Body.Eth1Data.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.Eth1Data.BlockHash")
	}
	graffiti, err := hexutil.Decode(signedBlindedBlock.Message.Body.Graffiti)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.Graffiti")
	}
	proposerSlashings, err := convertProposerSlashings(signedBlindedBlock.Message.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertAttesterSlashings(signedBlindedBlock.Message.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertAtts(signedBlindedBlock.Message.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertDeposits(signedBlindedBlock.Message.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertExits(signedBlindedBlock.Message.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	syncCommitteeBits, err := bytesutil.FromHexString(signedBlindedBlock.Message.Body.SyncAggregate.SyncCommitteeBits)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.SyncAggregate.SyncCommitteeBits")
	}
	syncCommitteeSig, err := hexutil.Decode(signedBlindedBlock.Message.Body.SyncAggregate.SyncCommitteeSignature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.SyncAggregate.SyncCommitteeSignature")
	}
	payloadParentHash, err := hexutil.Decode(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.ParentHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.ParentHash")
	}
	payloadFeeRecipient, err := hexutil.Decode(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.FeeRecipient)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.FeeRecipient")
	}
	payloadStateRoot, err := hexutil.Decode(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.StateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.StateRoot")
	}
	payloadReceiptsRoot, err := hexutil.Decode(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.ReceiptsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.ReceiptsRoot")
	}
	payloadLogsBloom, err := hexutil.Decode(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.LogsBloom)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.LogsBloom")
	}
	payloadPrevRandao, err := hexutil.Decode(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.PrevRandao)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.PrevRandao")
	}
	payloadBlockNumber, err := strconv.ParseUint(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.BlockNumber, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.BlockNumber")
	}
	payloadGasLimit, err := strconv.ParseUint(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.GasLimit, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.GasLimit")
	}
	payloadGasUsed, err := strconv.ParseUint(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.GasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.GasUsed")
	}
	payloadTimestamp, err := strconv.ParseUint(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.Timestamp, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.Timestamp")
	}
	payloadExtraData, err := hexutil.Decode(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.ExtraData)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.ExtraData")
	}
	payloadBaseFeePerGas, err := uint256ToHex(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.BaseFeePerGas)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.BaseFeePerGas")
	}
	payloadBlockHash, err := hexutil.Decode(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.BlockHash")
	}
	payloadTxsRoot, err := hexutil.Decode(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.TransactionsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.TransactionsRoot")
	}
	payloadWithdrawalsRoot, err := hexutil.Decode(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.WithdrawalsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayloadHeader.WithdrawalsRoot")
	}
	blsChanges, err := convertBlsChanges(signedBlindedBlock.Message.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}
	payloadDataGasUsed, err := strconv.ParseUint(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.DataGasUsed, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayload.DataGasUsed")
	}
	payloadExcessDataGas, err := strconv.ParseUint(signedBlindedBlock.Message.Body.ExecutionPayloadHeader.ExcessDataGas, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlindedBlock.Message.Body.ExecutionPayload.ExcessDataGas")
	}
	return &eth.SignedBlindedBeaconBlockDeneb{
		Block: &eth.BlindedBeaconBlockDeneb{
			Slot:          primitives.Slot(slot),
			ProposerIndex: primitives.ValidatorIndex(proposerIndex),
			ParentRoot:    parentRoot,
			StateRoot:     stateRoot,
			Body: &eth.BlindedBeaconBlockBodyDeneb{
				RandaoReveal: randaoReveal,
				Eth1Data: &eth.Eth1Data{
					DepositRoot:  depositRoot,
					DepositCount: depositCount,
					BlockHash:    blockHash,
				},
				Graffiti:          graffiti,
				ProposerSlashings: proposerSlashings,
				AttesterSlashings: attesterSlashings,
				Attestations:      atts,
				Deposits:          deposits,
				VoluntaryExits:    exits,
				SyncAggregate: &eth.SyncAggregate{
					SyncCommitteeBits:      syncCommitteeBits,
					SyncCommitteeSignature: syncCommitteeSig,
				},
				ExecutionPayloadHeader: &enginev1.ExecutionPayloadHeaderDeneb{
					ParentHash:       payloadParentHash,
					FeeRecipient:     payloadFeeRecipient,
					StateRoot:        payloadStateRoot,
					ReceiptsRoot:     payloadReceiptsRoot,
					LogsBloom:        payloadLogsBloom,
					PrevRandao:       payloadPrevRandao,
					BlockNumber:      payloadBlockNumber,
					GasLimit:         payloadGasLimit,
					GasUsed:          payloadGasUsed,
					Timestamp:        payloadTimestamp,
					ExtraData:        payloadExtraData,
					BaseFeePerGas:    payloadBaseFeePerGas,
					BlockHash:        payloadBlockHash,
					TransactionsRoot: payloadTxsRoot,
					WithdrawalsRoot:  payloadWithdrawalsRoot,
					DataGasUsed:      payloadDataGasUsed,
					ExcessDataGas:    payloadExcessDataGas,
				},
				BlsToExecutionChanges: blsChanges,
			},
		},
		Signature: sig,
	}, nil
}

func convertInternalToBlindedDenebBlock(b *eth.BlindedBeaconBlockDeneb) (*BlindedBeaconBlockDeneb, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}

	blsChanges, err := convertInternalBlsChanges(b.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}

	return &BlindedBeaconBlockDeneb{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BlindedBeaconBlockBodyDeneb{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
			ExecutionPayloadHeader: &ExecutionPayloadHeaderDeneb{
				ParentHash:       hexutil.Encode(b.Body.ExecutionPayloadHeader.ParentHash),
				FeeRecipient:     hexutil.Encode(b.Body.ExecutionPayloadHeader.FeeRecipient),
				StateRoot:        hexutil.Encode(b.Body.ExecutionPayloadHeader.StateRoot),
				ReceiptsRoot:     hexutil.Encode(b.Body.ExecutionPayloadHeader.ReceiptsRoot),
				LogsBloom:        hexutil.Encode(b.Body.ExecutionPayloadHeader.LogsBloom),
				PrevRandao:       hexutil.Encode(b.Body.ExecutionPayloadHeader.PrevRandao),
				BlockNumber:      fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.BlockNumber),
				GasLimit:         fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.GasLimit),
				GasUsed:          fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.GasUsed),
				Timestamp:        fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.Timestamp),
				ExtraData:        hexutil.Encode(b.Body.ExecutionPayloadHeader.ExtraData),
				BaseFeePerGas:    hexutil.Encode(b.Body.ExecutionPayloadHeader.BaseFeePerGas),
				BlockHash:        hexutil.Encode(b.Body.ExecutionPayloadHeader.BlockHash),
				TransactionsRoot: hexutil.Encode(b.Body.ExecutionPayloadHeader.TransactionsRoot),
				WithdrawalsRoot:  hexutil.Encode(b.Body.ExecutionPayloadHeader.WithdrawalsRoot),
				DataGasUsed:      fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.DataGasUsed),   // new in deneb TODO: rename to blob
				ExcessDataGas:    fmt.Sprintf("%d", b.Body.ExecutionPayloadHeader.ExcessDataGas), // new in deneb TODO: rename to blob
			},
			BlsToExecutionChanges: blsChanges, // new in capella
		},
	}, nil
}

func convertInternalToDenebBlock(b *eth.BeaconBlockDeneb) (*BeaconBlockDeneb, error) {
	if b == nil {
		return nil, errors.New("block is empty, nothing to convert.")
	}
	proposerSlashings, err := convertInternalProposerSlashings(b.Body.ProposerSlashings)
	if err != nil {
		return nil, err
	}
	attesterSlashings, err := convertInternalAttesterSlashings(b.Body.AttesterSlashings)
	if err != nil {
		return nil, err
	}
	atts, err := convertInternalAtts(b.Body.Attestations)
	if err != nil {
		return nil, err
	}
	deposits, err := convertInternalDeposits(b.Body.Deposits)
	if err != nil {
		return nil, err
	}
	exits, err := convertInternalExits(b.Body.VoluntaryExits)
	if err != nil {
		return nil, err
	}
	transactions := make([]string, len(b.Body.ExecutionPayload.Transactions))
	for i, tx := range b.Body.ExecutionPayload.Transactions {
		transactions[i] = hexutil.Encode(tx)
	}
	withdrawals := make([]*Withdrawal, len(b.Body.ExecutionPayload.Withdrawals))
	for i, w := range b.Body.ExecutionPayload.Withdrawals {
		withdrawals[i] = &Withdrawal{
			WithdrawalIndex:  fmt.Sprintf("%d", w.Index),
			ValidatorIndex:   fmt.Sprintf("%d", w.ValidatorIndex),
			ExecutionAddress: hexutil.Encode(w.Address),
			Amount:           fmt.Sprintf("%d", w.Amount),
		}
	}
	blsChanges, err := convertInternalBlsChanges(b.Body.BlsToExecutionChanges)
	if err != nil {
		return nil, err
	}

	return &BeaconBlockDeneb{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    hexutil.Encode(b.ParentRoot),
		StateRoot:     hexutil.Encode(b.StateRoot),
		Body: &BeaconBlockBodyDeneb{
			RandaoReveal: hexutil.Encode(b.Body.RandaoReveal),
			Eth1Data: &Eth1Data{
				DepositRoot:  hexutil.Encode(b.Body.Eth1Data.DepositRoot),
				DepositCount: fmt.Sprintf("%d", b.Body.Eth1Data.DepositCount),
				BlockHash:    hexutil.Encode(b.Body.Eth1Data.BlockHash),
			},
			Graffiti:          hexutil.Encode(b.Body.Graffiti),
			ProposerSlashings: proposerSlashings,
			AttesterSlashings: attesterSlashings,
			Attestations:      atts,
			Deposits:          deposits,
			VoluntaryExits:    exits,
			SyncAggregate: &SyncAggregate{
				SyncCommitteeBits:      hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeBits),
				SyncCommitteeSignature: hexutil.Encode(b.Body.SyncAggregate.SyncCommitteeSignature),
			},
			ExecutionPayload: &ExecutionPayloadDeneb{
				ParentHash:    hexutil.Encode(b.Body.ExecutionPayload.ParentHash),
				FeeRecipient:  hexutil.Encode(b.Body.ExecutionPayload.FeeRecipient),
				StateRoot:     hexutil.Encode(b.Body.ExecutionPayload.StateRoot),
				ReceiptsRoot:  hexutil.Encode(b.Body.ExecutionPayload.ReceiptsRoot),
				LogsBloom:     hexutil.Encode(b.Body.ExecutionPayload.LogsBloom),
				PrevRandao:    hexutil.Encode(b.Body.ExecutionPayload.PrevRandao),
				BlockNumber:   fmt.Sprintf("%d", b.Body.ExecutionPayload.BlockNumber),
				GasLimit:      fmt.Sprintf("%d", b.Body.ExecutionPayload.GasLimit),
				GasUsed:       fmt.Sprintf("%d", b.Body.ExecutionPayload.GasUsed),
				Timestamp:     fmt.Sprintf("%d", b.Body.ExecutionPayload.Timestamp),
				ExtraData:     hexutil.Encode(b.Body.ExecutionPayload.ExtraData),
				BaseFeePerGas: hexutil.Encode(b.Body.ExecutionPayload.BaseFeePerGas),
				BlockHash:     hexutil.Encode(b.Body.ExecutionPayload.BlockHash),
				Transactions:  transactions,
				Withdrawals:   withdrawals,
				DataGasUsed:   fmt.Sprintf("%d", b.Body.ExecutionPayload.DataGasUsed),   // new in deneb TODO: rename to blob
				ExcessDataGas: fmt.Sprintf("%d", b.Body.ExecutionPayload.ExcessDataGas), // new in deneb TODO: rename to blob
			},
			BlsToExecutionChanges: blsChanges, // new in capella
		},
	}, nil
}

func convertToSignedBlindedBlobSidecar(i int, signedBlob *SignedBlindedBlobSidecar) (*eth.SignedBlindedBlobSidecar, error) {
	blobSig, err := hexutil.Decode(signedBlob.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode signedBlob.Signature")
	}
	if signedBlob.Message == nil {
		return nil, fmt.Errorf("blobsidecar message was empty at index %d", i)
	}
	blockRoot, err := hexutil.Decode(signedBlob.Message.BlockRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.BlockRoot at index %d", i))
	}
	index, err := strconv.ParseUint(signedBlob.Message.Index, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.Index at index %d", i))
	}
	denebSlot, err := strconv.ParseUint(signedBlob.Message.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.Index at index %d", i))
	}
	blockParentRoot, err := hexutil.Decode(signedBlob.Message.BlockParentRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.BlockParentRoot at index %d", i))
	}
	proposerIndex, err := strconv.ParseUint(signedBlob.Message.ProposerIndex, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.ProposerIndex at index %d", i))
	}
	blobRoot, err := hexutil.Decode(signedBlob.Message.BlobRoot)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.BlobRoot at index %d", i))
	}
	kzgCommitment, err := hexutil.Decode(signedBlob.Message.KzgCommitment)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.KzgCommitment at index %d", i))
	}
	kzgProof, err := hexutil.Decode(signedBlob.Message.KzgProof)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not decode signedBlob.Message.KzgProof at index %d", i))
	}
	bsc := &eth.BlindedBlobSidecar{
		BlockRoot:       blockRoot,
		Index:           index,
		Slot:            primitives.Slot(denebSlot),
		BlockParentRoot: blockParentRoot,
		ProposerIndex:   primitives.ValidatorIndex(proposerIndex),
		BlobRoot:        blobRoot,
		KzgCommitment:   kzgCommitment,
		KzgProof:        kzgProof,
	}
	return &eth.SignedBlindedBlobSidecar{
		Message:   bsc,
		Signature: blobSig,
	}, nil
}

func convertInternalToBlindedBlobSidecar(b *eth.BlindedBlobSidecar) (*BlindedBlobSidecar, error) {
	if b == nil {
		return nil, errors.New("BlindedBlobSidecar is empty, nothing to convert.")
	}
	return &BlindedBlobSidecar{
		BlockRoot:       hexutil.Encode(b.BlockRoot),
		Index:           fmt.Sprintf("%d", b.Index),
		Slot:            fmt.Sprintf("%d", b.Slot),
		BlockParentRoot: hexutil.Encode(b.BlockParentRoot),
		ProposerIndex:   fmt.Sprintf("%d", b.ProposerIndex),
		BlobRoot:        hexutil.Encode(b.BlobRoot),
		KzgCommitment:   hexutil.Encode(b.KzgCommitment),
		KzgProof:        hexutil.Encode(b.KzgProof),
	}, nil
}

func convertInternalToBlobSidecar(b *eth.BlobSidecar) (*BlobSidecar, error) {
	if b == nil {
		return nil, errors.New("BlobSidecar is empty, nothing to convert.")
	}
	return &BlobSidecar{
		BlockRoot:       hexutil.Encode(b.BlockRoot),
		Index:           fmt.Sprintf("%d", b.Index),
		Slot:            fmt.Sprintf("%d", b.Slot),
		BlockParentRoot: hexutil.Encode(b.BlockParentRoot),
		ProposerIndex:   fmt.Sprintf("%d", b.ProposerIndex),
		Blob:            hexutil.Encode(b.Blob),
		KzgCommitment:   hexutil.Encode(b.KzgCommitment),
		KzgProof:        hexutil.Encode(b.KzgProof),
	}, nil
}

func convertProposerSlashings(src []*ProposerSlashing) ([]*eth.ProposerSlashing, error) {
	if src == nil {
		return nil, errors.New("nil b.Message.Body.ProposerSlashings")
	}
	proposerSlashings := make([]*eth.ProposerSlashing, len(src))
	for i, s := range src {
		h1Sig, err := hexutil.Decode(s.SignedHeader1.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader1.Signature", i)
		}
		h1Slot, err := strconv.ParseUint(s.SignedHeader1.Message.Slot, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader1.Message.Slot", i)
		}
		h1ProposerIndex, err := strconv.ParseUint(s.SignedHeader1.Message.ProposerIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader1.Message.ProposerIndex", i)
		}
		h1ParentRoot, err := hexutil.Decode(s.SignedHeader1.Message.ParentRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader1.Message.ParentRoot", i)
		}
		h1StateRoot, err := hexutil.Decode(s.SignedHeader1.Message.StateRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader1.Message.StateRoot", i)
		}
		h1BodyRoot, err := hexutil.Decode(s.SignedHeader1.Message.BodyRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader1.Message.BodyRoot", i)
		}
		h2Sig, err := hexutil.Decode(s.SignedHeader2.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader2.Signature", i)
		}
		h2Slot, err := strconv.ParseUint(s.SignedHeader2.Message.Slot, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader2.Message.Slot", i)
		}
		h2ProposerIndex, err := strconv.ParseUint(s.SignedHeader2.Message.ProposerIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader2.Message.ProposerIndex", i)
		}
		h2ParentRoot, err := hexutil.Decode(s.SignedHeader2.Message.ParentRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader2.Message.ParentRoot", i)
		}
		h2StateRoot, err := hexutil.Decode(s.SignedHeader2.Message.StateRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader2.Message.StateRoot", i)
		}
		h2BodyRoot, err := hexutil.Decode(s.SignedHeader2.Message.BodyRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.ProposerSlashings[%d].SignedHeader2.Message.BodyRoot", i)
		}
		proposerSlashings[i] = &eth.ProposerSlashing{
			Header_1: &eth.SignedBeaconBlockHeader{
				Header: &eth.BeaconBlockHeader{
					Slot:          primitives.Slot(h1Slot),
					ProposerIndex: primitives.ValidatorIndex(h1ProposerIndex),
					ParentRoot:    h1ParentRoot,
					StateRoot:     h1StateRoot,
					BodyRoot:      h1BodyRoot,
				},
				Signature: h1Sig,
			},
			Header_2: &eth.SignedBeaconBlockHeader{
				Header: &eth.BeaconBlockHeader{
					Slot:          primitives.Slot(h2Slot),
					ProposerIndex: primitives.ValidatorIndex(h2ProposerIndex),
					ParentRoot:    h2ParentRoot,
					StateRoot:     h2StateRoot,
					BodyRoot:      h2BodyRoot,
				},
				Signature: h2Sig,
			},
		}
	}
	return proposerSlashings, nil
}

func convertInternalProposerSlashings(src []*eth.ProposerSlashing) ([]*ProposerSlashing, error) {
	if src == nil {
		return nil, errors.New("proposer slashings are emtpy, nothing to convert.")
	}
	proposerSlashings := make([]*ProposerSlashing, len(src))
	for i, s := range src {
		proposerSlashings[i] = &ProposerSlashing{
			SignedHeader1: &SignedBeaconBlockHeader{
				Message: &BeaconBlockHeader{
					Slot:          fmt.Sprintf("%d", s.Header_1.Header.Slot),
					ProposerIndex: fmt.Sprintf("%d", s.Header_1.Header.ProposerIndex),
					ParentRoot:    hexutil.Encode(s.Header_1.Header.ParentRoot),
					StateRoot:     hexutil.Encode(s.Header_1.Header.StateRoot),
					BodyRoot:      hexutil.Encode(s.Header_1.Header.BodyRoot),
				},
				Signature: hexutil.Encode(s.Header_1.Signature),
			},
			SignedHeader2: &SignedBeaconBlockHeader{
				Message: &BeaconBlockHeader{
					Slot:          fmt.Sprintf("%d", s.Header_2.Header.Slot),
					ProposerIndex: fmt.Sprintf("%d", s.Header_2.Header.ProposerIndex),
					ParentRoot:    hexutil.Encode(s.Header_2.Header.ParentRoot),
					StateRoot:     hexutil.Encode(s.Header_2.Header.StateRoot),
					BodyRoot:      hexutil.Encode(s.Header_2.Header.BodyRoot),
				},
				Signature: hexutil.Encode(s.Header_2.Signature),
			},
		}
	}
	return proposerSlashings, nil
}

func convertAttesterSlashings(src []*AttesterSlashing) ([]*eth.AttesterSlashing, error) {
	if src == nil {
		return nil, errors.New("nil b.Message.Body.AttesterSlashings")
	}
	attesterSlashings := make([]*eth.AttesterSlashing, len(src))
	for i, s := range src {
		a1Sig, err := hexutil.Decode(s.Attestation1.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Signature", i)
		}
		a1AttestingIndices := make([]uint64, len(s.Attestation1.AttestingIndices))
		for j, ix := range s.Attestation1.AttestingIndices {
			attestingIndex, err := strconv.ParseUint(ix, 10, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.AttestingIndices[%d]", i, j)
			}
			a1AttestingIndices[j] = attestingIndex
		}
		a1Slot, err := strconv.ParseUint(s.Attestation1.Data.Slot, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.Slot", i)
		}
		a1CommitteeIndex, err := strconv.ParseUint(s.Attestation1.Data.Index, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.Index", i)
		}
		a1BeaconBlockRoot, err := hexutil.Decode(s.Attestation1.Data.BeaconBlockRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.BeaconBlockRoot", i)
		}
		a1SourceEpoch, err := strconv.ParseUint(s.Attestation1.Data.Source.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.Source.Epoch", i)
		}
		a1SourceRoot, err := hexutil.Decode(s.Attestation1.Data.Source.Root)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.Source.Root", i)
		}
		a1TargetEpoch, err := strconv.ParseUint(s.Attestation1.Data.Target.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.Target.Epoch", i)
		}
		a1TargetRoot, err := hexutil.Decode(s.Attestation1.Data.Target.Root)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation1.Data.Target.Root", i)
		}
		a2Sig, err := hexutil.Decode(s.Attestation2.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Signature", i)
		}
		a2AttestingIndices := make([]uint64, len(s.Attestation2.AttestingIndices))
		for j, ix := range s.Attestation2.AttestingIndices {
			attestingIndex, err := strconv.ParseUint(ix, 10, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.AttestingIndices[%d]", i, j)
			}
			a2AttestingIndices[j] = attestingIndex
		}
		a2Slot, err := strconv.ParseUint(s.Attestation2.Data.Slot, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.Slot", i)
		}
		a2CommitteeIndex, err := strconv.ParseUint(s.Attestation2.Data.Index, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.Index", i)
		}
		a2BeaconBlockRoot, err := hexutil.Decode(s.Attestation2.Data.BeaconBlockRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.BeaconBlockRoot", i)
		}
		a2SourceEpoch, err := strconv.ParseUint(s.Attestation2.Data.Source.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.Source.Epoch", i)
		}
		a2SourceRoot, err := hexutil.Decode(s.Attestation2.Data.Source.Root)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.Source.Root", i)
		}
		a2TargetEpoch, err := strconv.ParseUint(s.Attestation2.Data.Target.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.Target.Epoch", i)
		}
		a2TargetRoot, err := hexutil.Decode(s.Attestation2.Data.Target.Root)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.AttesterSlashings[%d].Attestation2.Data.Target.Root", i)
		}
		attesterSlashings[i] = &eth.AttesterSlashing{
			Attestation_1: &eth.IndexedAttestation{
				AttestingIndices: a1AttestingIndices,
				Data: &eth.AttestationData{
					Slot:            primitives.Slot(a1Slot),
					CommitteeIndex:  primitives.CommitteeIndex(a1CommitteeIndex),
					BeaconBlockRoot: a1BeaconBlockRoot,
					Source: &eth.Checkpoint{
						Epoch: primitives.Epoch(a1SourceEpoch),
						Root:  a1SourceRoot,
					},
					Target: &eth.Checkpoint{
						Epoch: primitives.Epoch(a1TargetEpoch),
						Root:  a1TargetRoot,
					},
				},
				Signature: a1Sig,
			},
			Attestation_2: &eth.IndexedAttestation{
				AttestingIndices: a2AttestingIndices,
				Data: &eth.AttestationData{
					Slot:            primitives.Slot(a2Slot),
					CommitteeIndex:  primitives.CommitteeIndex(a2CommitteeIndex),
					BeaconBlockRoot: a2BeaconBlockRoot,
					Source: &eth.Checkpoint{
						Epoch: primitives.Epoch(a2SourceEpoch),
						Root:  a2SourceRoot,
					},
					Target: &eth.Checkpoint{
						Epoch: primitives.Epoch(a2TargetEpoch),
						Root:  a2TargetRoot,
					},
				},
				Signature: a2Sig,
			},
		}
	}
	return attesterSlashings, nil
}

func convertInternalAttesterSlashings(src []*eth.AttesterSlashing) ([]*AttesterSlashing, error) {
	if src == nil {
		return nil, errors.New("AttesterSlashings is empty, nothing to convert.")
	}
	attesterSlashings := make([]*AttesterSlashing, len(src))
	for i, s := range src {
		a1AttestingIndices := make([]string, len(s.Attestation_1.AttestingIndices))
		for j, ix := range s.Attestation_1.AttestingIndices {
			a1AttestingIndices[j] = fmt.Sprintf("%d", ix)
		}
		a2AttestingIndices := make([]string, len(s.Attestation_2.AttestingIndices))
		for j, ix := range s.Attestation_2.AttestingIndices {
			a2AttestingIndices[j] = fmt.Sprintf("%d", ix)
		}
		attesterSlashings[i] = &AttesterSlashing{
			Attestation1: &IndexedAttestation{
				AttestingIndices: a1AttestingIndices,
				Data: &AttestationData{
					Slot:            fmt.Sprintf("%d", s.Attestation_1.Data.Slot),
					Index:           fmt.Sprintf("%d", s.Attestation_1.Data.CommitteeIndex),
					BeaconBlockRoot: hexutil.Encode(s.Attestation_1.Data.BeaconBlockRoot),
					Source: &Checkpoint{
						Epoch: fmt.Sprintf("%d", s.Attestation_1.Data.Source.Epoch),
						Root:  hexutil.Encode(s.Attestation_1.Data.Source.Root),
					},
					Target: &Checkpoint{
						Epoch: fmt.Sprintf("%d", s.Attestation_1.Data.Target.Epoch),
						Root:  hexutil.Encode(s.Attestation_1.Data.Target.Root),
					},
				},
				Signature: hexutil.Encode(s.Attestation_1.Signature),
			},
			Attestation2: &IndexedAttestation{
				AttestingIndices: a2AttestingIndices,
				Data: &AttestationData{
					Slot:            fmt.Sprintf("%d", s.Attestation_2.Data.Slot),
					Index:           fmt.Sprintf("%d", s.Attestation_2.Data.CommitteeIndex),
					BeaconBlockRoot: hexutil.Encode(s.Attestation_2.Data.BeaconBlockRoot),
					Source: &Checkpoint{
						Epoch: fmt.Sprintf("%d", s.Attestation_2.Data.Source.Epoch),
						Root:  hexutil.Encode(s.Attestation_2.Data.Source.Root),
					},
					Target: &Checkpoint{
						Epoch: fmt.Sprintf("%d", s.Attestation_2.Data.Target.Epoch),
						Root:  hexutil.Encode(s.Attestation_2.Data.Target.Root),
					},
				},
				Signature: hexutil.Encode(s.Attestation_2.Signature),
			},
		}
	}
	return attesterSlashings, nil
}

func convertAtts(src []*Attestation) ([]*eth.Attestation, error) {
	if src == nil {
		return nil, errors.New("nil b.Message.Body.Attestations")
	}
	atts := make([]*eth.Attestation, len(src))
	for i, a := range src {
		sig, err := hexutil.Decode(a.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Signature", i)
		}
		slot, err := strconv.ParseUint(a.Data.Slot, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.Slot", i)
		}
		committeeIndex, err := strconv.ParseUint(a.Data.Index, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.Index", i)
		}
		beaconBlockRoot, err := hexutil.Decode(a.Data.BeaconBlockRoot)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.BeaconBlockRoot", i)
		}
		sourceEpoch, err := strconv.ParseUint(a.Data.Source.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.Source.Epoch", i)
		}
		sourceRoot, err := hexutil.Decode(a.Data.Source.Root)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.Source.Root", i)
		}
		targetEpoch, err := strconv.ParseUint(a.Data.Target.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.Target.Epoch", i)
		}
		targetRoot, err := hexutil.Decode(a.Data.Target.Root)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Attestations[%d].Data.Target.Root", i)
		}
		atts[i] = &eth.Attestation{
			AggregationBits: []byte(a.AggregationBits),
			Data: &eth.AttestationData{
				Slot:            primitives.Slot(slot),
				CommitteeIndex:  primitives.CommitteeIndex(committeeIndex),
				BeaconBlockRoot: beaconBlockRoot,
				Source: &eth.Checkpoint{
					Epoch: primitives.Epoch(sourceEpoch),
					Root:  sourceRoot,
				},
				Target: &eth.Checkpoint{
					Epoch: primitives.Epoch(targetEpoch),
					Root:  targetRoot,
				},
			},
			Signature: sig,
		}
	}
	return atts, nil
}

func convertInternalAtts(src []*eth.Attestation) ([]*Attestation, error) {
	if src == nil {
		return nil, errors.New("Attestations are empty, nothing to convert.")
	}
	atts := make([]*Attestation, len(src))
	for i, a := range src {
		atts[i] = &Attestation{
			AggregationBits: hexutil.Encode(a.AggregationBits),
			Data: &AttestationData{
				Slot:            fmt.Sprintf("%d", a.Data.Slot),
				Index:           fmt.Sprintf("%d", a.Data.CommitteeIndex),
				BeaconBlockRoot: hexutil.Encode(a.Data.BeaconBlockRoot),
				Source: &Checkpoint{
					Epoch: fmt.Sprintf("%d", a.Data.Source.Epoch),
					Root:  hexutil.Encode(a.Data.Source.Root),
				},
				Target: &Checkpoint{
					Epoch: fmt.Sprintf("%d", a.Data.Target.Epoch),
					Root:  hexutil.Encode(a.Data.Target.Root),
				},
			},
			Signature: hexutil.Encode(a.Signature),
		}
	}
	return atts, nil
}

func convertDeposits(src []*Deposit) ([]*eth.Deposit, error) {
	if src == nil {
		return nil, errors.New("nil b.Message.Body.Deposits")
	}
	deposits := make([]*eth.Deposit, len(src))
	for i, d := range src {
		proof := make([][]byte, len(d.Proof))
		for j, p := range d.Proof {
			var err error
			proof[j], err = hexutil.Decode(p)
			if err != nil {
				return nil, errors.Wrapf(err, "could not decode b.Message.Body.Deposits[%d].Proof[%d]", i, j)
			}
		}
		pubkey, err := hexutil.Decode(d.Data.Pubkey)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Deposits[%d].Pubkey", i)
		}
		withdrawalCreds, err := hexutil.Decode(d.Data.WithdrawalCredentials)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Deposits[%d].WithdrawalCredentials", i)
		}
		amount, err := strconv.ParseUint(d.Data.Amount, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Deposits[%d].Amount", i)
		}
		sig, err := hexutil.Decode(d.Data.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.Deposits[%d].Signature", i)
		}
		deposits[i] = &eth.Deposit{
			Proof: proof,
			Data: &eth.Deposit_Data{
				PublicKey:             pubkey,
				WithdrawalCredentials: withdrawalCreds,
				Amount:                amount,
				Signature:             sig,
			},
		}
	}
	return deposits, nil
}

func convertInternalDeposits(src []*eth.Deposit) ([]*Deposit, error) {
	if src == nil {
		return nil, errors.New("deposits are empty, nothing to convert.")
	}
	deposits := make([]*Deposit, len(src))
	for i, d := range src {
		proof := make([]string, len(d.Proof))
		for j, p := range d.Proof {
			proof[j] = hexutil.Encode(p)
		}
		deposits[i] = &Deposit{
			Proof: proof,
			Data: &DepositData{
				Pubkey:                hexutil.Encode(d.Data.PublicKey),
				WithdrawalCredentials: hexutil.Encode(d.Data.WithdrawalCredentials),
				Amount:                fmt.Sprintf("%d", d.Data.Amount),
				Signature:             hexutil.Encode(d.Data.Signature),
			},
		}
	}
	return deposits, nil
}

func convertExits(src []*SignedVoluntaryExit) ([]*eth.SignedVoluntaryExit, error) {
	if src == nil {
		return nil, errors.New("nil b.Message.Body.VoluntaryExits")
	}
	exits := make([]*eth.SignedVoluntaryExit, len(src))
	for i, e := range src {
		sig, err := hexutil.Decode(e.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.VoluntaryExits[%d].Signature", i)
		}
		epoch, err := strconv.ParseUint(e.Message.Epoch, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.VoluntaryExits[%d].Epoch", i)
		}
		validatorIndex, err := strconv.ParseUint(e.Message.ValidatorIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.VoluntaryExits[%d].ValidatorIndex", i)
		}
		exits[i] = &eth.SignedVoluntaryExit{
			Exit: &eth.VoluntaryExit{
				Epoch:          primitives.Epoch(epoch),
				ValidatorIndex: primitives.ValidatorIndex(validatorIndex),
			},
			Signature: sig,
		}
	}
	return exits, nil
}

func convertInternalExits(src []*eth.SignedVoluntaryExit) ([]*SignedVoluntaryExit, error) {
	if src == nil {
		return nil, errors.New("VoluntaryExits are empty, nothing to convert.")
	}
	exits := make([]*SignedVoluntaryExit, len(src))
	for i, e := range src {
		exits[i] = &SignedVoluntaryExit{
			Message: &VoluntaryExit{
				Epoch:          fmt.Sprintf("%d", e.Exit.Epoch),
				ValidatorIndex: fmt.Sprintf("%d", e.Exit.ValidatorIndex),
			},
			Signature: hexutil.Encode(e.Signature),
		}
	}
	return exits, nil
}

func convertBlsChanges(src []*SignedBlsToExecutionChange) ([]*eth.SignedBLSToExecutionChange, error) {
	if src == nil {
		return nil, errors.New("nil b.Message.Body.BlsToExecutionChanges")
	}
	changes := make([]*eth.SignedBLSToExecutionChange, len(src))
	for i, ch := range src {
		sig, err := hexutil.Decode(ch.Signature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.BlsToExecutionChanges[%d].Signature", i)
		}
		index, err := strconv.ParseUint(ch.Message.ValidatorIndex, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.BlsToExecutionChanges[%d].Message.ValidatorIndex", i)
		}
		pubkey, err := hexutil.Decode(ch.Message.FromBlsPubkey)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.BlsToExecutionChanges[%d].Message.FromBlsPubkey", i)
		}
		address, err := hexutil.Decode(ch.Message.ToExecutionAddress)
		if err != nil {
			return nil, errors.Wrapf(err, "could not decode b.Message.Body.BlsToExecutionChanges[%d].Message.ToExecutionAddress", i)
		}
		changes[i] = &eth.SignedBLSToExecutionChange{
			Message: &eth.BLSToExecutionChange{
				ValidatorIndex:     primitives.ValidatorIndex(index),
				FromBlsPubkey:      pubkey,
				ToExecutionAddress: address,
			},
			Signature: sig,
		}
	}
	return changes, nil
}

func convertInternalBlsChanges(src []*eth.SignedBLSToExecutionChange) ([]*SignedBlsToExecutionChange, error) {
	if src == nil {
		return nil, errors.New("BlsToExecutionChanges are emtpy, nothing to convert.")
	}
	changes := make([]*SignedBlsToExecutionChange, len(src))
	for i, ch := range src {
		changes[i] = &SignedBlsToExecutionChange{
			Message: &BlsToExecutionChange{
				ValidatorIndex:     fmt.Sprintf("%d", ch.Message.ValidatorIndex),
				FromBlsPubkey:      hexutil.Encode(ch.Message.FromBlsPubkey),
				ToExecutionAddress: hexutil.Encode(ch.Message.ToExecutionAddress),
			},
			Signature: hexutil.Encode(ch.Signature),
		}
	}
	return changes, nil
}

func uint256ToHex(num string) ([]byte, error) {
	uint256, ok := new(big.Int).SetString(num, 10)
	if !ok {
		return nil, errors.New("could not parse Uint256")
	}
	bigEndian := uint256.Bytes()
	if len(bigEndian) > 32 {
		return nil, errors.New("number too big for Uint256")
	}
	return bytesutil2.ReverseByteOrder(bytesutil2.PadTo(bigEndian, 32)), nil
}
