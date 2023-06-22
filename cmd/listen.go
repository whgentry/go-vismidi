/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"os"

	"github.com/nsf/termbox-go"
	"github.com/spf13/cobra"
	"github.com/whgentry/go-vismidi/animation"
	"github.com/whgentry/go-vismidi/config"
	"github.com/whgentry/go-vismidi/control"
	"github.com/whgentry/go-vismidi/led"
	"github.com/whgentry/go-vismidi/midi"
)

// rootCmd represents the base command when called without any subcommands
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listens for midi events from the plugged in device and displays them in the terminal",
	Long: `This command will turn your current terminal window and turn it 
into a midi visualizer for an 88 key keyboard. You should just be able to 
plug in a keyboard with midi output `,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := termbox.Init()
		if err != nil {
			os.Exit(1)
		}
		defer termbox.Close()

		_, NumLEDPerCol := termbox.Size()

		// TODO: Load Settings
		config.LoadConfig(config.DefaultConfigPath)

		// Input and output structures
		animation.Initialize(NumLEDPerCol, midi.GetSettings().KeyCount)
		led.Initialize(NumLEDPerCol, midi.GetSettings().KeyCount)

		// Create Control Channels
		midiEventChan := make(chan midi.MIDIEvent, 100)
		animationFrameChan := make(chan animation.PixelStateFrame, 100)

		midiCB := control.NewIOBlock(
			nil,
			midiEventChan,
			midi.Inputs,
		)

		animationCB := control.NewIOBlock(
			midiEventChan,
			animationFrameChan,
			animation.Animations,
		)

		ledCB := control.NewIOBlock(
			animationFrameChan,
			nil,
			led.Displays,
		)

		// Start control blocks
		midiCB.Start(ctx)
		animationCB.Start(ctx)
		ledCB.Start(ctx)

		// TODO Add methods to controlblock to allow rotation through processor
		termbox.SetInputMode(termbox.InputEsc)
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyArrowRight, termbox.KeySpace:
					animationCB.Next()
				case termbox.KeyArrowLeft:
					animationCB.Previous()
				case termbox.KeyEsc, termbox.KeyCtrlC:
					os.Exit(0)
				}
			}
		}

	},
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-vismidi.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
