package postgres

import (
	"context"
	"errors"
	"server/internal/models/orm"
	"server/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresFriendRepository struct {
	db *gorm.DB
}

func NewPostgresFriendRepository(db *gorm.DB) repository.FriendRepository {
	return &PostgresFriendRepository{
		db: db,
	}
}

func (r *PostgresFriendRepository) GetFriendList(ctx context.Context, userID uuid.UUID) ([]orm.Friend, error) {
	var friends []orm.Friend
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&friends)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []orm.Friend{}, nil
		}
		return nil, result.Error
	}
	return friends, nil
}

func (r *PostgresFriendRepository) AddFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	friend := orm.Friend{
		UserID:   userID,
		FriendID: friendID,
	}
	result := r.db.WithContext(ctx).Create(&friend)
	return result.Error
}

func (r *PostgresFriendRepository) BlockFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&orm.Friend{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Update("is_blocked", true)
	return result.Error
}

func (r *PostgresFriendRepository) UnblockFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&orm.Friend{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Update("is_blocked", false)
	return result.Error
}

func (r *PostgresFriendRepository) CheckIfFriendExists(ctx context.Context, userID, friendID uuid.UUID) (bool, error) {
	var friend orm.Friend
	result := r.db.WithContext(ctx).Where("user_id = ? AND friend_id = ?", userID, friendID).First(&friend)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
