package models

import (
	"github.com/TheVovchenskiy/sportify-backend/pkg/common"

	"github.com/google/uuid"
)

type ProfileAPI struct {
	IsMyProfile  bool        `json:"isMyProfile"`
	IsTgOnlyUser bool        `json:"isTgOnlyUser"`
	Username     string      `json:"username"`
	FirstName    *string     `json:"first_name"`
	SecondName   *string     `json:"second_name"`
	PhotoURL     *string     `json:"photo_url"`
	Description  *string     `json:"description"`
	TgURL        *string     `json:"tg_url"`
	SportTypes   []SportType `json:"sport_types"`
}

func mapTgURL(tgUserID *int64, tgUsername string) *string {
	if tgUserID == nil {
		return nil
	}

	return common.Ref("https://t.me/" + tgUsername)
}

func MapUserFullToProfileAPI(userIDFromToken uuid.UUID, userFull *UserFull) *ProfileAPI {
	return &ProfileAPI{
		IsMyProfile:  userIDFromToken == userFull.ID,
		IsTgOnlyUser: userFull.TgID != nil,
		Username:     userFull.Username,
		FirstName:    userFull.FirstName,
		SecondName:   userFull.SecondName,
		PhotoURL:     userFull.PhotoURL,
		Description:  userFull.Description,
		TgURL:        mapTgURL(userFull.TgID, userFull.Username),
		SportTypes:   userFull.SportTypes,
	}
}
