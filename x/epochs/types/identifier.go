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

package types

import (
	"fmt"
	"strings"
)

const (
	// WeekEpochID defines the identifier for weekly epochs
	WeekEpochID = "week"
	// DayEpochID defines the identifier for daily epochs
	DayEpochID = "day"
	// HourEpochID defines the identifier for hourly epochs
	HourEpochID = "hour"
)

// ValidateEpochIdentifierInterface performs a stateless
// validation of the epoch ID interface.
func ValidateEpochIdentifierInterface(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return ValidateEpochIdentifierString(v)
}

// ValidateEpochIdentifierInterface performs a stateless
// validation of the epoch ID.
func ValidateEpochIdentifierString(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("blank epoch identifier: %s", s)
	}
	return nil
}
