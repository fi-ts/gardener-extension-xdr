package charts

import (
	"embed"
)

// InternalChart embeds the cortex chart in embed.FS
//
//go:embed internal
var InternalChart embed.FS

const (
	// CortexChartsPath is the path to the internal charts.
	CortexChartsPath = "cortex"
	CortexNamespace  = "kube-system"
	CortextName      = "cortex-agent"
)
