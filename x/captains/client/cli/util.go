package cli

import (
	"encoding/json"
	"fmt"

	"github.com/tabilabs/tabi/x/captains/types"
)

const (
	FlagOwner      = "owner"
	FlagReportType = "report_type"

	ReportTypeDigest = "digest"
	ReportTypeBatch  = "batch"
	ReportTypeEnd    = "end"
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

func parseReportType(reportType string) types.ReportType {
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
