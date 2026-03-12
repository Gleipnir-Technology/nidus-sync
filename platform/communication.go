package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/publicreport"
)

func NotificationCount(ctx context.Context, org *models.Organization, user *models.User) (result uint, err error) {
	count_nreports, err := publicreport.NuisanceReportForOrganizationCount(ctx, org.ID)
	if err != nil {
		return 0, fmt.Errorf("nuisance report query: %w", err)
	}
	result += count_nreports

	count_wreports, err := publicreport.WaterReportForOrganizationCount(ctx, org.ID)
	if err != nil {
		return 0, fmt.Errorf("water report query: %w", err)
	}
	result += count_wreports
	return result, nil
}
