package main

import (
	"fmt"
	"math"
	"os"
	"time"

	l7g "github.com/immesys/chirp-l7g"
)

//You can put some state variables and type definitions here
type path struct {
	Src int
	Dst int
}

type algorithmstate struct {
	//This is just example information copied from python
	//You can put anything you want here
	AnemometerID       string
	IsDuct             bool
	CalibrationPeriod  time.Duration
	IsCalibrating      bool
	PrevRelPhase       map[path]float64
	PrevAbsPhase       map[path]float64
	CurAbsPhase        map[path]float64
	CalibrationPhases  map[path][]float64
	CalibrationIndexes map[path][]int
	CalibratedIndex    map[path]int
	Paths              []path
}

var states map[string]*algorithmstate

func main() {
	fmt.Printf("Version 3.6 ====\n")
	localtty := ""
	if len(os.Args) > 1 {
		localtty = os.Args[1]
	}
	//Register and run our algorithm.
	err := l7g.RunDPA([]byte(ourEntity), Initialize, OnNewDataSet,
		//This is the algorithm vendor
		"ucberkeley",
		//This is the algorithm version
		"1.0",
		//This is the address of the local connection or blank to disable
		localtty)
	fmt.Printf("fatal error: %v\n", err)
}

func Initialize(emit l7g.Emitter) {
	states = make(map[string]*algorithmstate)
	//If you want the algorithm output to be on standard out as well
	//to allow for local use, do this
	emit.MirrorToStandardOutput(true)
}

// OnNewDataSet encapsulates the algorithm. You can store the emitter and
// use it asynchronously if required. You can see the documentation for the
// parameters at https://godoc.org/github.com/immesys/chirp-l7g
// The popHdr and data arrays are always length 4 and represent the data
// for four sets of readings with primary = 0..3. They are always contiguous
// in time. If the set is incomplete due to data loss, the missing elements
// will have nil in them.
// Sometimes the popHdr will be nil for an element even if the data is present
// this indicates a packet that was lost but reconstructed using forward
// error correction codes.
func OnNewDataSet(info *l7g.SetInfo, popHdr []*l7g.L7GHeader, data []*l7g.ChirpHeader, emit l7g.Emitter) {

	//We only want to process complete sets of data
	// if !info.Complete {
	// 	return
	// }

	//This string is the complete ID of this anemometer
	//You can use this as a key into buffers of historic state
	idstring := info.MAC
	state, found := states[idstring]
	if !found {
		//Initialize new algorithm state here
		//perhaps make it calibrate or something
		state = &algorithmstate{
			AnemometerID: idstring,
			IsDuct:       info.IsDuct,
		}
		states[idstring] = state
	}

	//Initialize the output data object
	//We will build up the data inside this as we process the input data
	outputdata := l7g.OutputData{
		Timestamp: info.TimeOfFirst.UnixNano(),
		Sensor:    info.MAC,
	}

	//Over all paths
	//fmt.Printf("data len is %d\n", len(data))
	// for src := 0; src < len(data); src++ {
	// 	for dst := 0; dst < len(data); dst++ {
	// 		if info.IsDuct6 {
	// 			if src < 3 && dst < 3 {
	// 				continue
	// 			}
	// 			if dst >= 3 && dst >= 3 {
	// 				continue
	// 			}
	// 		} else {
	// 			if src == dst {
	// 				//We don't use data from primary
	// 				continue
	// 			}
	// 		}
	// 		// spew.Dump(info)
	// 		// fmt.Printf("src=%d dst=%d\n", src, dst)
	// 		// spew.Dump(data)
	// 		p := path{src, dst}
	// 		maxIndex := data[src].MaxIndex[dst]
	// 		if maxIndex < 3 {
	// 			//This might indicate something is wrong
	// 			//It also means that our IQ values start from 0
	// 			//not from maxIndex-3
	// 		}
	// 		Ivalue2beforeMax := data[src].IValues[dst][1] //0 is 3 before
	// 		Qvalue2beforeMax := data[src].QValues[dst][1]
	//
	// 		//Use this as well as other data
	// 		_ = p
	// 		_ = maxIndex
	// 		_ = Ivalue2beforeMax
	// 		_ = Qvalue2beforeMax
	//
	// 		//To populate outputdata
	// 		outputdata.Tofs = append(outputdata.Tofs, l7g.TOFMeasure{
	// 			Src: src,
	// 			Dst: dst,
	// 			//Here you should actually put the time of flight in microseconds
	// 			Val: 555,
	// 		})
	// 		outputdata.Temperatures = append(outputdata.Temperatures, l7g.TempMeasure{
	// 			Src: src,
	// 			Dst: dst,
	// 			//Here you should actually put the temperature in celsius
	// 			Val: 25.5,
	// 		})
	// 	}
	// }

	//If you have a new velocity estimate, you should also add that here
	//These are in m/s
	outputdata.Velocities.Valid = true
	outputdata.Velocities.X = 1 //due north
	outputdata.Velocities.Y = 2 //due west
	outputdata.Velocities.Z = 3 //vertical
	//Also in m/s
	outputdata.Velocities.Mag = math.Sqrt(1 + 4 + 9)
	// Phi, the azimuthal angle is the degrees counterclockwize from North (X)
	outputdata.Velocities.Phi = 55
	// Theta, the polar angle is the degrees from vertical
	outputdata.Velocities.Theta = 66

	//If you have other information output from the algorithm you can
	//add that in here
	outputdata.Extradata = append(outputdata.Extradata, "the algorithm has not been filled in yet")

	//Emit the data on the data bus (and if MirrorToStandardOut is true, also
	//to standard output)
	emit.Data(outputdata)
}
