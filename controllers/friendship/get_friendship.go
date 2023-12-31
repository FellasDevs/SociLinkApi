package friendshipcontroller

import (
	"SociLinkApi/dto"
	"SociLinkApi/models"
	friendshiprepository "SociLinkApi/repository/friendship"
	userrepository "SociLinkApi/repository/user"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
)

func GetFriendship(context *gin.Context, db *gorm.DB) {
	nickname := context.Param("nickname")

	if nickname == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Apelido do amigo não informado",
		})
		return
	}

	friend := models.User{Nickname: nickname}
	if err := userrepository.GetUser(&friend, db); err != nil {
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

	uid, _ := context.Get("userId")
	userId := uid.(uuid.UUID)

	if friendship, err := friendshiprepository.GetFriendshipByUsers(userId, friend.ID, db); err != nil {
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
	} else {
		response := dto.GetFriendshipResponseDto{
			Friendship: dto.FriendshipToResponseDto(friendship),
		}

		context.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    response,
		})
	}
}
