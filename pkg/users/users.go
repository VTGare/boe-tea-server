package users

import (
	"context"
	"time"

	"github.com/VTGare/boe-tea-backend/internal/database"
	"go.uber.org/zap"
)

type User struct {
	ID         string       `json:"id" bson:"user_id"`
	DM         bool         `json:"dm" bson:"dm"`
	Crossport  bool         `json:"crossport" bson:"crossport"`
	Favourites []*Favourite `json:"favourites" bson:"new_favourites"`
	Groups     []*Group     `json:"groups" bson:"channel_groups"`
	CreatedAt  time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at" bson:"updated_at"`
}

type Favourite struct {
	ArtworkID string    `json:"artwork_id" bson:"artwork_id"`
	NSFW      bool      `json:"nsfw" bson:"nsfw"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Group struct {
	Name     string   `json:"name" bson:"name"`
	Parent   string   `json:"parent" bson:"parent"`
	Children []string `json:"children" bson:"children"`
}

type Service interface {
	InsertOne(context.Context, string) (*User, error)
	FindOne(context.Context, string) (*User, error)
	DeleteOne(context.Context, string) (*User, error)
	InsertFavourite(context.Context, *Favourite) (*Favourite, error)
	DeleteFavourite(context.Context, *Favourite) (*Favourite, error)
	InsertGroup(context.Context, *Group) (*Group, error)
	DeleteGroup(context.Context, string) (*Group, error)
	InsertToGroup(context.Context, string, string) (*Group, error)
	DeleteFromGroup(context.Context, string, string) (*Group, error)
}

type userService struct {
	db     *database.DB
	logger *zap.SugaredLogger
}
