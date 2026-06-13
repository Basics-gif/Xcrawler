package browser

import "github.com/playwright-community/playwright-go"

type RawCookie struct {
	Domain         string   `json:"domain"`
	ExpirationDate *float64 `json:"expirationDate"`
	HostOnly       bool     `json:"hostOnly"`
	HttpOnly       bool     `json:"httpOnly"`
	Name           string   `json:"name"`
	Path           string   `json:"path"`
	SameSite       string   `json:"sameSite"`
	Secure         bool     `json:"secure"`
	Session        bool     `json:"session"`
	StoreId        string   `json:"storeId"`
	Value          string   `json:"value"`
}

func toPlaywrightCookies(raw []RawCookie) []playwright.OptionalCookie {
	var cookies []playwright.OptionalCookie
	for _, c := range raw {
		cookie := playwright.OptionalCookie{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   playwright.String(c.Domain),
			Path:     playwright.String(c.Path),
			HttpOnly: playwright.Bool(c.HttpOnly),
			Secure:   playwright.Bool(c.Secure),
		}

		if c.ExpirationDate != nil {
			cookie.Expires = playwright.Float(*c.ExpirationDate)
		}

		switch c.SameSite {
		case "strict":
			cookie.SameSite = playwright.SameSiteAttributeStrict
		case "lax":
			cookie.SameSite = playwright.SameSiteAttributeLax
		default:
			cookie.SameSite = playwright.SameSiteAttributeNone
		}
		cookies = append(cookies, cookie)
	}

	return cookies

}
