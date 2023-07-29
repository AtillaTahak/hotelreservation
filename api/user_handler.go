package api

import (
	"errors"
	"fmt"
	"hotelreservation/db"
	"hotelreservation/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}
func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error{
	var (
		update types.UpdateUserParams
		userID = c.Params("id")
	)
	if err := c.BodyParser(&update); err != nil {
		return err
	}
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	if err = h.userStore.UpdateUser(c.Context(), filter, update); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"updated": userID})
}
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	err := h.userStore.DeleteUser(c.Context(),id)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"deleted": id})
}
func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	user, err := types.NewUserFromParams(&params)

	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	user, err := h.userStore.GetUserByID(c.Context(),id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}
		return err
	}
	return c.JSON(user)
}