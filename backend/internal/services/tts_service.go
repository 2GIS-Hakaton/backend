package services

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// TTSService генерирует аудио через Yandex SpeechKit
type TTSService struct {
	apiKey   string
	folderID string
	baseURL  string
	iamToken string
	voice    string
	client   *http.Client
}

func NewTTSService() *TTSService {
	return &TTSService{
		apiKey:   os.Getenv("YANDEX_API_KEY"),
		folderID: os.Getenv("YANDEX_FOLDER_ID"),
		baseURL:  "https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize",
		voice:    getVoice(),
		client:   &http.Client{},
	}
}

// GenerateAudio генерирует аудио из текста
func (s *TTSService) GenerateAudio(text, filename string) (string, int, error) {
	if s.apiKey == "" {
		return "", 0, fmt.Errorf("YANDEX_API_KEY not set")
	}

	// Формируем URL с правильно закодированными параметрами
	params := url.Values{}
	params.Add("text", text)
	params.Add("lang", "ru-RU")
	params.Add("voice", s.voice)
	params.Add("speed", "1.0")
	params.Add("format", "mp3")
	params.Add("emotion", "good")
	params.Add("folderId", s.folderID)

	fullURL := s.baseURL + "?" + params.Encode()

	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Аутентификация через API ключ
	req.Header.Set("Authorization", "Api-Key "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to call Yandex SpeechKit API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("Yandex SpeechKit API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Сохраняем аудио
	audioDir := "./audio"
	if err := os.MkdirAll(audioDir, 0755); err != nil {
		return "", 0, fmt.Errorf("failed to create audio directory: %w", err)
	}

	audioPath := filepath.Join(audioDir, filename+".mp3")
	file, err := os.Create(audioPath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create audio file: %w", err)
	}
	defer file.Close()

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read audio data: %w", err)
	}

	if _, err := file.Write(audioData); err != nil {
		return "", 0, fmt.Errorf("failed to write audio file: %w", err)
	}

	// Примерная длительность: 150 слов в минуту
	words := len(text) / 5 // примерное количество слов
	durationSeconds := (words * 60) / 150

	return audioPath, durationSeconds, nil
}

// getVoice возвращает голос из переменной окружения или дефолтный
func getVoice() string {
	voice := os.Getenv("YANDEX_VOICE")
	if voice == "" {
		// Дефолтный голос - женский, приятный
		return "alena"
	}
	return voice
}

// Доступные голоса Yandex SpeechKit:
// Женские: alena (мягкий), oksana (нейтральный), jane (эмоциональный)
// Мужские: filipp (нейтральный), ermil (эмоциональный), omazh (глубокий)
// Premium: marina (премиум женский), alexander (премиум мужской)
