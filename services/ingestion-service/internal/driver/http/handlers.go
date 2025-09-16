package http

import (
	"net/http"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/http/generated"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Server struct {
	channelUseCase input.ChannelInputPort
	videoUseCase   input.VideoInputPort
	systemUseCase  input.SystemInputPort
	keywordUseCase input.KeywordInputPort
}

func NewServer(
	channelUseCase input.ChannelInputPort,
	videoUseCase input.VideoInputPort,
	systemUseCase input.SystemInputPort,
) *Server {
	return &Server{
		channelUseCase: channelUseCase,
		videoUseCase:   videoUseCase,
		systemUseCase:  systemUseCase,
	}
}

// NewServerWithKeyword creates a new server with keyword support
func NewServerWithKeyword(
	channelUseCase input.ChannelInputPort,
	videoUseCase input.VideoInputPort,
	systemUseCase input.SystemInputPort,
	keywordUseCase input.KeywordInputPort,
) *Server {
	return &Server{
		channelUseCase: channelUseCase,
		videoUseCase:   videoUseCase,
		systemUseCase:  systemUseCase,
		keywordUseCase: keywordUseCase,
	}
}

func (s *Server) AdminCollectSubscriptions(c *gin.Context, params generated.AdminCollectSubscriptionsParams) {
	start := time.Now()

	result, err := s.videoUseCase.CollectSubscriptions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, generated.Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, generated.CollectSubscriptionsResponse{
		ChannelsProcessed: int32(result.ChannelsProcessed),
		VideosCollected:   int32(result.VideosCollected),
		VideosCreated:     int32(result.VideosCreated),
		Duration:          time.Since(start).String(),
	})
}

func (s *Server) AdminCollectTrending(c *gin.Context, params generated.AdminCollectTrendingParams) {
	start := time.Now()

	// TODO: Add genre_id parameter to OpenAPI spec and regenerate
	// For now, collect trending for all genres
	result, err := s.videoUseCase.CollectTrending(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, generated.Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, generated.CollectTrendingResponse{
		VideosCollected: int32(result.VideosCollected),
		VideosCreated:   int32(result.VideosCreated),
		VideosUpdated:   int32(result.VideosUpdated),
		Duration:        time.Since(start).String(),
	})
}

func (s *Server) AdminScheduleSnapshots(c *gin.Context, params generated.AdminScheduleSnapshotsParams) {
	start := time.Now()

	result, err := s.systemUseCase.ScheduleSnapshots(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, generated.Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, generated.ScheduleSnapshotsResponse{
		VideosProcessed: int32(result.VideosProcessed),
		TasksScheduled:  int32(result.TasksScheduled),
		Duration:        time.Since(start).String(),
	})
}

func (s *Server) AdminUpdateChannels(c *gin.Context, params generated.AdminUpdateChannelsParams) {
	start := time.Now()

	result, err := s.channelUseCase.UpdateChannels(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, generated.Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, generated.UpdateChannelsResponse{
		ChannelsProcessed: int32(result.ChannelsProcessed),
		ChannelsUpdated:   int32(result.ChannelsUpdated),
		Duration:          time.Since(start).String(),
	})
}

func (s *Server) TasksCreateSnapshot(c *gin.Context, params generated.TasksCreateSnapshotParams) {
	var req generated.CreateSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, generated.Error{
			Code:    "INVALID_REQUEST",
			Message: err.Error(),
		})
		return
	}

	// Parse video ID
	videoID, err := uuid.Parse(req.VideoId)
	if err != nil {
		c.JSON(http.StatusBadRequest, generated.Error{
			Code:    "INVALID_VIDEO_ID",
			Message: "Invalid video ID format",
		})
		return
	}

	// Create snapshot
	_, err = s.systemUseCase.CreateSnapshot(c.Request.Context(), &input.CreateSnapshotInput{
		VideoID:        videoID,
		CheckpointHour: int(req.CheckpointHour),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, generated.Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) WebSubVerify(c *gin.Context, params generated.WebSubVerifyParams) {
	// Verify the subscription request
	if params.HubMode == "subscribe" || params.HubMode == "unsubscribe" {
		// Return the challenge to confirm the subscription
		c.String(http.StatusOK, params.HubChallenge)
		return
	}

	c.Status(http.StatusNotFound)
}

func (s *Server) WebSubNotify(c *gin.Context, params generated.WebSubNotifyParams) {
	// TODO: Implement WebSub notification handling
	// This would parse the YouTube push notification and trigger appropriate actions
	// For now, just acknowledge receipt
	c.Status(http.StatusOK)
}
