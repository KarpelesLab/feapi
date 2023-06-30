package feapi

import (
	"bytes"
	"context"

	"github.com/KarpelesLab/pjson"
)

type Sortable []*SortableValue

type SortableValue struct {
	V any
	K []byte
}

func (s *SortableValue) MarshalJSON() ([]byte, error) {
	return pjson.Marshal(s.V)
}

func (s *SortableValue) MarshalContextJSON(ctx context.Context) ([]byte, error) {
	return pjson.MarshalContext(ctx, s.V)
}

func (s Sortable) Len() int {
	return len(s)
}

func (s Sortable) Less(i, j int) bool {
	return bytes.Compare(s[i].K, s[j].K) < 0
}

func (s Sortable) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
