package candashboard

import (
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
	mcapWriter, mcapErr := mcap.NewWriter(file, nil)

	if mcapErr != nil {
		fmt.Println("Error creating mcap file:", mcapErr)
		return nil
	}
	return mcapWriter
}

func NewSchema(id uint16, name string) *mcap.Schema {
	return &mcap.Schema{
		ID:   id,
		Name: name,
	}
}

func NewChannel(channelID uint16, schemaID mcap.Schema, name string) *mcap.Channel {
	return &mcap.Channel{
		ID:       channelID,
		SchemaID: schemaID.ID,
		Topic:    name,
	}
}

func NewMessage(channelID mcap.Channel, logtime uint64, pubtime uint64, data []byte) *mcap.Message {
	return &mcap.Message{
		ChannelID:   channelID.ID,
		LogTime:     logtime,
		PublishTime: pubtime,
		Data:        data,
	}
}
