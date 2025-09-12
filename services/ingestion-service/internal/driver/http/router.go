package http

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/http/generated"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	channelUseCase input.ChannelInputPort,
	videoUseCase input.VideoInputPort,
	systemUseCase input.SystemInputPort,
) *gin.Engine {
	router := gin.Default()

	server := NewServer(channelUseCase, videoUseCase, systemUseCase)

	// Register handlers with generated server interface
	generated.RegisterHandlers(router, server)

	return router
}
