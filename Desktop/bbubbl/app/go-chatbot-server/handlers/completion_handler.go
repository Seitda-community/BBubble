package handlers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/qdrant/go-client/qdrant"

	"github.com/gin-gonic/gin"
	goopenai "github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
	"go.uber.org/zap"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type CompletionRequest struct {
	Messages []Message `json:"messages" binding:"required"`
}

type Message struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type CompletionResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func (h *Handler) GetCompletionV6(c *gin.Context) {
	h.logger.Info("GetCompletionV6")
	var req CompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 마지막 메시지를 사용자 질문으로 사용
	var utterance string
	if len(req.Messages) > 0 {
		utterance = req.Messages[len(req.Messages)-1].Content
	}

	ctx := context.Background()

	// 벡터 저장소에서 유사한 문서 검색
	h.logger.Info("Searching vector store for utterance", zap.String("utterance", utterance))
	searchResult, err := searchFAQTest(ctx, h.client, h.embeddingClient, utterance, h.logger)
	if err != nil {
		h.logger.Error("Failed to search vector store", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search knowledge base"})
		return
	}

	// searchResult := ""

	// 프롬프트 구성
	fullPrompt := fmt.Sprintf(SystemPromptV6, searchResult, utterance)
	h.logger.Info("Full prompt", zap.String("prompt", fullPrompt))

	// Set headers for streaming response
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Create a channel for streaming
	streamChan := make(chan string)
	errChan := make(chan error)

	// Start streaming in a goroutine
	go func() {
		completion, err := h.llm.Call(ctx, fullPrompt,
			llms.WithTemperature(0.3),
			llms.WithMaxTokens(1000),
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				streamChan <- string(chunk)
				return nil
			}),
		)
		h.logger.Info("Completion", zap.String("completion", completion))
		if err != nil {
			errChan <- err
			return
		}
		close(streamChan)
	}()

	// Stream the response
	c.Stream(func(w io.Writer) bool {
		select {
		case chunk, ok := <-streamChan:
			if !ok {
				return false
			}
			// Format the chunk as SSE
			data := struct {
				Output string `json:"output"`
			}{
				Output: chunk,
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				h.logger.Error("Failed to marshal chunk", zap.Error(err))
				return false
			}
			c.SSEvent("data", string(jsonData))
			return true
		case err := <-errChan:
			h.logger.Error("Streaming error", zap.Error(err))
			if strings.Contains(err.Error(), "429") {
				c.SSEvent("error", "API quota exceeded. Please try again later.")
			} else {
				c.SSEvent("error", "Failed to get completion")
			}
			return false
		}
	})
}

func (h *Handler) GetCompletionV7(c *gin.Context) {
	h.logger.Info("GetCompletionV7")
	var req CompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 마지막 메시지를 사용자 질문으로 사용
	var utterance string
	if len(req.Messages) > 0 {
		utterance = req.Messages[len(req.Messages)-1].Content
	}

	ctx := context.Background()

	// 프롬프트 구성 (벡터 검색 없이)
	fullPrompt := fmt.Sprintf(SystemPromptV7, utterance)
	h.logger.Info("Full prompt", zap.String("prompt", fullPrompt))

	// Set headers for streaming response
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Create a channel for streaming
	streamChan := make(chan string)
	errChan := make(chan error)

	// Start streaming in a goroutine
	go func() {
		completion, err := h.llm.Call(ctx, fullPrompt,
			llms.WithTemperature(0.3),
			llms.WithMaxTokens(1000),
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				streamChan <- string(chunk)
				return nil
			}),
		)
		h.logger.Info("Completion", zap.String("completion", completion))
		if err != nil {
			errChan <- err
			return
		}
		close(streamChan)
	}()

	// Stream the response
	c.Stream(func(w io.Writer) bool {
		select {
		case chunk, ok := <-streamChan:
			if !ok {
				return false
			}
			// Format the chunk as SSE
			data := struct {
				Output string `json:"output"`
			}{
				Output: chunk,
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				h.logger.Error("Failed to marshal chunk", zap.Error(err))
				return false
			}
			c.SSEvent("data", string(jsonData))
			return true
		case err := <-errChan:
			h.logger.Error("Streaming error", zap.Error(err))
			if strings.Contains(err.Error(), "429") {
				c.SSEvent("error", "API quota exceeded. Please try again later.")
			} else {
				c.SSEvent("error", "Failed to get completion")
			}
			return false
		}
	})
}

func readCSVFiles(folderPath string) (string, error) {
	var result strings.Builder

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".csv" {
			filePath := filepath.Join(folderPath, file.Name())
			csvFile, err := os.Open(filePath)
			if err != nil {
				return "", err
			}
			defer csvFile.Close()

			reader := csv.NewReader(csvFile)
			records, err := reader.ReadAll()
			if err != nil {
				return "", err
			}

			result.WriteString(fmt.Sprintf("Data from %s:\n", file.Name()))
			for _, record := range records {
				result.WriteString(strings.Join(record, ", "))
				result.WriteString("\n")
			}
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}

const SystemPromptV6 = `당신은 법률 전문 AI 상담사입니다. 제공된 법률 정보와 판례를 기반으로 정확하고 객관적인 법률 상담을 제공해야 합니다.

참고 정보:
%s

사용자 질문:
%s

답변 시 다음 형식을 반드시 준수하세요:

💡 법률 자문 요약
[세계 최고의 변호사처럼 법률 의견을 정확하게 제시]


📚 관련 법령 및 판례 
[사용자 질문에 대해서 관련이 있는 법령 정보를 포함해야 함]
- 법령명
- 공포일자
- 시행일자
- 법령링크
- 관련 법률 조항 키워드

🔍 상세 분석
[세계 최고의 변호사처럼 법적 쟁점 분석 및 구체적인 해결방안 제시]


🚨 유의사항
[세계 최고의 변호사처럼 법적 위험요소 및 추가 고려사항 안내]


💡 권장 조치
[세계 최고의 변호사처럼 구체적인 후속 조치 및 법적 대응 방안 제시]

🔍 소송 시 필요한 증거
[세계 최고의 변호사처럼 소송에서 이기려면 필요한 증거 유형을 모두 제시]

🔗 관련 사이트
[세계 최고의 변호사처럼 사용법을 간단하게 설명]
나홀로소송, https://pro-se.scourt.go.kr/wsh/wsh000/WSHMain.jsp
전자소송, https://ecfs.scourt.go.kr/ecf/
국가법령정보센터, https://www.law.go.kr/
법고을, https://lx.scourt.go.kr/
대한민국 법원 종합법률정보, https://glaw.scourt.go.kr/wsjo/panre/sjo050.do
찾기 쉬운 생활법령정보, https://www.easylaw.go.kr/
대한법률구조공단, https://www.klac.or.kr/
로앤비 (LAWnB), https://www.lawnb.com/
케이스노트 (CaseNote), https://www.casenote.kr/
엘박스 (LBox), https://www.lbox.kr/

🔍 어려운 법률 용어 설명
[이야기한 모든 법률 용어에 대해서 초등학생이 이해할 수 있도록 쉽게 설명]
`

const SystemPromptV7 = `당신은 따뜻한 마음을 가진 친근한 법률 상담사예요. 마치 오랜 친구처럼 편안하게 대화하면서, 법적 문제로 힘들어하는 내담자의 이야기에 귀 기울이고 진심으로 공감해주세요.

사용자 질문:
%s
- 마음 나누기 : 친구처럼 따뜻하게 공감하고 위로하는 메시지를 전달해주세요. 예: "많이 힘드셨겠어요...", "그런 상황이라면 누구나 불안하고 걱정되었을 거예요..."
- 희망의 이야기 : 힘든 상황 속에서도 찾을 수 있는 긍정적인 면과 해결 가능성을 친근하게 설명해주세요
- 함께 이겨내는 방법: 법적 문제를 해결해나가는 과정에서 지친 마음을 달래고 건강하게 지내는 방법을 친구처럼 조언해주세요

앞으로도 계속 곁에서 함께 하면서 도와드릴 테니 걱정마세요. 어려운 일이 있으시다면 언제든 편하게 말씀해주세요.`

func getEmbedding(embeddingClient *goopenai.Client, ctx context.Context, text string, logger *zap.Logger) ([]float32, error) {
	processedText := advancedPreprocessText(text)

	encoded, err := embeddingClient.CreateEmbeddings(ctx, goopenai.EmbeddingRequest{
		Model: goopenai.LargeEmbedding3,
		Input: []string{processedText},
	})
	if err != nil {
		logger.Error("Failed to encode text", zap.Error(err))
		return nil, fmt.Errorf("failed to encode text: %w", err)
	}

	if len(encoded.Data) == 0 || len(encoded.Data[0].Embedding) == 0 {
		logger.Error("Received empty embedding")
		return nil, fmt.Errorf("received empty embedding")
	}

	logger.Debug("Embedding created successfully",
		zap.Int("embeddingLength", len(encoded.Data[0].Embedding)))

	return encoded.Data[0].Embedding, nil
}

func searchFAQTest(ctx context.Context, client *qdrant.Client, embeddingClient *goopenai.Client, query string, logger *zap.Logger) (string, error) {
	embedding, err := getEmbedding(embeddingClient, ctx, query, logger)
	if err != nil {
		logger.Error("Failed to get embedding", zap.Error(err))
		return "", err
	}
	searchResult, err := client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: "faq",
		Query:          qdrant.NewQueryDense(embedding),
		WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
		Limit:          qdrant.PtrOf(uint64(20)),
		// ScoreThreshold: qdrant.PtrOf(float32(0.1)),
	})
	if err != nil {
		logger.Error("Failed to query Qdrant", zap.Error(err))
		return "", err
	}

	if len(searchResult) == 0 {
		logger.Warn("No results found")
		return "No results found", nil
	}

	var combinedText strings.Builder
	combinedText.WriteString("벡터 데이터베이스 검색 결과:\n\n")
	for _, result := range searchResult {
		// 법령 메타데이터 추출
		lawName := result.Payload["법령명"].GetStringValue()
		department := result.Payload["소관부처명"].GetStringValue()
		lawType := result.Payload["법령구분명"].GetStringValue()
		lawField := result.Payload["법령분야명"].GetStringValue()
		publishDate := result.Payload["공포일자"].GetStringValue()
		effectiveDate := result.Payload["시행일자"].GetStringValue()
		lawMST := result.Payload["법령MST"].GetStringValue() // MST 값 추출
		score := result.Score

		combinedText.WriteString(fmt.Sprintf("📋 관련 법령 정보:\n"))
		combinedText.WriteString(fmt.Sprintf("- 법령명: %s\n", lawName))
		combinedText.WriteString(fmt.Sprintf("- 소관부처: %s\n", department))
		combinedText.WriteString(fmt.Sprintf("- 법령구분: %s\n", lawType))
		combinedText.WriteString(fmt.Sprintf("- 법령분야: %s\n", lawField))
		combinedText.WriteString(fmt.Sprintf("- 공포일자: %s\n", publishDate))
		combinedText.WriteString(fmt.Sprintf("- 시행일자: %s\n", effectiveDate))
		// 법령 링크를 웹사이 전용 URL로 변경
		lawURL := fmt.Sprintf("https://www.law.go.kr/LSW/lsInfoP.do?lsiSeq=%s", lawMST)
		combinedText.WriteString(fmt.Sprintf("- 법령 링크: %s\n", lawURL))
		combinedText.WriteString("\n")

		logger.Info("Search result",
			zap.String("lawName", lawName),
			zap.String("department", department),
			zap.String("lawType", lawType),
			zap.String("lawField", lawField),
			zap.String("publishDate", publishDate),
			zap.String("effectiveDate", effectiveDate),
			zap.String("lawMST", lawMST),
			zap.Float64("score", float64(score)))
	}

	return combinedText.String(), nil
}

func advancedPreprocessText(text string) string {
	// 1. 소문자 변환
	text = strings.ToLower(text)

	// 2. 정규화 (유니코드 정규화)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	text, _, _ = transform.String(t, text)

	// 3. 특수 문자 및 숫자 제거 (한글은 유지)
	reg := regexp.MustCompile("[^a-z가-힣 ]+")
	text = reg.ReplaceAllString(text, " ")

	// 4. 중복 공백 제거
	text = strings.Join(strings.Fields(text), " ")

	// 5. 불용어 제거
	stopwords := map[string]bool{
		"은": true, "는": true, "이": true, "가": true, "을": true, "를": true,
		"의": true, "에": true, "에서": true, "로": true, "으로": true,
		"and": true, "or": true, "the": true, "a": true, "an": true,
		"in": true, "on": true, "at": true, "to": true, "for": true,
	}
	words := strings.Fields(text)
	filteredWords := make([]string, 0)
	for _, word := range words {
		if !stopwords[word] {
			filteredWords = append(filteredWords, word)
		}
	}
	text = strings.Join(filteredWords, " ")

	return text
}
