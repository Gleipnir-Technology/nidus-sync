package mailer

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/lob"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/Gleipnir-Technology/nidus-sync/platform/pdf"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func ComplianceSend(ctx context.Context, row_id int32) error {
	bxn := db.PGInstance.BobDB
	compliance_req, err := models.FindComplianceReportRequest(ctx, bxn, row_id)
	if err != nil {
		return fmt.Errorf("find compliance report: %w", err)
	}
	log.Debug().Int32("id", row_id).Str("public_id", compliance_req.PublicID).Msg("working on mailer")

	if compliance_req.LeadID.IsNull() {
		return fmt.Errorf("no lead for compliance req %d", compliance_req.ID)
	}
	lead_id := compliance_req.LeadID.MustGet()
	lead, err := models.FindLead(ctx, bxn, lead_id)
	if err != nil {
		return fmt.Errorf("find lead: %w", err)
	}

	if lead.SiteID.IsNull() {
		return fmt.Errorf("no site for lead %d", lead.ID)
	}
	site_id := lead.SiteID.MustGet()
	site, err := models.FindSite(ctx, bxn, site_id)
	if err != nil {
		return fmt.Errorf("find site: %w", err)
	}

	address, err := models.FindAddress(ctx, bxn, site.AddressID)
	if err != nil {
		return fmt.Errorf("find address: %w", err)
	}

	organization, err := models.FindOrganization(ctx, bxn, site.OrganizationID)
	if err != nil {
		return fmt.Errorf("find address: %w", err)
	}
	if organization.LobAddressID.IsNull() {
		return fmt.Errorf("organization %d has no Lob Address ID", organization.ID)
	}

	path := fmt.Sprintf("/mailer/mode-3/%s/preview", compliance_req.PublicID)
	content, err := pdf.GeneratePDF(ctx, path)
	if err != nil {
		return fmt.Errorf("generate pdf: %w", err)
	}
	err = file.MailerFromReader(compliance_req.PublicID, bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("save pdf: %w", err)
	}

	// Do the part where we actually send to the mailer service
	if organization.LobAddressID.IsNull() {
		return fmt.Errorf("lob address for %d is null", organization.ID)
	}
	lob_address := organization.LobAddressID.MustGet()
	letter, err := sendMail(ctx, lob_address, compliance_req.PublicID, site, address, content)
	if err != nil {
		return fmt.Errorf("send mail: %w", err)
	}

	mailer_uuid, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("generate uuid: %w", err)
	}
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start txn: %w", err)
	}
	defer txn.Rollback(nil)
	mailer, err := models.CommsMailers.Insert(&models.CommsMailerSetter{
		AddressID:  omit.From(address.ID),
		Created:    omit.From(time.Now()),
		ExternalID: omit.From(letter.ID),
		// ID
		Recipient: omit.From(site.OwnerName),
		UUID:      omit.From(mailer_uuid),
	}).One(ctx, txn)
	if err != nil {
		return fmt.Errorf("create comms mailer: %w", err)
	}

	crrm, err := models.ComplianceReportRequestMailers.Insert(&models.ComplianceReportRequestMailerSetter{
		ComplianceReportRequestID: omit.From(compliance_req.ID),
		// ID
		MailerID: omit.From(mailer.ID),
	}).One(ctx, txn)
	if err != nil {
		return fmt.Errorf("create crrm: %w", err)
	}
	log.Info().Int32("id", crrm.ID).Msg("Created compliance report request mailer")
	txn.Commit(ctx)
	return nil
}

func sendMail(ctx context.Context, org_address_id string, public_id string, site *models.Site, address *models.Address, content []byte) (*lob.Letter, error) {
	key := config.LobAPIKey
	client := lob.NewLob(key)
	line1 := address.Number + " " + address.Street
	addr_req := lob.RequestAddressCreate{
		AddressLine1: line1,
		AddressCity:  address.Locality,
		AddressState: address.Region,
		AddressZip:   address.PostalCode,
		Name:         site.OwnerName,
	}
	addr_to, err := client.AddressCreate(ctx, addr_req)
	if err != nil {
		return nil, fmt.Errorf("create to addr: %w", err)
	}

	req := lob.RequestLetterCreate{
		To:      addr_to.ID,
		From:    org_address_id,
		File:    bytes.NewReader(content),
		Color:   true,
		UseType: "operational",
	}
	letter, err := client.LetterCreate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("letter create: %w", err)
	}

	return &letter, nil
}
