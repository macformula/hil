package main

import (
	"encoding/json"
	"fmt"
	"github.com/foxglove/mcap/go/mcap"
	"os"
	"time"
)

type MyData struct {
	Data []byte `json:"data"`
}

type mcapmessage struct {
	//Name string
	Data int
}

type plotter struct {
}

func main() {
	//buf := &bytes.Buffer{}
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006.01.02_15.04.05")
	fileName := fmt.Sprintf("file_%s.mcap", formattedTime)
	file, _ := os.Create(fileName)

	w, _ := mcap.NewWriter(file, &mcap.WriterOptions{
		Chunked: false,
	})

	defer func() {
		err := w.Close()
		if err != nil {
			fmt.Println("Error closing mcap file:", err)
			return
		}
	}()
	err := w.WriteHeader(&mcap.Header{})
	if err != nil {
		panic("FAILED")
	}
	err = w.WriteSchema(&mcap.Schema{
		ID:       1,
		Name:     "schema",
		Encoding: "jsonschema",
		Data:     []byte(`{"type":"object"}`),
	})
	if err != nil {
		panic("FAILED")
	}
	err = w.WriteChannel(&mcap.Channel{
		ID:              1,
		Topic:           "/test",
		MessageEncoding: "json",
		SchemaID:        0,
		Metadata: map[string]string{
			"callerid": "100", // cspell:disable-line
		},
	})
	if err != nil {
		panic("FAILED")
	}
	err = w.WriteChannel(&mcap.Channel{
		ID:              2,
		Topic:           "/test2",
		MessageEncoding: "json",
		SchemaID:        0,
	})
	if err != nil {
		panic("FAILED")
	}
	data := MyData{
		Data: []byte{5, 6, 7, 8},
		//Data: []byte{1},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic("FAILED")
	}
	err = w.WriteMessage(&mcap.Message{
		ChannelID: 1,
		Sequence:  0,
		//LogTime:   0,
		LogTime: uint64(time.Now().Nanosecond()),
		//PublishTime: 0,
		PublishTime: uint64(time.Now().Nanosecond()),
		Data:        jsonData,
	})
	if err != nil {
		panic("FAILED")
	}
	//data = MyData{
	//	Data: []byte{1, 2, 3, 4},
	//}
	//jsonData, err = json.Marshal(data)
	//err = w.WriteMessage(&mcap.Message{
	//	ChannelID:   2,
	//	Sequence:    0,
	//	LogTime:     100,
	//	PublishTime: 100,
	//	Data:        jsonData,
	//})
	//if err != nil {
	//	panic("FAILED")
	//}

	t_holder := uint64(time.Now().UnixNano())
	//t_hold := t_holder

	for i := 1; i <= 20; i++ {
		//channelid := uint16(1)
		//if i%2 == 1 {
		//	channelid = uint16(2)
		//}

		//what is happening here
		m := mcapmessage{
			//"Engine",
			i * 5}
		d, err := json.Marshal(m)

		//fmt.Println("m: ", m.Data)
		//fmt.Println("d: ", d)

		if err != nil {
			fmt.Println("Error marshalling message data into json format:", err)
			return
		}

		//binary.LittleEndian.PutUint64(bytearray, uint64(i*5))
		//t := uint64(time.Now().Nanosecond())

		t := uint64(time.Now().UnixNano())
		//t := uint64(time.Now().Second() + i*1000000)
		//fmt.Println("Current time ", i, ":", time.Now().Second())
		//t := uint64(i * 1000000000)
		err = w.WriteMessage(&mcap.Message{
			ChannelID:   2,
			LogTime:     t - t_holder, //Wrong time formatting
			PublishTime: t - t_holder,
			Data:        d,
		})
		if err != nil {
			fmt.Println("Error writing message to mcap file on iteration", i, " :", err)
			return
		}
		time.Sleep(1000 * time.Millisecond)
		//time.Sleep(3000 * time.Millisecond)
	}

	//fmt.Println(buf)
}
