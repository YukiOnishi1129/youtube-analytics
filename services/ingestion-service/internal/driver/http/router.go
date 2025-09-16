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

// SetupRouterWithKeyword creates a new router with keyword support
func SetupRouterWithKeyword(
	channelUseCase input.ChannelInputPort,
	videoUseCase input.VideoInputPort,
	systemUseCase input.SystemInputPort,
	keywordUseCase input.KeywordInputPort,
) *gin.Engine {
	router := gin.Default()

	server := NewServerWithKeyword(
		channelUseCase,
		videoUseCase,
		systemUseCase,
		keywordUseCase,
	)

	// Register handlers with generated server interface
	generated.RegisterHandlers(router, server)

	return router
}

// SetupRouterWithPresenters creates a new router with presenters for response formatting
func SetupRouterWithPresenters(
	channelUseCase input.ChannelInputPort,
	videoUseCase input.VideoInputPort,
	systemUseCase input.SystemInputPort,
	keywordUseCase input.KeywordInputPort,
	channelPresenter interface{},
	videoPresenter interface{},
	systemPresenter interface{},
	keywordPresenter interface{},
) *gin.Engine {
	router := gin.Default()

	// Create server with keyword support
	server := NewServerWithKeyword(
		channelUseCase,
		videoUseCase,
		systemUseCase,
		keywordUseCase,
	)

	// Register handlers with generated server interface
	generated.RegisterHandlers(router, server)

	return router
}
