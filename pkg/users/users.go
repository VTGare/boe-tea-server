package users

import (
	"context"

	"github.com/VTGare/boe-tea-backend/internal/database"
	"go.uber.org/zap"
)

type User struct {
}

type Service interface {
	InsertOne(context.Context, string) (*User, error)
	FindOne(context.Context, string) (*User, error)
	DeleteOne(context.Context, string) (*User, error)
	//ReplaceOne(context.Context, string) (*User, error)
	//AddFavourite(context.Context, int, bool) error
	//RemoveFavourite(context.Context, int, bool) error
}

type userService struct {
	db     *database.DB
	logger *zap.SugaredLogger
}
