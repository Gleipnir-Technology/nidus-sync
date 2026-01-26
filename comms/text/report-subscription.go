package text

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/nyaruka/phonenumbers"
	//"github.com/rs/zerolog/log"
)

func NewJobReportSubscriptionConfirmation(
	destination E164,
	report_id string,
	source E164) jobReportSubscription {
	return jobReportSubscription{
		dst:      destination,
		reportID: report_id,
		src:      source,
	}
}

type jobReportSubscription struct {
	dst      E164
	reportID string
	src      E164
}

func (j jobReportSubscription) content() string {
	return fmt.Sprintf("Thanks for submitting mosquito report %s. Text for any questions. We'll send you updates as we get them.", j.reportID)
}
func (j jobReportSubscription) destination() string {
	return phonenumbers.Format(&j.dst, phonenumbers.E164)
}
func (j jobReportSubscription) messageType() MessageType {
	return ReportSubscription
}
func (j jobReportSubscription) messageTypeName() string {
	return "report-subscription"
}
func (j jobReportSubscription) source() string {
	return phonenumbers.Format(&j.src, phonenumbers.E164)
}

func sendReportSubscription(ctx context.Context, job Job) error {
	j, ok := job.(jobReportSubscription)
	if !ok {
		return fmt.Errorf("job is not for report subscription confirmation")
	}

	sub, err := isSubscribed(ctx, job.destination())
	if err != nil {
		return fmt.Errorf("Failed to check if subscribed: %w", err)
	}
	if !sub {
		err = sendText(ctx, j.source(), j.destination(), j.content(), enums.CommsTextoriginWebsiteAction)
		if err != nil {
			return fmt.Errorf("Failed to send report subscription confirmation: %w", err)
		}
	} else {
		err = delayMessage(ctx, j.source(), j.destination(), j.content(), enums.CommsTextjobtypeReportConfirmation)
		if err != nil {
			return fmt.Errorf("Failed to delay report subscription message: %w", err)
		}
		err := ensureInitialText(ctx, j.source(), j.destination())
		if err != nil {
			return fmt.Errorf("Failed to ensure initial text has been sent: %w", err)
		}
	}
	return nil
}
