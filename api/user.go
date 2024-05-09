package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/izsal/simple_bank/db/sqlc"
	"github.com/izsal/simple_bank/util"
	"github.com/jackc/pgx/v5/pgconn"
)

// The `createUserRequest` struct is defining the structure of the request body expected when creating
// a new user. It contains fields such as `Username`, `Password`, `FullName`, and `Email`, each tagged
// with validation rules using the `binding` tag from the `gin` framework. These validation rules
// specify that certain fields are required (`binding:"required"`) and have specific constraints like
// being alphanumeric, having a minimum length, or being in email format. This struct helps in parsing
// and validating the incoming JSON request data before processing it further in the `createUser`
// function.
type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// The `createUserResponse` struct is defining the structure of the response that will be sent back to
// the client after a new user is successfully created. It contains fields such as `Username`,
// `FullName`, `Email`, `PasswordChangedAt`, and `CreatedAt` which represent the user's information.
// This struct is used to format the response data before sending it back as a JSON response in the
// `createUser` function.
type createUserResponse struct {
	Username          string
	FullName          string
	Email             string
	PasswordChangedAt time.Time
	CreatedAt         time.Time
}

// The `func (server *Server) createUser(ctx *gin.Context)` function is a method defined on the
// `Server` struct in the `api` package. This method is responsible for handling the creation of a new
// user based on the JSON request received in the `ctx` context.
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		// pqErr := err.(*pgconn.PgError)
		// log.Println(pqErr.Code)
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Code {
			case "23505":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt.Time,
		CreatedAt:         user.CreatedAt.Time,
	}
	ctx.JSON(http.StatusOK, rsp)
}
