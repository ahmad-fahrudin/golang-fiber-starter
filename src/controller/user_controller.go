package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/utils"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserController struct {
	UserService  service.UserService
	TokenService service.TokenService
}

func NewUserController(userService service.UserService, tokenService service.TokenService) *UserController {
	return &UserController{
		UserService:  userService,
		TokenService: tokenService,
	}
}

// @Tags         Users
// @Summary      Get all users with new pagination utility
// @Description  Example of using the new pagination utility with date filtering
// @Security BearerAuth
// @Produce      json
// @Param        page       query     int     false   "Page number"  default(1)
// @Param        limit      query     int     false   "Maximum number of users"    default(10)
// @Param        search     query     string  false  "Search by name or email or role"
// @Param        start_date query     string  false  "Filter by start date (YYYY-MM-DD)"
// @Param        end_date   query     string  false  "Filter by end date (YYYY-MM-DD)"
// @Router       /users/paginated [get]
// @Success      200  {object}  example.GetAllUserResponse
func (u *UserController) GetUsersWithPagination(c *fiber.Ctx) error {
	// Extract pagination parameters using utility
	paginationParams := utils.ExtractPaginationParams(c)

	result, err := u.UserService.GetUsersWithPagination(c, paginationParams)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithPaginate[model.User]{
			Code:         fiber.StatusOK,
			Status:       "success",
			Message:      "Get all users successfully",
			Results:      result.Results,
			Page:         result.Page,
			Limit:        result.Limit,
			TotalPages:   result.TotalPages,
			TotalResults: result.TotalResults,
		})
}

// @Tags         Users
// @Summary      Get a user
// @Description  Logged in users can fetch only their own user information. Only admins can fetch other users.
// @Security BearerAuth
// @Produce      json
// @Param        id  path  string  true  "User id"
// @Router       /users/{id} [get]
// @Success      200  {object}  example.GetUserResponse
func (u *UserController) GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("userId")

	if _, err := uuid.Parse(userID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := u.UserService.GetUserByID(c, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithUser{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Get user successfully",
			User:    *user,
		})
}

// @Tags         Users
// @Summary      Create a user
// @Description  Only admins can create other users.
// @Security BearerAuth
// @Produce      json
// @Param        request  body  validation.CreateUser  true  "Request body"
// @Router       /users [post]
// @Success      201  {object}  example.CreateUserResponse
func (u *UserController) CreateUser(c *fiber.Ctx) error {
	req := new(validation.CreateUser)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := u.UserService.CreateUser(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.SuccessWithUser{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Create user successfully",
			User:    *user,
		})
}

// @Tags         Users
// @Summary      Update a user
// @Description  Logged in users can only update their own information. Only admins can update other users.
// @Security BearerAuth
// @Produce      json
// @Param        id  path  string  true  "User id"
// @Param        request  body  validation.UpdateUser  true  "Request body"
// @Router       /users/{id} [patch]
// @Success      200  {object}  example.UpdateUserResponse
func (u *UserController) UpdateUser(c *fiber.Ctx) error {
	req := new(validation.UpdateUser)
	userID := c.Params("userId")

	if _, err := uuid.Parse(userID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := u.UserService.UpdateUser(c, req, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithUser{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Update user successfully",
			User:    *user,
		})
}

// @Tags         Users
// @Summary      Delete a user
// @Description  Logged in users can delete only themselves. Only admins can delete other users.
// @Security BearerAuth
// @Produce      json
// @Param        id  path  string  true  "User id"
// @Router       /users/{id} [delete]
// @Success      200  {object}  example.DeleteUserResponse
func (u *UserController) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("userId")

	if _, err := uuid.Parse(userID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := u.TokenService.DeleteAllToken(c, userID); err != nil {
		return err
	}

	if err := u.UserService.DeleteUser(c, userID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Delete user successfully",
		})
}
