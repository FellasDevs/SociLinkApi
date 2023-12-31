package usercontroller

import (
	"SociLinkApi/dto"
	"SociLinkApi/models"
	userrepository "SociLinkApi/repository/user"
	authservice "SociLinkApi/services/auth"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strings"
)

func EditUserInfo(context *gin.Context, db *gorm.DB) {
	uid, _ := context.Get("userId")

	var userInfo dto.EditUserInfoRequestDto
	if err := context.ShouldBindJSON(&userInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	user := models.User{ID: uid.(uuid.UUID)}
	if err := userrepository.GetUser(&user, db); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var fieldErrors []string

	if userInfo.Name != "" && userInfo.Name != user.Name {
		if len(userInfo.Name) < 6 {
			fieldErrors = append(fieldErrors, "Nome deve conter no mínimo 6 caracteres.")
		} else if len(userInfo.Name) > 50 {
			fieldErrors = append(fieldErrors, "Nome deve conter no máximo 50 caracteres.")
		} else {
			user.Name = userInfo.Name
		}
	}
	if userInfo.Nickname != "" && userInfo.Nickname != user.Nickname {
		if len(userInfo.Nickname) < 6 {
			fieldErrors = append(fieldErrors, "Nickname deve conter no mínimo 6 caracteres.")
		} else if len(userInfo.Nickname) > 50 {
			fieldErrors = append(fieldErrors, "Nickname deve conter no máximo 50 caracteres.")
		} else {
			user.Nickname = userInfo.Nickname
		}
	}
	if userInfo.Birthdate != "" {
		if birthdate, err := authservice.ParseBirthdate(userInfo.Birthdate); err != nil {
			fieldErrors = append(fieldErrors, "Data de nascimento inválida.")
		} else {
			user.Birthdate = birthdate
		}
	}
	if userInfo.Country != user.Country {
		if len(userInfo.Country) < 4 {
			fieldErrors = append(fieldErrors, "País deve conter no mínimo 4 caracteres.")
		} else if len(userInfo.Country) > 50 {
			fieldErrors = append(fieldErrors, "País deve conter no máximo 50 caracteres.")
		} else {
			user.Country = userInfo.Country
		}
	}
	if userInfo.City != user.City {
		if len(userInfo.City) < 4 {
			fieldErrors = append(fieldErrors, "Cidade deve conter no mínimo 4 caracteres.")
		} else if len(userInfo.City) > 50 {
			fieldErrors = append(fieldErrors, "Cidade deve conter no máximo 50 caracteres.")
		} else {
			user.City = userInfo.City
		}
	}
	if userInfo.Picture != user.Picture && userInfo.Picture != "" {
		if _, err := url.Parse(userInfo.Picture); err == nil {
			user.Picture = userInfo.Picture
		} else {
			fieldErrors = append(fieldErrors, "URL da foto de perfil inválida.")
		}
	}
	if userInfo.Banner != user.Banner && userInfo.Banner != "" {
		if _, err := url.Parse(userInfo.Banner); err == nil {
			user.Banner = userInfo.Banner
		} else {
			fieldErrors = append(fieldErrors, "URL do banner de perfil inválida.")
		}
	}

	if len(fieldErrors) > 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": strings.Join(fieldErrors, " "),
		})
		return
	}

	if err := userrepository.UpdateUser(&user, db); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.ConstraintName == "users_nickname_key" {
			context.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": err.Error(),
				"data": gin.H{
					"reason": "nickname",
				},
			})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	response := dto.EditUserInfoResponseDto{
		User: dto.UserToResponseDto(user),
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User info updated successfully",
		"data":    response,
	})
}
