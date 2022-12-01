package postgres_test

import (
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/post/storage/repo"
	"github.com/stretchr/testify/require"
)

func createComment(t *testing.T) *repo.Comment {
	Comment, err := strg.Comment().Create(&repo.Comment{
		PostId:      1,
		UserId:      1,
		Description: faker.Sentence(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, Comment)
	return Comment
}

func deleteComment(id int, t *testing.T) {
	err := strg.Comment().Delete(id)
	require.NoError(t, err)
}

func TestGetComment(t *testing.T) {
	n := createComment(t)
	note, err := strg.Comment().Get(n.Id)
	require.NoError(t, err)
	require.NotEmpty(t, note)

	deleteComment(note.Id, t)
}

func TestCreateComment(t *testing.T) {
	createComment(t)
}

func TestUpdateComment(t *testing.T) {
	n := createComment(t)
	
	n.PostId = 1
	n.UserId = 1
	n.Description = faker.Sentence()

	Comment, err := strg.Comment().Update(n)
	require.NoError(t, err)
	require.NotEmpty(t, Comment)

	deleteComment(Comment.Id, t)
}

func TestDeleteComment(t *testing.T) {
	u := createComment(t)
	deleteComment(u.Id, t)
}

func TestGetAllComment(t *testing.T) {
	u := createComment(t)
	n, _ := faker.RandomInt(100)
	_, err := strg.Comment().GetAll(repo.GetCommentQuery{
		Page:  n[0],
		Limit: n[0],
	})

	require.NoError(t, err)
	deleteComment(u.Id, t)
}
