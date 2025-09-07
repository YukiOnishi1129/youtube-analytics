# Testing Strategy

## Test Pyramid

```
         ╱╲
        ╱E2E╲       ← Few (staging smoke tests only)
       ╱──────╲
      ╱Service ╲    ← Service-scope external I/F tests
     ╱──────────╲
    ╱Integration ╲  ← Includes DB
   ╱──────────────╲
  ╱     Unit       ╲ ← Most (domain logic)
 ╱──────────────────╲
```

## Tests by Layer

### Unit Tests

Target: Pure business logic in the domain layer

```go
// domain/service/metric_calculator_test.go
func TestCalculateGrowthRate(t *testing.T) {
    tests := []struct {
        name     string
        baseline int64
        current  int64
        hours    int
        want     float64
    }{
        {"zero baseline", 0, 100, 24, 4.17},
        {"no growth", 100, 100, 24, 0.0},
        {"negative growth", 100, 50, 24, -2.08},
    }
    
    calc := NewMetricCalculator()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := calc.CalculateGrowthRate(tt.baseline, tt.current, tt.hours)
            assert.InDelta(t, tt.want, got, 0.01)
        })
    }
}
```

### Use Case Tests

Target: Application logic in the usecase layer  
Method: Replace port/output with mocks/fakes

```go
// usecase/collect_trending_videos_test.go
func TestCollectTrendingVideos(t *testing.T) {
    // Arrange mocks
    mockYouTube := &MockYouTubeGateway{}
    mockRepo := &MockVideoRepository{}
    mockTasks := &MockCloudTasksGateway{}
    
    mockYouTube.On("FetchTrending", ctx, []int{27, 28}).
        Return(sampleVideos, nil)
    mockRepo.On("Save", ctx, mock.Anything).
        Return(nil)
    mockTasks.On("ScheduleSnapshots", ctx, mock.Anything).
        Return(nil)
    
    uc := NewCollectTrendingVideosUseCase(
        mockYouTube, mockRepo, mockTasks,
    )
    
    // Act
    err := uc.Execute(ctx, CollectTrendingInput{
        Categories: []int{27, 28},
        Keywords:   []Keyword{...},
    })
    
    // Assert
    assert.NoError(t, err)
    mockYouTube.AssertExpectations(t)
    mockRepo.AssertExpectations(t)
    mockTasks.AssertExpectations(t)
}
```

### Integration Tests

Target: handler → usecase → repository integration  
Method: Start PostgreSQL with testcontainers-go

```go
// adapter/gateway/postgres_video_repository_test.go
func TestPostgresVideoRepository_Integration(t *testing.T) {
    // Start PostgreSQL container with testcontainers
    ctx := context.Background()
    container, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(5 * time.Second),
        ),
    )
    require.NoError(t, err)
    defer container.Terminate(ctx)
    
    // Run migrations
    connStr, _ := container.ConnectionString(ctx)
    db, _ := sql.Open("postgres", connStr)
    runMigrations(db)
    
    // Repository tests
    repo := NewPostgresVideoRepository(db)
    
    t.Run("Save and FindByID", func(t *testing.T) {
        video := &domain.Video{
            ID:        domain.VideoID("test123"),
            Title:     "Test Video",
            ChannelID: domain.ChannelID("ch123"),
        }
        
        // Save
        err := repo.Save(ctx, video)
        assert.NoError(t, err)
        
        // Fetch
        found, err := repo.FindByID(ctx, video.ID)
        assert.NoError(t, err)
        assert.Equal(t, video.Title, found.Title)
    })
}
```

### Contract Tests

Target: Backward compatibility of Protocol Buffers schemas

