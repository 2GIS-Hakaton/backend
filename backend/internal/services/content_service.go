package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dimmy-kor/audioguid/internal/models"
)

// ContentService генерирует контент через YandexGPT
type ContentService struct {
	apiKey     string
	folderID   string
	baseURL    string
	searchAPI  string
	searchUser string
	searchKey  string
	client     *http.Client
}

func NewContentService() *ContentService {
	return &ContentService{
		apiKey:     os.Getenv("YANDEX_API_KEY"),
		folderID:   os.Getenv("YANDEX_FOLDER_ID"),
		baseURL:    "https://llm.api.cloud.yandex.net/foundationModels/v1/completion",
		searchAPI:  "https://yandex.ru/search/xml",
		searchUser: os.Getenv("YANDEX_SEARCH_USER"),
		searchKey:  os.Getenv("YANDEX_SEARCH_KEY"),
		client:     &http.Client{},
	}
}

// GenerateDescription генерирует описание для места
func (s *ContentService) GenerateDescription(poi models.POI) (string, error) {
	// Сначала ищем дополнительную информацию через Yandex Search (опционально)
	additionalInfo := ""
	if s.searchUser != "" && s.searchKey != "" {
		info, err := s.searchForPOI(poi)
		if err == nil && info != "" {
			additionalInfo = info
		}
	}

	// Генерируем описание через YandexGPT
	prompt := s.buildPrompt(poi, additionalInfo)

	request := YandexGPTRequest{
		ModelURI: fmt.Sprintf("gpt://%s/yandexgpt/latest", s.folderID),
		CompletionOptions: CompletionOptions{
			Stream:      false,
			Temperature: 0.7,
			MaxTokens:   "1000",
		},
		Messages: []YandexMessage{
			{
				Role: "system",
				Text: "Ты - профессиональный экскурсовод с отличным знанием истории и умением увлекательно рассказывать. Твои рассказы точны, интересны и без ошибок. Ты специализируешься на истории Москвы и России.",
			},
			{
				Role: "user",
				Text: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Api-Key "+s.apiKey)
	req.Header.Set("x-folder-id", s.folderID)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call YandexGPT API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("YandexGPT API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result YandexGPTResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Result.Alternatives) == 0 {
		return "", fmt.Errorf("no alternatives in response")
	}

	return result.Result.Alternatives[0].Message.Text, nil
}

// searchForPOI ищет дополнительную информацию о месте через Yandex Search API
func (s *ContentService) searchForPOI(poi models.POI) (string, error) {
	// Формируем поисковый запрос
	query := fmt.Sprintf("%s Москва история", poi.Name)

	// Параметры запроса
	params := url.Values{}
	params.Set("user", s.searchUser)
	params.Set("key", s.searchKey)
	params.Set("query", query)
	params.Set("l10n", "ru")
	params.Set("sortby", "rlv") // сортировка по релевантности
	params.Set("maxpassages", "3")
	params.Set("groupby", "attr=d.mode=deep.groups-on-page=5.docs-in-group=1")

	searchURL := s.searchAPI + "?" + params.Encode()

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("search API error: status %d", resp.StatusCode)
	}

	// Читаем XML ответ (базовая обработка)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Извлекаем текстовые фрагменты (упрощенная версия)
	// В продакшене нужен полноценный XML парсер
	bodyStr := string(body)

	// Ищем теги <passage> с текстом
	passages := extractPassages(bodyStr)
	if len(passages) > 0 {
		return strings.Join(passages[:min(3, len(passages))], " "), nil
	}

	return "", nil
}

func (s *ContentService) buildPrompt(poi models.POI, additionalInfo string) string {
	epochNames := map[string]string{
		"medieval": "средневековье",
		"imperial": "имперский период",
		"soviet":   "советский период",
		"modern":   "современность",
	}

	categoryNames := map[string]string{
		"architecture": "архитектура",
		"history":      "история",
		"culture":      "культура",
		"religion":     "религия",
		"art":          "искусство",
	}

	prompt := fmt.Sprintf(`Создай увлекательный рассказ для аудиогида о следующем месте:

Название: %s
Описание: %s
Эпоха: %s
Категория: %s`,
		poi.Name,
		poi.Description,
		epochNames[poi.Epoch],
		categoryNames[poi.Category],
	)

	if poi.YearBuilt > 0 {
		prompt += fmt.Sprintf("\nГод постройки: %d", poi.YearBuilt)
	}

	if poi.Architect != "" {
		prompt += fmt.Sprintf("\nАрхитектор: %s", poi.Architect)
	}

	if poi.Style != "" {
		prompt += fmt.Sprintf("\nСтиль: %s", poi.Style)
	}

	if additionalInfo != "" {
		prompt += fmt.Sprintf("\n\nДополнительная информация из интернета:\n%s", additionalInfo)
	}

	prompt += `

Требования к рассказу:
- Длина: 2-3 минуты чтения (300-400 слов)
- Стиль: живой, увлекательный, но точный
- Включи интересные факты и легенды
- Без речевых ошибок и анахронизмов
- Используй яркие образы и эмоции
- Обращайся к слушателю на "вы"
- Начни с интригующей фразы
- Закончи призывом обратить внимание на детали

Формат: Чистый текст без заголовков и разметки.`

	return prompt
}

// Структуры для YandexGPT API

type YandexMessage struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type CompletionOptions struct {
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
	MaxTokens   string  `json:"maxTokens"`
}

type YandexGPTRequest struct {
	ModelURI          string            `json:"modelUri"`
	CompletionOptions CompletionOptions `json:"completionOptions"`
	Messages          []YandexMessage   `json:"messages"`
}

type YandexGPTResponse struct {
	Result struct {
		Alternatives []struct {
			Message YandexMessage `json:"message"`
			Status  string        `json:"status"`
		} `json:"alternatives"`
		Usage struct {
			InputTextTokens  string `json:"inputTextTokens"`
			CompletionTokens string `json:"completionTokens"`
			TotalTokens      string `json:"totalTokens"`
		} `json:"usage"`
		ModelVersion string `json:"modelVersion"`
	} `json:"result"`
}

// Вспомогательные функции

func extractPassages(xmlContent string) []string {
	// Упрощенная экстракция текста из XML
	// В продакшене использовать encoding/xml
	passages := []string{}

	start := 0
	for {
		idx := strings.Index(xmlContent[start:], "<passage>")
		if idx == -1 {
			break
		}
		start += idx + 9

		end := strings.Index(xmlContent[start:], "</passage>")
		if end == -1 {
			break
		}

		passage := xmlContent[start : start+end]
		// Убираем HTML теги
		passage = strings.ReplaceAll(passage, "<hlword>", "")
		passage = strings.ReplaceAll(passage, "</hlword>", "")
		passage = strings.TrimSpace(passage)

		if passage != "" {
			passages = append(passages, passage)
		}

		start += end + 10
	}

	return passages
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
