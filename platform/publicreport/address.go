package publicreport

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	//"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
)

func loadAddresses(ctx context.Context, txn bob.Executor, address_ids []int32) (results map[int32]types.Address, err error) {
	rows, err := models.Addresses.Query(
		models.SelectWhere.Addresses.ID.In(address_ids...),
	).All(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("query addresses: %w", err)
	}
	results = make(map[int32]types.Address, len(rows))
	for _, row := range rows {
		results[row.ID] = types.AddressFromModel(row)
	}
	return results, nil
}
