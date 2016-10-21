package plan

import (
	"io"

	"github.com/mvader/gitql/sql"
)

type Limit struct {
	UnaryNode
	size int64
}

func NewLimit(size int64, child sql.Node) *Limit {
	return &Limit{
		UnaryNode: UnaryNode{Child: child},
		size:      size,
	}
}

func (l *Limit) Schema() sql.Schema {
	return l.UnaryNode.Child.Schema()
}

func (l *Limit) Resolved() bool {
	return true
}

func (l *Limit) RowIter() (sql.RowIter, error) {
	li, err := l.Child.RowIter()
	if err != nil {
		return nil, err
	}
	return &limitIter{l, 0, li}, nil
}

type limitIter struct {
	l          *Limit
	currentPos int64
	childIter  sql.RowIter
}

func (li *limitIter) Next() (sql.Row, error) {
	if li.currentPos >= li.l.size {
		return nil, io.EOF
	}
	childRow, err := li.childIter.Next()
	li.currentPos++
	if err != nil {
		return nil, err
	}
	return childRow, nil
}