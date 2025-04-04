package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/hajimehoshi/go-mp3"
	"github.com/ebitengine/oto/v3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run player.go <file.mp3>")
		return
	}

	filePath := os.Args[1]

	// Open the MP3 file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the MP3 file
	decoded, err := mp3.NewDecoder(file)
	if err != nil {
		fmt.Println("Error decoding MP3:", err)
		return
	}

	// Create Oto context
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate: 44100,
		ChannelCount: 2,
		Format: oto.FormatSignedInt16LE,
	})// Sample rate: 44100 Hz, Stereo, 16-bit
	if err != nil {
		fmt.Println("Error creating Oto context:", err)
		return
	}
	<-ready // Wait until the context is ready

	// Create an Oto player
	player := ctx.NewPlayer(decoded)

	// Handle Ctrl+C to stop playback
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Play in a loop
	go func() {
		for {
			file.Seek(0, io.SeekStart) // Reset file pointer
			decoded, err = mp3.NewDecoder(file)
			if err != nil {
				fmt.Println("Error decoding MP3:", err)
				return
			}

			player = ctx.NewPlayer(decoded)
			player.Play()

			// Wait until playback is finished
			for player.IsPlaying() {
				continue
			}
		}
	}()

	// Wait for Ctrl+C
	<-sigChan
	fmt.Println("\nStopping playback...")
	player.Close()
}

