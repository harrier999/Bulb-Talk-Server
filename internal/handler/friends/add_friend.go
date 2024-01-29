package friends

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/internal/db/postgres_db"
	"server/internal/models/orm"
	"server/pkg/authenticator"
	"server/pkg/log"

	"github.com/google/uuid"
	"github.com/ttacon/libphonenumber"
)

type AddFriendRequest struct {
	PhoneNumber string `json:"phone_number"`
}

func AddFriend(w http.ResponseWriter, r *http.Request) {
	logger := log.NewColorLog()

	user_id, err := authenticator.GetUserID(r)
	if err != nil {
		logger.Info("Error while getting user_id from context", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var AddFriendData AddFriendRequest
	json.NewDecoder(r.Body).Decode(&AddFriendData)
	if err := validAddFriendData(AddFriendData); err != nil {
		logger.Info("Error while validating add friend data", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	friend_id, err := getUserIDByPhoneNumber(AddFriendData.PhoneNumber)
	if err != nil {
		logger.Info("Error while getting friend_id by phone_number", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if checkIfFriendExists(user_id, friend_id) {
		logger.Info("Friend already exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := AddFriendToDB(user_id, friend_id); err != nil {
		logger.Info("Error while adding friend to DB", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func validAddFriendData(addFriendData AddFriendRequest) error {
	if addFriendData.PhoneNumber == "" {
		return errors.New("phone_number is empty")
	}
	_, err := libphonenumber.Parse(addFriendData.PhoneNumber, "KR")
	if err != nil {
		return errors.New("phone_number is invalid")
	}

	return nil
}

func AddFriendToDB(user_id uuid.UUID, friend_id uuid.UUID) error {
	logger := log.NewColorLog()
	postgresClient := postgres_db.GetPostgresClient()
	user := orm.Friend{
		UserID:   user_id,
		FriendID: friend_id,
	}
	result := postgresClient.Create(&user)
	if result.Error != nil {
		logger.Info("Error while creating friend", "error", result.Error)
		return result.Error
	}

	return nil
}

func getUserIDByPhoneNumber(phoneNumber string) (uuid.UUID, error) {
	postgresClient := postgres_db.GetPostgresClient()
	user := orm.User{
		PhoneNumber: phoneNumber,
	}
	result := postgresClient.Where(&user).First(&user)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return uuid.UUID{}, errors.New("user not found")
		}
		return uuid.UUID{}, result.Error
	}
	return user.UserID, nil
}

func checkIfFriendExists(user_id uuid.UUID, friend_id uuid.UUID) bool {
	postgresClient := postgres_db.GetPostgresClient()
	friend := orm.Friend{
		UserID:   user_id,
		FriendID: friend_id,
	}
	result := postgresClient.Where(&friend).First(&friend)
	if result.Error != nil {
		return false
	}
	return true
}
