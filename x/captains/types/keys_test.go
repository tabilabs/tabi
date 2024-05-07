package types_test

//func TestKv(t *testing.T) {
//	mem := dbadapter.Store{DB: dbm.NewMemDB()}
//	meter := storetypes.NewGasMeter(10000)
//	gasStore := gaskv.NewStore(mem, meter, storetypes.KVGasConfig())
//	owner, err := sdk.AccAddressFromBech32("tabis1pytx0t95gfh0t3h89vt8nqadqq4hmwadkn7fk4")
//	require.NoError(t, err)
//
//	ownerStore := GetPrefixStoreKey(gasStore, owner)
//	ownerStore.Set([]byte("35971be6e9bb024a895582fe0e42e04848a86da550aaef0fccbfba86f99f617d"), types.PlaceHolder)
//
//	// for each
//	ownerStore2 := GetPrefixStoreKey(gasStore, owner)
//	ownerStoreIterator := ownerStore2.Iterator(nil, nil)
//	for ; ownerStoreIterator.Valid(); ownerStoreIterator.Next() {
//		fmt.Printf("key: %s\n", string(ownerStoreIterator.Key()))
//	}
//}
//
//func GetPrefixStoreKey(store storetypes.KVStore, owner sdk.AccAddress) prefix.Store {
//	key := types.NodeByOwnerStoreKey(owner)
//	return prefix.NewStore(store, key)
//}
