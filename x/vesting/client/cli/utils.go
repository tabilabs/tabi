// Copyright 2024 Tabi Foundation
// This file is part of the Tabi Network packages.
//
// Tabi is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Tabi packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

type VestingData struct {
	StartTime int64         `json:"start_time"`
	Periods   []InputPeriod `json:"periods"`
}

type InputPeriod struct {
	Coins  string `json:"coins"`
	Length int64  `json:"length_seconds"`
}

// readScheduleFile reads the file at path and unmarshals it to get the schedule.
// Returns start time, periods, and error.
func ReadScheduleFile(path string) (int64, sdkvesting.Periods, error) {
	contents, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return 0, nil, err
	}

	var data VestingData

	if err = json.Unmarshal(contents, &data); err != nil {
		return 0, nil, err
	}

	startTime := data.StartTime
	periods := make(sdkvesting.Periods, 0, len(data.Periods))

	for i, p := range data.Periods {
		if p.Length < 1 {
			return 0, nil, fmt.Errorf("invalid period length of %d in period %d, length must be greater than 0", p.Length, i)
		}

		amount, err := sdk.ParseCoinsNormalized(p.Coins)
		if err != nil {
			return 0, nil, err
		}

		period := sdkvesting.Period{Length: p.Length, Amount: amount}
		periods = append(periods, period)
	}

	return startTime, periods, nil
}
