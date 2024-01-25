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
package eip712

import (
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// WrapTxToTypedData wraps an Amino-encoded Cosmos Tx JSON SignDoc
// bytestream into an EIP712-compatible TypedData request.
func WrapTxToTypedData(
	chainID uint64,
	data []byte,
) (apitypes.TypedData, error) {
	messagePayload, err := createEIP712MessagePayload(data)
	message := messagePayload.message
	if err != nil {
		return apitypes.TypedData{}, err
	}

	types, err := createEIP712Types(messagePayload)
	if err != nil {
		return apitypes.TypedData{}, err
	}

	domain := createEIP712Domain(chainID)

	typedData := apitypes.TypedData{
		Types:       types,
		PrimaryType: txField,
		Domain:      domain,
		Message:     message,
	}

	return typedData, nil
}
