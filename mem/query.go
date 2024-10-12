package mem

import (
	"fold/util"
	"net/http"
)

type Condition struct {
	values   []string
	value    string
	any      bool
	multiple bool
}

type Query struct {
	cols []*Condition
	all  bool
}

func QueryForTable(r *http.Request, table *Table) Query {
	q, err := util.MapQuery(r)

	if err != nil || len(q) == 0 {
		return AllQuery()
	}

	return PrepareQuery(q, table)
}

func PrepareQuery(query map[string][]string, table *Table) Query {
	cols := table.DescribeColumns()
	result := NewQuery(len(cols))
	for i, col := range cols {
		val, ok := query[col.Name()]
		if ok && len(val) > 0 {
			if len(val) > 1 {
				result.ValuesAt(i, val)
			} else {
				result.ValueAt(i, val[0])
			}
		} else {
			result.AnyAt(i)
		}

	}
	return result.Finalize()
}

func AllQuery() Query {
	return Query{all: true}
}

func NewQuery(nCols int) *Query {
	return &Query{cols: make([]*Condition, nCols)}
}

func (q *Query) All() bool {
	return q.all
}

func (q *Query) AnyAt(i int) {
	q.cols[i] = &Condition{value: "*", any: true, multiple: false}
}

func (q *Query) ValueAt(i int, value string) {
	q.cols[i] = &Condition{value: value, any: false, multiple: false}
}

func (q *Query) ValuesAt(i int, values []string) {
	q.cols[i] = &Condition{values: values, any: false, multiple: true}
}

func (q *Query) Finalize() Query {
	if q.all {
		return *q
	}
	found := false
	for _, col := range q.cols {
		if !col.any {
			found = true
		}
	}
	if !found {
		q.all = true
	}
	return *q
}

func (q *Query) MatchesAt(i int, value string) bool {
	condition := q.cols[i]
	if condition.any {
		return true
	}

	if condition.multiple {
		for _, cond := range condition.values {
			if value == cond {
				return true
			}
		}
		return false
	}

	return condition.value == value
}

func (q *Query) Matches(row []Data) bool {
	for i, data := range row {
		if !q.MatchesAt(i, data.Str()) {
			return false
		}
	}

	return true
}
