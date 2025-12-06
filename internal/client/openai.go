// Package client provides a client for OpenAI's GPT-4 API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ctrl-vfr/astral-tui/internal/i18n"
	"github.com/ctrl-vfr/astral-tui/pkg/horoscope"
	"github.com/ctrl-vfr/astral-tui/pkg/position"
)

const (
	openaiURL    = "https://api.openai.com/v1/chat/completions"
	defaultModel = "gpt-4o-mini"
)

// OpenAIClient handles requests to OpenAI's GPT API.
type OpenAIClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewOpenAIClient creates a new OpenAI client.
// Uses OPENAI_API_KEY for authentication and ASTRAL_OPENAI_MODEL to override the default model (gpt-4o-mini).
func NewOpenAIClient() (*OpenAIClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	model := os.Getenv("ASTRAL_OPENAI_MODEL")
	if model == "" {
		model = defaultModel
	}

	return &OpenAIClient{
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{},
	}, nil
}

// GetInterpretation requests an astrological interpretation from GPT-4.
func (c *OpenAIClient) GetInterpretation(ctx context.Context, chart *horoscope.Chart, userContext string) (string, error) {
	userPrompt := buildUserPrompt(chart, userContext)

	reqBody := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: i18n.SystemPrompt()},
			{Role: "user", Content: userPrompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", openaiURL, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI returned status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func buildUserPrompt(chart *horoscope.Chart, userQuestion string) string {
	var sb strings.Builder

	// User question first - if no question, ask for general reading
	if userQuestion != "" {
		sb.WriteString(fmt.Sprintf("## %s: %s\n\n", i18n.T("PromptQuestion"), userQuestion))
	} else {
		sb.WriteString(fmt.Sprintf("## %s: %s\n\n", i18n.T("PromptQuestion"), i18n.T("PromptDefaultQuestion")))
	}

	// Today's date and transits
	now := time.Now()
	sb.WriteString(fmt.Sprintf("%s: %s (%s)\n\n",
		i18n.T("PromptToday"),
		now.Format("02/01/2006"),
		i18n.Weekday(int(now.Weekday()))))

	// Current planetary positions (transits)
	sb.WriteString(fmt.Sprintf("## %s:\n", i18n.T("PromptTransitsTitle")))
	todayPositions := position.CalculateAll(now)
	for _, pos := range todayPositions {
		if !pos.Body.IsMainPlanet() {
			continue
		}
		zodiac := horoscope.LongitudeToZodiac(pos.EclipticLongitude)
		sb.WriteString(fmt.Sprintf("- %s %s: %s %d°%d'%s\n",
			pos.Body.Symbol(), pos.Body.String(), zodiac.Sign.String(), zodiac.Degrees, zodiac.Minutes, retrogradeLabel(pos.Retrograde)))
	}

	// Birth chart data
	sb.WriteString(fmt.Sprintf("\n## %s:\n", i18n.T("PromptNatalTitle")))
	sb.WriteString(fmt.Sprintf("%s: %s\n", i18n.T("PromptBirthDate"), chart.DateTime.Format("02/01/2006 15:04")))
	sb.WriteString(fmt.Sprintf("%s: %s (%.4f, %.4f)\n\n", i18n.T("PromptLocation"), chart.Location, chart.Latitude, chart.Longitude))

	sb.WriteString(fmt.Sprintf("%s:\n", i18n.T("PromptPlanetPositions")))
	for _, pos := range chart.Positions {
		zodiac := horoscope.LongitudeToZodiac(pos.EclipticLongitude)
		sb.WriteString(fmt.Sprintf("- %s %s: %s %d°%d'%s\n",
			pos.Body.Symbol(), pos.Body.String(), zodiac.Sign.String(), zodiac.Degrees, zodiac.Minutes, retrogradeLabel(pos.Retrograde)))
	}

	sb.WriteString(fmt.Sprintf("\n%s:\n", i18n.T("PromptMajorAspects")))
	for _, aspect := range chart.Aspects {
		sb.WriteString(fmt.Sprintf("- %s %s %s %s %s (%s %.1f°)\n",
			aspect.Body1.Symbol(), aspect.Body1.String(), aspect.Type.String(), aspect.Body2.Symbol(), aspect.Body2.String(), i18n.T("PromptOrb"), aspect.Orb))
	}

	// Element distribution
	elements := calculateElements(chart.Positions)
	sb.WriteString(fmt.Sprintf("\n%s:\n", i18n.T("PromptElementDist")))
	sb.WriteString(fmt.Sprintf("- %s: "+i18n.T("ElementCount")+"\n", i18n.T("ElementFire"), elements[horoscope.Fire]))
	sb.WriteString(fmt.Sprintf("- %s: "+i18n.T("ElementCount")+"\n", i18n.T("ElementEarth"), elements[horoscope.Earth]))
	sb.WriteString(fmt.Sprintf("- %s: "+i18n.T("ElementCount")+"\n", i18n.T("ElementAir"), elements[horoscope.Air]))
	sb.WriteString(fmt.Sprintf("- %s: "+i18n.T("ElementCount")+"\n", i18n.T("ElementWater"), elements[horoscope.Water]))

	return sb.String()
}

func calculateElements(positions []position.Position) map[horoscope.Element]int {
	elements := make(map[horoscope.Element]int)
	for _, pos := range positions {
		if pos.Body > position.Pluto {
			continue
		}
		zodiac := horoscope.LongitudeToZodiac(pos.EclipticLongitude)
		elements[zodiac.Sign.Element()]++
	}
	return elements
}

func retrogradeLabel(isRetro bool) string {
	if isRetro {
		return i18n.T("PromptRetrograde")
	}
	return ""
}
