package service

import "github.com/kevin-fagan/learn-gin/internal/repository"

type ComponentService struct {
	linkRepository repository.LinkRepository
}

func NewComponentService(ctx *repository.SQLContext) *LinkService {
	return &LinkService{
		linkRepository: *repository.NewLinkRepository(ctx),
	}
}
