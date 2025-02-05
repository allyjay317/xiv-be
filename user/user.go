package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	database "github.com/alyjay/xiv-be/database"
	types "github.com/alyjay/xiv-be/types"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	var user types.User
	var characters []types.Character
	id := r.URL.Query().Get("id")

	db, err := database.GetDb()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Server Issue" + err.Error()))
		return
	}
	err = db.Get(&user, `SELECT * FROM users WHERE id=$1`, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No User Found"))
		return
	}

	err = db.Select(&characters, `SELECT character_id, name, avatar, portrait FROM characters WHERE user_id=$1`, id)
	if err != nil {

	}

	user.Characters = append(user.Characters, characters...)

	json.NewEncoder(w).Encode(user)
}

func GetAccentColor(user *types.DiscordUserResponse) (color string) {
	if user.AccentColor != 0 {
		color = fmt.Sprintf("#%x", user.AccentColor)
	} else if user.BannerColor != "" {
		color = user.BannerColor
	}
	return color
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	client_id, _ := os.LookupEnv("DISCORD_CLIENT_ID")
	client_secret, _ := os.LookupEnv("DISCORD_CLIENT_SECRET")
	redirect_uri, _ := os.LookupEnv("DISCORD_REDIRECT_URI")
	site_url, _ := os.LookupEnv("SITE_URL")
	code := r.URL.Query().Get("code")
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
		http.Redirect(w, r, site_url+"/error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Redirect(w, r, site_url+"/error", http.StatusInternalServerError)
		return
	}

	var result types.AuthTokenResponse
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		http.Redirect(w, r, site_url+"/error", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	req.Header.Set("Authorization", result.TokenType+" "+result.AccessToken)
	res, err := client.Do(req)
	if err != nil {
		http.Redirect(w, r, site_url+"/error", http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	var userResult types.DiscordUserResponse
	_ = json.NewDecoder(res.Body).Decode(&userResult)

	newUUID := uuid.NewString()

	db, err := database.GetDb()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Server Issue" + err.Error()))
		return
	}

	var user types.User
	user.ID = newUUID
	err = db.Get(&user, `SELECT * FROM users WHERE discord_id=$1`, userResult.Id)
	if err != nil {
		user.Username = userResult.Username
		user.DiscordId = userResult.Id
		user.Avatar = userResult.Avatar
		user.AccentColor = GetAccentColor(&userResult)
		_, err = db.Exec(`INSERT INTO users (id, username, discord_id, avatar, accent_color, auth_token, expires) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			user.ID,
			user.Username,
			user.DiscordId,
			user.Avatar,
			user.AccentColor,
			result.AccessToken,
			time.Now().UnixMilli()+int64(result.ExpiresIn), 0)
	} else {
		user.AccentColor = GetAccentColor(&userResult)
		user.Avatar = userResult.Avatar
		_, err = db.Exec(`UPDATE users SET 
		accent_color = $1, 
		avatar = $2, 
		auth_token = $3, 
		expires = $4 
		WHERE discord_id=$5 `,
			GetAccentColor(&userResult),
			userResult.Avatar,
			result.AccessToken,
			time.Now().UnixMilli()+int64(result.ExpiresIn),
			userResult.Id)
	}
	http.Redirect(w, r, site_url+"/login?id="+user.ID+"&token="+result.AccessToken+"&expires="+fmt.Sprint(result.ExpiresIn), http.StatusSeeOther)
}
