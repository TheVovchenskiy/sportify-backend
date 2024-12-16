package models

import (
	"fmt"
	"github.com/TheVovchenskiy/sportify-backend/pkg/common"
	"strings"

	"github.com/google/uuid"
)

type ProfileAPI struct {
	IsMyProfile  bool        `json:"isMyProfile"`
	IsTgOnlyUser bool        `json:"isTgOnlyUser"`
	UserID       string      `json:"user_id"`
	Username     string      `json:"username"`
	FirstName    *string     `json:"first_name"`
	SecondName   *string     `json:"second_name"`
	PhotoURL     *string     `json:"photo_url"`
	Description  *string     `json:"description"`
	TgURL        *string     `json:"tg_url"`
	SportTypes   []SportType `json:"sport_types"`
}

func MapTgURL(tgUserID *int64, tgUsername string) *string {
	if tgUserID == nil {
		return nil
	}

	return common.Ref("https://t.me/" + tgUsername)
}

func MapUserFullToProfileAPI(urlPrefixFile string, userIDFromToken uuid.UUID, userFull *UserFull) *ProfileAPI {
	if userFull.PhotoURL == nil {
		userFull.PhotoURL = common.Ref(urlPrefixFile + "user-default-avatar.png")
	}
	return &ProfileAPI{
		IsMyProfile:  userIDFromToken == userFull.ID,
		IsTgOnlyUser: userFull.TgID != nil,
		UserID:       userFull.ID.String(),
		Username:     userFull.Username,
		FirstName:    userFull.FirstName,
		SecondName:   userFull.SecondName,
		PhotoURL:     userFull.PhotoURL,
		Description:  userFull.Description,
		TgURL:        MapTgURL(userFull.TgID, userFull.Username),
		SportTypes:   userFull.SportTypes,
	}
}

type RequestUpdateProfile struct {
	FirstName   string      `json:"first_name"`
	SecondName  string      `json:"second_name"`
	PhotoURL    string      `json:"photo_url"`
	Description string      `json:"description"`
	SportTypes  []SportType `json:"sport_types"`
}

func (r *RequestUpdateProfile) Valid() error {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.SecondName = strings.TrimSpace(r.SecondName)
	r.Description = strings.TrimSpace(r.Description)

	if len(r.FirstName) > 256 {
		return fmt.Errorf("имя должно быть короче 256 символов")
	}

	if len(r.SecondName) > 256 {
		return fmt.Errorf("фамилия должна быть короче 256 символов")
	}

	if len(r.Description) > 4096 {
		return fmt.Errorf("описание должно быть короче 4096 символов")
	}

	return nil
}

type ResponseUpdateProfile struct {
	RedirectURL string `json:"redirect_url"`
}

func NewResponseUpdateProfile(url string, userID uuid.UUID) ResponseUpdateProfile {
	return ResponseUpdateProfile{RedirectURL: url + "/profiles/" + userID.String()}
}
