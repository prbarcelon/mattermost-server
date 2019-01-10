package querybuilder

import (
	"fmt"
	"reflect"
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
	args              map[string]interface{}
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

	if len(b.args) > 0 {
		clone.args = make(map[string]interface{}, len(b.args))
	}
	for key, value := range b.args {
		clone.args[key] = value
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

func (b *Builder) Join(sql string) *Builder {
	clone := b.Clone()
	clone.joinStatements = append(clone.joinStatements, "INNER JOIN "+sql)

	return clone
}

func (b *Builder) LeftJoin(sql string) *Builder {
	clone := b.Clone()
	clone.joinStatements = append(clone.joinStatements, "LEFT JOIN "+sql)

	return clone
}

func (b *Builder) RightJoin(sql string) *Builder {
	clone := b.Clone()
	clone.joinStatements = append(clone.joinStatements, "RIGHT JOIN "+sql)

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

	// It's a natural mistake to write the key prefixed with the colon: fix that.
	key = strings.TrimLeft(key, ":")

	if clone.args == nil {
		clone.args = make(map[string]interface{})
	}
	clone.args[key] = value

	// TODO: panic if empty array

	return clone
}

func (b *Builder) Offset(offset string) *Builder {
	clone := b.Clone()
	clone.offset = offset

	return clone
}

func (b *Builder) OffsetInt(offset int) *Builder {
	clone := b.Clone()

	// TODO: panic on Offset already set?
	clone.offset = ":Offset"
	clone = clone.Bind("Offset", offset)

	return clone
}

func (b *Builder) Limit(limit string) *Builder {
	clone := b.Clone()
	clone.limit = limit

	return clone
}

func (b *Builder) LimitInt(limit int) *Builder {
	clone := b.Clone()

	// TODO: panic on Limit already set?
	clone.limit = ":Limit"
	clone = clone.Bind("Limit", limit)

	return clone
}

func (b *Builder) Query() string {
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
	if len(b.limit) > 0 {
		query += " LIMIT " + b.limit
	}
	if len(b.offset) > 0 {
		query += " OFFSET " + b.offset
	}

	// Explode any array parameters in the generated query string.
	for key, value := range b.args {
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			valueLen := reflect.ValueOf(value).Len()

			keys := []string{}
			for index := 0; index < valueLen; index++ {
				keys = append(keys, fmt.Sprintf(":%s_%d", key, index))
			}

			query = strings.Replace(query, fmt.Sprintf(":%s", key), strings.Join(keys, ", "), -1)
		}
	}

	// func (b *Builder) BindStringArray(key string, values []string) *Builder {
	// 	clone := b.Clone()

	// 	valueKeys := make([]string, len(values))
	// 	for index, value := range values {
	// 		valueKey := fmt.Sprintf(":%s_%d", key, index)
	// 		valueKeys[index] = valueKey

	// 		builder.Bind(valueKey, value)

	// 		props["userId"+strconv.Itoa(index)] = userId
	// 		idQuery += ":userId" + strconv.Itoa(index)
	// 	}

	// 	clone.Bind(key,
	// }

	return query
}

func (b *Builder) String() string {
	return b.Query()
}

func (b *Builder) Args() map[string]interface{} {
	args := map[string]interface{}{}

	for key, value := range b.args {
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			switch rt.Elem().Kind() {
			case reflect.String:
				for index, v := range value.([]string) {
					args[fmt.Sprintf("%s_%d", key, index)] = v
				}

			case reflect.Int:
				for index, v := range value.([]int) {
					args[fmt.Sprintf("%s_%d", key, index)] = v
				}

			default:
				panic("unsupported argument array kind: " + rt.Elem().Kind().String())
			}
		default:
			args[key] = value
		}
	}

	return args
}
