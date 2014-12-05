package natbeat

import (
	"os"
	"time"

	"github.com/cloudfoundry/gunk/diegonats"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/restart"
)

func NewBackgroundHeartbeat(natsClient diegonats.NATSClient, natsAddress, natsUsername, natsPassword string, logger lager.Logger, registration RegistryMessage) ifrit.RunFunc {
	return func(signals <-chan os.Signal, ready chan<- struct{}) error {
		restarter := restart.Restarter{
			Runner: newBackgroundGroup(natsClient, natsAddress, natsUsername, natsPassword, logger, registration),
			Load: func(runner ifrit.Runner, err error) ifrit.Runner {
				return newBackgroundGroup(natsClient, natsAddress, natsUsername, natsPassword, logger, registration)
			},
		}
		// don't wait, start this thing in the background
		close(ready)
		return restarter.Run(signals, make(chan struct{}))
	}
}

func newBackgroundGroup(natsClient diegonats.NATSClient, natsAddress, natsUsername, natsPassword string, logger lager.Logger, registration RegistryMessage) ifrit.Runner {
	return grouper.NewOrdered(os.Interrupt, grouper.Members{
		{"nats_connection", diegonats.NewClientRunner(natsAddress, natsUsername, natsPassword, logger, natsClient)},
		{"router_heartbeat", New(client, registration, 50*time.Millisecond, logger)},
	})
}
