package handler

import (
	"multitenant/model"

	"github.com/gofiber/fiber/v2"
)

func (s *MultiTanentServer) CreateTransaction(ctx *fiber.Ctx) error {
	s.logger.Debug("CreateTransaction")
	crt := model.CrtTransaction{}
	if err := ctx.BodyParser(&crt); err != nil {
		return NewUnprocessableEntityResponse(ctx, err.Error())
	}

	if err := s.validator.Struct(crt); err != nil {
		return NewBadRequestResponse(ctx, err.Error())
	}

	t := crt.ToTransaction()
	s.engine.CreateTransaction(t)

	return NewSuccessNoContentResponse(ctx)
}
