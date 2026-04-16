package mailer

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/Gleipnir-Technology/nidus-sync/platform/pdf"
	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	lob "github.com/lob/lob-go"
	"github.com/rs/zerolog/log"
)

func ComplianceSend(ctx context.Context, txn bob.Executor, row_id int32) error {
	compliance_req, err := models.FindComplianceReportRequest(ctx, txn, row_id)
	if err != nil {
		return fmt.Errorf("find compliance report: %w", err)
	}
	log.Debug().Int32("id", row_id).Str("public_id", compliance_req.PublicID).Msg("working on mailer")

	if compliance_req.LeadID.IsNull() {
		return fmt.Errorf("no lead for compliance req %d", compliance_req.ID)
	}
	lead_id := compliance_req.LeadID.MustGet()
	lead, err := models.FindLead(ctx, txn, lead_id)
	if err != nil {
		return fmt.Errorf("find lead: %w", err)
	}

	if lead.SiteID.IsNull() {
		return fmt.Errorf("no site for lead %d", lead.ID)
	}
	site_id := lead.SiteID.MustGet()
	site, err := models.FindSite(ctx, txn, site_id)
	if err != nil {
		return fmt.Errorf("find site: %w", err)
	}

	address, err := models.FindAddress(ctx, txn, site.AddressID)
	if err != nil {
		return fmt.Errorf("find address: %w", err)
	}

	organization, err := models.FindOrganization(ctx, txn, site.OrganizationID)
	if err != nil {
		return fmt.Errorf("find address: %w", err)
	}
	if organization.LobAddressID.IsNull() {
		return fmt.Errorf("organization %d has no Lob Address ID", organization.ID)
	}

	content, err := pdf.GeneratePDF(ctx, compliance_req.PublicID)
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
	letter, err := sendMail(ctx, lob_address, compliance_req.PublicID, site, address)
	if err != nil {
		return fmt.Errorf("send mail: %w", err)
	}

	mailer_uuid, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("generate uuid: %w", err)
	}
	mailer, err := models.CommsMailers.Insert(&models.CommsMailerSetter{
		AddressID:  omit.From(address.ID),
		Created:    omit.From(time.Now()),
		ExternalID: omit.From(letter.Id),
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
	return nil
}

func sendMail(ctx context.Context, org_address_id string, public_id string, site *models.Site, address *models.Address) (*lob.Letter, error) {
	ctx_lob := context.WithValue(ctx, lob.ContextBasicAuth, lob.BasicAuth{UserName: config.LobAPIKey})
	config := lob.NewConfiguration()
	client := lob.NewAPIClient(config)

	from_addr, _, err := client.AddressesApi.Get(ctx_lob, org_address_id).Execute()
	if err != nil {
		return nil, fmt.Errorf("get from address '%s': %w", org_address_id)
	}

	var to_addr = *lob.NewAddressEditable()
	line1 := address.Number + " " + address.Street
	to_addr.SetAddressLine1(line1)
	to_addr.SetAddressLine2("")
	to_addr.SetAddressCity(address.Locality)
	to_addr.SetAddressState(address.Region)
	to_addr.SetAddressZip(address.PostalCode)
	to_addr.SetAddressCountry(lob.COUNTRYEXTENDED_US)
	addr_desc := fmt.Sprintf("site %d - %s", site.ID, site.OwnerName)
	to_addr.SetDescription(addr_desc)
	to_addr.SetName(site.OwnerName)

	var use_type lob.LtrUseType = lob.LTRUSETYPE_OPERATIONAL

	content_path := file.MailerPath(public_id)
	var ltr = lob.NewLetterEditable(true, to_addr, from_addr, content_path, *lob.NewNullableLtrUseType(&use_type))
	desc := fmt.Sprintf("Compliance request %s", public_id)
	ltr.SetDescription(desc)
	result, _, err := client.LettersApi.Create(ctx_lob).LetterEditable(*ltr).Execute()

	if err != nil {
		return nil, fmt.Errorf("create letter: %w", err)
	}
	log.Info().Str("id", result.Id).Msg("Created Lob letter")
	return result, nil
}
