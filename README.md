# Go ElevenLabs

A production-quality Go client library for the [ElevenLabs](https://elevenlabs.io) Text-to-Speech API.

## Features

- Idiomatic Go API with full context support
- Comprehensive coverage of ElevenLabs API endpoints
- Type-safe request/response handling
- Streaming audio support with `io.ReadCloser`
- Automatic authentication handling
- Proper error handling with custom error types
- Well-tested with unit tests
- Clean package structure for easy navigation

## Installation

```bash
go get github.com/dhia-gharsallaoui/go-elevenlabs
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/dhia-gharsallaoui/go-elevenlabs"
	"github.com/dhia-gharsallaoui/go-elevenlabs/tts"
)

func main() {
	// Create client with your API key
	client := elevenlabs.NewClient(os.Getenv("ELEVENLABS_API_KEY"))

	// Create TTS service
	ttsService := tts.NewService(client)

	// Prepare request
	req := tts.ConvertRequest{
		Text:    "Hello! This is a test of the ElevenLabs API.",
		ModelID: "eleven_flash_v2_5",
		VoiceSettings: &tts.VoiceSettings{
			Stability:       0.5,
			SimilarityBoost: 0.75,
		},
	}

	// Convert text to speech
	audio, err := ttsService.Convert(
		context.Background(),
		"21m00Tcm4TlvDq8ikWAM", // Rachel voice ID
		req,
		&tts.ConvertOptions{
			OutputFormat: "mp3_44100_128",
		},
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer audio.Close()

	// Save audio to file
	file, err := os.Create("output.mp3")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, audio)
	if err != nil {
		fmt.Printf("Error saving audio: %v\n", err)
		return
	}

	fmt.Println("Audio saved to output.mp3")
}
```

## Usage Examples

### Authentication

```go
import "github.com/dhia-gharsallaoui/go-elevenlabs"

// Basic client with API key
client := elevenlabs.NewClient("your-api-key")

// With custom configuration
client := elevenlabs.NewClient(
	"your-api-key",
	elevenlabs.WithBaseURL("https://api.elevenlabs.io"),
	elevenlabs.WithHTTPClient(&http.Client{
		Timeout: 60 * time.Second,
	}),
)
```

### Text-to-Speech

#### Basic Conversion

```go
import (
	"github.com/dhia-gharsallaoui/go-elevenlabs/tts"
)

ttsService := tts.NewService(client)

req := tts.ConvertRequest{
	Text:    "Your text here",
	ModelID: "eleven_flash_v2_5",
}

audio, err := ttsService.Convert(ctx, voiceID, req, nil)
if err != nil {
	// Handle error
}
defer audio.Close()

// Use the audio stream (io.ReadCloser)
```

#### Streaming Conversion

```go
// For real-time streaming
audio, err := ttsService.ConvertStream(ctx, voiceID, req, nil)
if err != nil {
	// Handle error
}
defer audio.Close()

// Stream audio data as it's generated
```

#### With Character Timestamps

```go
result, err := ttsService.ConvertWithTimestamps(ctx, voiceID, req, nil)
if err != nil {
	// Handle error
}

// Access timing information
fmt.Printf("Audio (base64): %s\n", result.AudioBase64)
fmt.Printf("Characters: %v\n", result.Alignment.Characters)
fmt.Printf("Timings: %v\n", result.Alignment.CharacterStartTimesSecs)
```

### Voice Management

#### List Available Voices

```go
import "github.com/dhia-gharsallaoui/go-elevenlabs/voices"

voicesService := voices.NewService(client)

result, err := voicesService.List(ctx)
if err != nil {
	// Handle error
}

for _, voice := range result.Voices {
	fmt.Printf("Voice: %s (ID: %s)\n", voice.Name, voice.VoiceID)
}
```

#### Get Voice Details

```go
voice, err := voicesService.Get(ctx, "voice-id")
if err != nil {
	// Handle error
}

fmt.Printf("Name: %s\n", voice.Name)
fmt.Printf("Category: %s\n", voice.Category)
fmt.Printf("Preview URL: %s\n", voice.PreviewURL)
```

#### Manage Voice Settings

```go
// Get voice settings
settings, err := voicesService.GetSettings(ctx, voiceID)
if err != nil {
	// Handle error
}

// Update voice settings
newSettings := voices.VoiceSettings{
	Stability:       0.7,
	SimilarityBoost: 0.85,
	Style:           0.5,
	UseSpeakerBoost: true,
}

err = voicesService.UpdateSettings(ctx, voiceID, newSettings)
if err != nil {
	// Handle error
}
```

#### Delete a Voice

```go
err := voicesService.Delete(ctx, voiceID)
if err != nil {
	// Handle error
}
```

### History Management

#### List Generation History

```go
import "github.com/dhia-gharsallaoui/go-elevenlabs/history"

historyService := history.NewService(client)

// List with options
result, err := historyService.List(ctx, &history.ListOptions{
	PageSize: 50,
	VoiceID:  "specific-voice-id", // Optional filter
})
if err != nil {
	// Handle error
}

for _, item := range result.History {
	fmt.Printf("Generated: %s at %d\n", item.Text, item.DateUnix)
}
```

#### Get History Item

```go
item, err := historyService.Get(ctx, historyItemID)
if err != nil {
	// Handle error
}

fmt.Printf("Text: %s\n", item.Text)
fmt.Printf("Voice: %s\n", item.VoiceName)
```

#### Download History Audio

```go
// Get audio for a single item
audio, err := historyService.GetAudio(ctx, historyItemID)
if err != nil {
	// Handle error
}
defer audio.Close()

// Save to file
file, _ := os.Create("history-audio.mp3")
defer file.Close()
io.Copy(file, audio)
```

#### Bulk Download

```go
// Download multiple items (returns ZIP if multiple IDs)
archive, err := historyService.Download(ctx, []string{
	"history-id-1",
	"history-id-2",
	"history-id-3",
})
if err != nil {
	// Handle error
}
defer archive.Close()

// Save ZIP file
file, _ := os.Create("history.zip")
defer file.Close()
io.Copy(file, archive)
```

#### Delete History Item

```go
err := historyService.Delete(ctx, historyItemID)
if err != nil {
	// Handle error
}
```

### Models

#### List Available Models

```go
import "github.com/dhia-gharsallaoui/go-elevenlabs/models"

modelsService := models.NewService(client)

result, err := modelsService.List(ctx)
if err != nil {
	// Handle error
}

for _, model := range result.Models {
	fmt.Printf("Model: %s (ID: %s)\n", model.Name, model.ModelID)
	fmt.Printf("  Can do TTS: %v\n", model.CanDoTextToSpeech)
	fmt.Printf("  Can use style: %v\n", model.CanUseStyle)
}
```

### User Information

#### Get User Info

```go
import "github.com/dhia-gharsallaoui/go-elevenlabs/user"

userService := user.NewService(client)

info, err := userService.GetInfo(ctx)
if err != nil {
	// Handle error
}

fmt.Printf("API Key: %s\n", info.XIAPIKey)
fmt.Printf("Tier: %s\n", info.Subscription.Tier)
```

#### Get Subscription Details

```go
sub, err := userService.GetSubscription(ctx)
if err != nil {
	// Handle error
}

fmt.Printf("Character Usage: %d / %d\n",
	sub.CharacterCount,
	sub.CharacterLimit,
)
fmt.Printf("Voice Limit: %d\n", sub.VoiceLimit)
fmt.Printf("Status: %s\n", sub.Status)
```

## Error Handling

The library provides structured error handling:

```go
audio, err := ttsService.Convert(ctx, voiceID, req, nil)
if err != nil {
	// Check if it's an API error
	if apiErr, ok := err.(*elevenlabs.APIError); ok {
		fmt.Printf("API Error (Status %d): %s\n",
			apiErr.StatusCode,
			apiErr.Message,
		)
		if apiErr.Detail != "" {
			fmt.Printf("Detail: %s\n", apiErr.Detail)
		}
	} else {
		// Other error (network, etc.)
		fmt.Printf("Error: %v\n", err)
	}
	return
}
```

## Package Structure

```
go-elevenlabs/
├── client.go          # Core client and HTTP handling
├── error.go           # Error types
├── tts/              # Text-to-speech operations
├── voices/           # Voice management
├── history/          # Generation history
├── models/           # Available models
└── user/             # User account info
```

## Best Practices

### Always Close Streams

```go
audio, err := ttsService.Convert(ctx, voiceID, req, nil)
if err != nil {
	return err
}
defer audio.Close() // Important!

// Use the audio stream
```

### Use Context for Timeouts

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

audio, err := ttsService.Convert(ctx, voiceID, req, nil)
```

### Reuse the Client

```go
// Create once, reuse across requests
client := elevenlabs.NewClient(apiKey)

// Use with different services
ttsService := tts.NewService(client)
voicesService := voices.NewService(client)
```

## API Coverage

### Implemented Packages

| Package | Endpoints | Description |
|---------|-----------|-------------|
| `tts` | 4 | Text-to-speech conversion (standard, streaming, with timestamps) |
| `voices` | 12 | Voice management (list, get, update, delete, settings) |
| `history` | 5 | Generation history (list, get, delete, download) |
| `models` | 1 | List available TTS models |
| `user` | 2 | User account and subscription info |

Additional packages can be easily added following the same patterns.

## Requirements

- Go 1.25 or higher
- ElevenLabs API key ([Get one here](https://elevenlabs.io))

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This library is provided as-is for use with the ElevenLabs API. Please refer to the [ElevenLabs Terms of Service](https://elevenlabs.io/terms) for API usage terms.

## Links

- [ElevenLabs Website](https://elevenlabs.io)
- [ElevenLabs API Documentation](https://elevenlabs.io/docs/api-reference)
- [Go Package Documentation](https://pkg.go.dev/github.com/dhia-gharsallaoui/go-elevenlabs)
