package postgres

import (
	"context"
	"server/internal/models/orm"
	"server/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) repository.UserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (orm.User, error) {
	var user orm.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	return user, result.Error
}

func (r *PostgresUserRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (orm.User, error) {
	var user orm.User
	result := r.db.WithContext(ctx).Where("phone_number = ?", phoneNumber).First(&user)
	return user, result.Error
}

func (r *PostgresUserRepository) Create(ctx context.Context, user orm.User) (orm.User, error) {
	result := r.db.WithContext(ctx).Create(&user)
	return user, result.Error
}

func (r *PostgresUserRepository) Update(ctx context.Context, user orm.User) error {
	result := r.db.WithContext(ctx).Save(&user)
	return result.Error
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&orm.User{}, id)
	return result.Error
}
