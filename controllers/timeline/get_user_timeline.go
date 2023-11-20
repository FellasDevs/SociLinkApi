package timeline

import (
	"SociLinkApi/dto"
	"SociLinkApi/models"
	frienshiprepository "SociLinkApi/repository/frienship"
	postrepository "SociLinkApi/repository/post"
	userrepository "SociLinkApi/repository/user"
	authtypes "SociLinkApi/types/auth"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
)

func GetUserTimeline(context *gin.Context, db *gorm.DB) {
	nickname := context.Param("nick")

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

	visibility := authtypes.Public

	uid, exists := context.Get("userId")
	if exists {
		userId := uid.(uuid.UUID)

		if userId == user.ID {
			visibility = authtypes.Private
		} else if _, err := frienshiprepository.GetFriendshipByUsers(userId, user.ID, db); err == nil {
			visibility = authtypes.Friends
		}
	}

	if posts, err := postrepository.GetPostsByUser(user.ID, visibility, db); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
	} else {
		var response dto.GetUserTimelineResponseDto

		response.User = dto.UserResponseDto{
			Id:        user.ID.String(),
			Name:      user.Name,
			Nickname:  user.Nickname,
			Birthdate: user.Birthdate.String(),
			Country:   user.Country,
			City:      user.City,
			Picture:   user.Picture,
			Banner:    user.Banner,
			CreatedAt: user.CreatedAt.String(),
		}

		response.Posts = make([]dto.PostResponseDto, len(posts))

		for i, post := range posts {
			response.Posts[i] = dto.PostResponseDto{
				Id: post.ID.String(),
				User: dto.UserResponseDto{
					Id:        post.User.ID.String(),
					Name:      post.User.Name,
					Nickname:  post.User.Nickname,
					Birthdate: post.User.Birthdate.String(),
					Country:   post.User.Country,
					City:      post.User.City,
					Picture:   post.User.Picture,
					Banner:    post.User.Banner,
					CreatedAt: post.User.CreatedAt.String(),
				},
				Content:    post.Content,
				Images:     post.Images,
				Visibility: post.Visibility,
			}
		}

		context.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "posts recuperados com sucesso",
			"data":    response,
		})
	}
}