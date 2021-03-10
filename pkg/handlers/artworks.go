package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/VTGare/boe-tea-backend/internal/jsonutils"
	"github.com/VTGare/boe-tea-backend/pkg/artworks"
	"github.com/VTGare/boe-tea-backend/pkg/artworks/options"
	"github.com/gofiber/fiber/v2"
)

func registerArtworks(r fiber.Router, env *Env) {
	r.Get("/artworks", allArtworks(env))
	r.Get("/artworks/:id", findOneArtwork(env))
	r.Post("/artworks", insertOneArtwork(env))   //TODO: protected route
	r.Delete("/artworks", deleteOneArtwork(env)) //TODO: protected route
}

func allArtworks(env *Env) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var filt options.Filter
		err := ctx.QueryParser(&filt)
		if err != nil {
			return err
		}

		if last := ctx.Query("last"); last != "" {
			dur, err := time.ParseDuration(last)
			if err != nil {
				return err
			}

			filt.Time = dur
		}

		opts := options.DefaultFind()
		opts.Filter = &filt

		artworks, err := env.Artworks.FindMany(ctx.Context(), opts)
		if err != nil {
			return err
		}

		return ctx.JSON(artworks)
	}
}

func findOneArtwork(env *Env) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		id, err := strconv.Atoi(ctx.Params("id"))
		if err != nil {
			return fiber.NewError(http.StatusBadRequest, "Parsing error")
		}

		artwork, err := env.Artworks.FindOne(ctx.Context(), &options.FilterOne{ID: id})
		if err != nil {
			return err
		}

		return ctx.JSON(artwork)
	}
}

func insertOneArtwork(env *Env) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var req artworks.ArtworkInsert
		err := jsonutils.DecodeJSON(ctx, &req)
		if err != nil {
			return err
		}

		artwork, err := env.Artworks.InsertOne(ctx.Context(), &req)
		if err != nil {
			return err
		}

		return ctx.JSON(artwork)
	}
}

func deleteOneArtwork(env *Env) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var req options.FilterOne
		err := ctx.QueryParser(&req)
		if err != nil {
			return err
		}

		artwork, err := env.Artworks.DeleteOne(ctx.Context(), &req)
		if err != nil {
			return err
		}

		return ctx.JSON(artwork)
	}
}
