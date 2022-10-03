package keeper

import (
	"github.com/NicholasDotSol/duality/x/dex/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetLimitOrderPoolTotalSharesMap set a specific limitOrderPoolTotalSharesMap in the store from its index
func (k Keeper) SetLimitOrderPoolTotalSharesMap(ctx sdk.Context, limitOrderPoolTotalSharesMap types.LimitOrderPoolTotalSharesMap) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LimitOrderPoolTotalSharesMapKeyPrefix))
	b := k.cdc.MustMarshal(&limitOrderPoolTotalSharesMap)
	store.Set(types.LimitOrderPoolTotalSharesMapKey(
		limitOrderPoolTotalSharesMap.Count,
	), b)
}

// GetLimitOrderPoolTotalSharesMap returns a limitOrderPoolTotalSharesMap from its index
func (k Keeper) GetLimitOrderPoolTotalSharesMap(
	ctx sdk.Context,
	count string,

) (val types.LimitOrderPoolTotalSharesMap, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LimitOrderPoolTotalSharesMapKeyPrefix))

	b := store.Get(types.LimitOrderPoolTotalSharesMapKey(
		count,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveLimitOrderPoolTotalSharesMap removes a limitOrderPoolTotalSharesMap from the store
func (k Keeper) RemoveLimitOrderPoolTotalSharesMap(
	ctx sdk.Context,
	count string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LimitOrderPoolTotalSharesMapKeyPrefix))
	store.Delete(types.LimitOrderPoolTotalSharesMapKey(
		count,
	))
}

// GetAllLimitOrderPoolTotalSharesMap returns all limitOrderPoolTotalSharesMap
func (k Keeper) GetAllLimitOrderPoolTotalSharesMap(ctx sdk.Context) (list []types.LimitOrderPoolTotalSharesMap) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LimitOrderPoolTotalSharesMapKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LimitOrderPoolTotalSharesMap
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
