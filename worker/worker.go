package worker

import (
	// "context"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"git.progwebtech.com/code/worksmart/store"
	"github.com/centrifugal/gocent"
	"github.com/tidwall/gjson"
	"github.com/osiloke/fcm/queues/machinery"
)

type Config struct {
	Realtime struct {
		Domain string `json:"domain"`
	} `json:"realtime"`
}

type Params struct {
	BodyTemplate    string `json:"bodyTemplate"`
	ChannelTemplate string `json:"channelTemplate"`
	TitleTemplate   string `json:"titleTemplate"`
}

type FCMData struct {
	Data       map[string]interface{}
	Method     string
	GroupName  string
	Owner      string
	StoreTitle string
	StoreID    string `json:"StoreId`
}

func doFCM(c *gocent.Client) func(addonConfig, params, data, traceID string) error {
	return func(addonConfig, params, data, traceID string) error {
		config := Config{}
		ctx := context.Background()
		err := json.Unmarshal([]byte(addonConfig), &config)
		if err != nil {
			return err
		}
		fcmParams := Params{}
		err = json.Unmarshal([]byte(params), &fcmParams)
		if err != nil {
			return err
		}
		// fcmData := FCMData{}
		// err = json.Unmarshal([]byte(data), &fcmData)
		// if err != nil {
		// 	return err
		// }
		// name, err := StringFromFCMData(fcmParams.ChannelTemplate, &fcmData)
		// if err != nil {
		// 	return err
		// }
		storeID := gjson.Get(data, "StoreId")
		if !storeID.Exists() {
			return errors.New("StoreId is missing")
		}
		channelName := gjson.Get(data, fcmParams.ChannelTemplate)
		if !channelName.Exists() {
			return errors.New("channel name cannot be retrieved from template, " + fcmParams.ChannelTemplate)
		}
		channel := store.GenerateStoreCentrifugoChannelFromStore(channelName.String(), storeID.String())
		fmt.Println("gen channel", channel)
		// ctx := context.Background()
		storeEntry := gjson.Get(data, "Data")
		fmt.Println("sending", storeEntry.String(), "to", channel)
		channels, _ := c.Channels(ctx)
		fmt.Printf("Channels: %v\n", channels)
		err = c.Publish(ctx, channel, []byte(storeEntry.String()))
		if err != nil {
			log.Println("failed publishing")
		}
		return err

	}
}

// Worker a fcm worker that sends messages to centrifuge
type Worker struct {
	Addr    string        `help:"centrifuge web address"`
	Key     string        `help:"centrifuge key"`
	Timeout time.Duration `help:"gocent timeout"`
	ID      string        `help:"worker id"`
	Build   string        `help:"build"`
}

// Run run the worker
func (w *Worker) Run() error {
	c := gocent.New(gocent.Config{
		Addr: w.Addr,
		Key:  w.Key,
	})
	return machinery.Worker(w.ID, map[string]interface{}{
		"notification": doFCM(c),
	})
}

// NewWorker new worker
func NewWorker(build string) *Worker {
	centrifugoSecret := os.Getenv("CENTRIFUGO_SECRET")
	if centrifugoSecret == "" {
		centrifugoSecret = "rsecret"
	}
	centrifugoURL := os.Getenv("CENTRIFUGO_URL")
	if centrifugoURL == "" {
		centrifugoURL = "rsecret"
	}
	return &Worker{Addr: centrifugoURL, Key: centrifugoSecret, Timeout: 5 * time.Second, Build: build}
}
