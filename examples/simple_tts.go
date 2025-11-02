package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dhia-gharsallaoui/go-elevenlabs"
	"github.com/dhia-gharsallaoui/go-elevenlabs/tts"
)

func main() {
	// Get credentials from environment variables
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	voiceID := os.Getenv("ELEVENLABS_VOICE_ID")

	if apiKey == "" {
		fmt.Println("Error: ELEVENLABS_API_KEY environment variable not set")
		os.Exit(1)
	}
	if voiceID == "" {
		fmt.Println("Error: ELEVENLABS_VOICE_ID environment variable not set")
		os.Exit(1)
	}

	// Create client
	client := elevenlabs.NewClient(apiKey)
	ttsService := tts.NewService(client)

	// Get text from user or use default
	text := "Hello! This is a simple test of the ElevenLabs Go client library."
	if len(os.Args) > 1 {
		text = strings.Join(os.Args[1:], " ")
	}

	fmt.Printf("Converting: \"%s\"\n", text)

	// Convert text to speech
	req := tts.ConvertRequest{
		Text:    text,
		ModelID: "eleven_flash_v2_5",
		VoiceSettings: &tts.VoiceSettings{
			Stability:       0.5,
			SimilarityBoost: 0.75,
		},
	}

	audio, err := ttsService.Convert(
		context.Background(),
		voiceID,
		req,
		&tts.ConvertOptions{
			OutputFormat: "mp3_44100_128",
		},
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer audio.Close()

	// Save to file
	outputFile := "output.mp3"
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	written, err := io.Copy(file, audio)
	if err != nil {
		fmt.Printf("Error saving audio: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Success! Audio saved to %s (%d bytes)\n", outputFile, written)
}
