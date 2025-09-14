package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/presenter"
	pb "github.com/YukiOnishi1129/youtube-analytics/services/pkg/pb/ingestion/v1"
)

// grpcPresenter implements presenter.GRPCPresenter
type grpcPresenter struct{}

// NewGRPCPresenter creates a new gRPC presenter
func NewGRPCPresenter() presenter.GRPCPresenter {
	return &grpcPresenter{}
}

// PresentChannel presents a single channel for gRPC
func (p *grpcPresenter) PresentChannel(channel *domain.Channel) interface{} {
	return &pb.GetChannelResponse{
		Channel: domainChannelToProto(channel),
	}
}

// PresentChannels presents multiple channels for gRPC
func (p *grpcPresenter) PresentChannels(channels []*domain.Channel) interface{} {
	items := make([]*pb.Channel, len(channels))
	for i, channel := range channels {
		items[i] = domainChannelToProto(channel)
	}
	return &pb.ListChannelsResponse{
		Channels: items,
	}
}

// PresentVideo presents a single video for gRPC
func (p *grpcPresenter) PresentVideo(video *domain.Video) interface{} {
	return &pb.GetVideoResponse{
		Video: domainVideoToProto(video),
	}
}

// PresentVideos presents multiple videos for gRPC
func (p *grpcPresenter) PresentVideos(videos []*domain.Video) interface{} {
	items := make([]*pb.Video, len(videos))
	for i, video := range videos {
		items[i] = domainVideoToProto(video)
	}
	return &pb.ListVideosResponse{
		Videos: items,
	}
}

// PresentSnapshot presents a single snapshot for gRPC
func (p *grpcPresenter) PresentSnapshot(snapshot *domain.VideoSnapshot) interface{} {
	return &pb.GetSnapshotResponse{
		Snapshot: domainSnapshotToProto(snapshot),
	}
}

// PresentSnapshots presents multiple snapshots for gRPC
func (p *grpcPresenter) PresentSnapshots(snapshots []*domain.VideoSnapshot) interface{} {
	items := make([]*pb.VideoSnapshot, len(snapshots))
	for i, snapshot := range snapshots {
		items[i] = domainSnapshotToProto(snapshot)
	}
	return &pb.ListSnapshotsResponse{
		Snapshots: items,
	}
}

// PresentKeyword presents a single keyword for gRPC
func (p *grpcPresenter) PresentKeyword(keyword *domain.Keyword) interface{} {
	return &pb.GetKeywordResponse{
		Keyword: domainKeywordToProto(keyword),
	}
}

// PresentKeywords presents multiple keywords for gRPC
func (p *grpcPresenter) PresentKeywords(keywords []*domain.Keyword) interface{} {
	items := make([]*pb.Keyword, len(keywords))
	for i, keyword := range keywords {
		items[i] = domainKeywordToProto(keyword)
	}
	return &pb.ListKeywordsResponse{
		Keywords: items,
	}
}

// PresentScheduleSnapshotsResult presents the result of scheduling snapshots
func (p *grpcPresenter) PresentScheduleSnapshotsResult(result interface{}) interface{} {
	return &pb.ScheduleSnapshotsResponse{
		VideosProcessed: 0, // Not used in gRPC - this is HTTP admin endpoint
		TasksScheduled:  0,
		DurationMs:      0,
	}
}

// PresentUpdateChannelsResult presents the result of updating channels
func (p *grpcPresenter) PresentUpdateChannelsResult(result interface{}) interface{} {
	return &pb.UpdateChannelsResponse{
		ChannelsProcessed: 0, // Not used in gRPC - this is HTTP admin endpoint
		ChannelsUpdated:   0,
		DurationMs:        0,
	}
}

// PresentCollectTrendingResult presents the result of collecting trending videos
func (p *grpcPresenter) PresentCollectTrendingResult(result interface{}) interface{} {
	return &pb.CollectTrendingResponse{
		VideosProcessed: 0, // Not used in gRPC - this is HTTP admin endpoint
		VideosAdded:     0,
		DurationMs:      0,
	}
}

// PresentCollectSubscriptionsResult presents the result of collecting subscription videos
func (p *grpcPresenter) PresentCollectSubscriptionsResult(result interface{}) interface{} {
	return &pb.CollectSubscriptionsResponse{
		ChannelsProcessed: 0, // Not used in gRPC - this is HTTP admin endpoint
		VideosProcessed:   0,
		VideosAdded:       0,
		DurationMs:        0,
	}
}

// PresentDeleted presents a successful deletion
func (p *grpcPresenter) PresentDeleted() interface{} {
	return nil // gRPC uses empty response
}

// PresentError presents an error for gRPC
func (p *grpcPresenter) PresentError(err error) error {
	switch err {
	case domain.ErrChannelNotFound, domain.ErrVideoNotFound, domain.ErrSnapshotNotFound, domain.ErrKeywordNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrKeywordDuplicate:
		return status.Error(codes.AlreadyExists, err.Error())
	case domain.ErrInvalidInput:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

// Helper functions to convert domain models to proto

func domainChannelToProto(channel *domain.Channel) *pb.Channel {
	pbChannel := &pb.Channel{
		Id:               string(channel.ID),
		YoutubeChannelId: string(channel.YouTubeChannelID),
		Title:            channel.Title,
		ThumbnailUrl:     channel.ThumbnailURL,
		Subscribed:       channel.Subscribed,
		CreatedAt:        timestamppb.New(channel.CreatedAt),
	}
	if channel.UpdatedAt != nil {
		pbChannel.UpdatedAt = timestamppb.New(*channel.UpdatedAt)
	}
	return pbChannel
}

func domainVideoToProto(video *domain.Video) *pb.Video {
	pbVideo := &pb.Video{
		Id:               string(video.ID),
		YoutubeVideoId:   string(video.YouTubeVideoID),
		YoutubeChannelId: string(video.YouTubeChannelID),
		Title:            video.Title,
		PublishedAt:      timestamppb.New(video.PublishedAt),
		CategoryId:       int32(video.CategoryID),
		CreatedAt:        timestamppb.New(video.CreatedAt),
	}
	if video.UpdatedAt != nil {
		pbVideo.UpdatedAt = timestamppb.New(*video.UpdatedAt)
	}
	return pbVideo
}

func domainSnapshotToProto(snapshot *domain.VideoSnapshot) *pb.VideoSnapshot {
	return &pb.VideoSnapshot{
		Id:                string(snapshot.ID),
		VideoId:           string(snapshot.VideoID),
		CheckpointHour:    int32(snapshot.CheckpointHour),
		MeasuredAt:        timestamppb.New(snapshot.MeasuredAt),
		ViewsCount:        snapshot.ViewsCount,
		LikesCount:        snapshot.LikesCount,
		SubscriptionCount: snapshot.SubscriptionCount,
		Source:            string(snapshot.Source),
		CreatedAt:         timestamppb.New(snapshot.CreatedAt),
	}
}

func domainKeywordToProto(keyword *domain.Keyword) *pb.Keyword {
	pbKeyword := &pb.Keyword{
		Id:          string(keyword.ID),
		Name:        keyword.Name,
		FilterType:  string(keyword.FilterType),
		Pattern:     keyword.Pattern,
		CreatedAt:   timestamppb.New(keyword.CreatedAt),
	}
	if keyword.Description != nil {
		pbKeyword.Description = *keyword.Description
	}
	if keyword.UpdatedAt != nil {
		pbKeyword.UpdatedAt = timestamppb.New(*keyword.UpdatedAt)
	}
	return pbKeyword
}
