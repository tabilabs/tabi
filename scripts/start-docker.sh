#!/bin/bash

KEY="dev0"
CHAINID="tabi_9789-1"
MONIKER="mymoniker"
DATA_DIR=$(mktemp -d -t tabi-datadir.XXXXX)

echo "create and add new keys"
./tabid keys add $KEY --home $DATA_DIR --no-backup --chain-id $CHAINID --algo "eth_secp256k1" --keyring-backend test
echo "init Tabi with moniker=$MONIKER and chain-id=$CHAINID"
./tabid init $MONIKER --chain-id $CHAINID --home $DATA_DIR
echo "prepare genesis: Allocate genesis accounts"
./tabid add-genesis-account \
"$(./tabid keys show $KEY -a --home $DATA_DIR --keyring-backend test)" 1000000000000000000atabi,1000000000000000000stake \
--home $DATA_DIR --keyring-backend test
echo "prepare genesis: Sign genesis transaction"
./tabid gentx $KEY 1000000000000000000stake --keyring-backend test --home $DATA_DIR --keyring-backend test --chain-id $CHAINID
echo "prepare genesis: Collect genesis tx"
./tabid collect-gentxs --home $DATA_DIR
echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./tabid validate-genesis --home $DATA_DIR

echo "starting tabi node $i in background ..."
./tabid start --pruning=nothing --rpc.unsafe \
--keyring-backend test --home $DATA_DIR \
>$DATA_DIR/node.log 2>&1 & disown

echo "started tabi node"
tail -f /dev/null