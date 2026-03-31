package pdf

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
)

func GeneratePDF(ctx context.Context, path string) ([]byte, error) {
	// create context
	chrome_ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// capture pdf
	var buf []byte
	url := fmt.Sprintf("http://%s%s", config.Bind, path)
	log.Info().Str("url", url).Msg("Getting with headless chrome")
	if err := chromedp.Run(chrome_ctx, printToPDF(url, &buf)); err != nil {
		return nil, fmt.Errorf("print to pdf: %w", err)
	}

	return buf, nil
}

// print a specific pdf page.
func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
