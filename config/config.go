package config

import (
	"context"
	"os"

	"github.com/whgentry/go-vismidi/animation"
	"github.com/whgentry/go-vismidi/led"
	"github.com/whgentry/go-vismidi/midi"
	"gopkg.in/yaml.v3"
)

var DefaultConfigPath = "./config/visadadfmidi.yaml"

type Settings struct {
	animation animation.Settings
	led       led.Settings
	midi      midi.Settings
}

func LoadConfig(filepath string) error {
	data, err := os.ReadFile(filepath)

	settings := Settings{}
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, settings)
	if err != nil {
		return err
	}

	animation.ApplySettings(settings.animation)
	led.ApplySettings(settings.led)
	midi.ApplySettings(settings.midi)

	return nil
}

func StartConfigWatch(ctx context.Context) {

}
