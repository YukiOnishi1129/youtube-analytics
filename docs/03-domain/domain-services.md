# Domain Services & Policies

Implementation of logic spanning multiple aggregates and business rules.

## PatternBuilder (Pattern Builder)

### Purpose
Generate regular expression patterns from keyword names to absorb notation variations

### Implementation Example
```go
type PatternBuilder interface {
    BuildPattern(name string) string
}

type patternBuilder struct{}

func (pb *patternBuilder) BuildPattern(name string) string {
    // Escape basic string
    escaped := regexp.QuoteMeta(name)
    
    // Build pattern to absorb notation variations
    pattern := ""
    for _, char := range escaped {
        switch {
        case isAlphabetic(char):
            // For alphabetic chars, match both uppercase and lowercase
            upper := strings.ToUpper(string(char))
            lower := strings.ToLower(string(char))
            pattern += fmt.Sprintf("[%s%s]", upper, lower)
            
        case isNumeric(char):
            // For numeric chars, match both half-width and full-width
            halfWidth := string(char)
            fullWidth := toFullWidth(halfWidth)
            pattern += fmt.Sprintf("[%s%s]", halfWidth, fullWidth)
            
        case isPunctuation(char):
            // For punctuation, also match spaces or other punctuation
            pattern += "[\\s\\-_\\.]*"
            
        default:
            pattern += string(char)
        }
    }
    
    // Consider word boundaries
    return fmt.Sprintf("\\b%s\\b", pattern)
}
```

### Usage Examples
```go
// Input: "Next.js"
// Output: "\\b[Nn][Ee][Xx][Tt][\\s\\-_\\.]*[Jj][Ss]\\b"
// Matches: "Next.js", "next.js", "NEXT.JS", "Next js", "next-js"

// Input: "React"
// Output: "\\b[Rr][Ee][Aa][Cc][Tt]\\b"
// Matches: "React", "react", "REACT"

// Input: "Go language"
// Output: "\\b[Gg][Oo] language\\b"
// Matches: "Go language", "go language", "GO language"
```

## MetricsCalculator (Metrics Calculator)

### Purpose
Calculate various metrics from snapshot data

### Implementation

#### Growth Rate Calculation
```go
func (mc *MetricsCalculator) Growth(
    valueAtX int64,
    valueAt0 int64,
    hours int,
) float64 {
    if hours == 0 {
        return 0
    }
    return float64(valueAtX - valueAt0) / float64(hours)
}

// Usage example
viewGrowthRate := mc.Growth(views24h, views0h, 24)
// views24h=10000, views0h=1000, hours=24
// → (10000-1000)/24 = 375 views/hour
```

#### Engagement Level Calculation
```go
func (mc *MetricsCalculator) ViewsPerSub(
    viewsX int64,
    subsX int64,
) float64 {
    if subsX == 0 {
        return 0
    }
    return float64(viewsX) / float64(subsX)
}

// Usage example
relativeViews := mc.ViewsPerSub(views24h, subs24h)
// views24h=10000, subs24h=1000
// → 10000/1000 = 10.0 (each subscriber watched 10 times)
```

#### Wilson Score Lower Bound
```go
func (mc *MetricsCalculator) WilsonLowerBound(
    likes int64,
    views int64,
) float64 {
    if views == 0 {
        return 0
    }
    
    n := float64(views)
    p := float64(likes) / n
    
    if n < 30 {
        // Sample size too small
        return 0
    }
    
    // Wilson score calculation (95% confidence interval)
    z := 1.96 // Z-value for 95% confidence interval
    denominator := 1 + z*z/n
    
    centre := p + z*z/(2*n)
    deviation := z * math.Sqrt((p*(1-p)+z*z/(4*n))/n)
    
    lowerBound := (centre - deviation) / denominator
    
    if lowerBound < 0 {
        return 0
    }
    return lowerBound
}

// Usage example
quality := mc.WilsonLowerBound(likes24h, views24h)
// likes24h=100, views24h=1000
// → 0.082 (confidence lower bound of like rate ~8.2%)
```

#### LPS Shrunk Score
```go
func (mc *MetricsCalculator) LpsShrunk(
    likes int64,
    subs int64,
) float64 {
    const (
        SCALE  = 1000
        OFFSET = 500
    )
    
    numerator := float64(likes + OFFSET)
    denominator := float64(subs + SCALE)
    
    return numerator / denominator
}

// Usage example
heat := mc.LpsShrunk(likes24h, subs24h)
// likes24h=100, subs24h=1000
// → (100+500)/(1000+1000) = 0.3
// Fair evaluation even for channels with few subscribers
```

## ExclusionPolicy (Exclusion Policy)

### Purpose
Determine exclude_from_ranking based on checkpoint-specific minimum thresholds

### Implementation
```go
type ExclusionThresholds struct {
    MinViews int64
    MinSubs  int64
}

var defaultThresholds = map[int]ExclusionThresholds{
    3:   {MinViews: 100,  MinSubs: 100},
    6:   {MinViews: 300,  MinSubs: 100},
    12:  {MinViews: 500,  MinSubs: 100},
    24:  {MinViews: 1000, MinSubs: 100},
    48:  {MinViews: 2000, MinSubs: 100},
    72:  {MinViews: 3000, MinSubs: 100},
    168: {MinViews: 5000, MinSubs: 100},
}

func (ep *ExclusionPolicy) ShouldExclude(
    checkpointHour int,
    viewsCount int64,
    subsCount int64,
) bool {
    threshold, exists := defaultThresholds[checkpointHour]
    if !exists {
        return false
    }
    
    // Exclude if any threshold is not met
    if viewsCount < threshold.MinViews {
        return true
    }
    if subsCount < threshold.MinSubs {
        return true
    }
    
    return false
}
```

### Usage Examples
```go
// 24-hour point with 800 views, 50 subscribers
// → Subscriber count below threshold (100), so exclude
exclude := ep.ShouldExclude(24, 800, 50) // true

// 24-hour point with 2000 views, 500 subscribers
// → Both above thresholds, so include in ranking
exclude := ep.ShouldExclude(24, 2000, 500) // false
```

## FilterService (Filter Service)

### Purpose
Evaluate video titles with keyword filters

### Implementation
```go
type FilterResult int

const (
    FilterResultExclude FilterResult = -1
    FilterResultNeutral FilterResult = 0
    FilterResultInclude FilterResult = 1
)

type FilterService interface {
    // Evaluate title with filters
    Filter(title string, keywords []FilterKeyword) FilterResult
}

type filterService struct{}

func (fs *filterService) Filter(
    title string,
    keywords []FilterKeyword,
) FilterResult {
    // Exclude filters take precedence
    for _, keyword := range keywords {
        if !keyword.Enabled {
            continue
        }
        
        if keyword.FilterType == "exclude" {
            matched, _ := regexp.MatchString(keyword.Pattern, title)
            if matched {
                return FilterResultExclude
            }
        }
    }
    
    // Check include filters
    for _, keyword := range keywords {
        if !keyword.Enabled {
            continue
        }
        
        if keyword.FilterType == "include" {
            matched, _ := regexp.MatchString(keyword.Pattern, title)
            if matched {
                return FilterResultInclude
            }
        }
    }
    
    // Matches neither
    return FilterResultNeutral
}
```

### Usage Examples
```go
keywords := []FilterKeyword{
    {Name: "React", FilterType: "include", Pattern: "[Rr][Ee][Aa][Cc][Tt]", Enabled: true},
    {Name: "singing cover", FilterType: "exclude", Pattern: "singing cover", Enabled: true},
}

// "React Tutorial" → Include (matches "React" keyword)
result := fs.Filter("React Tutorial", keywords)

// "React Singing Cover" → Exclude (exclusion takes precedence)
result := fs.Filter("React Singing Cover", keywords)

// "Vue.js Tutorial" → Neutral (matches neither)
result := fs.Filter("Vue.js Tutorial", keywords)
```

## TaskIDGenerator (Task ID Generator)

### Purpose
Generate deterministic TaskIDs for idempotency guarantee

### Implementation
```go
type TaskIDGenerator interface {
    GenerateSnapshotTaskID(videoId string, checkpointHour int) string
    GenerateTrendingTaskID(category int, timestamp time.Time) string
}

type taskIDGenerator struct{}

func (g *taskIDGenerator) GenerateSnapshotTaskID(
    videoId string,
    checkpointHour int,
) string {
    // Generate deterministic ID from video ID and checkpoint
    data := fmt.Sprintf("snapshot:%s:%d", videoId, checkpointHour)
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:16]) // 128bit
}

func (g *taskIDGenerator) GenerateTrendingTaskID(
    category int,
    timestamp time.Time,
) string {
    // Generate deterministic ID from category and time
    // Round to minute to allow short-term duplicates
    rounded := timestamp.Truncate(time.Minute)
    data := fmt.Sprintf("trending:%d:%s", category, rounded.Format(time.RFC3339))
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:16])
}
```

### Usage Examples
```go
// Generate snapshot task ID
taskId := gen.GenerateSnapshotTaskID("abc123", 24)
// Same input always generates same ID
// → "a3f5b2c1d4e6f789" (example)

// Generate trending task ID
taskId := gen.GenerateTrendingTaskID(27, time.Now())
// Duplicate execution within 1 minute gets same ID
// → "b2c3d4e5f6a7b8c9" (example)
```

## Domain Policy Application Points

| Policy | Application Point | Purpose |
|---------|------------------|---------|
| PatternBuilder | FilterKeyword Aggregate | Pattern generation during keyword registration |
| MetricsCalculator | VideoMetrics Aggregate | Metric calculation from snapshots |
| ExclusionPolicy | VideoMetrics Aggregate | Low-sample determination |
| FilterService | Application Layer | Include/exclude determination during video collection |
| TaskIDGenerator | Application Layer | ID generation during Cloud Tasks registration |
