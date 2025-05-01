package service

import (
	"context"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository"
	"net/url"
)

type VideoService struct {
	metaRepo    *repository.VideoMetaRepository
	managerRepo *repository.VideoManagerRepository
}

func NewVideoService(metaRepo *repository.VideoMetaRepository, mngerRepo *repository.VideoManagerRepository) *VideoService {
	return &VideoService{
		metaRepo:    metaRepo,
		managerRepo: mngerRepo,
	}
}

func (s *VideoService) CreateVideoEntry(ctx context.Context, vidMeta model.Video) (video model.Video, url url.URL, err error) {
	video, err = s.metaRepo.CreateVideoMeta(ctx, vidMeta)

	//TODO: Add proper error handling for video creation
	if err != nil {
		return
	}

	// Disallow upload if the video is not in a pending state or if the retry count is greater than 3.
	if video.RetryCount > 3 || video.Status != model.VideoStatusPending {
		err = fluxerrors.ErrVideoUploadNotAllowed
		return
	}

	ptrURL, err := s.managerRepo.GenerateVideoUploadURL(video.ID, video.Slug)

	if err != nil {
		err = fluxerrors.ErrVideoURLGenerationFailed
		return
	}

	url = *ptrURL

	return
}
