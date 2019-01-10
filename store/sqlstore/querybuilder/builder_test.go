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
	assert.Equal(t, map[string]interface{}{"UserId": "id"}, builder.Bindings())
}
