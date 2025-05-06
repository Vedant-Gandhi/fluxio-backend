package service

import (
	"context"
	"fluxio-backend/pkg/constants"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository"
	"net/url"
	"strings"
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
	if video.RetryCount > constants.MaxVideoURLRegenerateRetryCount || video.Status != model.VideoStatusPending {
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

// Handles the meta update after the video file is uploaded.
func (s *VideoService) UpdateUploadStatus(ctx context.Context, id model.VideoID, storagePath string) (err error) {
	if strings.EqualFold(storagePath, "") {
		err = fluxerrors.ErrMalformedStoragePath
		return
	}

	if strings.EqualFold(id.String(), "") {
		err = fluxerrors.ErrInvalidVideoID
		return
	}

	existData, err := s.metaRepo.GetProcessingDetailsByID(ctx, id)

	if err != nil {
		if err == fluxerrors.ErrVideoNotFound {
			return
		}

		err = fluxerrors.ErrUnknown
		return
	}

	// If video status is not pending or retry limit is over end it.
	if existData.Status != model.VideoStatusPending || existData.RetryCount > constants.MaxVideoURLRegenerateRetryCount {
		err = fluxerrors.ErrInvalidVideoStatus
		return
	}

	err = s.metaRepo.UpdateProcessingDetails(ctx, id, model.VideoStatusProcessing, storagePath)

	if err != nil {
		if err == fluxerrors.ErrInvalidVideoID || err == fluxerrors.ErrMalformedStoragePath {
			return
		}

		err = fluxerrors.ErrVideoMetaUpdateFailed
		return
	}

	return
}
