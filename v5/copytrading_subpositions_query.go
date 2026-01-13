package okx

import (
	"net/url"
	"strconv"
)

type copyTradingSubpositionsQuery struct {
	instType   string
	uniqueCode string
	instId     string
	after      string
	before     string
	limit      *int
}

func (q copyTradingSubpositionsQuery) values() url.Values {
	v := url.Values{}
	if q.instType != "" {
		v.Set("instType", q.instType)
	}
	if q.uniqueCode != "" {
		v.Set("uniqueCode", q.uniqueCode)
	}
	if q.instId != "" {
		v.Set("instId", q.instId)
	}
	if q.after != "" {
		v.Set("after", q.after)
	}
	if q.before != "" {
		v.Set("before", q.before)
	}
	if q.limit != nil {
		v.Set("limit", strconv.Itoa(*q.limit))
	}

	if len(v) == 0 {
		return nil
	}
	return v
}
