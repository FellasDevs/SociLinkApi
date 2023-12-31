package postcontroller

import (
	"SociLinkApi/dto"
	"SociLinkApi/models"
	likerepository "SociLinkApi/repository/like"
	postrepository "SociLinkApi/repository/post"
	authtypes "SociLinkApi/types/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreatePost(context *gin.Context, db *gorm.DB) {
	var postData dto.CreatePostRequestDto

	if err := context.ShouldBindJSON(&postData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var fieldErrors []string

	if postData.Content == "" {
		fieldErrors = append(fieldErrors, "Conteúdo não pode ser vazio.")
	}
	if postData.Visibility != "public" && postData.Visibility != "private" && postData.Visibility != "friends" {
		fieldErrors = append(fieldErrors, "Visibilidade deve ser public, private ou friends.")
	}

	if len(fieldErrors) > 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": strings.Join(fieldErrors, " "),
		})
		return
	}

	visibility := authtypes.Public

	if postData.Visibility == "private" {
		visibility = authtypes.Private
	} else if postData.Visibility == "friends" {
		visibility = authtypes.Friends
	}

	uid, _ := context.Get("userId")
	userId := uid.(uuid.UUID)

	var originalPostId *uuid.UUID = nil

	// Check if original post exists
	if postData.OriginalPostId != "" {
		if originalPostUuid, err := uuid.Parse(postData.OriginalPostId); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Id do post original inválido.",
			})
			return
		} else {
			originalPost := models.Post{ID: originalPostUuid}

			// Get post user is referencing
			err = postrepository.GetPost(&originalPost, &userId, db)
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Post original não encontrado.",
				})
				return
			}

			// If post user is referencing is a repost, get the original post
			if originalPost.OriginalPostID != nil {
				originalPost = models.Post{ID: *originalPost.OriginalPostID}
				err = postrepository.GetPost(&originalPost, &userId, db)
				if err != nil {
					context.JSON(http.StatusBadRequest, gin.H{
						"success": false,
						"message": "Post original não encontrado.",
					})
					return
				}
			}

			originalPostId = &originalPost.ID
		}
	}

	post := models.Post{
		UserID:         userId,
		OriginalPostID: originalPostId,
		Content:        postData.Content,
		Images:         postData.Images,
		Visibility:     string(visibility),
	}

	if err := postrepository.CreatePost(&post, db); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	response := dto.CreatePostResponseDto{
		Post: dto.PostToResponseDto(post, 0, false),
	}

	if post.OriginalPostID != nil {
		originalPost := models.Post{
			ID: *post.OriginalPostID,
		}

		err := postrepository.GetPost(&originalPost, &userId, db)

		if err == nil {
			likes, _ := likerepository.GetPostLikes(post.ID, db)

			userLikedPost := false
			for _, like := range likes {
				if like.UserID == userId {
					userLikedPost = true
					break
				}
			}

			originalPostResponseDto := dto.PostToResponseDto(originalPost, len(likes), userLikedPost)
			response.Post.OriginalPost = &originalPostResponseDto
		}
	}

	context.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Post criado com sucesso!",
		"data":    response,
	})
}
