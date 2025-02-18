package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	database "github.com/alyjay/xiv-be/database"
	"github.com/alyjay/xiv-be/externalapi"
	"github.com/alyjay/xiv-be/response"
	types "github.com/alyjay/xiv-be/types"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	user, err := database.GetUser(id)
	if err != nil {
		response.NotFoundError(w, "No User Found")
		return
	}

	characters, err := database.GetCharacters(id)

	if err == nil {
		user.Characters = append(user.Characters, characters...)
	}

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
	site_url, _ := os.LookupEnv("SITE_URL")
	code := r.URL.Query().Get("code")
	var auth types.AuthTokenResponse
	var userResult types.DiscordUserResponse
	err := externalapi.RequestAuthToken(&auth, code)
	if err != nil {
		http.Redirect(w, r, site_url+"/error", http.StatusInternalServerError)
		return
	}

	err = externalapi.GetUserInfo(&userResult, auth)

	if err != nil {
		http.Redirect(w, r, site_url+"/error", http.StatusInternalServerError)
		return
	}

	user := types.User{
		ID:          uuid.NewString(),
		Username:    userResult.Username,
		DiscordId:   userResult.Id,
		Avatar:      userResult.Avatar,
		AccentColor: GetAccentColor(&userResult),
		AuthToken:   auth.AccessToken,
		Expires:     time.Now().UnixMilli() + int64(auth.ExpiresIn),
	}
	id, err := database.InsertUser(user)
	if err != nil {
		http.Redirect(w, r, site_url+"/error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, site_url+"/login?id="+id+"&token="+auth.AccessToken+"&expires="+fmt.Sprint(auth.ExpiresIn), http.StatusSeeOther)
}
