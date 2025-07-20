package components

import (
	"github.com/kevin-fagan/go-links/internal/db"
	"github.com/kevin-fagan/go-links/internal/links"
	"github.com/kevin-fagan/go-links/internal/logs"
	"github.com/kevin-fagan/go-links/internal/tags"
)

type Service struct {
	tags  tags.Repository
	logs  logs.Repository
	links links.Repository
}

func NewService(ctx *db.SQLiteContext) *Service {
	return &Service{
		tags:  *tags.NewRepository(ctx),
		logs:  *logs.NewRepository(ctx),
		links: *links.NewRepository(ctx),
	}
}
