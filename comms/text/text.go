package text

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/config"
)

func SendText(ctx context.Context, source string, destination string, message string) (string, error) {
	switch config.TextProvider {
	case "voipms":
		return sendTextVoipms(ctx, destination, message)
	case "twilio":
		return sendTextTwilio(ctx, source, destination, message)
	}
	return "", fmt.Errorf("Unsupported provider '%s'", config.TextProvider)
}
