package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/post/api/models"
	"github.com/post/storage/repo"
)

// @Router /comments/{id} [get]
// @Summary Get comment by id
// @Description Get comment by id
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Comment
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := h.storage.Comment().Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Comment{
		Id:          resp.Id,
		PostId:      resp.PostId,
		UserId:      resp.UserId,
		Description: resp.Description,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
		User: &models.UserProfile{
			Id:              resp.UserId,
			FirstName:       resp.User.FirstName,
			LastName:        resp.User.LastName,
			Email:           resp.User.Email,
			ProfileImageUrl: resp.User.ProfileImageUrl,
		},
	})
}

// @Security ApiKeyAuth
// @Router /comments [post]
// @Summary Create a comment
// @Description Create a comment
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body models.CreateComment true "comment"
// @Success 201 {object} models.Comment
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateComment(c *gin.Context) {
	var (
		req models.CreateComment
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	usr, err := h.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := h.storage.Comment().Create(&repo.Comment{
		PostId:      req.PostId,
		UserId:      usr.UserId,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Comment{
		Id:          resp.Id,
		PostId:      resp.PostId,
		UserId:      resp.UserId,
		Description: resp.Description,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
		User: &models.UserProfile{
			Id:              usr.UserId,
			FirstName:       usr.FirstName,
			LastName:        usr.LastName,
			Email:           usr.Email,
			ProfileImageUrl: usr.ProfileImageUrl,
		},
	})
}

// @Router /comments [get]
// @Summary Get all comments
// @Description Get all comments
// @Tags comments
// @Accept json
// @Produce json
// @Param filter query models.GetAllCommentsParams false "Filter"
// @Success 200 {object} models.GetAllResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllComment(c *gin.Context) {
	req, err := commentsParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	result, err := h.storage.Comment().GetAll(repo.GetCommentQuery{
		Page:       req.Page,
		Limit:      req.Limit,
		PostId:     req.PostID,
		SortByDate: req.SortByDate,
		UserId:     req.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, commentsResponse(h, result))
}

func commentsParams(c *gin.Context) (*models.GetAllCommentsParams, error) {
	var (
		limit          int = 10
		page           int = 1
		err            error
		sortByDate     string
		PostId, UserId int
	)

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("sort_by_date") != "" &&
		(c.Query("sort_by_date") == "desc" || c.Query("sort_by_date") == "asc" || c.Query("sort_by_date") == "none") {
		sortByDate = c.Query("sort_by_date")
	}

	if c.Query("post_id") != "" {
		PostId, err = strconv.Atoi(c.Query("post_id"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("user_id") != "" {
		UserId, err = strconv.Atoi(c.Query("user_id"))
		if err != nil {
			return nil, err
		}
	}

	return &models.GetAllCommentsParams{
		Limit:      limit,
		Page:       page,
		SortByDate: sortByDate,
		PostID:     PostId,
		UserID:     UserId,
	}, nil
}

func commentsResponse(h *handlerV1, data *repo.GetAllCommentsResult) *models.GetAllCommentsResponse {
	response := models.GetAllCommentsResponse{
		Comments: make([]*models.Comment, 0),
		Count:    data.Count,
	}

	for _, comment := range data.Comments {
		usr, _ := h.storage.User().GetUserProfileInfo(comment.UserId)
		comment.User = repo.UserProfile{
			Id:              comment.UserId,
			FirstName:       usr.FirstName,
			LastName:        usr.LastName,
			Email:           usr.Email,
			ProfileImageUrl: usr.ProfileImageUrl,
		}

		p := parseCommentModel(comment)
		response.Comments = append(response.Comments, &p)
	}

	return &response
}

func parseCommentModel(Comment *repo.Comment) models.Comment {
	return models.Comment{
		Id:          Comment.Id,
		PostId:      Comment.PostId,
		UserId:      Comment.UserId,
		Description: Comment.Description,
		CreatedAt:   Comment.CreatedAt,
		UpdatedAt:   Comment.UpdatedAt,
		User: &models.UserProfile{
			Id:              Comment.UserId,
			FirstName:       Comment.User.FirstName,
			LastName:        Comment.User.LastName,
			Email:           Comment.User.Email,
			ProfileImageUrl: Comment.User.ProfileImageUrl,
		},
	}
}

// @Security ApiKeyAuth
// @Summary Update a comment
// @Description Update a comments
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param comment body models.UpdateComment true "comment"
// @Success 200 {object} models.Comment
// @Failure 500 {object} models.ErrorResponse
// @Router /comments/{id} [put]
func (h *handlerV1) UpdateComment(ctx *gin.Context) {
	var b models.UpdateComment
	err := ctx.ShouldBindJSON(&b)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if payload.UserType != repo.UserTypeSuperadmin || payload.UserId!= h.storage.Comment().GetUserInfo(id) {
		ctx.JSON(http.StatusForbidden, errorResponse(ErrForbidden))
		return
	}


	b.Id = id
	comment, err := h.storage.Comment().Update(&repo.Comment{
		Id:          b.Id,
		Description: b.Description,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	profil, _ := h.storage.User().GetUserProfileInfo(comment.UserId)
	comment.User.Id = profil.Id
	comment.User.FirstName = profil.FirstName
	comment.User.LastName = profil.LastName
	comment.User.Email = profil.Email
	comment.User.ProfileImageUrl = profil.ProfileImageUrl

	ctx.JSON(http.StatusOK, comment)
}

// @Security ApiKeyAuth
// @Summary Delete a comment
// @Description Delete a comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Failure 500 {object} models.ErrorResponse
// @Router /comments/{id} [delete]
func (h *handlerV1) DeleteComment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to convert",
		})
		return
	}
	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if payload.UserType != repo.UserTypeSuperadmin || payload.UserId!= h.storage.Comment().GetUserInfo(id) {
		ctx.JSON(http.StatusForbidden, errorResponse(ErrForbidden))
		return
	}

	err = h.storage.Comment().Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to Delete method",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successful delete method",
	})
}
