package postcontroller

import (
	"SociLinkApi/dto"
	"SociLinkApi/models"
	likerepository "SociLinkApi/repository/like"
	postrepository "SociLinkApi/repository/post"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
)

func SearchPosts(context *gin.Context, db *gorm.DB) {
	uid, exist := context.Get("userId")
	var userId *uuid.UUID
	if exist {
		id := uid.(uuid.UUID)
		userId = &id
	}

	var params dto.SearchPostRequestDto
	if err := context.ShouldBindQuery(&params); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if posts, err := postrepository.SearchPosts(params.Search, userId, params.PaginationRequestDto, db); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
	} else {
		response := dto.SearchPostResponseDto{
			Posts: make([]dto.PostResponseDto, len(posts)),
		}

		for i, post := range posts {
			likes, _ := likerepository.GetPostLikes(post.ID, db)

			userLikedPost := false
			if userId != nil {
				for _, like := range likes {
					if like.UserID == *userId {
						userLikedPost = true
						break
					}
				}
			}

			response.Posts[i] = dto.PostToResponseDto(post, len(likes), userLikedPost)

			if post.OriginalPostID != nil {
				originalPost := models.Post{
					ID: *post.OriginalPostID,
				}

				err = postrepository.GetPost(&originalPost, userId, db)

				if err == nil {
					likes, _ = likerepository.GetPostLikes(post.ID, db)

					userLikedPost = false
					for _, like := range likes {
						if like.UserID == *userId {
							userLikedPost = true
							break
						}
					}

					originalPostResponseDto := dto.PostToResponseDto(originalPost, len(likes), userLikedPost)
					response.Posts[i].OriginalPost = &originalPostResponseDto
				}
			}
		}

		context.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "posts encontrados",
			"data":    response,
		})
	}
}
