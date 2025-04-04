package player

import (
	"fmt"
	"os"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

// PlaySound plays the given MP3 file using oto.
func PlaySound(filePath string) error {
	// Open the MP3 file
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer f.Close()

	// Decode the MP3 file using go-mp3
	decoder, err := mp3.NewDecoder(f)
	if err != nil {
		return fmt.Errorf("could not decode mp3 file: %v", err)
	}

	// Create a new audio context with the correct options
	options := oto.NewContextOptions{
		SampleRate: decoder.SampleRate(),
		NumChannels: 2,     // Stereo
		BitDepth: 2,        // 16-bit audio
		BufferSize: 1024,   // Buffer size
	}

	ctx, err := oto.NewContextWithOptions(options)
	if err != nil {
		return fmt.Errorf("could not create audio context: %v", err)
	}
	defer ctx.Close()

	// Create a new player for the audio stream
	player := ctx.NewPlayer(decoder)
	defer player.Close()

	// Play the sound
	player.Play()

	// Wait until the player finishes playing
	for player.IsPlaying() {
		// You can add a small sleep here to avoid hogging CPU resources
	}

	// Cleanup
	return nil
}

