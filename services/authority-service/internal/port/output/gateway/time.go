package gateway

import "time"

// Clock abstracts time for use cases to enable deterministic tests.
type Clock interface {
    Now() time.Time
}

