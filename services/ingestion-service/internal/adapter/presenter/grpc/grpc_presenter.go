package grpc

import (
	"fmt"
	
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
func (p *grpcPresenter) PresentChannel(channel *domain.Channel) *pb.Channel {
	return domainChannelToProto(channel)
}

// PresentChannels presents multiple channels for gRPC
func (p *grpcPresenter) PresentChannels(channels []*domain.Channel) []*pb.Channel {
	items := make([]*pb.Channel, len(channels))
	for i, channel := range channels {
		items[i] = domainChannelToProto(channel)
	}
	return items
}

// PresentVideo presents a single video for gRPC
func (p *grpcPresenter) PresentVideo(video *domain.Video) *pb.Video {
	return domainVideoToProto(video)
}

// PresentVideos presents multiple videos for gRPC
func (p *grpcPresenter) PresentVideos(videos []*domain.Video) []*pb.Video {
	items := make([]*pb.Video, len(videos))
	for i, video := range videos {
		items[i] = domainVideoToProto(video)
	}
	return items
}

// PresentSnapshot presents a single snapshot for gRPC
func (p *grpcPresenter) PresentSnapshot(snapshot *domain.VideoSnapshot) *pb.VideoSnapshot {
	return domainSnapshotToProto(snapshot)
}

// PresentSnapshots presents multiple snapshots for gRPC
func (p *grpcPresenter) PresentSnapshots(snapshots []*domain.VideoSnapshot) []*pb.VideoSnapshot {
	items := make([]*pb.VideoSnapshot, len(snapshots))
	for i, snapshot := range snapshots {
		items[i] = domainSnapshotToProto(snapshot)
	}
	return items
}

// PresentKeyword presents a single keyword for gRPC
func (p *grpcPresenter) PresentKeyword(keyword *domain.Keyword) *pb.Keyword {
	return domainKeywordToProto(keyword)
}

// PresentKeywords presents multiple keywords for gRPC
func (p *grpcPresenter) PresentKeywords(keywords []*domain.Keyword) []*pb.Keyword {
	items := make([]*pb.Keyword, len(keywords))
	for i, keyword := range keywords {
		items[i] = domainKeywordToProto(keyword)
	}
	return items
}

// PresentScheduleSnapshotsResult presents the result of scheduling snapshots
func (p *grpcPresenter) PresentScheduleSnapshotsResult(count int32) *pb.ScheduleSnapshotsResponse {
	return &pb.ScheduleSnapshotsResponse{
		VideosProcessed: count,
		TasksScheduled:  count,
		DurationMs:      0,
	}
}

// PresentUpdateChannelsResult presents the result of updating channels
func (p *grpcPresenter) PresentUpdateChannelsResult(count int32) *pb.UpdateChannelsResponse {
	return &pb.UpdateChannelsResponse{
		ChannelsProcessed: count,
		ChannelsUpdated:   count,
		DurationMs:        0,
	}
}

// PresentCollectTrendingResult presents the result of collecting trending videos
func (p *grpcPresenter) PresentCollectTrendingResult(count int32) *pb.CollectTrendingResponse {
	return &pb.CollectTrendingResponse{
		VideosProcessed: count,
		VideosAdded:     count,
		DurationMs:      0,
	}
}

// PresentCollectSubscriptionsResult presents the result of collecting subscription videos
func (p *grpcPresenter) PresentCollectSubscriptionsResult(count int32) *pb.CollectSubscriptionsResponse {
	return &pb.CollectSubscriptionsResponse{
		ChannelsProcessed: count,
		VideosProcessed:   count,
		VideosAdded:       count,
		DurationMs:        0,
	}
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
		GenreId:     string(keyword.GenreID),
		Name:        keyword.Name,
		FilterType:  string(keyword.FilterType),
		Pattern:     keyword.Pattern,
		TargetField: keyword.TargetField,
		Enabled:     keyword.Enabled,
		CreatedAt:   timestamppb.New(keyword.CreatedAt),
	}
	if keyword.Description != nil {
		pbKeyword.Description = *keyword.Description
	}
	if keyword.UpdatedAt != nil {
		pbKeyword.UpdatedAt = timestamppb.New(*keyword.UpdatedAt)
	}
	if keyword.DeletedAt != nil {
		pbKeyword.DeletedAt = timestamppb.New(*keyword.DeletedAt)
	}
	return pbKeyword
}

// PresentGenre presents a single genre for gRPC
func (p *grpcPresenter) PresentGenre(genre *domain.Genre) *pb.Genre {
	return domainGenreToProto(genre)
}

// PresentGenres presents multiple genres for gRPC
func (p *grpcPresenter) PresentGenres(genres []*domain.Genre) []*pb.Genre {
	items := make([]*pb.Genre, len(genres))
	for i, genre := range genres {
		items[i] = domainGenreToProto(genre)
	}
	return items
}

// PresentYouTubeCategory presents a single YouTube category for gRPC
func (p *grpcPresenter) PresentYouTubeCategory(category *domain.YouTubeCategory) *pb.YouTubeCategory {
	return domainYouTubeCategoryToProto(category)
}

// PresentYouTubeCategories presents multiple YouTube categories for gRPC
func (p *grpcPresenter) PresentYouTubeCategories(categories []*domain.YouTubeCategory) []*pb.YouTubeCategory {
	items := make([]*pb.YouTubeCategory, len(categories))
	for i, category := range categories {
		items[i] = domainYouTubeCategoryToProto(category)
	}
	return items
}

// PresentVideoGenre presents a single video-genre association for gRPC
func (p *grpcPresenter) PresentVideoGenre(videoGenre *domain.VideoGenre) *pb.VideoGenre {
	return domainVideoGenreToProto(videoGenre)
}

// PresentVideoGenres presents multiple video-genre associations for gRPC
func (p *grpcPresenter) PresentVideoGenres(videoGenres []*domain.VideoGenre) []*pb.VideoGenre {
	items := make([]*pb.VideoGenre, len(videoGenres))
	for i, vg := range videoGenres {
		items[i] = domainVideoGenreToProto(vg)
	}
	return items
}

// PresentAuditLog presents a single audit log for gRPC
func (p *grpcPresenter) PresentAuditLog(auditLog *domain.AuditLog) *pb.AuditLog {
	return domainAuditLogToProto(auditLog)
}

// PresentAuditLogs presents multiple audit logs for gRPC
func (p *grpcPresenter) PresentAuditLogs(auditLogs []*domain.AuditLog) []*pb.AuditLog {
	items := make([]*pb.AuditLog, len(auditLogs))
	for i, log := range auditLogs {
		items[i] = domainAuditLogToProto(log)
	}
	return items
}

// PresentBatchJob presents a single batch job for gRPC
func (p *grpcPresenter) PresentBatchJob(batchJob *domain.BatchJob) *pb.BatchJob {
	return domainBatchJobToProto(batchJob)
}

// PresentBatchJobs presents multiple batch jobs for gRPC
func (p *grpcPresenter) PresentBatchJobs(batchJobs []*domain.BatchJob) []*pb.BatchJob {
	items := make([]*pb.BatchJob, len(batchJobs))
	for i, job := range batchJobs {
		items[i] = domainBatchJobToProto(job)
	}
	return items
}

// Helper functions for new domain objects

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

func domainVideoGenreToProto(vg *domain.VideoGenre) *pb.VideoGenre {
	return &pb.VideoGenre{
		VideoId:   string(vg.VideoID),
		GenreId:   string(vg.GenreID),
		CreatedAt: timestamppb.New(vg.CreatedAt),
	}
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
