package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"net/http"
)

/* Notification data model */
type Frame struct {
	Icon     string `json:"icon,omitempty"`
	Text     string `json:"text,omitempty"`
	GoalData struct {
		Start   int    `json:"start"`
		Current int    `json:"current"`
		End     int    `json:"end"`
		Unit    string `json:"unit"`
	} `json:"goalData,omitempty"`
	ChartData []int `json:"chartData,omitempty"`
}

type Sound struct {
	Category string `json:"category"`
	ID       string `json:"id"`
	Repeat   int    `json:"repeat"`
}

type Model struct {
	Frames []Frame `json:"frames"`
	Sound Sound `json:"sound"`
	Cycles int `json:"cycles"`
}

type Notification struct {
	Priority string `json:"priority"`
	IconType string `json:"icon_type"`
	LifeTime int    `json:"lifeTime"`
	Model    Model `json:"model"`
}

/* Handler */

type HandlerConfig struct {
	sensu.PluginConfig
	Ip      string
	Key  string
	EntityIcon string
	
	OkIcon string
	OkSound string
	
	WarningIcon string
	WarningSound string
	
	CriticalIcon string
	CriticalSound string
}

var (
	config = HandlerConfig{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-lametric-handler",
			Short:    "The Sensu Go lametric notifications handler",
			Timeout:  10,
			Keyspace: "sensu.io/plugins/lametric/config",
		},
	}

	configOptions = []*sensu.PluginConfigOption{
		{
			Path:      "ip",
			Env:       "SENSU_LAMETRIC_IP",
			Argument:  "ip",
			Shorthand: "i",
			Default:   "",
			Usage:     "The lametric ip to send notifications to, defaults to value of SENSU_LAMETRIC_IP env variable",
			Value:     &config.Ip,
		},
		{
			Path:      "key",
			Env:       "SENSU_LAMETRIC_KEY",
			Argument:  "key",
			Shorthand: "k",
			Default:   "",
			Usage:     "The lametric auth key, defaults to value of SENSU_LAMETRIC_KEY env variable",
			Value:     &config.Key,
		},
		{
			Path:      "entity-icon",
			Env:       "SENSU_LAMETRIC_ENTITY_ICON",
			Argument:  "entity-icon",
			Shorthand: "e",
			Default:   "i31916",
			Usage:     "The entity notification icon",
			Value:     &config.EntityIcon,
		},
		{
			Path:      "ok-icon",
			Env:       "SENSU_LAMETRIC_OK_ICON",
			Argument:  "ok-icon",
			Shorthand: "o",
			Default:   "a25939",
			Usage:     "The ok state notification icon",
			Value:     &config.OkIcon,
		},
		{
			Path:      "ok-sound",
			Env:       "SENSU_LAMETRIC_OK_SOUND",
			Argument:  "ok-sound",
			Shorthand: "O",
			Default:   "positive1",
			Usage:     "The ok state notification sound",
			Value:     &config.OkSound,
		},
		{
			Path:      "warning-icon",
			Env:       "SENSU_LAMETRIC_warning_ICON",
			Argument:  "warning-icon",
			Shorthand: "w",
			Default:   "a7756",
			Usage:     "The warning state notification icon",
			Value:     &config.WarningIcon,
		},
		{
			Path:      "warning-sound",
			Env:       "SENSU_LAMETRIC_warning_SOUND",
			Argument:  "warning-sound",
			Shorthand: "W",
			Default:   "negative5",
			Usage:     "The warning state notification sound",
			Value:     &config.WarningSound,
		},
		{
			Path:      "critical-icon",
			Env:       "SENSU_LAMETRIC_critical_ICON",
			Argument:  "critical-icon",
			Shorthand: "c",
			Default:   "a2715",
			Usage:     "The critical state notification icon",
			Value:     &config.CriticalIcon,
		},
		{
			Path:      "critical-sound",
			Env:       "SENSU_LAMETRIC_critical_SOUND",
			Argument:  "critical-sound",
			Shorthand: "C",
			Default:   "negative1",
			Usage:     "The critical state notification sound",
			Value:     &config.CriticalSound,
		},
	}
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	goHandler := sensu.NewGoHandler(&config.PluginConfig, configOptions, checkArgs, sendMessage)
	goHandler.Execute()
}

func checkArgs(_ *corev2.Event) error {
	if len(config.Ip) == 0 {
		return fmt.Errorf("--ip or SENSU_LAMETRIC_IP environment variable is required")
	}
	if len(config.Key) == 0 {
		return fmt.Errorf("--key or SENSU_LAMETRIC_KEY environment variable is required")
	}

	return nil
}

func messageStatus(event *corev2.Event) (string, string) {
	switch event.Check.Status {
	case 0:
		return config.OkIcon, config.OkSound
	case 2:
		return config.CriticalIcon, config.CriticalSound
	default:
		return config.WarningIcon, config.WarningSound
	}
}

func sendMessage(event *corev2.Event) error {
	icon, sound := messageStatus(event)

	n := &Notification{
		Priority: "info",
		IconType: "info",
		Model: Model{
			Frames: []Frame{
				{
					Text: event.Entity.Name,
					Icon: config.EntityIcon,
				},
				{
					Text: event.Check.Name,
					Icon: icon,
				},
			},
			Sound:  Sound{
				Category: "notifications",
				ID: sound,
				Repeat: 1,
			},
			Cycles: 1,
		},
	}

	dat, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://%s/api/v2/device/notifications", config.Ip),
		bytes.NewReader(dat))
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth("dev", config.Key)

	hc := http.Client{}
	_, err = hc.Do(req)
	if err != nil {
		panic(err)
	}

	return err
}
