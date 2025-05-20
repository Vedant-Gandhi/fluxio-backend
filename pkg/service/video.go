package service

import (
	"context"
	"encoding/json"
	"fluxio-backend/pkg/common/schema"
	"fluxio-backend/pkg/constants"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository"
	"fluxio-backend/pkg/utils"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type VideoService struct {
	videRepo *repository.VideoRepository
	l        schema.Logger
}

func NewVideoService(videRepo *repository.VideoRepository, logger schema.Logger) *VideoService {
	return &VideoService{
		videRepo: videRepo,
		l:        logger,
	}
}

func (s *VideoService) AddVideo(ctx context.Context, vidMeta model.Video, mimeType string) (video model.Video, url url.URL, err error) {
	logger := s.l.With("title", vidMeta.Title)

	if !utils.CheckVideoMimeTypeValidity(mimeType) {
		err = fluxerrors.ErrInvalidVideoExtension
		logger.Error("Video Format is required", err)
		return
	}

	splitMime := strings.SplitN(mimeType, "/", 2)

	// Set the format type.
	vidMeta.Format = splitMime[1]

	video, err = s.videRepo.CreateVideoMeta(ctx, vidMeta)
	if err != nil {
		logger.Error("Failed to create video metadata", err)
		return
	}

	logger = logger.With("video_id", video.ID.String())

	// Disallow upload if the video is not in a pending state or if the retry count is greater than 3.
	if video.RetryCount > constants.MaxVideoURLRegenerateRetryCount || video.Status != model.VideoStatusUploadPending {
		err = fluxerrors.ErrVideoUploadNotAllowed
		logger.Error("Video upload not allowed", err)
		return
	}

	// Generate the upload URL for the video.
	ptrURL, err := s.videRepo.GenerateUnProcessedVideoUploadURL(ctx, video.ID, video.Slug, mimeType)
	if err != nil {
		err = fluxerrors.ErrVideoURLGenerationFailed
		logger.Error("Failed to generate video upload URL", err)
		_ = s.videRepo.IncrementVideoRetryCount(ctx, video.ID)

		// Prevent returning the video if the url generation fails.
		video = model.Video{}
		return
	}

	url = *ptrURL
	logger.Info("Video metadata created successfully", "slug", video.Slug)
	return
}

// Handles the meta update after the video file is uploaded.
func (s *VideoService) UpdateUploadStatus(ctx context.Context, slug string, params model.Video) (err error) {
	logger := s.l.With("slug", slug)

	if strings.EqualFold(params.StoragePath, "") {
		err = fluxerrors.ErrMalformedStoragePath
		logger.Error("Invalid storage path", err)
		return
	}

	if strings.EqualFold(slug, "") {
		err = fluxerrors.ErrInvalidVideoSlug
		logger.Error("Invalid video slug", err)
		return
	}

	existData, err := s.videRepo.GetProcessingDetailsBySlug(ctx, slug)
	if err != nil {
		if err == fluxerrors.ErrVideoNotFound {
			logger.Error("Video not found", err)
			return
		}

		logger.Error("Failed to get video processing details", err)
		err = fluxerrors.ErrUnknown
		return
	}

	logger = logger.With("video_id", existData.ID.String())

	// If video status is not pending or retry limit is over end it.
	if existData.Status != model.VideoStatusUploadPending || existData.RetryCount > constants.MaxVideoURLRegenerateRetryCount {
		err = fluxerrors.ErrInvalidVideoStatus
		logger.Error("Invalid video status for update", err)
		return
	}

	err = s.videRepo.UpdateMeta(ctx, existData.ID, model.VideoStatusProcessing, model.Video{
		StoragePath: params.StoragePath,
	})

	if err != nil {
		if err == fluxerrors.ErrInvalidVideoID || err == fluxerrors.ErrMalformedStoragePath {
			logger.Error("Failed to update video metadata", err)
			return
		}

		logger.Error("Video metadata update failed", err)
		err = fluxerrors.ErrVideoMetaUpdateFailed
		return
	}

	logger.Info("Video upload status updated to processing")
	return
}

// Performs all the post upload processing for the video.
func (s *VideoService) PerformPostUploadProcessing(ctx context.Context, slug string) (err error) {
	logger := s.l.With("slug", slug)
	logger.Info("Starting post-upload processing")

	videoMeta, err := s.videRepo.GetVideoBySlug(ctx, slug)
	if err != nil {
		if err == fluxerrors.ErrVideoNotFound {
			logger.Error("Video not found for processing", err)
			return
		}
		logger.Error("Failed to get video by slug", err)
		err = fluxerrors.ErrUnknown
		return
	}

	logger = logger.With("video_id", videoMeta.ID.String())

	if videoMeta.Status != model.VideoStatusProcessing {
		err = fluxerrors.ErrInvalidVideoStatus
		logger.Error("Invalid video status for processing", err)
		return
	}

	downloadURL, err := s.videRepo.GetUnProcessedVideoDownloadURL(ctx, videoMeta.Slug)
	if err != nil {
		if err == fluxerrors.ErrVideoURLGenerationFailed {
			logger.Error("Failed to generate video download URL", err)
			return
		}
		logger.Error("Unknown error getting video download URL", err)
		err = fluxerrors.ErrUnknown
		return
	}

	// Extract the whole video meta data like size, type, width, height,etc
	logger.Info("Extracting video metadata")
	rawProbe, err := ffmpeg_go.Probe(downloadURL.String())
	if err != nil {
		logger.Error("Failed to probe video using ffmpeg", err)
		err = fluxerrors.ErrVideoPhysicalMetaExtractionFailed
		return
	}

	var probe model.FFProbeOutput
	err = json.Unmarshal([]byte(rawProbe), &probe)
	if err != nil {
		logger.Error("Failed to parse ffprobe output", err)
		err = fluxerrors.ErrVideoPhysicalMetaExtractionFailed
		return
	}

	if probe.Format.NbStreams != 2 {
		logger.Error("Unsupported video stream count", nil)
		err = fluxerrors.ErrVideoStreamCountNotSupported
		return
	}

	videoStream := model.FFProbeStream{}
	audioStream := model.FFProbeStream{}

	// Get the streams from the probe
	if probe.Streams[0].CodecType == "video" {
		videoStream = probe.Streams[0]
		audioStream = probe.Streams[1]
	} else {
		audioStream = probe.Streams[0]
		videoStream = probe.Streams[1]
	}

	updateData := model.Video{}
	updateData.AudioCodec = audioStream.CodecName

	sampleRate, err := strconv.Atoi(audioStream.SampleRate)
	if err != nil {
		logger.Error("Failed to parse audio sample rate", err)
		err = fluxerrors.ErrVideoPhysicalMetaExtractionFailed
		return
	}

	updateData.AudioSampleRate = uint32(sampleRate)
	updateData.Width = uint32(videoStream.Width)
	updateData.Height = uint32(videoStream.Height)
	updateData.Format = videoStream.CodecName

	duration, err := strconv.ParseFloat(probe.Format.Duration, 64)
	if err != nil {
		logger.Error("Failed to parse video duration", err)
		err = fluxerrors.ErrVideoPhysicalMetaExtractionFailed
		return
	}

	updateData.Length = uint64(math.Ceil(duration))

	size, err := strconv.ParseFloat(probe.Format.Size, 64)
	if err != nil {
		logger.Error("Failed to parse video size", err)
		err = fluxerrors.ErrVideoPhysicalMetaExtractionFailed
		return
	}

	calcPrec := math.Pow(10, float64(constants.VidSizeDecimalPrecision)) // Stores the power precision to round the size.
	updateData.Size = float32(math.Round(size*calcPrec) / calcPrec)      // Round the size to decimal places.

	err = s.videRepo.UpdateInternalStatus(ctx, videoMeta.ID, model.VidInternalStatusMetaExtracted)
	if err != nil {
		if err == fluxerrors.ErrVideoNotFound {
			logger.Error("Video not found when updating status", err)
			return
		}
		logger.Error("Failed to update video internal status", err)
		err = fluxerrors.ErrVideoMetaUpdateFailed
		return
	}

	// Create thumbnails for the video and store them in the db
	logger.Info("Starting thumbnail generation")
	thumbnailTempDir, err := os.MkdirTemp(os.TempDir(), "fluxio-thumbnails-*")
	if err != nil {
		logger.Error("Failed to create temporary directory for thumbnails", err)
		err = fluxerrors.ErrThumbnailCreationFailed
		return
	}

	defer os.RemoveAll(thumbnailTempDir)

	thumbnailWidth := 1280
	thumbnailHeight := 720
	thumbnailFormat := "jpg"

	timestamps := s.generateDistinctTimestamps(updateData.Length)

	successThumbnailCount := 0
	client := &http.Client{}

	// Generate three thumbnails
	for idx, timestamp := range timestamps {
		// We need to convert the timestamp to ffmpeg format of HH:MM:SS
		timestampSeconds := timestamp
		hours := timestampSeconds / 3600
		minutes := (timestampSeconds % 3600) / 60
		seconds := timestampSeconds % 60
		timeStr := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

		opPath := path.Join(thumbnailTempDir, fmt.Sprintf("%s-%s.%s", videoMeta.Slug, fmt.Sprint(timestamp), thumbnailFormat))

		// We pass the URL so the ffmpeg will smartly use HTTP Range requests to get the exact frame.
		err = ffmpeg_go.Input(downloadURL.String(), ffmpeg_go.KwArgs{
			"ss":      timeStr, // The Timestamp to extract the thumbnail from
			"y":       "",      // Overwrite the output file if exists
			"timeout": "40",    // Timeout for whole op execution
		}).Output(opPath, ffmpeg_go.KwArgs{
			"vframes": 1,                                                                                                                                                                                // How many frames to output
			"s":       fmt.Sprintf("%dx%d", thumbnailWidth, thumbnailHeight),                                                                                                                            // Pass the thumbnail dimensions here
			"q:v":     3,                                                                                                                                                                                // Quality of the thumbnail
			"vf":      fmt.Sprintf("thumbnail,scale=w=%[1]s:h=%[2]s:force_original_aspect_ratio=decrease,pad=%[1]s:%[2]s:(ow-iw)/2:(oh-ih)/2", fmt.Sprint(thumbnailWidth), fmt.Sprint(thumbnailHeight)), // Apply thumbnail filter, scale, maintain aspect ratio
		}).OverWriteOutput().Run()

		// Perform cleanup of the temporary file
		defer os.Remove(opPath)

		if err != nil {
			// Silence the error if ffmpeg fails to generate the thumbnail.
			err = nil
			continue
		}

		fileStat, err := os.Stat(opPath)
		if err != nil {
			// If the file does not exist, we can skip this thumbnail.
			err = nil
			continue
		}

		thumbnail := model.Thumbnail{
			VideoID:   videoMeta.ID,
			Width:     uint16(thumbnailWidth),
			Height:    uint16(thumbnailHeight),
			Size:      uint32(fileStat.Size() / 1024), // Size in KB
			Format:    thumbnailFormat,
			TimeStamp: timestamp,
			IsDefault: idx == 0, // Set the first thumbnail as default
		}

		url, err := s.videRepo.GenerateThumbnailUploadURL(ctx, thumbnail.VideoID, thumbnail.TimeStamp, thumbnailFormat)
		if err != nil {
			err = nil
			continue
		}

		thumbFile, err := os.Open(opPath)
		if err != nil {
			err = nil
			continue
		}

		defer thumbFile.Close()

		uploadReq, err := http.NewRequest(http.MethodPut, url.String(), thumbFile)
		if err != nil {
			err = nil
			continue
		}

		uploadReq.Header.Set("Content-Type", fmt.Sprintf("image/%s", thumbnailFormat))
		uploadReq.ContentLength = fileStat.Size()

		resp, err := client.Do(uploadReq)
		if err != nil {
			err = nil
			continue
		}

		if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) {
			resp.Body.Close()
			continue
		}

		resp.Body.Close()

		thumbnail.StoragePath = fmt.Sprintf("%s.%s", utils.CreateURLSafeThumbnailFileName(thumbnail.VideoID.String(), fmt.Sprint(thumbnail.TimeStamp)), thumbnailFormat)

		_, err = s.videRepo.CreateThumbnail(ctx, thumbnail)
		if err != nil {
			err = nil // Ignore the error if thumbnail creation fails.
			continue
		}

		successThumbnailCount++
	}

	logger.Info("Video processing completed successfully", "thumbnails_created", successThumbnailCount)
	return
}

func (v *VideoService) generateDistinctTimestamps(videoDuration uint64) []uint64 {
	if videoDuration < 4 {
		if videoDuration < 2 {
			return []uint64{videoDuration}
		}
		step := videoDuration / 3
		return []uint64{step, 2 * step, videoDuration}
	}

	// Divide video into 3 segments and pick a random time from each segment
	segmentDuration := videoDuration / 4 // Use 4 segments to avoid the very beginning and end

	timestamps := make([]uint64, 3)

	// First thumbnail from first quarter (excluding first 5% of video)
	minTime1 := uint64(float64(videoDuration) * 0.05)
	maxTime1 := segmentDuration
	timestamps[0] = minTime1 + uint64(rand.Int63n(int64(maxTime1-minTime1)))

	// Second thumbnail from middle section
	minTime2 := segmentDuration * 1
	maxTime2 := segmentDuration * 2
	timestamps[1] = minTime2 + uint64(rand.Int63n(int64(maxTime2-minTime2)))

	// Third thumbnail from later section (avoiding last 5% of video)
	minTime3 := segmentDuration * 2
	maxTime3 := uint64(float64(videoDuration) * 0.95)
	timestamps[2] = minTime3 + uint64(rand.Int63n(int64(maxTime3-minTime3)))

	return timestamps
}
