package postgres_test

import (
	"fmt"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/post/storage/repo"
	"github.com/stretchr/testify/require"
)

func createPost(t *testing.T) *repo.Post {
	n, _ := faker.RandomInt(10)
	Post, err := strg.Post().Create(&repo.Post{
		Title:       faker.Sentence(),
		Description: faker.Sentence(),
		ImageUrl:    faker.URL(),
		UserId:      1,
		CategoryId:  1,
		ViewsCount:  n[0],
	})
	fmt.Println(Post.UserId)
	require.NoError(t, err)
	require.NotEmpty(t, Post)
	return Post
}

func deletePost(id int, t *testing.T) {
	err := strg.Post().Delete(id)
	require.NoError(t, err)
}

func TestGetPost(t *testing.T) {
	n := createPost(t)
	note, err := strg.Post().Get(n.Id)
	require.NoError(t, err)
	require.NotEmpty(t, note)

	deletePost(note.Id, t)
}

func TestCreatePost(t *testing.T) {
	createPost(t)
}

func TestUpdatePost(t *testing.T) {
	n := createPost(t)
	num, _ := faker.RandomInt(10)

	//UserId and CategoryId are connected with user. So, I gave values
	n.UserId = 1
	n.CategoryId = 1

	n.Title = faker.Sentence()
	n.Description = faker.Sentence()
	n.ImageUrl = faker.URL()
	n.ViewsCount = num[0]

	Post, err := strg.Post().Update(n)
	require.NoError(t, err)
	require.NotEmpty(t, Post)

	deletePost(Post.Id, t)
}

func TestDeletePost(t *testing.T) {
	u := createPost(t)
	deletePost(u.Id, t)
}

func TestGetAllPost(t *testing.T) {
	u := createPost(t)
	n, _ := faker.RandomInt(100)
	_, err := strg.Post().GetAll(repo.GetPostQuery{
		Page:  n[0],
		Limit: n[0],
	})

	require.NoError(t, err)
	deletePost(u.Id, t)
}
