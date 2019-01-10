package querybuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	b := &Builder{}
	assert.Equal(t, "", b.String())

	b = b.Select("u.*")
	assert.Equal(t, "SELECT u.*", b.String())

	b = b.Select("b.UserId IS NOT NULL AS IsBot")
	assert.Equal(t, "SELECT u.*, b.UserId IS NOT NULL AS IsBot", b.String())

	b = (&Builder{}).Select("*").From("Users u")
	assert.Equal(t, "SELECT * FROM Users u", b.String())

	b = b.From("Bots b")
	assert.Equal(t, "SELECT * FROM Users u, Bots b", b.String())

	b = (&Builder{}).Select("*").From("Users u")
	b = b.LeftJoin("Bots b ON ( b.UserId = u.Id )")
	assert.Equal(t, "SELECT * FROM Users u LEFT JOIN Bots b ON ( b.UserId = u.Id )", b.String())
}

func TestClone(t *testing.T) {
	b := &Builder{}
	b = b.Select("u.*")
	assert.Equal(t, "SELECT u.*", b.String())

	clone := b.Clone()
	assert.Equal(t, "SELECT u.*", clone.String())
}

func TestRandom(t *testing.T) {
	builder := New()
	builder = builder.Select("u.*")
	builder = builder.Select("b.UserId IS NOT NULL AS IsBot")
	builder = builder.From("Users u")
	builder = builder.LeftJoin("Bots b ON ( b.UserId = u.Id )")
	builder = builder.Where("Id = :UserId").Bind("UserId", "id")

	assert.Equal(t, "SELECT u.*, b.UserId IS NOT NULL AS IsBot FROM Users u LEFT JOIN Bots b ON ( b.UserId = u.Id ) WHERE Id = :UserId", builder.String())
	assert.Equal(t, map[string]interface{}{"UserId": "id"}, builder.Args())
}

func TestBind(t *testing.T) {
	// t.Run("empty array binding", func(t *testing.T) {
	// 	builder := New().Select("*").From("Users").Where("Id IN (:Ids)").Bind("Ids", []string{})
	// 	assert.Equal(t, "SELECT * FROM Users WHERE Id IN ()", builder.String())
	// 	assert.Equal(t, map[string]interface{}{}, builder.Args())
	// })

	t.Run("array with multiple string elements", func(t *testing.T) {
		builder := New().Select("*").From("Users").Where("Id IN (:Ids)").Bind("Ids", []string{"id1", "id2", "id3"})
		assert.Equal(t, "SELECT * FROM Users WHERE Id IN (:Ids_0, :Ids_1, :Ids_2)", builder.String())
		assert.Equal(t, map[string]interface{}{
			"Ids_0": "id1",
			"Ids_1": "id2",
			"Ids_2": "id3",
		}, builder.Args())
	})

	t.Run("array with multiple integer elements", func(t *testing.T) {
		builder := New().Select("*").From("Users").Where("Count IN (:Ids)").Bind("Ids", []int{1, 2, 3})
		assert.Equal(t, "SELECT * FROM Users WHERE Count IN (:Ids_0, :Ids_1, :Ids_2)", builder.String())
		assert.Equal(t, map[string]interface{}{
			"Ids_0": 1,
			"Ids_1": 2,
			"Ids_2": 3,
		}, builder.Args())
	})
}
