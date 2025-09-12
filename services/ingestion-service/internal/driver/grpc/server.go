package grpc

import (
	"fmt"
	"log"
	"net"

	// pb "github.com/YukiOnishi1129/youtube-analytics/services/pkg/pb/proto/ingestion/v1"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
	// "google.golang.org/protobuf/types/known/timestamppb"
)

// Server implements the gRPC server for ingestion service
// Note: The actual gRPC implementation is commented out until proto generation is properly set up
type Server struct {
	// pb.UnimplementedIngestionServiceServer // Commented out until proto generation is fixed
	channelUseCase input.ChannelInputPort
	videoUseCase   input.VideoInputPort
	systemUseCase  input.SystemInputPort
	keywordUseCase input.KeywordInputPort // Optional keyword use case
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

// Start starts the gRPC server
func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	// pb.RegisterIngestionServiceServer(grpcServer, s) // Commented out until proto generation is fixed

	log.Printf("gRPC server starting on port %d (note: no services registered until proto generation is fixed)", port)
	return grpcServer.Serve(lis)
}

/*
// The following methods will be uncommented once proto generation is properly set up

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
	result, err := s.videoUseCase.CollectTrending(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to collect trending videos: %v", err))
	}

	return &pb.CollectTrendingResponse{
		VideosProcessed: int32(result.VideosProcessed),
		VideosAdded:     int32(result.VideosAdded),
		DurationMs:      result.Duration.Milliseconds(),
	}, nil
}

func (s *Server) CollectSubscriptions(ctx context.Context, req *pb.CollectSubscriptionsRequest) (*pb.CollectSubscriptionsResponse, error) {
	result, err := s.videoUseCase.CollectSubscriptions(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to collect subscription videos: %v", err))
	}

	return &pb.CollectSubscriptionsResponse{
		ChannelsProcessed: int32(result.ChannelsProcessed),
		VideosProcessed:   int32(result.VideosProcessed),
		VideosAdded:       int32(result.VideosAdded),
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

	snapshot, err := s.systemUseCase.CreateSnapshot(ctx, videoID, int(req.CheckpointHour))
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

// Helper functions to convert between domain and proto types

func domainChannelToProto(channel *domain.Channel) *pb.Channel {
	proto := &pb.Channel{
		Id:               string(channel.ID),
		YoutubeChannelId: string(channel.YouTubeChannelID),
		Title:            channel.Title,
		ThumbnailUrl:     channel.ThumbnailURL,
		Subscribed:       channel.Subscribed,
	}

	// CreatedAt is not a pointer
	proto.CreatedAt = timestamppb.New(channel.CreatedAt)

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
*/
