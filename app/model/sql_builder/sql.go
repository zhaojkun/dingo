package SQL

import (
	"bytes"
)

func Select(what string) *selectBuilder {
	return newSelect(what)
}

// A simple SQL statement builder
type selectBuilder struct {
	what    string
	from    string
	where   []string
	orderBy string
	limit   string
	offset  string
}

func newSelect(what string) *selectBuilder {
	b := new(selectBuilder)
	b.what = what
	return b
}

func (b *selectBuilder) Select(what string) *selectBuilder {
	b.what = what
	return b
}

func (b *selectBuilder) Copy() *selectBuilder {
	var new selectBuilder
	new = *b
	return &new
}

func (b *selectBuilder) From(from string) *selectBuilder {
	b.from = from
	return b
}

func (b *selectBuilder) Where(stmt ...string) *selectBuilder {
	b.where = append(b.where, stmt...)
	return b
}

func (b *selectBuilder) OrderBy(stmt string) *selectBuilder {
	b.orderBy = stmt
	return b
}

func (b *selectBuilder) Limit(limit string) *selectBuilder {
	b.limit = limit
	return b
}

func (b *selectBuilder) Offset(offset string) *selectBuilder {
	b.offset = offset
	return b
}

func (b *selectBuilder) SQL() string {
	sqlBuf := new(bytes.Buffer)
	sqlBuf.WriteString("SELECT " + b.what + " FROM " + b.from)
	if len(b.where) != 0 {
		buf := new(bytes.Buffer)
		for i, s := range b.where {
			if i != 0 {
				buf.WriteString(" AND ")
			}
			buf.WriteString(s)
		}
		sqlBuf.WriteString(" WHERE ")
		sqlBuf.WriteString(buf.String())
	}
	if b.orderBy != "" {
		sqlBuf.WriteString(" ORDER BY ")
		sqlBuf.WriteString(b.orderBy)
	}
	if b.limit != "" {
		sqlBuf.WriteString(" LIMIT ")
		sqlBuf.WriteString(b.limit)
	}
	if b.offset != "" {
		sqlBuf.WriteString(" OFFSET ")
		sqlBuf.WriteString(b.offset)
	}
	return sqlBuf.String()
}
