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

/*
Package ante defines the SDK auth module's AnteHandler as well as an internal
AnteHandler for an Ethereum transaction (i.e MsgEthereumTx).

During CheckTx, the transaction is passed through a series of
pre-message execution validation checks such as signature and account
verification in addition to minimum fees being checked. Otherwise, during
DeliverTx, the transaction is simply passed to the EVM which will also
perform the same series of checks. The distinction is made in CheckTx to
prevent spam and DoS attacks.
*/
package ante
