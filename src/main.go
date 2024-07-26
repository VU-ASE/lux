package main

import (
	"os"
	"time"

	pb_module_outputs "github.com/VU-ASE/rovercom/packages/go/outputs"
	pb_core_messages "github.com/VU-ASE/rovercom/packages/go/core"
	servicerunner "github.com/VU-ASE/roverlib/src"
	"github.com/d2r2/go-bh1750"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	zmq "github.com/pebbe/zmq4"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func run(
	serviceInfo servicerunner.ResolvedService,
	sysMan servicerunner.SystemManagerInfo,
	initialTuning *pb_core_messages.TuningState) error {

	// Set up logging
	err := logger.ChangePackageLogLevel("bh1750", logger.InfoLevel)
	if err != nil {
		return err
	}
	err = logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
	if err != nil {
		return err
	}

	// Fetch address to output on from service yaml
	addr, err := serviceInfo.GetOutputAddress("lux-output")
	if err != nil {
		return err
	}
	log.Info().Msgf("Publishing lux data to address: '%s'", addr)

	// Create ZMQ socket to publish on
	publisher, _ := zmq.NewSocket(zmq.PUB)
	defer publisher.Close()
	err = publisher.Bind(addr)
	if err != nil {
		return err
	}

	// Create i2c device to use
	i2c, err := i2c.NewI2C(0x23, 5)
	if err != nil {
		log.Err(err)
	}
	defer i2c.Close()

	// Initialize driver using the just configured i2c device
	sensor := bh1750.NewBH1750()
	err = sensor.Reset(i2c)
	if err != nil {
		log.Err(err)
	}

	// Configure driver
	err = sensor.ChangeSensivityFactor(i2c, sensor.GetDefaultSensivityFactor())
	if err != nil {
		log.Err(err)
	}
	resolution := bh1750.HighResolution

	// Main loop
	for {
		// Read sensor data
		amb, err := sensor.MeasureAmbientLight(i2c, resolution)
		if err != nil {
			log.Err(err)
		}

		// Create a pb message and wrap it
		msgData := &pb_module_outputs.LuxSensorOutput{
			Lux: int32(amb),
		}
		wrapperData := &pb_module_outputs.SensorOutput{
			SensorId:  1,
			Timestamp: uint64(time.Now().UnixMilli()),
			SensorOutput: &pb_module_outputs.SensorOutput_LuxOutput{
				LuxOutput: msgData,
			},
		}

		// Marshal the message and send it
		msgBytes, err := proto.Marshal(wrapperData)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal protobuf")
			continue
		}
		log.Debug().Msgf("Ambient light (%s) = %v lx", resolution, amb)

		_, err = publisher.SendBytes(msgBytes, 0)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send zmq message")
			continue
		}

		// Don't overload the sensor
		time.Sleep(time.Millisecond * 150)
	}
}

func tuningCallback(newtuning *pb_core_messages.TuningState) {
	log.Warn().Msg("Tuning state changed, but this module does not have any tunable parameters")
}

func onTerminate(sig os.Signal) {
	log.Info().Msg("Terminating")
}

func main() {
	servicerunner.Run(run, tuningCallback, onTerminate, false)
}
