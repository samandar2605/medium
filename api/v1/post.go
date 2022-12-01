package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/post/api/models"
	"github.com/post/storage/repo"
)

// @Router /posts/{id} [get]
// @Summary Get post by id
// @Description Get post by id
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Post
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err = h.storage.Post().ViewsInc(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	resp, err := h.storage.Post().Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	usr, _ := h.storage.User().GetUserProfileInfo(resp.UserId)
	c.JSON(http.StatusOK, models.Post{
		Id:          resp.Id,
		Title:       resp.Title,
		Description: resp.Description,
		ImageUrl:    resp.ImageUrl,
		UserId:      resp.UserId,
		CategoryId:  resp.CategoryId,
		UpdatedAt:   resp.UpdatedAt,
		ViewsCount:  resp.ViewsCount,
		CreatedAt:   resp.CreatedAt,
		User: models.UserProfile{
			Id:              resp.UserId,
			FirstName:       usr.FirstName,
			LastName:        usr.LastName,
			Email:           usr.Email,
			ProfileImageUrl: usr.ProfileImageUrl,
		},
	})
}

// @Security ApiKeyAuth
// @Router /posts [post]
// @Summary Create a post
// @Description Create a post
// @Tags post
// @Accept json
// @Produce json
// @Param post body models.CreatePost true "post"
// @Success 201 {object} models.Post
// @Failure 500 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreatePost(c *gin.Context) {
	var (
		req models.CreatePost
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

	image, _ := h.storage.User().GetUserProfileInfo(usr.UserId)
	resp, err := h.storage.Post().Create(&repo.Post{
		Title:       req.Title,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		UserId:      usr.UserId,
		CategoryId:  req.CategoryId,
		User: repo.UserProfile{
			Id:              usr.UserId,
			FirstName:       usr.FirstName,
			LastName:        usr.LastName,
			Email:           usr.Email,
			ProfileImageUrl:  image.ProfileImageUrl,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Post{
		Id:          resp.Id,
		Title:       resp.Title,
		Description: resp.Description,
		ImageUrl:    resp.ImageUrl,
		UserId:      resp.UserId,
		CategoryId:  resp.CategoryId,
		ViewsCount:  resp.ViewsCount,
		UpdatedAt:   resp.UpdatedAt,
		CreatedAt:   resp.CreatedAt,
		User:        models.UserProfile(resp.User),
	})
}

// @Router /posts [get]
// @Summary Get all posts
// @Description Get all posts
// @Tags post
// @Accept json
// @Produce json
// @Param filter query models.GetAllPostsParams false "Filter"
// @Success 200 {object} models.GetAllPostsResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllPost(c *gin.Context) {
	req, err := postsParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	result, err := h.storage.Post().GetAll(repo.GetPostQuery{
		Page:       req.Page,
		Limit:      req.Limit,
		CategoryID: req.CategoryId,
		UserID:     req.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, postsResponse(h, result))
}

func postsParams(c *gin.Context) (*models.GetAllPostsParams, error) {
	var (
		limit              int = 10
		page               int = 1
		err                error
		CategoryId, UserId int
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

	if c.Query("category_id") != "" {
		CategoryId, err = strconv.Atoi(c.Query("category_id"))
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

	return &models.GetAllPostsParams{
		Limit:      limit,
		Page:       page,
		CategoryId: CategoryId,
		UserID:     UserId,
	}, nil
}

func postsResponse(h *handlerV1, data *repo.GetAllPostResult) *models.GetAllPostsResponse {
	response := models.GetAllPostsResponse{
		Posts: make([]*models.Post, 0),
		Count: data.Count,
	}

	for _, post := range data.Post {
		usr, _ := h.storage.User().GetUserProfileInfo(post.UserId)
		if err := h.storage.Post().ViewsInc(post.Id); err != nil {
			return nil
		}
		post.User = repo.UserProfile{
			Id:              post.UserId,
			FirstName:       usr.FirstName,
			LastName:        usr.LastName,
			Email:           usr.Email,
			ProfileImageUrl: usr.ProfileImageUrl,
		}
		post.ViewsCount++
		p := parsePostModel(post)
		response.Posts = append(response.Posts, &p)
	}

	return &response
}

func parsePostModel(post *repo.Post) models.Post {

	return models.Post{
		Id:          post.Id,
		Title:       post.Title,
		Description: post.Description,
		ImageUrl:    post.ImageUrl,
		UserId:      post.UserId,
		CategoryId:  post.CategoryId,
		UpdatedAt:   post.UpdatedAt,
		ViewsCount:  post.ViewsCount,
		CreatedAt:   post.CreatedAt,
		User: models.UserProfile{
			Id:              post.UserId,
			FirstName:       post.User.FirstName,
			LastName:        post.User.LastName,
			Email:           post.User.Email,
			ProfileImageUrl: post.User.ProfileImageUrl,
		},
	}
}

// @Security ApiKeyAuth
// @Summary Update a post
// @Description Update a post
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param user body models.CreatePost true "post"
// @Success 200 {object} models.Post
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id} [put]
func (h *handlerV1) UpdatePost(ctx *gin.Context) {
	var b models.Post

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

	usr, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	b.Id = id
	post, err := h.storage.Post().Update(&repo.Post{
		Id:          b.Id,
		Title:       b.Title,
		Description: b.Description,
		ImageUrl:    b.ImageUrl,
		UserId:      usr.UserId,
		CategoryId:  b.CategoryId,
		ViewsCount:  b.ViewsCount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	profil, _ := h.storage.User().GetUserProfileInfo(post.UserId)
	post.User.Id = profil.Id
	post.User.FirstName = profil.FirstName
	post.User.LastName = profil.LastName
	post.User.Email = profil.Email
	post.User.ProfileImageUrl = profil.ProfileImageUrl

	ctx.JSON(http.StatusOK, post)
}

// @Security ApiKeyAuth
// @Summary Delete a posts
// @Description Delete a posts
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id} [delete]
func (h *handlerV1) DeletePost(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to convert",
		})
		return
	}

	err = h.storage.Post().Delete(id)
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
