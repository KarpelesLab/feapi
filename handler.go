package feapi

import (
	"context"

	"github.com/KarpelesLab/apirouter"
)

func Handle(ctx context.Context, path, verb string) (any, error) {
	c := apirouter.New(ctx, path, verb)
	return c.Call()
}
