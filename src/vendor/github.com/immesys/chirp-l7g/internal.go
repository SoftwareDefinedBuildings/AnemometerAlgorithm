package chirpl7g

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/BTrDB/btrdb"
	"gopkg.in/immesys/bw2bind.v5"
)

const serverVK = "MT3dKUYB8cnIfsbnPrrgy8Cb_8whVKM-Gtg2qd79Xco="
const serverIP = "54.215.229.112:28590"
const btrdbIP = "128.32.37.198:4410"

type dataProcessingAlgorithm struct {
	BWCL          *bw2bind.BW2Client
	Vendor        string
	Algorithm     string
	Process       func(info *SetInfo, popHdr []*L7GHeader, h []*ChirpHeader, e Emitter)
	Initialize    func(e Emitter)
	Uncorrectable map[string]int
	Total         map[string]int
	Correctable   map[string]int
	LastRawInput  map[string]RawInputData
	EmitToStdout  bool
	DB            *btrdb.BTrDB
	DBChan        chan *OutputData
}

func (a *dataProcessingAlgorithm) MirrorToStandardOutput(v bool) {
	a.EmitToStdout = v
}
func runDPA(entitycontents []byte, iz func(e Emitter), cb func(info *SetInfo, popHdr []*L7GHeader, h []*ChirpHeader, e Emitter), vendor string, algorithm string) error {

	a := dataProcessingAlgorithm{}

	db, err := btrdb.Connect(context.Background(), btrdbIP)
	if err != nil {
		panic(err)
	}

	a.DB = db
	a.DBChan = make(chan *OutputData, 10000)
	a.BWCL = nil
	a.Process = cb
	a.Initialize = iz
	a.Vendor = vendor
	a.LastRawInput = make(map[string]RawInputData)
	a.Algorithm = algorithm

	a.Uncorrectable = make(map[string]int)
	a.Total = make(map[string]int)

	listener, err := net.Listen("tcp6", ":4000")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("listen err: %v\n", err)
			continue
		}
		go a.handleIncomingData(conn)
	}
}

func (a *dataProcessingAlgorithm) handleIncomingData(c net.Conn) {

	batchInfo := make(map[string]*SetInfo)
	batchL7G := make(map[string][]*L7GHeader)
	batchChirp := make(map[string][]*ChirpHeader)
	lastseq := make(map[string]int)

	for {
		barr := make([]byte, 70)
		_, err := io.ReadFull(c, barr)
		if err != nil {
			fmt.Printf("Read error: %v\n", err)
			return
		}
		//m.Dump()
		//parts := strings.Split(m.URI, "/")
		site := "tcp"

		//po := m.GetOnePODF(bw2bind.PODFL7G1Raw).(bw2bind.MsgPackPayloadObject)
		h := L7GHeader{}
		h.Srcmac = "tcpdirect"
		h.Payload = barr
		if h.Payload[0] != 9 {
			fmt.Printf("Skipping l7g packet type %d\n", h.Payload[0])
			continue
		}

		// if h.Payload[1] > 20 {
		// 	//Skip the xor packets for now
		// 	fmt.Println("Skipping xor packet", h.Payload[0])
		// 	continue
		// }

		ch := ChirpHeader{}
		isAnemometer := loadChirpHeader(h.Payload, &ch)
		if !isAnemometer {
			continue
		}

		did := h.Srcmac
		numasics := 4
		if ch.Build%10 == 7 {
			numasics = 6
		}
		lastseqi, ok := lastseq[did]
		if !ok {
			lastseqi = int(ch.Seqno - 10)
		}
		uncorrectablei, ok := a.Uncorrectable[h.Srcmac]
		if !ok {
			uncorrectablei = 0
		}
		lastseqi &= 0xFFFF
		lastsegment := lastseqi / 4
		currentsegment := ch.Seqno / 4
		if ch.Build%10 == 7 {
			lastsegment = lastseqi / 8
			currentsegment = ch.Seqno / 8
		}

		if int(currentsegment) != int(lastsegment) {
			//Send the last segment if it is not nil
			l7g := batchL7G[did]
			hdr := batchChirp[did]
			info := batchInfo[did]
			if !(info == nil || hdr == nil || l7g == nil) {
				//Maybe we have to send
				mustsend := false
				for i := 0; i < numasics; i++ {
					if hdr[i] != nil {
						mustsend = true
					}
				}
				if mustsend {
					complete := true
					var t time.Time
					hast := false
					for i := 0; i < numasics; i++ {
						if l7g[i] == nil {
							complete = false
						} else if !hast {
							t = time.Unix(0, l7g[i].Brtime)
							hast = true
						}
					}
					info.Complete = complete
					info.TimeOfFirst = t
					ri := RawInputData{
						SetInfo:      info,
						L7GHeaders:   l7g,
						ChirpHeaders: hdr,
					}
					a.LastRawInput[did] = ri
					a.Process(info, l7g, hdr, a)
				}
			}
			batchL7G[did] = make([]*L7GHeader, numasics)
			batchChirp[did] = make([]*ChirpHeader, numasics)
			batchInfo[did] = &SetInfo{
				Site:    site,
				MAC:     did,
				Build:   ch.Build,
				IsDuct:  ch.Build%10 == 5,
				IsRoom:  ch.Build%10 == 0,
				IsDuct6: ch.Build%10 == 7,
			}
		}

		//Save the sequence number
		lastseqi++
		lastseqi &= 0xFFFF
		if int(ch.Seqno) != lastseqi {
			//TODO this is not valid for 6 channel
			uncorrectablei++
		}
		lastseqi = int(ch.Seqno)
		lastseq[h.Srcmac] = lastseqi
		a.Uncorrectable[h.Srcmac] = uncorrectablei

		//Save the total packets
		totali, ok := a.Total[h.Srcmac]
		if !ok {
			totali = 0
		}
		totali++
		a.Total[h.Srcmac] = totali

		//Save this header
		batchL7G[did][ch.Primary] = &h
		batchChirp[did][ch.Primary] = &ch
	}
}
func (a *dataProcessingAlgorithm) Data(od OutputData) {
	od.Vendor = a.Vendor
	od.Algorithm = a.Algorithm

	od.Uncorrectable, _ = a.Uncorrectable[od.Sensor]
	od.Total, _ = a.Total[od.Sensor]
	od.Correctable, _ = a.Correctable[od.Sensor]
	od.RawInput = a.LastRawInput[od.Sensor]
	//spew.Dump(a.LastRawInput)
	//spew.Dump(od)
	// URI := fmt.Sprintf("ucberkeley/anem/%s/%s/%s/s.anemometer/%s/i.anemometerdata/signal/feed", od.Vendor, od.Algorithm, od.RawInput.SetInfo.Site, od.Sensor)
	// po, err := bw2bind.CreateMsgPackPayloadObject(bw2bind.PONumChirpFeed, od)
	// if err != nil {
	// 	panic(err)
	// }
	// doPersist := false
	// if od.Total%200 < 5 {
	// 	doPersist = true
	// }
	// err = a.BWCL.Publish(&bw2bind.PublishParams{
	// 	URI:            URI,
	// 	AutoChain:      true,
	// 	PayloadObjects: []bw2bind.PayloadObject{po},
	// 	Persist:        doPersist,
	// })
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Got publish error: %v\n", err)
	// } else {
	// 	//fmt.Println("Publish ok")
	// }
	select {
	case a.DBChan <- &od:
	default:
		fmt.Println("WARNING: DROPPING DATA BEFORE WRITING TO DB\n")
	}

	for _, ch := range od.RawInput.ChirpHeaders {
		for i := 0; i < 4; i++ {
			fmt.Printf("FROM %d TO %d MI=%2d %5d + %5di | %5d + %5di | %5d + %5di | %5d + %5di |\n", ch.Primary, i, ch.MaxIndex,
				ch.QValues[i][0], ch.IValues[i][0], ch.QValues[i][1], ch.IValues[i][1], ch.QValues[i][2], ch.IValues[i][2], ch.QValues[i][3], ch.IValues[i][3])
		}
	}
	/*if a.EmitToStdout {
		data, err := json.Marshal(od)
		if err != nil {
			panic(err)
		}
		fmt.Println("MIRROR_STDOUT:" + string(data))
	}*/
}

type StreamKey struct {
	Collection string
	Name       string
}

func (a *dataProcessingAlgorithm) processDBChannel() {

	streamcache := make(map[StreamKey]*btrdb.Stream)
	for od := range a.DBChan {
		_ = od
		_ = streamcache
		// fo
		// for src := 0; src < 4; src++{
		// 	for dst := 0; dst < 4; dst++ {
		// 		if src == dest {
		// 			continue
		// 		}
		// 		sk := StreamKey{
		// 			Collection:fmt.Sprintf("anem/%s/%s/tof", od.Algorithm, od.Sensor),
		// 			Name: fmt.Sprintf("%d_to_%d", src, dst),
		// 		}
		//
		// 		stream, ok := streamcache[sk]
		// 		if !ok {
		// 			//TODO create
		// 		}
		// 		stream.InsertTV(context.Background(), float64(od.Timestamp), float64(od.
		// 	}
		// }

	}
}

/*
typedef struct __attribute__((packed))
{
  uint8_t   l7type;     // 0
  uint8_t   type;       // 1
  uint16_t  seqno;      // 2:3
  uint8_t   primary;    // 4
  uint8_t   buildnum;   // 5
  int16_t   acc_x;      // 6:7
  int16_t   acc_y;      // 8:9
  int16_t   acc_z;      // 10:11
  int16_t   mag_x;      // 12:13
  int16_t   mag_y;      // 14:15
  int16_t   mag_z;      // 16:17
  int16_t   hdc_temp;   // 18:19
  int16_t   hdc_hum;    // 20:21
  uint8_t   max_index[3]; // 22:24
  uint8_t   parity;    // 25
  uint16_t  cal_res;   // 26:27
  //Packed IQ data for 4 pairs
  //M-3, M-2, M-1, M
  uint8_t data[3][16];  //28:75
} measure_set_t; //76 bytes
*/
var xormap map[int][]byte

func init() {
	xormap = make(map[int][]byte)
}
func check_xor(seqno int, arr []byte) {
	cmp := make([]byte, len(arr))
	for i := 0; i < 4; i++ {
		thisbuf, ok := xormap[seqno-i]
		if !ok {
			return
		}
		for x := 0; x < len(arr); x++ {
			cmp[x] ^= thisbuf[x]
		}
	}
	for x := 0; x < len(arr); x++ {
		cmp[x] ^= arr[x]
	}
	//fmt.Printf("XOR RESULT: %x\n", cmp)
}
func loadChirpHeader(arr []byte, h *ChirpHeader) bool {
	//Drop the type info we added
	if arr[0] != 9 {
		return false
	}

	h.Type = int(arr[1])
	h.Seqno = binary.LittleEndian.Uint16(arr[2:])

	if h.Type > 15 {
		//check_xor(int(h.Seqno), arr)
		return false
	}

	var paritycheck uint8 = 0
	for _, b := range arr {
		paritycheck ^= b
	}
	if paritycheck != 0 {
		fmt.Fprintf(os.Stderr, "Received packet failed parity check, ignoring %d\n", paritycheck)
		//expected_seqnos := arr[2] ^ (arr[2] - 1) ^ (arr[2] - 2) ^ (arr[2] - 3)
		//fmt.Printf("B: %d\n", paritycheck^arr[2])
		//fmt.Printf("C: %d\n", expected_seqnos)
		return false
	}
	xormap[int(h.Seqno)] = arr
	h.Build = int(arr[5])
	numasics := 4
	if h.Build%10 == 7 {
		numasics = 6
	}
	h.CalPulse = 160
	h.Primary = arr[4]
	h.CalRes = make([]int, numasics)
	for i := 0; i < numasics; i++ {
		if i == (int(h.Primary)) {
			h.CalRes[i] = int(binary.LittleEndian.Uint16(arr[26:]))
		} else {
			h.CalRes[i] = -1
		}
	}

	h.MaxIndex = make([]int, numasics)
	h.IValues = make([][]int, numasics)
	h.QValues = make([][]int, numasics)
	offset := 0
	//fmt.Printf("arr: %x\n", arr)

	for i := 0; i < numasics; i++ {
		if h.Build%10 == 5 || h.Build%10 == 0 {
			if i == int(h.Primary) {
				h.MaxIndex[i] = -1
				continue
			}
		}
		if h.Build%10 == 7 {
			if int(h.Primary) < 3 && i < 3 {
				h.MaxIndex[i] = -1
				continue
			}
			if int(h.Primary) >= 3 && i >= 3 {
				h.MaxIndex[i] = -1
				continue
			}
		}
		h.MaxIndex[i] = int(arr[22+offset])
		h.IValues[i] = make([]int, 4)
		h.QValues[i] = make([]int, 4)
		for j := 0; j < 4; j++ {
			h.QValues[i][j] = int(int16(binary.LittleEndian.Uint16(arr[28+offset*16+j*4:])))
			h.IValues[i][j] = int(int16(binary.LittleEndian.Uint16(arr[28+offset*16+j*4+2:])))
		}
		offset++
	}

	f_acc_x := int16(binary.LittleEndian.Uint16(arr[6:]))
	f_acc_y := int16(binary.LittleEndian.Uint16(arr[8:]))
	f_acc_z := int16(binary.LittleEndian.Uint16(arr[10:]))
	f_mag_x := int16(binary.LittleEndian.Uint16(arr[12:]))
	f_mag_y := int16(binary.LittleEndian.Uint16(arr[14:]))
	f_mag_z := int16(binary.LittleEndian.Uint16(arr[16:]))
	//f_hdc_tmp := int16(binary.LittleEndian.Uint16(arr[18:]))
	//f_hdc_rh := binary.LittleEndian.Uint16(arr[20:])

	if f_acc_x >= 8192 {
		f_acc_x -= 16384
	}
	if f_acc_y >= 8192 {
		f_acc_y -= 16384
	}
	if f_acc_z >= 8192 {
		f_acc_z -= 16384
	}
	h.Accelerometer = []float64{
		float64(f_acc_x) * 0.244,
		float64(f_acc_y) * 0.244,
		float64(f_acc_z) * 0.244,
	}
	h.Magnetometer = []float64{
		float64(f_mag_x) * 0.1,
		float64(f_mag_y) * 0.1,
		float64(f_mag_z) * 0.1,
	}
	//h.Humidity = float64(f_hdc_rh) / 100.0
	//h.Temperature = float64(f_hdc_tmp) / 100.0
	//fmt.Printf("Received frame:\n")
	//spew.Dump(h)
	return true
}
