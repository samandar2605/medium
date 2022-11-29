package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/post/api/models"
	"github.com/post/storage/repo"
)

// @Router /categories/{id} [get]
// @Summary Get category by id
// @Description Get category by id
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.Category
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println("Error at GetCategory 1")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	resp, err := h.storage.Category().Get(id)
	if err != nil {
		fmt.Println("Error at GetCategory 2")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Category{
		Id:        resp.Id,
		Title:     resp.Title,
		CreatedAt: resp.CreatedAt,
	})
}

// @Security ApiKeyAuth
// @Router /categories [post]
// @Summary Create a category
// @Description Create a category
// @Tags category
// @Accept json
// @Produce json
// @Param category body models.CreateCategory true "Category"
// @Success 201 {object} models.Category
// @Failure 500 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateCategory(c *gin.Context) {
	var (
		req models.CreateCategory
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		fmt.Println("Error at CreateCategory 1")

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := h.storage.Category().Create(&repo.Category{
		Title: req.Title,
	})
	if err != nil {
		fmt.Println("Error at CreateCategory 2")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Category{
		Id:        resp.Id,
		Title:     resp.Title,
		CreatedAt: resp.CreatedAt,
	})
}

// @Router /categories [get]
// @Summary Get all categories
// @Description Get all categories
// @Tags category
// @Accept json
// @Produce json
// @Param filter query models.GetAllCategoryParams false "Filter"
// @Success 200 {object} models.GetAllCategoriesResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllCategories(c *gin.Context) {
	req, err := categoryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := h.storage.Category().GetAll(repo.GetCategoryQuery{
		Page:   req.Page,
		Limit:  req.Limit,
		Search: req.Search,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, categoryResponse(result))
}

func categoryParams(c *gin.Context) (*models.GetAllCategoryParams, error) {
	var (
		limit int = 10
		page  int = 1
		err   error
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

	return &models.GetAllCategoryParams{
		Page:   page,
		Limit:  limit,
		Search: c.Query("search"),
	}, nil
}

func categoryResponse(data *repo.GetAllCategoriesResult) *models.GetAllCategoriesResponse {
	response := models.GetAllCategoriesResponse{
		Categories: make([]*models.Category, 0),
		Count:      data.Count,
	}

	for _, category := range data.Categories {
		p := parseCategoryModel(category)
		response.Categories = append(response.Categories, &p)
	}

	return &response
}

func parseCategoryModel(Category *repo.Category) models.Category {
	return models.Category{
		Id:        Category.Id,
		Title:     Category.Title,
		CreatedAt: Category.CreatedAt,
	}
}

// @Security ApiKeyAuth
// @Summary Update a Category
// @Description Update a Category
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param user body models.CreateCategory true "Category"
// @Success 200 {object} models.Category
// @Failure 500 {object} models.ErrorResponse
// @Router /categories/{id} [put]
func (h *handlerV1) UpdateCategory(ctx *gin.Context) {
	var b models.Category

	err := ctx.ShouldBindJSON(&b)
	if err != nil {
		fmt.Println("Error at UpdateCategory 1")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fmt.Println("Error at UpdateCategory 2")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	b.Id = id
	category, err := h.storage.Category().Update(&repo.Category{
		Id:    b.Id,
		Title: b.Title,
	})
	fmt.Println(category)
	if err != nil {
		fmt.Println("Error at UpdateCategory 3")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create category",
		})
		return
	}

	ctx.JSON(http.StatusOK, category)
}

// @Security ApiKeyAuth
// @Summary Delete a categories
// @Description Delete a categories
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Failure 500 {object} models.ErrorResponse
// @Router /categories/{id} [delete]
func (h *handlerV1) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fmt.Println("Error at DeleteCategory 1")

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = h.storage.Category().Delete(id)
	if err != nil {
		fmt.Println("Error at DeleteCategory 2")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to Delete method",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successful delete method",
	})
}
