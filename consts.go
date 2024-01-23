package violin

import "math"

const (
	DefaultViolinMaxWorkerNum      = math.MaxInt32
	DefaultViolinWorkerIdleTimeout = 0
)

const (
	_ uint32 = iota
	statusInitialized
	statusPlaying
	statusCleaning
	statusShutdown
)
