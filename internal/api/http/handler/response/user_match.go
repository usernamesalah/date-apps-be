package response

import "date-apps-be/internal/model"

type UserMatchResponse struct {
	QuotaLeft int     `json:"quota_left"`
	Users     []*User `json:"users"`
}

type User struct {
	UserUID string `json:"user_uid"`
	Name    string `json:"name"`
}

func NewUserMatchResponse(users []*model.User, quotaLeft int) UserMatchResponse {
	var userMatchResponse []*User
	for _, user := range users {
		userMatchResponse = append(userMatchResponse, &User{
			UserUID: user.UID,
			Name:    user.Name,
		})
	}

	return UserMatchResponse{
		QuotaLeft: quotaLeft,
		Users:     userMatchResponse,
	}
}
