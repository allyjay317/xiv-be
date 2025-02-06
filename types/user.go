package types

type User struct {
	ID          string      `json:"id,omitempty" db:"id"`
	Username    string      `json:"username" db:"username"`
	DiscordId   string      `json:"discord_id" db:"discord_id"`
	Avatar      string      `json:"avatar" db:"avatar"`
	AccentColor string      `json:"accent_color" db:"accent_color"`
	Characters  []Character `json:"characters"`
	AuthToken   string      `json:"auth_token" db:"auth_token,omitempty"`
	Expires     int64       `json:"expires" db:"expires,omitempty"`
}

type AuthTokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type DiscordUserResponse struct {
	Id                   string `json:"id"`
	Username             string `json:"username"`
	Avatar               string `json:"avatar"`
	Discriminator        string `json:"discriminator"`
	PublicFlags          int    `json:"public_flags"`
	Flags                int    `json:"flags"`
	Banner               string `json:"banner"`
	AccentColor          int    `json:"accent_color"`
	GlobalName           string `json:"global_name"`
	AvatarDecorationData map[string]struct {
		SkuId string `json:"sku_id"`
		Asset string `json:"asset"`
	} `json:"avatar_decoration_data"`
	BannerColor  string              `json:"banner_color"`
	Clan         map[string]struct{} `json:"clan"`
	PrimaryGuild map[string]struct{} `json:"primary_guild"`
	MfaEnabled   bool                `json:"mfa_enabled"`
	Locale       string              `json:"locale"`
	PremiumType  int                 `json:"premium_type"`
}
