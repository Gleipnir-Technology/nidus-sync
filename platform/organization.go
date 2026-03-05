package platform

import (
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type Organization struct {
	ID   int32
	Name string
}

func NewOrganization(org *models.Organization) Organization {
	return Organization{
		ID:   org.ID,
		Name: org.Name,
	}
}
