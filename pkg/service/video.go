package service

import (
	"context"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository"
	"net/url"
)

const (
	MAX_VIDEO_RETRY_COUNT = 4
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

	if err != nil {
		return
	}

	// Disallow upload if the video is not in a pending state or if the retry count is greater than 3.
	if video.RetryCount > MAX_VIDEO_RETRY_COUNT || video.Status != model.VideoStatusPending {
		err = fluxerrors.ErrVideoUploadNotAllowed
		return
	}

	// Generate the upload URL for the video.
	ptrURL, err := s.managerRepo.GenerateVideoUploadURL(ctx, video.ID, video.Slug)

	if err != nil {
		err = fluxerrors.ErrVideoURLGenerationFailed
		_ = s.metaRepo.IncrementVideoRetryCount(ctx, video.ID)

		// Prevent returning the video if the url generation fails.
		video = model.Video{}
		return
	}

	url = *ptrURL

	return
}
