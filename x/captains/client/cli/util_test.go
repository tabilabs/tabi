package cli

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

func TestParseReport(t *testing.T) {
	for _, tc := range []struct {
		name       string
		prepare    func() ([]byte, error)
		contents   []byte
		reportType string
		expectErr  bool
	}{
		{
			name: "digest",
			prepare: func() ([]byte, error) {
				digest := types.ReportDigest{
					EpochId:                  0,
					TotalBatchCount:          0,
					TotalNodeCount:           0,
					MaximumNodeCountPerBatch: 0,
					GlobalOnOperationRatio:   sdk.Dec{},
				}
				return json.Marshal(digest)
			},
			reportType: ReportTypeDigest,
		},
		{
			name:       "batch",
			reportType: ReportTypeBatch,
			prepare: func() ([]byte, error) {
				batch := types.ReportBatch{
					EpochId: 1,
					BatchId: 0,
					Nodes:   nil,
				}
				return json.Marshal(batch)
			},
		},
		{
			name:       "end",
			reportType: ReportTypeEnd,
			prepare: func() ([]byte, error) {
				end := types.ReportEnd{
					EpochId: 1,
				}
				return json.Marshal(end)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bz, err := tc.prepare()
			if err != nil {
				t.Error(err)
			}
			tc.contents = bz
			_, err = parseReport(tc.contents, tc.reportType)
			if tc.expectErr {
				t.Error(err)
			}
		})
	}
}
