package postgres_test

import (
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/post/storage/repo"
	"github.com/stretchr/testify/require"
)

func createCategory(t *testing.T) *repo.Category {
	Category, err := strg.Category().Create(&repo.Category{
		Title: faker.Sentence(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, Category)
	return Category
}

func deleteCategory(id int, t *testing.T) {
	err := strg.Category().Delete(id)
	require.NoError(t, err)
}

func TestGetCategory(t *testing.T) {
	n := createCategory(t)
	note, err := strg.Category().Get(n.Id)
	require.NoError(t, err)
	require.NotEmpty(t, note)
	deleteCategory(note.Id, t)
}

func TestCreateCategory(t *testing.T) {
	createCategory(t)
}

func TestUpdateCategory(t *testing.T) {
	n := createCategory(t)

	n.Title = faker.Sentence()

	Category, err := strg.Category().Update(n)
	require.NoError(t, err)
	require.NotEmpty(t, Category)

	deleteCategory(Category.Id, t)
}

func TestDeleteCategory(t *testing.T) {
	u := createCategory(t)
	deleteCategory(u.Id, t)
}

func TestGetAllCategory(t *testing.T) {
	u := createCategory(t)
	n, _ := faker.RandomInt(100)
	_, err := strg.Category().GetAll(repo.GetCategoryQuery{
		Page:  n[0],
		Limit: n[1],
	})

	require.NoError(t, err)
	deleteCategory(u.Id, t)
}
