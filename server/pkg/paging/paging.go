package paging

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	DefaultLimit     = 100
	MaxLimit         = 1000
	MaxCursorLen     = 200
	DirectionForward = "forward"
	SortAsc          = "asc"
	SortByID         = "id"
	SortByName       = "name"
	SortByUpdated    = "updated"
)

var (
	ErrInvalidCursor = errors.New("invalid cursor")
)

type Options struct {
	PrevCursor string
	NextCursor string
	Direction  string // "forward", "backward"
	Limit      int
	Sort       string // "asc" or "desc"
	SortBy     string // "id", "name", "updated"
}

func (o *Options) Validate() {
	if o.Limit <= 0 {
		o.Limit = DefaultLimit
	} else if o.Limit > MaxLimit {
		o.Limit = MaxLimit
	}

	if o.Direction != DirectionForward {
		o.Direction = "backward"
	}

	if o.Sort != SortAsc {
		o.Sort = "desc"
	}

	if o.SortBy != SortByID && o.SortBy != SortByName && o.SortBy != SortByUpdated {
		o.SortBy = SortByName
	}
}

// Cursor holds the column values that make up the cursor.
type Cursor struct {
	ID      string
	Name    string
	Updated time.Time
}

func (o *Options) CreateCursor(cursor *Cursor) string {
	if cursor == nil {
		return ""
	}

	if o.SortBy == SortByID {
		return fmt.Sprintf("%s", cursor.ID)
	}

	if o.SortBy == SortByName {
		return fmt.Sprintf("%s,%s", cursor.ID, cursor.Name)
	}

	return fmt.Sprintf("%s,%s", cursor.ID, cursor.Updated.Format(time.RFC3339))
}

func (o *Options) GetCursor() (*Cursor, error) {
	cursorStr := o.NextCursor
	if o.Direction != DirectionForward {
		cursorStr = o.PrevCursor
	}

	if cursorStr == "" {
		return nil, nil
	}

	if len(cursorStr) > MaxCursorLen {
		return nil, fmt.Errorf("%w: cursor too long: %d: %s", ErrInvalidCursor, len(cursorStr), cursorStr[:MaxCursorLen/2]+"...")
	}

	vals := strings.Split(cursorStr, ",")

	id := vals[0]

	cursor := &Cursor{ID: id}

	if o.SortBy == SortByID {
		if len(vals) != 1 {
			return nil, fmt.Errorf("%w: id", ErrInvalidCursor)
		}

		return cursor, nil
	}

	if len(vals) != 2 {
		return nil, fmt.Errorf("%w: non-id", ErrInvalidCursor)
	}

	if o.SortBy == SortByName {
		cursor.Name = vals[1]
	} else if o.SortBy == SortByUpdated {
		updated, err := time.Parse(time.RFC3339, vals[1])
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidCursor, err)
		}

		cursor.Updated = updated
	}

	return cursor, nil
}

type Page[TItem any] struct {
	Items      []TItem `json:"items"`
	PrevCursor string  `json:"prevCursor"` // Cursor to the first item in the list
	NextCursor string  `json:"nextCursor"` // Cursor to the last item in the list
	HasMore    bool    `json:"hasMore"`
}
