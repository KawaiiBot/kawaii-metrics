package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	NormalCommandsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kawaiibot",
			Name: "normal_commands_total",
			Help: "Total number of normal commands used, partitioned by command name",
		},
		[]string{"name"},
	)
	NSFWCommandsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kawaiibot",
			Name: "nsfw_commands_total",
			Help: "Total number of NSFW commands used, partitioned by command name",
		},
		[]string{"name"},
	)
)

func init() {
	prometheus.MustRegister(NormalCommandsTotal)
	prometheus.MustRegister(NSFWCommandsTotal)
}
