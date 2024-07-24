package main

import (
	pb_systemmanager_messages "github.com/VU-ASE/pkg-CommunicationDefinitions/v2/packages/go/systemmanager"
	servicerunner "github.com/VU-ASE/pkg-ServiceRunner/v2/src"

	"github.com/rs/zerolog/log"
)

func run(
	service servicerunner.ResolvedService,
	sysMan servicerunner.SystemManagerInfo,
	initialTuning *pb_systemmanager_messages.TuningState) error {

	log.Info().Str("Planet", "Earth").Msg("Hello world")

	//TODO: Implement the service logic here. Likely this will involve creating a pub/sub and some main logic.
	//      The de facto standard is to have some read (zmq/IO), some handling logic (may be several items),
	//      and some write (zmq/IO). The go routines typically communicate via channels.

	return nil
}

func onTuningState(newtuning *pb_systemmanager_messages.TuningState) {
	log.Info().Str("Value", newtuning.String()).Msg("Received tuning state from system manager")
	//TODO: Update this service based on the new tuning state
}

func main() {
	servicerunner.Run(run, onTuningState, false)
}
