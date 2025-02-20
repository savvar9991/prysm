package electra_test

import (
	"context"
	"testing"

	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/electra"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/blocks"
	enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
	"github.com/prysmaticlabs/prysm/v5/testing/util"
)

func TestProcessOperationsWithNilRequests(t *testing.T) {
	tests := []struct {
		name      string
		modifyBlk func(blockElectra *ethpb.SignedBeaconBlockElectra)
		errMsg    string
	}{
		{
			name: "Nil deposit request",
			modifyBlk: func(blk *ethpb.SignedBeaconBlockElectra) {
				blk.Block.Body.ExecutionRequests.Deposits = []*enginev1.DepositRequest{nil}
			},
			errMsg: "nil deposit request",
		},
		{
			name: "Nil withdrawal request",
			modifyBlk: func(blk *ethpb.SignedBeaconBlockElectra) {
				blk.Block.Body.ExecutionRequests.Withdrawals = []*enginev1.WithdrawalRequest{nil}
			},
			errMsg: "nil withdrawal request",
		},
		{
			name: "Nil consolidation request",
			modifyBlk: func(blk *ethpb.SignedBeaconBlockElectra) {
				blk.Block.Body.ExecutionRequests.Consolidations = []*enginev1.ConsolidationRequest{nil}
			},
			errMsg: "nil consolidation request",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			st, ks := util.DeterministicGenesisStateElectra(t, 128)
			blk, err := util.GenerateFullBlockElectra(st, ks, util.DefaultBlockGenConfig(), 1)
			require.NoError(t, err)

			tc.modifyBlk(blk)

			b, err := blocks.NewSignedBeaconBlock(blk)
			require.NoError(t, err)

			require.NoError(t, st.SetSlot(1))

			_, err = electra.ProcessOperations(context.Background(), st, b.Block())
			require.ErrorContains(t, tc.errMsg, err)
		})
	}
}
