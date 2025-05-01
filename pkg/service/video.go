package service

import (
	"context"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository"
	"net/url"
)

type VideoService struct {
	metaRepo *repository.VideoMetaRepository
}

func NewVideoService(metaRepo *repository.VideoMetaRepository) *VideoService {
	return &VideoService{
		metaRepo: metaRepo,
	}
}

func (s *VideoService) CreateVideoEntry(ctx context.Context, vidMeta model.Video) (video *model.Video, url url.URL, err error) {
	s.metaRepo.CreateVideoMeta(ctx, vidMeta)
	return
}
