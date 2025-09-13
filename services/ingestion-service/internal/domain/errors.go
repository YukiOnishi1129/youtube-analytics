package domain

import "errors"

// Domain errors
var (
	// Channel errors
	ErrChannelNotFound      = errors.New("channel not found")
	ErrChannelAlreadyExists = errors.New("channel already exists")
	ErrInvalidChannelID     = errors.New("invalid channel ID")

	// Video errors
	ErrVideoNotFound      = errors.New("video not found")
	ErrVideoAlreadyExists = errors.New("video already exists")
	ErrInvalidVideoID     = errors.New("invalid video ID")

	// Snapshot errors
	ErrSnapshotNotFound      = errors.New("snapshot not found")
	ErrSnapshotAlreadyExists = errors.New("snapshot already exists")
	ErrInvalidCheckpoint     = errors.New("invalid checkpoint hour")

	// Keyword errors
	ErrKeywordNotFound  = errors.New("keyword not found")
	ErrKeywordDuplicate = errors.New("keyword already exists")

	// General errors
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("not found")
)
