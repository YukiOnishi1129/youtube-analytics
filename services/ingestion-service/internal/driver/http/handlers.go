package http

import (
	"net/http"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/http/generated"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/gin-gonic/gin"
)

type Server struct {
	channelUseCase input.ChannelInputPort
	videoUseCase   input.VideoInputPort
	systemUseCase  input.SystemInputPort
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
		VideosAdded:       int32(result.VideosAdded),
		VideosProcessed:   int32(result.VideosProcessed),
		Duration:          time.Since(start).String(),
	})
}

func (s *Server) AdminCollectTrending(c *gin.Context, params generated.AdminCollectTrendingParams) {
	start := time.Now()

	result, err := s.videoUseCase.CollectTrending(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, generated.Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, generated.CollectTrendingResponse{
		VideosAdded:     int32(result.VideosAdded),
		VideosProcessed: int32(result.VideosProcessed),
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
