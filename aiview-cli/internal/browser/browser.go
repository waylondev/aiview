// Package browser provides browser automation for cookie extraction.
package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

// GetCookies opens a browser and extracts cookies after user login.
func GetCookies(url string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var cookies string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(30*time.Second), // Wait for user to login
		chromedp.Evaluate(`document.cookie`, &cookies),
	)
	if err != nil {
		return "", fmt.Errorf("browser automation failed: %w", err)
	}

	return cookies, nil
}
