package douyin

// Credential holds Douyin authentication data.
type Credential struct {
	SessionID string `json:"session_id"`
	TTwid     string `json:"tt_webid"`
	Cookie    string `json:"cookie"`
}

// IsValid checks if the credential has basic required fields.
func (c *Credential) IsValid() bool {
	return c.Cookie != ""
}