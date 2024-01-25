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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	tabitypes "github.com/tabi/tabi/v1/types"
)

// NewRevenue returns an instance of Revenue. If the provided withdrawer
// address is empty, it sets the value to an empty string.
func NewRevenue(contract common.Address, deployer, withdrawer sdk.AccAddress) Revenue {
	withdrawerAddr := ""
	if len(withdrawer) > 0 {
		withdrawerAddr = withdrawer.String()
	}

	return Revenue{
		ContractAddress:   contract.String(),
		DeployerAddress:   deployer.String(),
		WithdrawerAddress: withdrawerAddr,
	}
}

// GetContractAddr returns the contract address
func (fs Revenue) GetContractAddr() common.Address {
	return common.HexToAddress(fs.ContractAddress)
}

// GetDeployerAddr returns the contract deployer address
func (fs Revenue) GetDeployerAddr() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(fs.DeployerAddress)
}

// GetWithdrawerAddr returns the account address to where the funds proceeding
// from the fees will be received. If the withdraw address is not defined, it
// defaults to the deployer address.
func (fs Revenue) GetWithdrawerAddr() sdk.AccAddress {
	if fs.WithdrawerAddress == "" {
		return nil
	}

	return sdk.MustAccAddressFromBech32(fs.WithdrawerAddress)
}

// Validate performs a stateless validation of a Revenue
func (fs Revenue) Validate() error {
	if err := tabitypes.ValidateNonZeroAddress(fs.ContractAddress); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(fs.DeployerAddress); err != nil {
		return err
	}

	if fs.WithdrawerAddress != "" {
		if _, err := sdk.AccAddressFromBech32(fs.WithdrawerAddress); err != nil {
			return err
		}
	}

	return nil
}
