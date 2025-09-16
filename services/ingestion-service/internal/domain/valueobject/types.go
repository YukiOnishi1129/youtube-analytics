package valueobject

import "time"

// UUID represents a v7 UUID
type UUID string

// YouTubeVideoID represents a YouTube video ID
type YouTubeVideoID string

// YouTubeChannelID represents a YouTube channel ID
type YouTubeChannelID string

// CategoryID represents a YouTube video category ID
type CategoryID int

// CheckpointHour represents a checkpoint hour (0,3,6,12,24,48,72,168)
type CheckpointHour int

// Valid checkpoint hours
const (
	CheckpointHour0   CheckpointHour = 0
	CheckpointHour3   CheckpointHour = 3
	CheckpointHour6   CheckpointHour = 6
	CheckpointHour12  CheckpointHour = 12
	CheckpointHour24  CheckpointHour = 24
	CheckpointHour48  CheckpointHour = 48
	CheckpointHour72  CheckpointHour = 72
	CheckpointHour168 CheckpointHour = 168
)

// IsValid checks if the checkpoint hour is valid
func (c CheckpointHour) IsValid() bool {
	switch c {
	case CheckpointHour0, CheckpointHour3, CheckpointHour6, CheckpointHour12,
		CheckpointHour24, CheckpointHour48, CheckpointHour72, CheckpointHour168:
		return true
	default:
		return false
	}
}

// AllCheckpointHours returns all valid checkpoint hours
func AllCheckpointHours() []CheckpointHour {
	return []CheckpointHour{
		CheckpointHour0,
		CheckpointHour3,
		CheckpointHour6,
		CheckpointHour12,
		CheckpointHour24,
		CheckpointHour48,
		CheckpointHour72,
		CheckpointHour168,
	}
}

// GetCheckpointHoursAfter returns checkpoint hours after the given hour
func GetCheckpointHoursAfter(hour CheckpointHour) []CheckpointHour {
	var result []CheckpointHour
	allHours := AllCheckpointHours()
	
	for _, h := range allHours {
		if h > hour {
			result = append(result, h)
		}
	}
	
	return result
}

// FilterType represents the type of filter (include/exclude)
type FilterType string

const (
	FilterTypeInclude FilterType = "include"
	FilterTypeExclude FilterType = "exclude"
)

// IsValid checks if the filter type is valid
func (f FilterType) IsValid() bool {
	switch f {
	case FilterTypeInclude, FilterTypeExclude:
		return true
	default:
		return false
	}
}

// Source represents the source of a snapshot
type Source string

const (
	SourceWebSub Source = "websub"
	SourceTask   Source = "task"
	SourceManual Source = "manual"
)

// IsValid checks if the source is valid
func (s Source) IsValid() bool {
	switch s {
	case SourceWebSub, SourceTask, SourceManual:
		return true
	default:
		return false
	}
}

// GenerateUUID generates a new UUID
func GenerateUUID() UUID {
	// This is a placeholder - in production, use a proper UUID generator
	return UUID("generated-uuid")
}

// VideoMeta represents video metadata for registration
type VideoMeta struct {
	YouTubeVideoID   YouTubeVideoID
	YouTubeChannelID YouTubeChannelID
	Title            string
	PublishedAt      time.Time
	CategoryID       CategoryID
}