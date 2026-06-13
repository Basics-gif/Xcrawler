package browser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/playwright-community/playwright-go"
)

type Manager struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	Context playwright.BrowserContext
	Page    playwright.Page
}

func New(cookiePath string) (*Manager, error) {
	data, err := os.ReadFile(cookiePath)
	if err != nil {
		return nil, fmt.Errorf("could not read session file: %v", err)
	}

	var rawCookies []RawCookie
	if err = json.Unmarshal(data, &rawCookies); err != nil {
		return nil, fmt.Errorf("could not parse cookes files: %v", err)
	}

	cookies := toPlaywrightCookies(rawCookies)

	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not start playwright: %v", err)
	}

	br, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		pw.Stop()
		return nil, fmt.Errorf("could not launch browser: %v", err)
	}

	ctx, err := br.NewContext()
	if err != nil {
		br.Close()
		pw.Stop()
		return nil, fmt.Errorf("could not create context: %v", err)
	}

	if err = ctx.AddCookies(cookies); err != nil {
		ctx.Close()
		br.Close()
		pw.Stop()
		return nil, fmt.Errorf("could not add cookies: %v", err)
	}

	page, err := ctx.NewPage()
	if err != nil {
		ctx.Close()
		br.Close()
		pw.Stop()
		return nil, fmt.Errorf("could not create page: %v", err)
	}

	return &Manager{
		pw:      pw,
		browser: br,
		Context: ctx,
		Page:    page,
	}, nil

}

func (m *Manager) Close() {
	m.Context.Close()
	m.browser.Close()
	m.pw.Stop()
}
