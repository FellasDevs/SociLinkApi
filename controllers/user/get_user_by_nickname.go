package usercontroller

import (
	"SociLinkApi/dto"
	"SociLinkApi/models"
	userrepository "SociLinkApi/repository/user"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetUserByNickname(context *gin.Context, db *gorm.DB) {
	var params dto.GetUserByNicknameRequestDto

	if err := context.ShouldBindQuery(&params); err != nil || params.Nickname == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Apelido do usuário deve ser informado",
		})
		return
	}

	nickname := params.Nickname

	user := models.User{Nickname: nickname}
	if err := userrepository.GetUser(&user, db); err != nil {
		var statusCode int

		if errors.Is(err, gorm.ErrRecordNotFound) {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}

		context.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	response := dto.GetUserByNicknameResponseDto{
		User: dto.UserToResponseDto(user),
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "usuário obtido com sucesso",
		"data":    response,
	})
}
