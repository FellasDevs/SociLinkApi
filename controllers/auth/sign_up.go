package authcontroller

import (
	"SociLinkApi/dto"
	"SociLinkApi/models"
	userrepository "SociLinkApi/repository/user"
	authservice "SociLinkApi/services/auth"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

func SignUpController(context *gin.Context, db *gorm.DB) {
	var userInfo dto.SignUpRequestDto

	if err := context.ShouldBindJSON(&userInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var fieldErrors []string

	if len(userInfo.Name) < 6 {
		fieldErrors = append(fieldErrors, "Nome deve conter no mínimo 6 caracteres.")
	} else if len(userInfo.Name) > 50 {
		fieldErrors = append(fieldErrors, "Nome deve conter no máximo 50 caracteres.")
	}
	if len(userInfo.Nickname) < 6 {
		fieldErrors = append(fieldErrors, "Nickname deve conter no mínimo 6 caracteres.")
	} else if len(userInfo.Nickname) > 50 {
		fieldErrors = append(fieldErrors, "Nickname deve conter no máximo 50 caracteres.")
	}
	if _, err := mail.ParseAddress(userInfo.Email); err != nil {
		fieldErrors = append(fieldErrors, "Email inválido.")
	} else if len(userInfo.Email) > 50 {
		fieldErrors = append(fieldErrors, "Email deve conter no máximo 50 caracteres.")
	}
	if len(userInfo.Password) < 6 {
		fieldErrors = append(fieldErrors, "Senha deve conter no mínimo 6 caracteres.")
	} else if len(userInfo.Password) > 50 {
		fieldErrors = append(fieldErrors, "Senha deve conter no máximo 50 caracteres.")
	}
	if userInfo.Birthdate == "" {
		fieldErrors = append(fieldErrors, "Data de nascimento não informada.")
	}

	if len(fieldErrors) > 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": strings.Join(fieldErrors, " "),
		})
		return
	}

	birthdate, err := authservice.ParseBirthdate(userInfo.Birthdate)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	password, err := authservice.EncryptPassword(userInfo.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	user := models.User{
		Name:      userInfo.Name,
		Nickname:  userInfo.Nickname,
		Email:     userInfo.Email,
		Password:  password,
		Birthdate: birthdate,
	}

	err = userrepository.CreateUser(&user, db)
	if err != nil {
		var pgErr *pgconn.PgError

		if !errors.As(err, &pgErr) {
			context.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		reason := ""

		switch pgErr.ConstraintName {
		case "users_email_key":
			reason = "email"
		case "users_nickname_key":
			reason = "nickname"
		}

		context.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": err.Error(),
			"data": gin.H{
				"reason": reason,
			},
		})
		return
	}

	token, err := authservice.CreateAuthToken(user.ID, time.Hour*24)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	response := dto.SignUpResponseDto{
		User:      dto.UserToResponseDto(user),
		AuthToken: token,
	}

	context.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Usuário criado com sucesso!",
		"data":    response,
	})
}
