package publicreport

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
)

func loadAddresses(ctx context.Context, txn db.Tx, address_ids []int64) (results map[int32]types.Address, err error) {
	addresses, err := querypublic.AddressesFromIDs(ctx, txn, address_ids)
	if err != nil {
		return nil, fmt.Errorf("query addresses: %w", err)
	}
	results = make(map[int32]types.Address, len(addresses))
	for _, row := range addresses {
		results[row.ID] = types.AddressFromModel(row)
	}
	return results, nil
}
