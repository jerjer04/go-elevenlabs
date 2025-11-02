package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/dhia-gharsallaoui/go-elevenlabs"
	"github.com/dhia-gharsallaoui/go-elevenlabs/history"
	"github.com/dhia-gharsallaoui/go-elevenlabs/models"
	"github.com/dhia-gharsallaoui/go-elevenlabs/tts"
	"github.com/dhia-gharsallaoui/go-elevenlabs/user"
	"github.com/dhia-gharsallaoui/go-elevenlabs/voices"
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

	fmt.Println("🚀 Testing ElevenLabs Go Client Library")
	fmt.Println("=====================================\n")

	// Create client
	client := elevenlabs.NewClient(apiKey)
	ctx := context.Background()

	// Test 1: Get User Info
	fmt.Println("📊 Test 1: Getting user information...")
	userService := user.NewService(client)
	userInfo, err := userService.GetInfo(ctx)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
	} else {
		fmt.Printf("✅ Success!\n")
		fmt.Printf("   Subscription Tier: %s\n", userInfo.Subscription.Tier)
		fmt.Printf("   Character Usage: %d / %d\n",
			userInfo.Subscription.CharacterCount,
			userInfo.Subscription.CharacterLimit)
		fmt.Printf("   Voice Limit: %d\n\n", userInfo.Subscription.VoiceLimit)
	}

	// Test 2: List Available Models
	fmt.Println("🤖 Test 2: Listing available models...")
	modelsService := models.NewService(client)
	modelsList, err := modelsService.List(ctx)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
	} else {
		fmt.Printf("✅ Success! Found %d models\n", len(modelsList.Models))
		if len(modelsList.Models) > 0 {
			fmt.Printf("   First model: %s (ID: %s)\n\n",
				modelsList.Models[0].Name,
				modelsList.Models[0].ModelID)
		}
	}

	// Test 3: List Voices
	fmt.Println("🎤 Test 3: Listing available voices...")
	voicesService := voices.NewService(client)
	voicesList, err := voicesService.List(ctx)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
	} else {
		fmt.Printf("✅ Success! Found %d voices\n", len(voicesList.Voices))
		// Find the voice we're using
		for _, v := range voicesList.Voices {
			if v.VoiceID == voiceID {
				fmt.Printf("   Using voice: %s (Category: %s)\n\n", v.Name, v.Category)
				break
			}
		}
	}

	// Test 4: Text-to-Speech Conversion
	fmt.Println("🔊 Test 4: Converting text to speech...")
	ttsService := tts.NewService(client)

	req := tts.ConvertRequest{
		Text:    "Hello! This is a test of the ElevenLabs Go client library. It seems to be working perfectly!",
		ModelID: "eleven_flash_v2_5",
		VoiceSettings: &tts.VoiceSettings{
			Stability:       0.5,
			SimilarityBoost: 0.75,
		},
	}

	audio, err := ttsService.Convert(
		ctx,
		voiceID,
		req,
		&tts.ConvertOptions{
			OutputFormat: "mp3_44100_128",
		},
	)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
	} else {
		defer audio.Close()

		// Save to file
		outputFile := "test_output.mp3"
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("❌ Failed to create file: %v\n\n", err)
		} else {
			defer file.Close()

			written, err := io.Copy(file, audio)
			if err != nil {
				fmt.Printf("❌ Failed to save audio: %v\n\n", err)
			} else {
				fmt.Printf("✅ Success!\n")
				fmt.Printf("   Audio saved to: %s\n", outputFile)
				fmt.Printf("   File size: %d bytes\n\n", written)
			}
		}
	}

	// Test 5: Get Recent History
	fmt.Println("📜 Test 5: Fetching generation history...")
	historyService := history.NewService(client)
	historyList, err := historyService.List(ctx, &history.ListOptions{
		PageSize: 5,
	})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n\n", err)
	} else {
		fmt.Printf("✅ Success! Found %d recent items\n", len(historyList.History))
		if len(historyList.History) > 0 {
			fmt.Printf("   Most recent: \"%s\" (Voice: %s)\n\n",
				historyList.History[0].Text[:min(50, len(historyList.History[0].Text))]+"...",
				historyList.History[0].VoiceName)
		}
	}

	// Summary
	fmt.Println("=====================================")
	fmt.Println("✨ All tests completed!")
	fmt.Println("=====================================")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
