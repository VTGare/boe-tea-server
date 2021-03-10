package jsonutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

//DecodeJSON decodes fiber.Ctx request body to dst.
func DecodeJSON(ctx *fiber.Ctx, dst interface{}) error {
	if ctx.Get("Content-Type") != "application/json" {
		return fiber.NewError(http.StatusUnsupportedMediaType, "Content-Type header is not application/json")
	}

	dec := json.NewDecoder(bytes.NewReader(ctx.Body()))
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var (
			syntaxError    *json.SyntaxError
			unmarshalError *json.UnmarshalTypeError
		)

		switch {
		case errors.As(err, &syntaxError) || errors.Is(err, io.ErrUnexpectedEOF):
			return fiber.NewError(http.StatusBadRequest, "Request body contains badly-formatted JSON.")

		case errors.As(err, &unmarshalError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalError.Field, unmarshalError.Offset)
			return fiber.NewError(http.StatusBadRequest, msg)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)

			return fiber.NewError(http.StatusBadRequest, msg)

		case errors.Is(err, io.EOF):
			return fiber.NewError(http.StatusBadRequest, "Request body is empty")

		case err.Error() == "http: request body too large":
			return fiber.NewError(http.StatusRequestEntityTooLarge, "Request body shouldn't be larger than 1MB.")

		default:
			return err
		}
	}

	return nil
}
