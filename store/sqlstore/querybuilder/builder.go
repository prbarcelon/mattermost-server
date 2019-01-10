package querybuilder

import (
	"strings"
)

type Builder struct {
	selectStatements  []string
	fromStatements    []string
	joinStatements    []string
	whereStatements   []string
	orderByStatements []string
	offset            string
	limit             string
	bindings          map[string]interface{}
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) Clone() *Builder {
	clone := &Builder{
		selectStatements:  make([]string, len(b.selectStatements)),
		fromStatements:    make([]string, len(b.fromStatements)),
		joinStatements:    make([]string, len(b.joinStatements)),
		whereStatements:   make([]string, len(b.whereStatements)),
		orderByStatements: make([]string, len(b.orderByStatements)),
		offset:            b.offset,
		limit:             b.limit,
	}
	copy(clone.selectStatements, b.selectStatements)
	copy(clone.fromStatements, b.fromStatements)
	copy(clone.joinStatements, b.joinStatements)
	copy(clone.whereStatements, b.whereStatements)
	copy(clone.orderByStatements, b.orderByStatements)

	if len(b.bindings) > 0 {
		clone.bindings = make(map[string]interface{}, len(b.bindings))
	}
	for key, value := range b.bindings {
		clone.bindings[key] = value
	}

	return clone
}

func (b *Builder) Select(sql string) *Builder {
	clone := b.Clone()
	clone.selectStatements = append(clone.selectStatements, sql)

	return clone
}

func (b *Builder) From(sql string) *Builder {
	clone := b.Clone()
	clone.fromStatements = append(clone.fromStatements, sql)

	return clone
}

func (b *Builder) LeftJoin(sql string) *Builder {
	clone := b.Clone()
	clone.joinStatements = append(clone.joinStatements, "LEFT JOIN "+sql)

	return clone
}

func (b *Builder) Where(sql string) *Builder {
	clone := b.Clone()
	clone.whereStatements = append(clone.whereStatements, sql)

	return clone
}

func (b *Builder) OrderBy(sql string) *Builder {
	clone := b.Clone()
	clone.orderByStatements = append(clone.orderByStatements, sql)

	return clone
}

func (b *Builder) Bind(key string, value interface{}) *Builder {
	clone := b.Clone()

	if clone.bindings == nil {
		clone.bindings = make(map[string]interface{})
	}
	clone.bindings[key] = value

	return clone
}

func (b *Builder) Offset(offset string) *Builder {
	clone := b.Clone()
	clone.offset = offset

	return clone
}

func (b *Builder) Limit(limit string) *Builder {
	clone := b.Clone()
	clone.limit = limit

	return clone
}

func (b *Builder) String() string {
	var query string

	if len(b.selectStatements) > 0 {
		query += "SELECT " + strings.Join(b.selectStatements, ", ")
	}
	if len(b.fromStatements) > 0 {
		query += " FROM " + strings.Join(b.fromStatements, ", ")
	}
	if len(b.joinStatements) > 0 {
		query += " " + strings.Join(b.joinStatements, " ")
	}
	if len(b.whereStatements) > 0 {
		query += " WHERE " + strings.Join(b.whereStatements, " AND ")
	}
	if len(b.orderByStatements) > 0 {
		query += " ORDER BY " + strings.Join(b.orderByStatements, ", ")
	}
	if len(b.offset) > 0 {
		query += " OFFSET " + b.offset
	}
	if len(b.limit) > 0 {
		query += " LIMIT " + b.limit
	}

	return query
}

func (b *Builder) Bindings() map[string]interface{} {
	return b.bindings
}
