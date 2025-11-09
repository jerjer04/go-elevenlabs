# Go ElevenLabs

[![Go Reference](https://pkg.go.dev/badge/github.com/dhia-gharsallaoui/go-elevenlabs.svg)](https://pkg.go.dev/github.com/dhia-gharsallaoui/go-elevenlabs)
[![Go Report Card](https://goreportcard.com/badge/github.com/dhia-gharsallaoui/go-elevenlabs)](https://goreportcard.com/report/github.com/dhia-gharsallaoui/go-elevenlabs)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dhia-gharsallaoui/go-elevenlabs)](https://github.com/dhia-gharsallaoui/go-elevenlabs)

A production-grade Go client library for the [ElevenLabs](https://elevenlabs.io) Text-to-Speech API. Built with idiomatic Go practices, comprehensive error handling, and full support for streaming audio generation.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
  - [Authentication](#authentication)
  - [Text-to-Speech](#text-to-speech)
  - [Voice Management](#voice-management)
  - [History Management](#history-management)
  - [Models](#models)
  - [User Information](#user-information)
- [Error Handling](#error-handling)
- [API Coverage](#api-coverage)
- [Examples](#examples)
- [Best Practices](#best-practices)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Idiomatic Go Design** - Clean, type-safe API following Go best practices
- **Full API Coverage** - Support for TTS, voices, history, models, and user endpoints
- **Context Support** - Built-in context propagation for cancellation and timeouts
- **Streaming Audio** - Efficient audio streaming with `io.ReadCloser` interface
- **Type Safety** - Strongly-typed requests and responses
- **Zero Dependencies** - Uses only Go standard library
- **Production Ready** - Comprehensive error handling and unit tests
- **Well Documented** - Extensive documentation and examples

## Installation

```bash
go get github.com/dhia-gharsallaoui/go-elevenlabs
```

**Requirements:**
- Go 1.25 or higher
- ElevenLabs API key ([Sign up here](https://elevenlabs.io))

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
	// Initialize client with your API key
	client := elevenlabs.NewClient(os.Getenv("ELEVENLABS_API_KEY"))

	// Create TTS service
	ttsService := tts.NewService(client)

	// Configure request
	req := tts.ConvertRequest{
		Text:    "Hello! Welcome to ElevenLabs Text-to-Speech API.",
		ModelID: "eleven_flash_v2_5",
		VoiceSettings: &tts.VoiceSettings{
			Stability:       0.5,
			SimilarityBoost: 0.75,
		},
	}

	// Generate speech
	audio, err := ttsService.Convert(
		context.Background(),
		"21m00Tcm4TlvDq8ikWAM", // Voice ID (Rachel)
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

	// Save to file
	file, err := os.Create("output.mp3")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, audio); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("✓ Audio saved to output.mp3")
}
```

## Usage

### Authentication

Create a client with your ElevenLabs API key:

```go
import "github.com/dhia-gharsallaoui/go-elevenlabs"

// Basic initialization
client := elevenlabs.NewClient(os.Getenv("ELEVENLABS_API_KEY"))

// Advanced configuration
client := elevenlabs.NewClient(
	apiKey,
	elevenlabs.WithBaseURL("https://api.elevenlabs.io"),
	elevenlabs.WithHTTPClient(&http.Client{
		Timeout: 60 * time.Second,
	}),
)
```

### Text-to-Speech

#### Standard Conversion

Convert text to speech with customizable voice settings:

```go
import "github.com/dhia-gharsallaoui/go-elevenlabs/tts"

ttsService := tts.NewService(client)

req := tts.ConvertRequest{
	Text:    "Your text here",
	ModelID: "eleven_flash_v2_5",
	VoiceSettings: &tts.VoiceSettings{
		Stability:       0.5,
		SimilarityBoost: 0.75,
		Style:           0.0,
		UseSpeakerBoost: true,
	},
}

audio, err := ttsService.Convert(ctx, voiceID, req, &tts.ConvertOptions{
	OutputFormat: "mp3_44100_128",
})
if err != nil {
	log.Fatal(err)
}
defer audio.Close()
```

#### Streaming Conversion

For real-time audio streaming:

```go
audioStream, err := ttsService.ConvertStream(ctx, voiceID, req, opts)
if err != nil {
	log.Fatal(err)
}
defer audioStream.Close()

// Process audio chunks as they arrive
io.Copy(outputWriter, audioStream)
```

#### With Character Timestamps

Get character-level timing information:

```go
result, err := ttsService.ConvertWithTimestamps(ctx, voiceID, req, nil)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Audio (base64): %s\n", result.AudioBase64)
fmt.Printf("Characters: %v\n", result.Alignment.Characters)
fmt.Printf("Start times: %v\n", result.Alignment.CharacterStartTimesSecs)
```

### Voice Management

#### List Voices

Retrieve all available voices:

```go
import "github.com/dhia-gharsallaoui/go-elevenlabs/voices"

voicesService := voices.NewService(client)

result, err := voicesService.List(ctx)
if err != nil {
	log.Fatal(err)
}

for _, voice := range result.Voices {
	fmt.Printf("%-30s ID: %s Category: %s\n",
		voice.Name, voice.VoiceID, voice.Category)
}
```

#### Get Voice Details

```go
voice, err := voicesService.Get(ctx, voiceID)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Name: %s\n", voice.Name)
fmt.Printf("Category: %s\n", voice.Category)
fmt.Printf("Preview: %s\n", voice.PreviewURL)
```

#### Update Voice Settings

```go
settings := voices.VoiceSettings{
	Stability:       0.7,
	SimilarityBoost: 0.85,
	Style:           0.5,
	UseSpeakerBoost: true,
}

if err := voicesService.UpdateSettings(ctx, voiceID, settings); err != nil {
	log.Fatal(err)
}
```

#### Delete Voice

```go
if err := voicesService.Delete(ctx, voiceID); err != nil {
	log.Fatal(err)
}
```

### History Management

#### List Generation History

```go
import "github.com/dhia-gharsallaoui/go-elevenlabs/history"

historyService := history.NewService(client)

result, err := historyService.List(ctx, &history.ListOptions{
	PageSize: 50,
	VoiceID:  "optional-voice-id-filter",
})
if err != nil {
	log.Fatal(err)
}

for _, item := range result.History {
	fmt.Printf("[%s] %s - %s\n",
		time.Unix(item.DateUnix, 0).Format(time.RFC3339),
		item.VoiceName,
		item.Text[:50])
}
```

#### Download History Audio

```go
// Single item
audio, err := historyService.GetAudio(ctx, historyItemID)
if err != nil {
	log.Fatal(err)
}
defer audio.Close()

// Bulk download (returns ZIP for multiple items)
archive, err := historyService.Download(ctx, []string{
	"history-id-1",
	"history-id-2",
})
if err != nil {
	log.Fatal(err)
}
defer archive.Close()
```

#### Delete History Item

```go
if err := historyService.Delete(ctx, historyItemID); err != nil {
	log.Fatal(err)
}
```

### Models

#### List Available Models

```go
import "github.com/dhia-gharsallaoui/go-elevenlabs/models"

modelsService := models.NewService(client)

result, err := modelsService.List(ctx)
if err != nil {
	log.Fatal(err)
}

for _, model := range result.Models {
	fmt.Printf("%-30s ID: %s\n", model.Name, model.ModelID)
	fmt.Printf("  TTS: %v | Voice Conversion: %v | Style: %v\n",
		model.CanDoTextToSpeech,
		model.CanDoVoiceConversion,
		model.CanUseStyle)
}
```

### User Information

#### Get Account Information

```go
import "github.com/dhia-gharsallaoui/go-elevenlabs/user"

userService := user.NewService(client)

info, err := userService.GetInfo(ctx)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Tier: %s\n", info.Subscription.Tier)
fmt.Printf("Characters: %d / %d\n",
	info.Subscription.CharacterCount,
	info.Subscription.CharacterLimit)
```

#### Get Subscription Details

```go
sub, err := userService.GetSubscription(ctx)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Status: %s\n", sub.Status)
fmt.Printf("Voice Limit: %d\n", sub.VoiceLimit)
fmt.Printf("Next Reset: %s\n",
	time.Unix(sub.NextCharacterCountResetUnix, 0).Format(time.RFC3339))
```

## Error Handling

The library provides structured error handling with detailed error information:

```go
audio, err := ttsService.Convert(ctx, voiceID, req, nil)
if err != nil {
	// Check for API-specific errors
	if apiErr, ok := err.(*elevenlabs.APIError); ok {
		fmt.Printf("API Error (HTTP %d): %s\n", apiErr.StatusCode, apiErr.Message)
		if apiErr.Detail != "" {
			fmt.Printf("Details: %s\n", apiErr.Detail)
		}
		return
	}

	// Handle other errors (network, context cancellation, etc.)
	fmt.Printf("Error: %v\n", err)
	return
}
```

## API Coverage

The library provides comprehensive coverage of the ElevenLabs API:

| Package | Endpoints | Description |
|---------|-----------|-------------|
| `tts` | 4 | Text-to-speech conversion (standard, streaming, with timestamps) |
| `voices` | 12 | Voice management (CRUD operations, settings, sharing) |
| `history` | 5 | Generation history (list, retrieve, delete, download) |
| `models` | 1 | List available TTS models and their capabilities |
| `user` | 2 | User account information and subscription details |

**Total:** 24 API endpoints implemented and tested.

## Examples

The [examples](./examples) directory contains sample programs demonstrating library usage:

- **[simple_tts](./examples/simple_tts)** - Basic text-to-speech conversion
- **[test_library](./examples/test_library)** - Comprehensive API testing suite

Run examples:

```bash
export ELEVENLABS_API_KEY="your-api-key"
export ELEVENLABS_VOICE_ID="your-voice-id"

# Simple TTS
go run ./examples/simple_tts "Your text here"

# Comprehensive test
go run ./examples/test_library
```

## Best Practices

### Always Close Audio Streams

```go
audio, err := ttsService.Convert(ctx, voiceID, req, nil)
if err != nil {
	return err
}
defer audio.Close() // ← Important!
```

### Use Context for Timeout Control

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

audio, err := ttsService.Convert(ctx, voiceID, req, nil)
```

### Reuse Client Instances

```go
// Initialize once
client := elevenlabs.NewClient(apiKey)

// Reuse across services
ttsService := tts.NewService(client)
voicesService := voices.NewService(client)
historyService := history.NewService(client)
```

### Handle Rate Limits

```go
audio, err := ttsService.Convert(ctx, voiceID, req, nil)
if err != nil {
	if apiErr, ok := err.(*elevenlabs.APIError); ok {
		if apiErr.StatusCode == 429 {
			// Implement exponential backoff
			time.Sleep(time.Second * 5)
			// Retry request
		}
	}
}
```

## Package Structure

```
go-elevenlabs/
├── client.go          # Core HTTP client with authentication
├── error.go           # Error types and handling
├── tts/              # Text-to-speech operations
├── voices/           # Voice management
├── history/          # Generation history
├── models/           # Available models
├── user/             # User account operations
└── examples/         # Usage examples
    ├── simple_tts/
    └── test_library/
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development

```bash
# Clone repository
git clone https://github.com/dhia-gharsallaoui/go-elevenlabs.git
cd go-elevenlabs

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Build all packages
go build ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built for the [ElevenLabs](https://elevenlabs.io) Text-to-Speech API
- Inspired by idiomatic Go client library design patterns

## Links

- **ElevenLabs**: [https://elevenlabs.io](https://elevenlabs.io)
- **API Documentation**: [https://elevenlabs.io/docs/api-reference](https://elevenlabs.io/docs/api-reference)
- **Go Package**: [https://pkg.go.dev/github.com/dhia-gharsallaoui/go-elevenlabs](https://pkg.go.dev/github.com/dhia-gharsallaoui/go-elevenlabs)

---

<div align="center">
  <sub>Built with ❤️ using Go</sub>
</div>
