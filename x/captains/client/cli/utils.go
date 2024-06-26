package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tabilabs/tabi/x/captains/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	FlagOwner      = "owner"
	FlagReportType = "report-type"

	ReportTypeDigest = "digest"
	ReportTypeBatch  = "batch"
	ReportTypeEnd    = "end"

	draftReportDigestFileName = "draft_report_digest.json"
	draftReportBatchFileName  = "draft_report_batch.json"
	draftReportEndFileName    = "draft_report_end.json"
)

func parseReport(contents []byte, reportType string) (any, error) {
	if len(contents) == 0 {
		return nil, fmt.Errorf("report is empty")
	}

	switch reportType {
	case ReportTypeDigest:
		var report types.ReportDigest
		err := json.Unmarshal(contents, &report)
		if err != nil {
			return nil, err
		}
		return &report, nil
	case ReportTypeBatch:
		var report types.ReportBatch
		err := json.Unmarshal(contents, &report)
		if err != nil {
			return nil, err
		}
		return &report, nil
	case ReportTypeEnd:
		var report types.ReportEnd
		err := json.Unmarshal(contents, &report)
		if err != nil {
			return nil, err
		}
		return &report, nil
	}

	return nil, fmt.Errorf("report type %s is not supported\n", reportType)
}

// parseReportType parses the report type string and returns the corresponding ReportType
func parseReportType(reportType string) types.ReportType {
	reportType = strings.TrimSpace(reportType)
	switch reportType {
	case ReportTypeDigest:
		return types.ReportType_REPORT_TYPE_DIGEST
	case ReportTypeBatch:
		return types.ReportType_REPORT_TYPE_BATCH
	case ReportTypeEnd:
		return types.ReportType_REPORT_TYPE_END
	}
	return types.ReportType_REPORT_TYPE_UNSPECIFIED
}

// writeFile writes the input to the file
func writeFile(fileName string, input any) error {
	raw, err := json.MarshalIndent(input, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal proposal: %w", err)
	}

	if err := os.WriteFile(fileName, raw, 0o600); err != nil {
		return err
	}

	return nil
}

func draftReport(reportType string) error {
	switch parseReportType(reportType) {
	case types.ReportType_REPORT_TYPE_DIGEST:
		report := types.ReportDigest{
			EpochId:                  1,
			TotalBatchCount:          1,
			TotalNodeCount:           1,
			MaximumNodeCountPerBatch: 1,
			GlobalOnOperationRatio:   sdk.Dec{},
		}
		return writeFile(draftReportDigestFileName, report)
	case types.ReportType_REPORT_TYPE_BATCH:
		report := types.ReportBatch{
			EpochId:   1,
			BatchId:   1,
			NodeCount: 1,
			Nodes: []types.NodePowerOnRatio{
				{
					NodeId:           "node-id",
					OnOperationRatio: sdk.Dec{},
				},
			},
		}
		return writeFile(draftReportBatchFileName, report)
	case types.ReportType_REPORT_TYPE_END:
		report := types.ReportEnd{
			EpochId: 1,
		}
		return writeFile(draftReportEndFileName, report)
	}
	return fmt.Errorf("report type %s is not supported\n", reportType)
}
