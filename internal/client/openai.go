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

	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

const openaiURL = "https://api.openai.com/v1/chat/completions"

const systemPrompt = `Tu es un oracle cosmique complètement perché, mi-astrologue mi-voyant extralucide.

MISSION CRITIQUE: L'utilisateur te pose une QUESTION SPECIFIQUE. Tu DOIS y répondre en combinant:
- Les TRANSITS DU JOUR (où sont les planètes MAINTENANT) → timing, énergie du moment
- Le THÈME NATAL (positions à la naissance) → personnalité, tendances profondes

COMMENT REPONDRE:
1. Lis la question posée
2. Regarde les TRANSITS (positions d'aujourd'hui) pour le timing et l'énergie actuelle
3. Compare avec le THÈME NATAL pour voir comment ça résonne avec la personne
4. Cite des positions SPECIFIQUES des deux pour justifier ta réponse

Exemple:
- Question: "Dois-je changer de travail?"
- Transit: Mars en Sagittaire aujourd'hui
- Natal: Soleil en Taureau
- Réponse: "Mars galope en Sagittaire aujourd'hui, il t'insuffle une soif d'aventure! Mais ton Soleil natal en Taureau te rappelle: ne quitte pas le navire sans avoir rempli tes poches de provisions..."

TON STYLE:
- Oracle déjanté qui canalise des entités astrales farfelues
- Métaphores cosmiques absurdes mais conseils étrangement pertinents
- CITE au moins 1-2 transits ET 1-2 positions natales
- Prédictions décalées et avertissements mystérieux rigolos

REGLES:
- Réponds DIRECTEMENT à la question (pas de généralités!)
- UTILISE les transits ET le natal, pas juste l'un ou l'autre
- Les planètes RÉTROGRADES sont importantes: elles indiquent des énergies tournées vers l'intérieur, des révisions, des retards ou des leçons du passé
- Français, 300-400 mots max
- Répond en markdown`

type OpenAIClient struct {
	apiKey     string
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

func NewOpenAIClient() (*OpenAIClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	return &OpenAIClient{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}, nil
}

func (c *OpenAIClient) GetInterpretation(ctx context.Context, chart *horoscope.Chart, userContext string) (string, error) {
	userPrompt := buildUserPrompt(chart, userContext)

	reqBody := chatRequest{
		Model: "gpt-4o",
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
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
	defer resp.Body.Close()

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
		sb.WriteString(fmt.Sprintf("🔮 QUESTION: %s\n\n", userQuestion))
	} else {
		sb.WriteString("🔮 QUESTION: Donne-moi une lecture cosmique générale pour aujourd'hui basée sur mon thème.\n\n")
	}

	// Today's date and transits
	now := time.Now()
	sb.WriteString(fmt.Sprintf("Aujourd'hui: %s (%s)\n\n",
		now.Format("02/01/2006"),
		frenchWeekday(now.Weekday())))

	// Current planetary positions (transits)
	sb.WriteString("TRANSITS DU JOUR (positions actuelles):\n")
	todayPositions := position.CalculateAll(now)
	for _, pos := range todayPositions {
		if !pos.Body.IsMainPlanet() {
			continue
		}
		zodiac := horoscope.LongitudeToZodiac(pos.EclipticLongitude)
		sb.WriteString(fmt.Sprintf("- %s en %s à %d°%d'%s\n",
			pos.Body.String(), zodiac.Sign.String(), zodiac.Degrees, zodiac.Minutes, retrogradeLabel(pos.Retrograde)))
	}

	// Birth chart data
	sb.WriteString("\nTHÈME NATAL (positions à la naissance):\n")
	sb.WriteString(fmt.Sprintf("Date de naissance: %s\n", chart.DateTime.Format("02/01/2006 15:04")))
	sb.WriteString(fmt.Sprintf("Lieu: %s (%.4f, %.4f)\n\n", chart.Location, chart.Latitude, chart.Longitude))

	sb.WriteString("Positions planétaires:\n")
	for _, pos := range chart.Positions {
		zodiac := horoscope.LongitudeToZodiac(pos.EclipticLongitude)
		sb.WriteString(fmt.Sprintf("- %s en %s à %d°%d'%s\n",
			pos.Body.String(), zodiac.Sign.String(), zodiac.Degrees, zodiac.Minutes, retrogradeLabel(pos.Retrograde)))
	}

	sb.WriteString("\nAspects majeurs:\n")
	for _, aspect := range chart.Aspects {
		sb.WriteString(fmt.Sprintf("- %s %s %s (orbe %.1f°)\n",
			aspect.Body1.String(), aspect.Type.String(), aspect.Body2.String(), aspect.Orb))
	}

	// Element distribution
	elements := calculateElements(chart.Positions)
	sb.WriteString("\nRépartition des éléments:\n")
	sb.WriteString(fmt.Sprintf("- Feu: %d planètes\n", elements[horoscope.Fire]))
	sb.WriteString(fmt.Sprintf("- Terre: %d planètes\n", elements[horoscope.Earth]))
	sb.WriteString(fmt.Sprintf("- Air: %d planètes\n", elements[horoscope.Air]))
	sb.WriteString(fmt.Sprintf("- Eau: %d planètes\n", elements[horoscope.Water]))

	return sb.String()
}

func frenchWeekday(w time.Weekday) string {
	days := []string{"Dimanche", "Lundi", "Mardi", "Mercredi", "Jeudi", "Vendredi", "Samedi"}
	return days[w]
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
		return " (RÉTROGRADE)"
	}
	return ""
}
