package customErrors

import (
	"Backend-trainee-assignment-winter-2025/internal/models"
)

var (
	ErrBindFailed      = models.ErrorResponse{Errors: "data bind failed, check fields"}
	ErrValidate        = models.ErrorResponse{Errors: "cannot validate struct"}
	ErrCreateUser      = models.ErrorResponse{Errors: "failed to create user"}
	ErrUnauthorized    = models.ErrorResponse{Errors: "user unauthorized"}
	ErrInvalidPassword = models.ErrorResponse{Errors: "invalid password"}
	ErrInternal        = models.ErrorResponse{Errors: "internal server error"}
	ErrMerchNotFound   = models.ErrorResponse{Errors: "merch item not found"}
	ErrUserNotFound    = models.ErrorResponse{Errors: "user not found"}
	ErrSelfTransfer    = models.ErrorResponse{Errors: "self-transfer is not allowed"}
)
