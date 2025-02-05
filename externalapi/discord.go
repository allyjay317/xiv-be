package externalapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/alyjay/xiv-be/types"
)

func RequestAuthToken(result *types.AuthTokenResponse, code string) (err error) {
	var client_id, _ = os.LookupEnv("DISCORD_CLIENT_ID")
	var client_secret, _ = os.LookupEnv("DISCORD_CLIENT_SECRET")
	var redirect_uri, _ = os.LookupEnv("DISCORD_REDIRECT_URI")
	pUrl := "https://discord.com/api/oauth2/token"
	b := types.AuthTokenRequest{
		ClientId:     client_id,
		ClientSecret: client_secret,
		Code:         code,
		GrantType:    "client_credentials",
		RedirectURI:  redirect_uri,
		Scope:        "identify",
	}
	payloadBuf := new(bytes.Buffer)

	data := url.Values{}
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirect_uri)

	json.NewEncoder(payloadBuf).Encode(b)

	resp, err := http.Post(pUrl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(body), &result)

	return err
}

func GetUserInfo(userResult *types.DiscordUserResponse, auth types.AuthTokenResponse) (err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	req.Header.Set("Authorization", auth.TokenType+" "+auth.AccessToken)
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	_ = json.NewDecoder(res.Body).Decode(&userResult)

	return err
}
