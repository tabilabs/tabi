package keeper_test

import (
	"math"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

// Reporter submits the report as per epoch phase.
type CaptainsReporter struct {
	GlobalOnOperationRatio   sdk.Dec
	MaximumNodeCountPerBatch uint64
	TotalBatchCount          uint64
	TotalNodeCount           uint64
}

func NewCaptainsReporter(gor sdk.Dec, mnc uint64) *CaptainsReporter {
	return &CaptainsReporter{
		GlobalOnOperationRatio:   gor,
		MaximumNodeCountPerBatch: mnc,
		TotalBatchCount:          0,
		TotalNodeCount:           0,
	}
}

// SubmitDigest submits the digest report for the epoch.
func (cpr *CaptainsReporter) SubmitDigest(suite *IntegrationTestSuite, es *EpochState) {
	nc := uint64(len(es.Nodes))
	cpr.TotalNodeCount = nc
	cpr.TotalBatchCount = uint64(math.Ceil(float64(nc) / float64(cpr.MaximumNodeCountPerBatch)))

	digest := types.ReportDigest{
		EpochId:                  es.Epoch,
		TotalBatchCount:          cpr.TotalBatchCount,
		TotalNodeCount:           cpr.TotalNodeCount,
		MaximumNodeCountPerBatch: cpr.MaximumNodeCountPerBatch,
		GlobalOnOperationRatio:   cpr.GlobalOnOperationRatio,
	}
	anyVal, err := cdctypes.NewAnyWithValue(&digest)
	suite.Require().NoError(err)
	suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		Report:     anyVal,
		ReportType: types.ReportType_REPORT_TYPE_DIGEST,
	})
	suite.Require().NoError(err)
}

// SubmitBatches submits the batches for the epoch.
// NOTE: this func causes height increment by TotalBatchCount-1.
func (cpr *CaptainsReporter) SubmitBatches(suite *IntegrationTestSuite, es *EpochState) {
	suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		ReportType: types.ReportType_REPORT_TYPE_BATCH,
	})

	for i := 1; i < int(cpr.TotalBatchCount); i++ {
		batch := types.ReportBatch{
			EpochId:   es.Epoch,
			BatchId:   uint64(i),
			NodeCount: cpr.MaximumNodeCountPerBatch,
			Nodes:     es.Nodes.PowerOnRatios((i-1)*int(cpr.MaximumNodeCountPerBatch), i*int(cpr.MaximumNodeCountPerBatch)),
		}
		anyVal, err := cdctypes.NewAnyWithValue(&batch)
		suite.Require().NoError(err)

		_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
			Authority:  accounts[0].String(),
			Report:     anyVal,
			ReportType: types.ReportType_REPORT_TYPE_BATCH,
		})
		suite.Require().NoError(err)
		suite.Commit()
	}

	//  the last batch
	batch := types.ReportBatch{
		EpochId:   es.Epoch,
		BatchId:   cpr.TotalBatchCount,
		NodeCount: cpr.TotalNodeCount - (cpr.TotalBatchCount-1)*cpr.MaximumNodeCountPerBatch,
		Nodes:     es.Nodes.PowerOnRatios((int(cpr.TotalBatchCount)-1)*int(cpr.MaximumNodeCountPerBatch), int(cpr.TotalNodeCount)),
	}
	anyVal, err := cdctypes.NewAnyWithValue(&batch)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		Report:     anyVal,
		ReportType: types.ReportType_REPORT_TYPE_BATCH,
	})
	suite.Require().NoError(err)
}

// SubmitEnd submits the end report for the epoch.
func (cpr *CaptainsReporter) SubmitEnd(suite *IntegrationTestSuite, es *EpochState) {
	end1 := types.ReportEnd{
		EpochId: es.Epoch,
	}
	anyVal, err := cdctypes.NewAnyWithValue(&end1)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		ReportType: types.ReportType_REPORT_TYPE_END,
		Report:     anyVal,
	})
	suite.Require().NoError(err)
}
