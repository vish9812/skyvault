package common

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultPagingLimit = 100
	MaxPagingLimit     = 1000
	SortAsc            = "asc"
	SortDesc           = "desc"
	DefaultSort        = SortAsc
	SortByID           = "id"
	SortByName         = "name"
	SortByUpdated      = "updated"
	DefaultSortBy      = SortByName
	DirectionNext      = "next"
	DirectionPrev      = "prev"
	DefaultDirection   = DirectionNext
)

var (
	ErrInvalidCursor = errors.New("invalid cursor")
)

type PagingOptions struct {
	PrevCursor string
	NextCursor string
	Direction  string
	Limit      int
	Sort       string // "asc" or "desc"
	SortBy     string // "id", "name", "updated"
}

func (o *PagingOptions) Validate() {
	if o.Direction != DirectionPrev && o.Direction != DirectionNext {
		o.Direction = DefaultDirection
	}

	if o.Limit <= 0 {
		o.Limit = DefaultPagingLimit
	} else if o.Limit > MaxPagingLimit {
		o.Limit = MaxPagingLimit
	}

	if o.Sort != SortAsc && o.Sort != SortDesc {
		o.Sort = DefaultSort
	}

	if o.SortBy != SortByID && o.SortBy != SortByName && o.SortBy != SortByUpdated {
		o.SortBy = DefaultSortBy
	}
}

// Cursor holds the column values that make up the cursor.
type Cursor struct {
	ID      int64
	Name    string
	Updated time.Time
}

func (o *PagingOptions) GetCursor() (*Cursor, error) {
	cursorStr := o.NextCursor
	if o.Direction == DirectionPrev {
		cursorStr = o.PrevCursor
	}

	if cursorStr == "" {
		return nil, nil
	}

	if len(cursorStr) > 200 {
		return nil, fmt.Errorf("cursor too long: %w", ErrInvalidCursor)
	}

	vals := strings.Split(cursorStr, ",")

	id, err := strconv.ParseInt(vals[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrInvalidCursor)
	}

	cursor := &Cursor{ID: id}

	if o.SortBy == SortByID {
		if len(vals) != 1 {
			return nil, fmt.Errorf("id: %w", ErrInvalidCursor)
		}

		return cursor, nil
	}

	if len(vals) != 2 {
		return nil, fmt.Errorf("non-id: %w", ErrInvalidCursor)
	}

	if o.SortBy == SortByName {
		cursor.Name = vals[1]
	} else if o.SortBy == SortByUpdated {
		updated, err := time.Parse(time.RFC3339, vals[1])
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, ErrInvalidCursor)
		}

		cursor.Updated = updated
	}

	return cursor, nil
}

func (o *PagingOptions) CreateCursor(cursor *Cursor) string {
	if cursor == nil {
		return ""
	}

	if o.SortBy == SortByID {
		return fmt.Sprintf("%d", cursor.ID)
	}

	if o.SortBy == SortByName {
		return fmt.Sprintf("%d,%s", cursor.ID, cursor.Name)
	}

	return fmt.Sprintf("%d,%s", cursor.ID, cursor.Updated.Format(time.RFC3339))
}

type PagedItems[TItem any] struct {
	Items      []TItem `json:"items"`
	PrevCursor string  `json:"prevCursor"` // Cursor to the first item in the list
	NextCursor string  `json:"nextCursor"` // Cursor to the last item in the list
	HasMore    bool    `json:"hasMore"`
}

// ResetCursorIfNoMore resets either next or previous cursor, if there are no more items in the list.
func (pi *PagedItems[TItem]) ResetCursorIfNoMore(pagingOpt *PagingOptions) {
	if !pi.HasMore {
		if pagingOpt.Direction == DirectionNext {
			pi.NextCursor = ""
		} else {
			pi.PrevCursor = ""
		}
	}
}
