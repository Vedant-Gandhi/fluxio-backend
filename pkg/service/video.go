package service

import (
	"context"
	"encoding/json"
	"fluxio-backend/pkg/constants"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository"
	"math"
	"net/url"
	"strconv"
	"strings"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type VideoService struct {
	videRepo *repository.VideoRepository
}

func NewVideoService(videRepo *repository.VideoRepository) *VideoService {
	return &VideoService{
		videRepo: videRepo,
	}
}

func (s *VideoService) CreateVideoEntry(ctx context.Context, vidMeta model.Video) (video model.Video, url url.URL, err error) {
	video, err = s.videRepo.CreateVideoMeta(ctx, vidMeta)

	if err != nil {
		return
	}

	// Disallow upload if the video is not in a pending state or if the retry count is greater than 3.
	if video.RetryCount > constants.MaxVideoURLRegenerateRetryCount || video.Status != model.VideoStatusPending {
		err = fluxerrors.ErrVideoUploadNotAllowed
		return
	}

	// Generate the upload URL for the video.
	ptrURL, err := s.videRepo.GenerateVideoUploadURL(ctx, video.ID, video.Slug)

	if err != nil {
		err = fluxerrors.ErrVideoURLGenerationFailed
		_ = s.videRepo.IncrementVideoRetryCount(ctx, video.ID)

		// Prevent returning the video if the url generation fails.
		video = model.Video{}
		return
	}

	url = *ptrURL

	return
}

// Handles the meta update after the video file is uploaded.
func (s *VideoService) UpdateUploadStatus(ctx context.Context, slug string, params model.UpdateVideoMeta) (err error) {
	if strings.EqualFold(params.StoragePath, "") {
		err = fluxerrors.ErrMalformedStoragePath
		return
	}

	if strings.EqualFold(slug, "") {
		err = fluxerrors.ErrInvalidVideoSlug
		return
	}

	existData, err := s.videRepo.GetProcessingDetailsBySlug(ctx, slug)

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

	err = s.videRepo.UpdateMeta(ctx, existData.ID, model.VideoStatusProcessing, model.UpdateVideoMeta{
		StoragePath: params.StoragePath,
	})

	if err != nil {
		if err == fluxerrors.ErrInvalidVideoID || err == fluxerrors.ErrMalformedStoragePath {
			return
		}

		err = fluxerrors.ErrVideoMetaUpdateFailed
		return
	}

	return
}

// Performs all the post upload processing for the video.
func (s *VideoService) PerformPostUploadProcessing(ctx context.Context, slug string) (err error) {

	videoMeta, err := s.videRepo.GetVideoBySlug(ctx, slug)

	if err != nil {

		if err == fluxerrors.ErrVideoNotFound {
			return
		}

		err = fluxerrors.ErrUnknown
		return
	}

	downloadURL, err := s.videRepo.GetVideoTemporaryDownloadURL(ctx, videoMeta.Slug)

	if err != nil {

		if err == fluxerrors.ErrVideoURLGenerationFailed {
			return
		}

		err = fluxerrors.ErrUnknown
		return
	}

	// Extract the whole video meta data like size, type, width, height,etc
	rawProbe, err := ffmpeg_go.Probe(downloadURL.String())

	if err != nil {

		err = fluxerrors.ErrVideoPhysicalMetaExtractionFailed
		return
	}

	var probe model.FFProbeOutput

	err = json.Unmarshal([]byte(rawProbe), &probe)

	if err != nil {

		err = fluxerrors.ErrVideoPhysicalMetaExtractionFailed
		return
	}

	if probe.Format.NbStreams != 2 {
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

	updateData := model.UpdateVideoMeta{}

	updateData.AudioCodec = audioStream.CodecName

	sampleRate, err := strconv.Atoi(audioStream.SampleRate)

	if err != nil {
		err = fluxerrors.ErrVideoPhysicalMetaExtractionFailed
		return
	}

	updateData.AudioSampleRate = uint32(sampleRate)

	updateData.Width = uint32(videoStream.Width)
	updateData.Height = uint32(videoStream.Height)
	updateData.Format = videoStream.CodecName

	duration, err := strconv.ParseFloat(probe.Format.Duration, 64)

	if err != nil {
		err = fluxerrors.ErrVideoPhysicalMetaExtractionFailed
		return
	}

	updateData.Length = uint64(math.Ceil(duration))

	size, err := strconv.ParseFloat(probe.Format.Size, 64)

	if err != nil {
		err = fluxerrors.ErrVideoPhysicalMetaExtractionFailed
		return
	}

	calcPrec := math.Pow(10, float64(constants.VidSizeDecimalPrecision)) // Stores the power precision to round the size.
	updateData.Size = float32(math.Round(size*calcPrec) / calcPrec)      // Round the size to decimal places.

	err = s.videRepo.UpdateMeta(ctx, videoMeta.ID, model.VideoStatusMetaExtracted, updateData)
	if err != nil {
		if err == fluxerrors.ErrVideoNotFound {
			return
		}

		err = fluxerrors.ErrVideoMetaUpdateFailed
		return
	}

	// Create thumbnails for the video and store them in the db
	return
}

/**
{
  "streams": [
    {
      "index": 0,
      "codec_name": "h264",
      "codec_long_name": "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
      "profile": "Constrained Baseline",
      "codec_type": "video",
      "codec_tag_string": "avc1",
      "codec_tag": "0x31637661",
      "width": 480,
      "height": 270,
      "coded_width": 480,
      "coded_height": 270,
      "closed_captions": 0,
      "has_b_frames": 0,
      "sample_aspect_ratio": "1:1",
      "display_aspect_ratio": "16:9",
      "pix_fmt": "yuv420p",
      "level": 30,
      "chroma_location": "left",
      "refs": 1,
      "is_avc": "true",
      "nal_length_size": "4",
      "r_frame_rate": "30/1",
      "avg_frame_rate": "30/1",
      "time_base": "1/30",
      "start_pts": 0,
      "start_time": "0.000000",
      "duration_ts": 901,
      "duration": "30.033333",
      "bit_rate": "301201",
      "bits_per_raw_sample": "8",
      "nb_frames": "901",
      "disposition": {
        "default": 1,
        "dub": 0,
        "original": 0,
        "comment": 0,
        "lyrics": 0,
        "karaoke": 0,
        "forced": 0,
        "hearing_impaired": 0,
        "visual_impaired": 0,
        "clean_effects": 0,
        "attached_pic": 0,
        "timed_thumbnails": 0
      },
      "tags": {
        "creation_time": "2015-08-07T09:13:02.000000Z",
        "language": "und",
        "handler_name": "L-SMASH Video Handler",
        "vendor_id": "[0][0][0][0]",
        "encoder": "AVC Coding"
      }
    },
    {
      "index": 1,
      "codec_name": "aac",
      "codec_long_name": "AAC (Advanced Audio Coding)",
      "profile": "LC",
      "codec_type": "audio",
      "codec_tag_string": "mp4a",
      "codec_tag": "0x6134706d",
      "sample_fmt": "fltp",
      "sample_rate": "48000",
      "channels": 2,
      "channel_layout": "stereo",
      "bits_per_sample": 0,
      "r_frame_rate": "0/0",
      "avg_frame_rate": "0/0",
      "time_base": "1/48000",
      "start_pts": 0,
      "start_time": "0.000000",
      "duration_ts": 1465280,
      "duration": "30.526667",
      "bit_rate": "112000",
      "nb_frames": "1431",
      "disposition": {
        "default": 1,
        "dub": 0,
        "original": 0,
        "comment": 0,
        "lyrics": 0,
        "karaoke": 0,
        "forced": 0,
        "hearing_impaired": 0,
        "visual_impaired": 0,
        "clean_effects": 0,
        "attached_pic": 0,
        "timed_thumbnails": 0
      },
      "tags": {
        "creation_time": "2015-08-07T09:13:02.000000Z",
        "language": "und",
        "handler_name": "L-SMASH Audio Handler",
        "vendor_id": "[0][0][0][0]"
      }
    }
  ],
  "format": {
    "filename": "http://127.0.0.1:9000/videos/new-with-slu1g-adc53916?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=LAiAfReWLSdL4OiKoKW6%2F20250509%2Fsa-east-1%2Fs3%2Faws4_request&X-Amz-Date=20250509T124227Z&X-Amz-Expires=3600&X-Amz-SignedHeaders=host&X-Amz-Signature=e5b9c3166090a83ec1773ce5fe1e0ff0c73f15f2ad2fdbb912a284950b884897",
    "nb_streams": 2,
    "nb_programs": 0,
    "format_name": "mov,mp4,m4a,3gp,3g2,mj2",
    "format_long_name": "QuickTime / MOV",
    "start_time": "0.000000",
    "duration": "30.526667",
    "size": "1570024",
    "bit_rate": "411449",
    "probe_score": 100,
    "tags": {
      "major_brand": "mp42",
      "minor_version": "0",
      "compatible_brands": "mp42mp41isomavc1",
      "creation_time": "2015-08-07T09:13:02.000000Z"
    }
  }
}
*/
