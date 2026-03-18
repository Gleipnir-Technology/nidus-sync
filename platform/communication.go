package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/publicreport"
)

func NotificationCount(ctx context.Context, org *models.Organization, user *models.User) (result uint, err error) {
	count_reports, err := publicreport.ReportsForOrganizationCount(ctx, org.ID)
	if err != nil {
		return 0, fmt.Errorf("report query: %w", err)
	}
	return uint(count_reports), nil
}
