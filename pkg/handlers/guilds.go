package handlers

import (
	"github.com/VTGare/boe-tea-backend/internal/jsonutils"
	"github.com/VTGare/boe-tea-backend/pkg/guilds"
	"github.com/gofiber/fiber/v2"
)

func registerGuilds(r fiber.Router, env *Env) {
	r.Get("/guilds/:id", findOneGuild(env))
	r.Get("/guilds", allGuilds(env))          //TODO: protected route
	r.Post("/guilds/:id", insertGuild(env))   //TODO: protected route
	r.Delete("/guilds/:id", deleteGuild(env)) //TODO: protected route
	r.Put("/guilds", replaceGuild(env))       //TODO: protected route
}

func allGuilds(env *Env) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		guilds, err := env.Guilds.All(ctx.Context())
		if err != nil {
			panic(err)
		}

		return ctx.JSON(guilds)
	}
}

func insertGuild(env *Env) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")

		g, err := env.Guilds.InsertOne(ctx.Context(), id)
		if err != nil {
			return err
		}

		return ctx.JSON(g)
	}
}

func findOneGuild(env *Env) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")

		g, err := env.Guilds.FindOne(ctx.Context(), id)
		if err != nil {
			return err
		}

		return ctx.JSON(g)
	}
}

func deleteGuild(env *Env) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")

		g, err := env.Guilds.DeleteOne(ctx.Context(), id)
		if err != nil {
			return err
		}

		return ctx.JSON(g)
	}
}

func replaceGuild(env *Env) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		g := &guilds.Guild{}
		err := jsonutils.DecodeJSON(ctx, g)
		if err != nil {
			return err
		}

		g, err = env.Guilds.ReplaceOne(ctx.Context(), g)
		if err != nil {
			return err
		}

		return ctx.JSON(g)
	}
}
