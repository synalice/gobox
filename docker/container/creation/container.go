package container

import (
	"time"
)

type Container struct {
	ID        string
	TimeLimit time.Duration
}
