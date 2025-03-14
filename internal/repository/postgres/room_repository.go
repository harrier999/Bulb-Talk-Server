package postgres

import (
	"context"
	"server/internal/models/orm"
	"server/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresRoomRepository struct {
	db *gorm.DB
}

func NewPostgresRoomRepository(db *gorm.DB) repository.RoomRepository {
	return &PostgresRoomRepository{
		db: db,
	}
}

func (r *PostgresRoomRepository) Create(ctx context.Context, room orm.Room) (orm.Room, error) {
	result := r.db.WithContext(ctx).Create(&room)
	return room, result.Error
}

func (r *PostgresRoomRepository) FindByID(ctx context.Context, id uuid.UUID) (orm.Room, error) {
	var room orm.Room
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&room)
	return room, result.Error
}

func (r *PostgresRoomRepository) GetUserRooms(ctx context.Context, userID uuid.UUID) ([]orm.Room, error) {
	var rooms []orm.Room
	result := r.db.WithContext(ctx).Table("rooms").
		Select("rooms.*").
		Joins("left join room_users on rooms.id = room_users.room_id").
		Where("room_users.user_id = ?", userID).
		Find(&rooms)
	return rooms, result.Error
}

func (r *PostgresRoomRepository) AddUserToRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	roomUser := orm.RoomUser{
		RoomID: roomID,
		UserID: userID,
	}
	result := r.db.WithContext(ctx).Create(&roomUser)
	return result.Error
}

func (r *PostgresRoomRepository) RemoveUserFromRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("room_id = ? AND user_id = ?", roomID, userID).Delete(&orm.RoomUser{})
	return result.Error
}

func (r *PostgresRoomRepository) CreateRoomWithUsers(ctx context.Context, roomName string, userIDs []uuid.UUID) (uuid.UUID, error) {

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return uuid.Nil, tx.Error
	}

	room := orm.Room{
		RoomName: roomName,
	}
	if err := tx.Create(&room).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	roomUsers := make([]orm.RoomUser, len(userIDs))
	for i, userID := range userIDs {
		roomUsers[i] = orm.RoomUser{
			RoomID: room.ID,
			UserID: userID,
		}
	}
	if err := tx.Create(&roomUsers).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return uuid.Nil, err
	}

	return room.ID, nil
}
