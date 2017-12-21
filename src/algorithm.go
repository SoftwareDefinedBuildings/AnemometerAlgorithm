package main

import (
	l7g "github.com/immesys/chirp-l7g"
)

func main() {
	//Register and run our algorithm.
	l7g.RunDPA([]byte(ourEntity), Initialize, OnNewData,
		//This is the algorithm vendor
		"ucberkeley",
		//This is the algorithm version
		"1.0")
}

func Initialize(emit l7g.Emitter) {
	//We actually do not do any initialization in this implementation, but if
	//you want to, you can do it here.
}

// OnNewData encapsulates the algorithm. You can store the emitter and
// use it asynchronously if required. You can see the documentation for the
// parameters at https://godoc.org/github.com/immesys/chirp-l7g
func OnNewData(popHdr *l7g.L7GHeader, h *l7g.ChirpHeader, emit l7g.Emitter) {
	/*
		// Define some magic constants for the algorithm
		magic_count_tx := -4

		// Create our output data set. For this reference implementation,
		// we emit one TOF measurement for every raw TOF sample (no averaging)
		// so the timestamp is simply the raw timestamp obtained from the
		// Border Router. We also identify the sensor simply from the mac address
		// (this is fine for most cases)
		odata := l7g.OutputData{
			Timestamp: popHdr.Brtime,
			Sensor:    popHdr.Srcmac,
		}

		// For each of the four measurements in the data set
		for set := 0; set < 4; set++ {
			// For now, ignore the data read from the ASIC in TXRX
			if int(h.Primary) == set {
				continue
			}

			// alias the data for readability. This is the 70 byte dataset
			// read from the ASIC
			data := h.Data[set]

			//The first six bytes of the data
			tof_sf := binary.LittleEndian.Uint16(data[0:2])
			tof_est := binary.LittleEndian.Uint16(data[2:4])
			intensity := binary.LittleEndian.Uint16(data[4:6])

			//Load the complex numbers
			iz := make([]int16, 16)
			qz := make([]int16, 16)
			for i := 0; i < 16; i++ {
				qz[i] = int16(binary.LittleEndian.Uint16(data[6+4*i:]))
				iz[i] = int16(binary.LittleEndian.Uint16(data[6+4*i+2:]))
			}

			//Find the largest complex magnitude (as a square). We do this like this
			//because it more closely mirror how it would be done on an embedded device
			// (actually because I copied the known-good firestorm implementation)
			magsqr := make([]uint64, 16)
			magmax := uint64(0)
			for i := 0; i < 16; i++ {
				magsqr[i] = uint64(int64(qz[i])*int64(qz[i]) + int64(iz[i])*int64(iz[i]))
				if magsqr[i] > magmax {
					magmax = magsqr[i]
				}
			}

			//Find the first index to be greater than half the max (quarter the square)
			quarter := magmax / 4
			less_idx := 0
			greater_idx := 0
			for i := 0; i < 16; i++ {
				if magsqr[i] < quarter {
					less_idx = i
				}
				if magsqr[i] > quarter {
					greater_idx = i
					break
				}
			}

			//Convert the squares into normal floating point
			less_val := math.Sqrt(float64(magsqr[less_idx]))
			greater_val := math.Sqrt(float64(magsqr[greater_idx]))
			half_val := math.Sqrt(float64(quarter))
			//CalPulse is in microseconds
			freq := float64(tof_sf) / 2048 * float64(h.CalRes[set]) / (float64(h.CalPulse) / 1000)
			//Linearly interpolate the index (the index is related to time of flight because it is regularly sampled)
			lerp_idx := float64(less_idx) + (half_val-less_val)/(greater_val-less_val)
			//Fudge the result with magic_count_tx and turn into time of flight
			//The *2 here is debugging, its not real
			tof := (lerp_idx + float64(magic_count_tx)) / freq * 8 * 2

			//We print these just for fun / debugging, but this is not actually emitting the data
			fmt.Printf("SEQ %d ASIC %d primary=%d\n", h.Seqno, set, h.Primary)
			fmt.Println("lerp_idx: ", lerp_idx)
			fmt.Println("tof_sf: ", tof_sf)
			fmt.Println("freq: ", freq)
			fmt.Printf("tof: %.2f us\n", tof*1000000)
			fmt.Println("intensity: ", intensity)
			fmt.Println("tof chip estimate: ", tof_est)
			fmt.Println("tof 50us estimate: ", lerp_idx*50)
			fmt.Println("data: ")
			for i := 0; i < 16; i++ {
				fmt.Printf(" [%2d] %6d + %6di (%.2f)\n", i, qz[i], iz[i], math.Sqrt(float64(magsqr[i])))
			}
			fmt.Println(".")

			//Append this time of flight to the output data set
			//For more "real" implementations, this would likely
			//be a rolling-window smoothed time of flight. You do not have
			//to base this value on just the data from this set and
			//you do not have to emit every time either (downsampling is ok)
			odata.Tofs = append(odata.Tofs, l7g.TOFMeasure{
				Src: int(h.Primary),
				Dst: set,
				Val: tof * 1000000,
			})
		} //end for each of the four measurements

		// Now we would also emit the velocities. I imagine this would use
		// the averaged/corrected time of flights that are emitted above
		// (when they are actually averaged/corrected)
		// For now, just a placeholder
		odata.Velocities = append(odata.Velocities, l7g.VelocityMeasure{X: 42, Y: 43, Z: 44})

		//Emit the data on the SASC bus
		emit.Data(odata)

	*/
}
