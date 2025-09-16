package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/YukiOnishi1129/youtube-analytics/services/pkg/pb/ingestion/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the gRPC server for ingestion service
type Server struct {
	pb.UnimplementedIngestionServiceServer
	channelUseCase         input.ChannelInputPort
	videoUseCase           input.VideoInputPort
	systemUseCase          input.SystemInputPort
	keywordUseCase         input.KeywordInputPort         // Optional keyword use case
	genreUseCase           input.GenreInputPort           // Genre use case
	youtubeCategoryUseCase input.YouTubeCategoryInputPort // YouTube category use case
	videoGenreUseCase      input.VideoGenreInputPort      // Video-Genre use case
	auditLogUseCase        input.AuditLogInputPort        // Audit log use case
	batchJobUseCase        input.BatchJobInputPort        // Batch job use case
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

// NewServerWithAllUseCases creates a new server with all use cases
func NewServerWithAllUseCases(
	channelUseCase input.ChannelInputPort,
	videoUseCase input.VideoInputPort,
	systemUseCase input.SystemInputPort,
	keywordUseCase input.KeywordInputPort,
	genreUseCase input.GenreInputPort,
	youtubeCategoryUseCase input.YouTubeCategoryInputPort,
	videoGenreUseCase input.VideoGenreInputPort,
	auditLogUseCase input.AuditLogInputPort,
	batchJobUseCase input.BatchJobInputPort,
) *Server {
	return &Server{
		channelUseCase:         channelUseCase,
		videoUseCase:           videoUseCase,
		systemUseCase:          systemUseCase,
		keywordUseCase:         keywordUseCase,
		genreUseCase:           genreUseCase,
		youtubeCategoryUseCase: youtubeCategoryUseCase,
		videoGenreUseCase:      videoGenreUseCase,
		auditLogUseCase:        auditLogUseCase,
		batchJobUseCase:        batchJobUseCase,
	}
}

// Start starts the gRPC server
func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterIngestionServiceServer(grpcServer, s)

	log.Printf("gRPC server starting on port %d", port)
	return grpcServer.Serve(lis)
}

// GetChannel gets a channel by ID
func (s *Server) GetChannel(ctx context.Context, req *pb.GetChannelRequest) (*pb.GetChannelResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "channel id is required")
	}

	channelID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid channel id format")
	}

	channel, err := s.channelUseCase.GetChannel(ctx, channelID)
	if err != nil {
		if err == domain.ErrChannelNotFound {
			return nil, status.Error(codes.NotFound, "channel not found")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get channel: %v", err))
	}

	return &pb.GetChannelResponse{
		Channel: domainChannelToProto(channel),
	}, nil
}

func (s *Server) ListChannels(ctx context.Context, req *pb.ListChannelsRequest) (*pb.ListChannelsResponse, error) {
	channels, err := s.channelUseCase.ListChannels(ctx, req.SubscribedOnly)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list channels: %v", err))
	}

	protoChannels := make([]*pb.Channel, len(channels))
	for i, channel := range channels {
		protoChannels[i] = domainChannelToProto(channel)
	}

	return &pb.ListChannelsResponse{
		Channels: protoChannels,
	}, nil
}

func (s *Server) CollectTrending(ctx context.Context, req *pb.CollectTrendingRequest) (*pb.CollectTrendingResponse, error) {
	var genreID *string
	if req.GenreId != "" {
		genreID = &req.GenreId
	}

	result, err := s.videoUseCase.CollectTrending(ctx, genreID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to collect trending videos: %v", err))
	}

	return &pb.CollectTrendingResponse{
		VideosProcessed: int32(result.VideosCollected),
		VideosAdded:     int32(result.VideosCreated),
		DurationMs:      result.Duration.Milliseconds(),
		GenreCode:       result.GenreCode,
	}, nil
}

func (s *Server) CollectSubscriptions(ctx context.Context, req *pb.CollectSubscriptionsRequest) (*pb.CollectSubscriptionsResponse, error) {
	result, err := s.videoUseCase.CollectSubscriptions(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to collect subscription videos: %v", err))
	}

	return &pb.CollectSubscriptionsResponse{
		ChannelsProcessed: int32(result.ChannelsProcessed),
		VideosProcessed:   int32(result.VideosCollected),
		VideosAdded:       int32(result.VideosCreated),
		DurationMs:        result.Duration.Milliseconds(),
	}, nil
}

func (s *Server) CreateSnapshot(ctx context.Context, req *pb.CreateSnapshotRequest) (*pb.CreateSnapshotResponse, error) {
	if req.VideoId == "" {
		return nil, status.Error(codes.InvalidArgument, "video_id is required")
	}

	videoID, err := uuid.Parse(req.VideoId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid video_id format")
	}

	snapshot, err := s.systemUseCase.CreateSnapshot(ctx, &input.CreateSnapshotInput{
		VideoID:        videoID,
		CheckpointHour: int(req.CheckpointHour),
	})
	if err != nil {
		if err == domain.ErrVideoNotFound {
			return nil, status.Error(codes.NotFound, "video not found")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create snapshot: %v", err))
	}

	return &pb.CreateSnapshotResponse{
		Snapshot: domainSnapshotToProto(snapshot),
	}, nil
}

func (s *Server) ScheduleSnapshots(ctx context.Context, req *pb.ScheduleSnapshotsRequest) (*pb.ScheduleSnapshotsResponse, error) {
	result, err := s.systemUseCase.ScheduleSnapshots(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to schedule snapshots: %v", err))
	}

	return &pb.ScheduleSnapshotsResponse{
		VideosProcessed: int32(result.VideosProcessed),
		TasksScheduled:  int32(result.TasksScheduled),
		DurationMs:      result.Duration.Milliseconds(),
	}, nil
}

func (s *Server) UpdateChannels(ctx context.Context, req *pb.UpdateChannelsRequest) (*pb.UpdateChannelsResponse, error) {
	result, err := s.channelUseCase.UpdateChannels(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update channels: %v", err))
	}

	return &pb.UpdateChannelsResponse{
		ChannelsProcessed: int32(result.ChannelsProcessed),
		ChannelsUpdated:   int32(result.ChannelsUpdated),
		DurationMs:        result.Duration.Milliseconds(),
	}, nil
}

// Unimplemented methods
func (s *Server) SubscribeChannel(ctx context.Context, req *pb.SubscribeChannelRequest) (*pb.SubscribeChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeChannel not implemented")
}

func (s *Server) UnsubscribeChannel(ctx context.Context, req *pb.UnsubscribeChannelRequest) (*pb.UnsubscribeChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnsubscribeChannel not implemented")
}

func (s *Server) GetVideo(ctx context.Context, req *pb.GetVideoRequest) (*pb.GetVideoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVideo not implemented")
}

func (s *Server) ListVideos(ctx context.Context, req *pb.ListVideosRequest) (*pb.ListVideosResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVideos not implemented")
}

func (s *Server) GetSnapshot(ctx context.Context, req *pb.GetSnapshotRequest) (*pb.GetSnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSnapshot not implemented")
}

func (s *Server) ListSnapshots(ctx context.Context, req *pb.ListSnapshotsRequest) (*pb.ListSnapshotsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSnapshots not implemented")
}

// Genre operations

func (s *Server) ListGenres(ctx context.Context, req *pb.ListGenresRequest) (*pb.ListGenresResponse, error) {
	if s.genreUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "genre use case not available")
	}

	var genres []*domain.Genre
	var err error

	if req.EnabledOnly {
		genres, err = s.genreUseCase.ListEnabledGenres(ctx)
	} else {
		genres, err = s.genreUseCase.ListGenres(ctx)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoGenres := make([]*pb.Genre, len(genres))
	for i, genre := range genres {
		protoGenres[i] = domainGenreToProto(genre)
	}

	return &pb.ListGenresResponse{
		Genres:     protoGenres,
		TotalCount: int32(len(genres)),
	}, nil
}

func (s *Server) GetGenre(ctx context.Context, req *pb.GetGenreRequest) (*pb.GetGenreResponse, error) {
	if s.genreUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "genre use case not available")
	}

	genreID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid genre ID")
	}

	genre, err := s.genreUseCase.GetGenre(ctx, genreID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "genre not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetGenreResponse{
		Genre: domainGenreToProto(genre),
	}, nil
}

func (s *Server) GetGenreByCode(ctx context.Context, req *pb.GetGenreByCodeRequest) (*pb.GetGenreByCodeResponse, error) {
	if s.genreUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "genre use case not available")
	}

	genre, err := s.genreUseCase.GetGenreByCode(ctx, req.Code)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "genre not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetGenreByCodeResponse{
		Genre: domainGenreToProto(genre),
	}, nil
}

func (s *Server) CreateGenre(ctx context.Context, req *pb.CreateGenreRequest) (*pb.CreateGenreResponse, error) {
	if s.genreUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "genre use case not available")
	}

	// Convert category IDs
	categoryIDs := make([]int, len(req.CategoryIds))
	for i, id := range req.CategoryIds {
		categoryIDs[i] = int(id)
	}

	genre, err := s.genreUseCase.CreateGenre(ctx, &input.CreateGenreInput{
		Code:        req.Code,
		Name:        req.Name,
		Language:    req.Language,
		RegionCode:  req.RegionCode,
		CategoryIDs: categoryIDs,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateGenreResponse{
		Genre: domainGenreToProto(genre),
	}, nil
}

func (s *Server) UpdateGenre(ctx context.Context, req *pb.UpdateGenreRequest) (*pb.UpdateGenreResponse, error) {
	if s.genreUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "genre use case not available")
	}

	genreID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid genre ID")
	}

	// Convert category IDs
	categoryIDs := make([]int, len(req.CategoryIds))
	for i, id := range req.CategoryIds {
		categoryIDs[i] = int(id)
	}

	genre, err := s.genreUseCase.UpdateGenre(ctx, &input.UpdateGenreInput{
		GenreID:     genreID,
		Name:        req.Name,
		CategoryIDs: categoryIDs,
	})
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "genre not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateGenreResponse{
		Genre: domainGenreToProto(genre),
	}, nil
}

func (s *Server) EnableGenre(ctx context.Context, req *pb.EnableGenreRequest) (*pb.EnableGenreResponse, error) {
	if s.genreUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "genre use case not available")
	}

	genreID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid genre ID")
	}

	genre, err := s.genreUseCase.EnableGenre(ctx, genreID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "genre not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EnableGenreResponse{
		Genre: domainGenreToProto(genre),
	}, nil
}

func (s *Server) DisableGenre(ctx context.Context, req *pb.DisableGenreRequest) (*pb.DisableGenreResponse, error) {
	if s.genreUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "genre use case not available")
	}

	genreID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid genre ID")
	}

	genre, err := s.genreUseCase.DisableGenre(ctx, genreID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "genre not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DisableGenreResponse{
		Genre: domainGenreToProto(genre),
	}, nil
}

// YouTube Category operations

func (s *Server) ListYouTubeCategories(ctx context.Context, req *pb.ListYouTubeCategoriesRequest) (*pb.ListYouTubeCategoriesResponse, error) {
	if s.youtubeCategoryUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "youtube category use case not available")
	}

	var categories []*domain.YouTubeCategory
	var err error

	if req.AssignableOnly {
		categories, err = s.youtubeCategoryUseCase.ListAssignableYouTubeCategories(ctx)
	} else {
		categories, err = s.youtubeCategoryUseCase.ListYouTubeCategories(ctx)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoCategories := make([]*pb.YouTubeCategory, len(categories))
	for i, category := range categories {
		protoCategories[i] = domainYouTubeCategoryToProto(category)
	}

	return &pb.ListYouTubeCategoriesResponse{
		Categories: protoCategories,
		TotalCount: int32(len(categories)),
	}, nil
}

func (s *Server) GetYouTubeCategory(ctx context.Context, req *pb.GetYouTubeCategoryRequest) (*pb.GetYouTubeCategoryResponse, error) {
	if s.youtubeCategoryUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "youtube category use case not available")
	}

	category, err := s.youtubeCategoryUseCase.GetYouTubeCategory(ctx, int(req.Id))
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "youtube category not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetYouTubeCategoryResponse{
		Category: domainYouTubeCategoryToProto(category),
	}, nil
}

func (s *Server) UpdateYouTubeCategory(ctx context.Context, req *pb.UpdateYouTubeCategoryRequest) (*pb.UpdateYouTubeCategoryResponse, error) {
	if s.youtubeCategoryUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "youtube category use case not available")
	}

	category, err := s.youtubeCategoryUseCase.UpdateYouTubeCategory(ctx, &input.UpdateYouTubeCategoryInput{
		CategoryID: int(req.Id),
		Name:       req.Title,
		Assignable: req.Assignable,
	})
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "youtube category not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateYouTubeCategoryResponse{
		Category: domainYouTubeCategoryToProto(category),
	}, nil
}

// Keyword operations (updated with genre support)

func (s *Server) GetKeyword(ctx context.Context, req *pb.GetKeywordRequest) (*pb.GetKeywordResponse, error) {
	if s.keywordUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "keyword use case not available")
	}

	keywordID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid keyword ID")
	}

	keyword, err := s.keywordUseCase.GetKeyword(ctx, keywordID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "keyword not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetKeywordResponse{
		Keyword: domainKeywordToProto(keyword),
	}, nil
}

func (s *Server) ListKeywords(ctx context.Context, req *pb.ListKeywordsRequest) (*pb.ListKeywordsResponse, error) {
	if s.keywordUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "keyword use case not available")
	}

	keywords, err := s.keywordUseCase.ListKeywords(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Filter enabled only if requested
	filteredKeywords := keywords
	if req.EnabledOnly {
		filteredKeywords = make([]*domain.Keyword, 0)
		for _, k := range keywords {
			if k.Enabled {
				filteredKeywords = append(filteredKeywords, k)
			}
		}
	}

	protoKeywords := make([]*pb.Keyword, len(filteredKeywords))
	for i, keyword := range filteredKeywords {
		protoKeywords[i] = domainKeywordToProto(keyword)
	}

	return &pb.ListKeywordsResponse{
		Keywords:   protoKeywords,
		TotalCount: int32(len(filteredKeywords)),
	}, nil
}

func (s *Server) ListKeywordsByGenre(ctx context.Context, req *pb.ListKeywordsByGenreRequest) (*pb.ListKeywordsByGenreResponse, error) {
	if s.keywordUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "keyword use case not available")
	}

	genreID, err := uuid.Parse(req.GenreId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid genre ID")
	}

	keywords, err := s.keywordUseCase.ListKeywordsByGenre(ctx, genreID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Filter enabled only if requested
	filteredKeywords := keywords
	if req.EnabledOnly {
		filteredKeywords = make([]*domain.Keyword, 0)
		for _, k := range keywords {
			if k.Enabled {
				filteredKeywords = append(filteredKeywords, k)
			}
		}
	}

	protoKeywords := make([]*pb.Keyword, len(filteredKeywords))
	for i, keyword := range filteredKeywords {
		protoKeywords[i] = domainKeywordToProto(keyword)
	}

	return &pb.ListKeywordsByGenreResponse{
		Keywords:   protoKeywords,
		TotalCount: int32(len(filteredKeywords)),
	}, nil
}

func (s *Server) CreateKeyword(ctx context.Context, req *pb.CreateKeywordRequest) (*pb.CreateKeywordResponse, error) {
	if s.keywordUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "keyword use case not available")
	}

	genreID, err := uuid.Parse(req.GenreId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid genre ID")
	}

	keyword, err := s.keywordUseCase.CreateKeyword(ctx, &input.CreateKeywordInput{
		GenreID:     genreID,
		Name:        req.Name,
		FilterType:  req.FilterType,
		Pattern:     req.Pattern,
		TargetField: req.TargetField,
		Description: &req.Description,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateKeywordResponse{
		Keyword: domainKeywordToProto(keyword),
	}, nil
}

func (s *Server) UpdateKeyword(ctx context.Context, req *pb.UpdateKeywordRequest) (*pb.UpdateKeywordResponse, error) {
	if s.keywordUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "keyword use case not available")
	}

	keywordID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid keyword ID")
	}

	keyword, err := s.keywordUseCase.UpdateKeyword(ctx, &input.UpdateKeywordInput{
		KeywordID:   keywordID,
		Name:        req.Name,
		FilterType:  req.FilterType,
		Pattern:     req.Pattern,
		TargetField: req.TargetField,
		Enabled:     req.Enabled,
		Description: &req.Description,
	})
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "keyword not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateKeywordResponse{
		Keyword: domainKeywordToProto(keyword),
	}, nil
}

func (s *Server) EnableKeyword(ctx context.Context, req *pb.EnableKeywordRequest) (*pb.EnableKeywordResponse, error) {
	if s.keywordUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "keyword use case not available")
	}

	keywordID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid keyword ID")
	}

	keyword, err := s.keywordUseCase.EnableKeyword(ctx, keywordID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "keyword not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EnableKeywordResponse{
		Keyword: domainKeywordToProto(keyword),
	}, nil
}

func (s *Server) DisableKeyword(ctx context.Context, req *pb.DisableKeywordRequest) (*pb.DisableKeywordResponse, error) {
	if s.keywordUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "keyword use case not available")
	}

	keywordID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid keyword ID")
	}

	keyword, err := s.keywordUseCase.DisableKeyword(ctx, keywordID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "keyword not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DisableKeywordResponse{
		Keyword: domainKeywordToProto(keyword),
	}, nil
}

func (s *Server) DeleteKeyword(ctx context.Context, req *pb.DeleteKeywordRequest) (*pb.DeleteKeywordResponse, error) {
	if s.keywordUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "keyword use case not available")
	}

	keywordID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid keyword ID")
	}

	if err := s.keywordUseCase.DeleteKeyword(ctx, keywordID); err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "keyword not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteKeywordResponse{}, nil
}

// Helper functions to convert between domain and proto types

func domainChannelToProto(channel *domain.Channel) *pb.Channel {
	proto := &pb.Channel{
		Id:                string(channel.ID),
		YoutubeChannelId:  string(channel.YouTubeChannelID),
		Title:             channel.Title,
		ThumbnailUrl:      channel.ThumbnailURL,
		Description:       channel.Description,
		Country:           channel.Country,
		ViewCount:         channel.ViewCount,
		SubscriptionCount: channel.SubscriptionCount,
		VideoCount:        channel.VideoCount,
		Subscribed:        channel.Subscribed,
		CreatedAt:         timestamppb.New(channel.CreatedAt),
	}

	if channel.UpdatedAt != nil {
		proto.UpdatedAt = timestamppb.New(*channel.UpdatedAt)
	}
	if channel.DeletedAt != nil {
		proto.DeletedAt = timestamppb.New(*channel.DeletedAt)
	}

	return proto
}

func domainSnapshotToProto(snapshot *domain.VideoSnapshot) *pb.VideoSnapshot {
	proto := &pb.VideoSnapshot{
		Id:                string(snapshot.ID),
		VideoId:           string(snapshot.VideoID),
		CheckpointHour:    int32(snapshot.CheckpointHour),
		MeasuredAt:        timestamppb.New(snapshot.MeasuredAt),
		ViewsCount:        snapshot.ViewsCount,
		LikesCount:        snapshot.LikesCount,
		SubscriptionCount: snapshot.SubscriptionCount,
		Source:            string(snapshot.Source),
	}

	// CreatedAt is not a pointer
	proto.CreatedAt = timestamppb.New(snapshot.CreatedAt)

	return proto
}

// Video-Genre operations

func (s *Server) ListVideoGenres(ctx context.Context, req *pb.ListVideoGenresRequest) (*pb.ListVideoGenresResponse, error) {
	if s.videoGenreUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "video-genre use case not available")
	}

	videoID, err := uuid.Parse(req.VideoId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid video ID")
	}

	videoGenres, err := s.videoGenreUseCase.GetVideoGenres(ctx, videoID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoVideoGenres := make([]*pb.VideoGenre, len(videoGenres))
	for i, vg := range videoGenres {
		protoVideoGenres[i] = &pb.VideoGenre{
			VideoId:   string(vg.VideoID),
			GenreId:   string(vg.GenreID),
			CreatedAt: timestamppb.New(vg.CreatedAt),
		}
	}

	return &pb.ListVideoGenresResponse{
		VideoGenres: protoVideoGenres,
	}, nil
}

func (s *Server) AssignVideoToGenre(ctx context.Context, req *pb.AssignVideoToGenreRequest) (*pb.AssignVideoToGenreResponse, error) {
	if s.videoGenreUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "video-genre use case not available")
	}

	videoID, err := uuid.Parse(req.VideoId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid video ID")
	}

	genreID, err := uuid.Parse(req.GenreId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid genre ID")
	}

	videoGenre, err := s.videoGenreUseCase.AssociateVideoWithGenre(ctx, videoID, genreID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AssignVideoToGenreResponse{
		VideoGenre: &pb.VideoGenre{
			VideoId:   string(videoGenre.VideoID),
			GenreId:   string(videoGenre.GenreID),
			CreatedAt: timestamppb.New(videoGenre.CreatedAt),
		},
	}, nil
}

func (s *Server) RemoveVideoFromGenre(ctx context.Context, req *pb.RemoveVideoFromGenreRequest) (*pb.RemoveVideoFromGenreResponse, error) {
	if s.videoGenreUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "video-genre use case not available")
	}

	videoID, err := uuid.Parse(req.VideoId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid video ID")
	}

	genreID, err := uuid.Parse(req.GenreId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid genre ID")
	}

	if err := s.videoGenreUseCase.DisassociateVideoFromGenre(ctx, videoID, genreID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveVideoFromGenreResponse{}, nil
}

// Audit log operations

func (s *Server) ListAuditLogs(ctx context.Context, req *pb.ListAuditLogsRequest) (*pb.ListAuditLogsResponse, error) {
	if s.auditLogUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "audit log use case not available")
	}

	// TODO: Implement pagination
	pageSize := int(req.PageSize)
	if pageSize == 0 {
		pageSize = 100
	}

	var logs []*domain.AuditLog
	var err error

	if req.ActorId != "" {
		actorID, err := uuid.Parse(req.ActorId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid actor ID")
		}
		logs, err = s.auditLogUseCase.ListAuditLogsByActor(ctx, actorID, pageSize, 0)
	} else if req.ResourceType != "" && req.ResourceId != "" {
		logs, err = s.auditLogUseCase.ListAuditLogsByResource(ctx, req.ResourceType, req.ResourceId, pageSize, 0)
	} else {
		logs, err = s.auditLogUseCase.ListRecentAuditLogs(ctx, pageSize)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoLogs := make([]*pb.AuditLog, len(logs))
	for i, log := range logs {
		protoLogs[i] = domainAuditLogToProto(log)
	}

	return &pb.ListAuditLogsResponse{
		AuditLogs:  protoLogs,
		TotalCount: int32(len(logs)),
	}, nil
}

func (s *Server) GetAuditLog(ctx context.Context, req *pb.GetAuditLogRequest) (*pb.GetAuditLogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuditLog not implemented")
}

// Batch job operations

func (s *Server) ListBatchJobs(ctx context.Context, req *pb.ListBatchJobsRequest) (*pb.ListBatchJobsResponse, error) {
	if s.batchJobUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "batch job use case not available")
	}

	// TODO: Implement pagination and filtering
	pageSize := 100

	var jobType *domain.JobType
	if req.JobType != "" {
		jt := domain.JobType(req.JobType)
		jobType = &jt
	}

	var jobTypeStr *string
	if jobType != nil {
		jt := string(*jobType)
		jobTypeStr = &jt
	}
	jobs, err := s.batchJobUseCase.ListRecentBatchJobs(ctx, jobTypeStr, pageSize)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoJobs := make([]*pb.BatchJob, len(jobs))
	for i, job := range jobs {
		protoJobs[i] = domainBatchJobToProto(job)
	}

	return &pb.ListBatchJobsResponse{
		BatchJobs:  protoJobs,
		TotalCount: int32(len(jobs)),
	}, nil
}

func (s *Server) GetBatchJob(ctx context.Context, req *pb.GetBatchJobRequest) (*pb.GetBatchJobResponse, error) {
	if s.batchJobUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "batch job use case not available")
	}

	jobID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid job ID")
	}

	job, err := s.batchJobUseCase.GetBatchJob(ctx, jobID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "batch job not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetBatchJobResponse{
		BatchJob: domainBatchJobToProto(job),
	}, nil
}

// System operations

func (s *Server) CollectTrendingByGenre(ctx context.Context, req *pb.CollectTrendingByGenreRequest) (*pb.CollectTrendingByGenreResponse, error) {
	if s.videoUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "video use case not available")
	}

	result, err := s.videoUseCase.CollectTrendingByGenre(ctx, req.GenreId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CollectTrendingByGenreResponse{
		GenreCode:       result.GenreCode,
		VideosProcessed: int32(result.VideosCollected),
		VideosAdded:     int32(result.VideosCreated),
		DurationMs:      result.Duration.Milliseconds(),
	}, nil
}

func (s *Server) CollectAllTrending(ctx context.Context, req *pb.CollectAllTrendingRequest) (*pb.CollectAllTrendingResponse, error) {
	if s.videoUseCase == nil {
		return nil, status.Error(codes.Unimplemented, "video use case not available")
	}

	result, err := s.videoUseCase.CollectAllTrending(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	genreResults := make([]*pb.CollectTrendingByGenreResponse, len(result.GenreResults))
	for i, gr := range result.GenreResults {
		genreResults[i] = &pb.CollectTrendingByGenreResponse{
			GenreCode:       gr.GenreCode,
			VideosProcessed: int32(gr.VideosCollected),
			VideosAdded:     int32(gr.VideosCreated),
			DurationMs:      gr.Duration.Milliseconds(),
		}
	}

	return &pb.CollectAllTrendingResponse{
		GenresProcessed: int32(result.GenresProcessed),
		TotalVideos:     int32(result.TotalCollected),
		TotalAdded:      int32(result.TotalCreated),
		GenreResults:    genreResults,
		DurationMs:      result.Duration.Milliseconds(),
	}, nil
}

// Add helper functions for new domain types

func domainGenreToProto(genre *domain.Genre) *pb.Genre {
	categoryIDs := make([]int32, len(genre.CategoryIDs))
	for i, id := range genre.CategoryIDs {
		categoryIDs[i] = int32(id)
	}

	return &pb.Genre{
		Id:          string(genre.ID),
		Code:        genre.Code,
		Name:        genre.Name,
		Language:    genre.Language,
		RegionCode:  genre.RegionCode,
		CategoryIds: categoryIDs,
		Enabled:     genre.Enabled,
	}
}

func domainYouTubeCategoryToProto(category *domain.YouTubeCategory) *pb.YouTubeCategory {
	return &pb.YouTubeCategory{
		Id:         int32(category.ID),
		Title:      category.Name,
		Assignable: category.Assignable,
	}
}

func domainKeywordToProto(keyword *domain.Keyword) *pb.Keyword {
	proto := &pb.Keyword{
		Id:          string(keyword.ID),
		GenreId:     string(keyword.GenreID),
		Name:        keyword.Name,
		FilterType:  string(keyword.FilterType),
		Pattern:     keyword.Pattern,
		TargetField: keyword.TargetField,
		Enabled:     keyword.Enabled,
		CreatedAt:   timestamppb.New(keyword.CreatedAt),
	}

	if keyword.Description != nil {
		proto.Description = *keyword.Description
	}

	if keyword.UpdatedAt != nil {
		proto.UpdatedAt = timestamppb.New(*keyword.UpdatedAt)
	}

	if keyword.DeletedAt != nil {
		proto.DeletedAt = timestamppb.New(*keyword.DeletedAt)
	}

	return proto
}

func domainAuditLogToProto(log *domain.AuditLog) *pb.AuditLog {
	proto := &pb.AuditLog{
		Id:           string(log.ID),
		ActorId:      string(log.ActorID),
		ActorEmail:   log.ActorEmail,
		Action:       log.Action,
		ResourceType: log.ResourceType,
		ResourceId:   log.ResourceID,
		UserAgent:    log.UserAgent,
		CreatedAt:    timestamppb.New(log.CreatedAt),
	}

	if log.IPAddress != nil {
		proto.IpAddress = log.IPAddress.String()
	}

	// Convert old/new values
	if log.OldValues != nil {
		proto.OldValues = make(map[string]string)
		for k, v := range log.OldValues {
			proto.OldValues[k] = fmt.Sprintf("%v", v)
		}
	}

	if log.NewValues != nil {
		proto.NewValues = make(map[string]string)
		for k, v := range log.NewValues {
			proto.NewValues[k] = fmt.Sprintf("%v", v)
		}
	}

	return proto
}

func domainBatchJobToProto(job *domain.BatchJob) *pb.BatchJob {
	proto := &pb.BatchJob{
		Id:           string(job.ID),
		JobType:      string(job.JobType),
		Status:       string(job.Status),
		ErrorMessage: job.ErrorMessage,
		CreatedAt:    timestamppb.New(job.CreatedAt),
	}

	if job.StartedAt != nil {
		proto.StartedAt = timestamppb.New(*job.StartedAt)
	}

	if job.CompletedAt != nil {
		proto.CompletedAt = timestamppb.New(*job.CompletedAt)
	}

	// Convert parameters
	if job.Parameters != nil {
		proto.Parameters = make(map[string]string)
		for k, v := range job.Parameters {
			proto.Parameters[k] = fmt.Sprintf("%v", v)
		}
	}

	// Convert statistics
	if job.Statistics != nil {
		proto.Statistics = make(map[string]string)
		for k, v := range job.Statistics {
			proto.Statistics[k] = fmt.Sprintf("%v", v)
		}
	}

	return proto
}
