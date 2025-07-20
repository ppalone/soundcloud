package soundcloud

import "strconv"

type searchOptions struct {
	q      string
	limit  int
	offset int
}

type SearchOption func(o *searchOptions)

func defaultSearchOptions() *searchOptions {
	return &searchOptions{
		q:      "",
		limit:  20,
		offset: 0,
	}
}

func (o *searchOptions) build() map[string]string {
	p := make(map[string]string)
	p["q"] = o.q
	p["limit"] = strconv.Itoa(o.limit)
	p["offset"] = strconv.Itoa(o.offset)

	return p
}

func WithLimit(limit int) SearchOption {
	return func(o *searchOptions) {
		o.limit = limit
	}
}

func WithOffset(offset int) SearchOption {
	return func(o *searchOptions) {
		o.offset = offset
	}
}
