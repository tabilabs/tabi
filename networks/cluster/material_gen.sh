### build dockerfile

rm -rf ./chainData/tabi
../../build/tabid testnet init-files --v 8 -o ./chainData/tabi --keyring-backend=test --starting-ip-address 192.167.10.2 --chain-id tabi_9789-1

  for i in {0..7}; do
    APP_TOML="./chainData/tabi/node$i/tabid/config/app.toml"
    GENESIS="./chainData/tabi/node$i/tabid/config/genesis.json"
    TMP_GENESIS="./chainData/tabi/node$i/tabid/config/tmp_genesis.json"
    DATADIR="./chainData/tabi/node$i/tabid/data"

  	# Set gas limit in genesis
	  jq '.consensus_params["block"]["max_gas"]="30000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' 's/127.0.0.1/0.0.0.0/g' "$APP_TOML"
    else
        sed -i 's/127.0.0.1/0.0.0.0/g' "$APP_TOML"
    fi
    rm -rf $DATADIR
    mkdir -p $DATADIR
    cp ./priv_validator_state.json $DATADIR
  done