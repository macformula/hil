package main

import (
	"encoding/binary"
	"fmt"
	"github.com/foxglove/mcap/go/mcap"
	"os"
	"time"
)

func initWriter() *mcap.Writer {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006.01.02_15.04.05")
	fileName := fmt.Sprintf("file_%s.txt", formattedTime)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil
	}
	mcapwriter, mcaperr := mcap.NewWriter(file, &mcap.WriterOptions{
		Chunked: false,
	})

	if mcaperr != nil {
		fmt.Println("Error creating mcap file:", mcaperr)
		return nil
	}
	return mcapwriter
}

func main() {
	writer := initWriter()
	defer func() {
		err := writer.Close()
		if err != nil {
			fmt.Println("Error closing mcap file:", err)
			return
		}
	}()
	err := writer.WriteHeader(&mcap.Header{})
	if err != nil {
		fmt.Println("Error closing mcap file:", err)
		return
	}

	err = writer.WriteSchema(&mcap.Schema{
		ID:       1,
		Name:     "schema",
		Encoding: "msg",
		Data:     []byte{},
	})

	if err != nil {
		fmt.Println("Error creating mcap schema:", err)
		return
	}

	err = writer.WriteChannel(&mcap.Channel{
		ID:              1,
		SchemaID:        1,
		Topic:           "Engine Speed (RPM)",
		MessageEncoding: "msg",
	})

	if err != nil {
		fmt.Println("Error creating Engine Speed channel:", err)
		return
	}

	err = writer.WriteChannel(&mcap.Channel{
		ID:              2,
		SchemaID:        1,
		Topic:           "Torque (N-m)",
		MessageEncoding: "msg",
	})

	if err != nil {
		fmt.Println("Error creating Torque channel:", err)
		return
	}

	bytearray := make([]byte, 8)

	for i := 0; i < 20; i++ {
		channelid := uint16(1)
		if i%2 == 1 {
			channelid = uint16(2)
		}
		binary.LittleEndian.PutUint64(bytearray, uint64(i*5))
		err = writer.WriteMessage(&mcap.Message{
			ChannelID:   channelid,
			LogTime:     uint64(time.Now().Unix()),
			PublishTime: uint64(time.Now().Unix()),
			Data:        bytearray,
		})
	}
}
