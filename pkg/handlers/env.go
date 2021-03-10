package handlers

import (
	"github.com/VTGare/boe-tea-backend/pkg/artworks"
	"github.com/VTGare/boe-tea-backend/pkg/guilds"
	"github.com/gofiber/fiber/v2"
)

type Env struct {
	Artworks artworks.Service
	Guilds   guilds.Service
}

func Register(r fiber.Router, env *Env) {
	registerArtworks(r, env)
	registerGuilds(r, env)
}
