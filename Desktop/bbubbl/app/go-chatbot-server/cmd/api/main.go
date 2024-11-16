package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"

	"go-chatbot-server/config"
	"go-chatbot-server/db"
	sqlc "go-chatbot-server/db/sqlc"
	"go-chatbot-server/router"
	"go-chatbot-server/server"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/qdrant/go-client/qdrant"
	goopenai "github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms/openai"
	"go.uber.org/zap"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"google.golang.org/api/option"
)

const createMode = false

// const createMode = true
const testMode = false

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env 파일을 찾을 수 없습니다. 기본 환경 변수를 사용합니다.", err)
	}

	logger, err := config.InitLogger()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	dbConn, err := db.Connect(logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbConn.Close()

	queries := sqlc.New(dbConn)

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		logger.Fatal("GEMINI_API_KEY not set in .env file")
	}
	ctx := context.Background()
	geminiClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		logger.Fatal("Failed to create Gemini client", zap.Error(err))
	}
	defer geminiClient.Close()

	llm, err := openai.New(
		openai.WithModel("gpt-4o-mini"),
	)
	if err != nil {
		logger.Fatal("Failed to create OpenAI client", zap.Error(err))
	}

	// Add this after creating the OpenAI client
	embeddingClient := goopenai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Remove the tokenizer initialization
	// tk := pretrained.BertBaseUncased()

	// Add this line near the top of the main function

	logger.Info("Creating qdrant client...")
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		logger.Fatal("Failed to create Qdrant client", zap.Error(err))
	}

	if createMode {
		// Remove all collections
		client.DeleteCollection(ctx, "faq")
		client.DeleteCollection(ctx, "faq_keywords")
		client.DeleteCollection(ctx, "company")
		client.DeleteCollection(ctx, "inspection_schedules")
		client.DeleteCollection(ctx, "regional_data")

		// Create collections for each CSV file
		// collections := []string{"faq", "faq_keywords", "company", "inspection_schedules", "regional_data"}
		collections := []string{"faq"}
		for _, collection := range collections {
			client.CreateCollection(ctx, &qdrant.CreateCollection{
				CollectionName: collection,
				VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
					Size:     3072,
					Distance: qdrant.Distance_Cosine,
				}),
			})
		}

		// Process and embed CSV files
		csvFiles := map[string]string{
			"faq": "./csv/faq.csv", // 원격 서버에는 현재 이것만 학습함
			// "faq_keywords": "./csv/faq_keywords.csv",
			// "company":              "./csv/company.csv",
			// "inspection_schedules": "./csv/inspection_schedules.csv",
			// "regional_data": "./csv/regional_data.csv",
		}

		// Update the processAndEmbedCSV function calls
		for collection, filePath := range csvFiles {
			err := processAndEmbedCSV(ctx, embeddingClient, client, collection, filePath, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to process and embed %s", filePath), zap.Error(err))
			}
		}
	} else {
		logger.Info("Skipping collection creation and data embedding")
	}

	// testMode 변수를 사용하여 테스트 실행 여부 제어
	if testMode {
		err = testSearchAndCreateCSV(ctx, client, embeddingClient, "./csv/faq.csv", "./csv/search_results.csv", logger)
		if err != nil {
			logger.Error("Failed to test search and create CSV", zap.Error(err))
		} else {
			logger.Info("Search test completed and results saved to search_results.csv")
		}
	} else {
		logger.Info("Skipping search test")
	}

	r := router.New(queries, logger, dbConn, geminiClient, llm, client, embeddingClient)

	srv := server.New(r.Engine(), logger, ":8080")
	logger.Info("Starting server on :8080")
	if err := srv.Run(); err != nil {
		logger.Fatal("cannot start server", zap.Error(err))
	}
}

func processAndEmbedCSV(ctx context.Context, embeddingClient *goopenai.Client, qdrantClient *qdrant.Client, collection, filePath string, logger *zap.Logger) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true // Handle quotes in CSV

	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header from %s: %w", filePath, err)
	}

	var records []map[string]string
	var id uint64 = 1
	// 10개 데이터만 처리
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record from %s: %w", filePath, err)
		}

		// Create a map for each record
		recordMap := make(map[string]string)
		for i, value := range record {
			recordMap[header[i]] = value
		}
		records = append(records, recordMap)

		id++
	}

	var points []*qdrant.PointStruct

	// Update the switch statement to use embeddingClient instead of tk
	switch collection {
	case "faq":
		points = processFAQ(ctx, records, embeddingClient)
	case "faq_keywords":
		points = processFAQKeywords(ctx, records, embeddingClient)
	case "company":
		points = processCompany(ctx, records, embeddingClient)
	case "inspection_schedules":
		points = processInspectionSchedules(ctx, records, embeddingClient)
	case "regional_data":
		points = processRegionalData(ctx, records, embeddingClient)
	default:
		return fmt.Errorf("unknown collection: %s", collection)
	}

	// Upsert points to Qdrant
	_, err = qdrantClient.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collection,
		Points:         points,
	})
	if err != nil {
		return fmt.Errorf("failed to upsert points to collection %s: %w", collection, err)
	}

	logger.Info(fmt.Sprintf("Successfully embedded %d records into collection %s", len(points), collection))
	return nil
}

func processFAQ(ctx context.Context, records []map[string]string, embeddingClient *goopenai.Client) []*qdrant.PointStruct {
	var points []*qdrant.PointStruct
	// maxRecords := 50
	// if len(records) > maxRecords {
	// 	records = records[:maxRecords]
	// }
	for i, record := range records {
		fullText := fmt.Sprintf("%s %s",
			record["법령명"], record["소관부처명"])
		embedding, err := getEmbedding(ctx, embeddingClient, fullText)
		fmt.Println(len(embedding), fullText)
		if err != nil {
			log.Printf("Failed to get embedding for record %d: %v", i, err)
			continue
		}
		point := &qdrant.PointStruct{
			Id:      qdrant.NewIDNum(uint64(i + 1)),
			Vectors: qdrant.NewVectorsDense(embedding),
			Payload: map[string]*qdrant.Value{
				"순번":     qdrant.NewValueString(record["순번"]),
				"법령MST":  qdrant.NewValueString(record["법령MST"]),
				"소관부처코드": qdrant.NewValueString(record["소관부처코드"]),
				"소관부처명":  qdrant.NewValueString(record["소관부처명"]),
				"법령ID":   qdrant.NewValueString(record["법령ID"]),
				"법령명":    qdrant.NewValueString(record["법령명"]),
				"공포일자":   qdrant.NewValueString(record["공포일자"]),
				"공포번호":   qdrant.NewValueString(record["공포번호"]),
				"시행일자":   qdrant.NewValueString(record["시행일자"]),
				"법령구분코드": qdrant.NewValueString(record["법령구분코드"]),
				"법령구분명":  qdrant.NewValueString(record["법령구분명"]),
				"법령분야코드": qdrant.NewValueString(record["법령분야코드"]),
				"법령분야명":  qdrant.NewValueString(record["법령분야명"]),
			},
		}
		points = append(points, point)
	}
	return points
}

func processFAQKeywords(ctx context.Context, records []map[string]string, embeddingClient *goopenai.Client) []*qdrant.PointStruct {
	var points []*qdrant.PointStruct
	for i, record := range records {
		fullText := fmt.Sprintf("%s %s %s", record["업종"], record["핵심키워드"], record["질문"])
		embedding, err := getEmbedding(ctx, embeddingClient, fullText)
		if err != nil {
			log.Printf("Failed to get embedding for record %d: %v", i, err)
			continue
		}
		point := &qdrant.PointStruct{
			Id:      qdrant.NewIDNum(uint64(i + 1)),
			Vectors: qdrant.NewVectorsDense(embedding),
			Payload: map[string]*qdrant.Value{
				"업종":    qdrant.NewValueString(record["업종"]),
				"핵심키워드": qdrant.NewValueString(record["핵심키워드"]),
				"질문":    qdrant.NewValueString(record["질문"]),
			},
		}
		points = append(points, point)
	}
	return points
}

func processCompany(ctx context.Context, records []map[string]string, embeddingClient *goopenai.Client) []*qdrant.PointStruct {
	var points []*qdrant.PointStruct
	for i, record := range records {
		fullText := fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s",
			record["분야필터"], record["기관명"], record["대표자명"], record["상세보기소재지"],
			record["도/시"], record["시/군/구"], record["전화"], record["팩스"],
			record["비고"], record["유효기간"])
		embedding, err := getEmbedding(ctx, embeddingClient, fullText)
		if err != nil {
			log.Printf("Failed to get embedding for record %d: %v", i, err)
			continue
		}
		point := &qdrant.PointStruct{
			Id:      qdrant.NewIDNum(uint64(i + 1)),
			Vectors: qdrant.NewVectorsDense(embedding),
			Payload: map[string]*qdrant.Value{
				"분야필터":    qdrant.NewValueString(record["분야필터"]),
				"기관명":     qdrant.NewValueString(record["기관명"]),
				"대표자명":    qdrant.NewValueString(record["대표자명"]),
				"상세보기소재지": qdrant.NewValueString(record["상세보기소재지"]),
				"도/시":     qdrant.NewValueString(record["도/시"]),
				"시/군/구":   qdrant.NewValueString(record["시/군/구"]),
				"전화":      qdrant.NewValueString(record["전화"]),
				"팩스":      qdrant.NewValueString(record["팩스"]),
				"비고":      qdrant.NewValueString(record["비고"]),
				"유효기간":    qdrant.NewValueString(record["유효기간"]),
			},
		}
		points = append(points, point)
	}
	return points
}

func processInspectionSchedules(ctx context.Context, records []map[string]string, embeddingClient *goopenai.Client) []*qdrant.PointStruct {
	var points []*qdrant.PointStruct
	for i, record := range records {
		fullText := fmt.Sprintf("%s %s %s %s %s %s %s %s",
			record["대분류"], record["중분류"], record["소분류"], record["제품"],
			record["검사항목"], record["검사금액"], record["검사주기"], record["비고"])
		embedding, err := getEmbedding(ctx, embeddingClient, fullText)
		if err != nil {
			log.Printf("Failed to get embedding for record %d: %v", i, err)
			continue
		}
		point := &qdrant.PointStruct{
			Id:      qdrant.NewIDNum(uint64(i + 1)),
			Vectors: qdrant.NewVectorsDense(embedding),
			Payload: map[string]*qdrant.Value{
				"대분류":  qdrant.NewValueString(record["대분류"]),
				"중분류":  qdrant.NewValueString(record["중분류"]),
				"소분류":  qdrant.NewValueString(record["소분류"]),
				"제품":   qdrant.NewValueString(record["제품"]),
				"검사항목": qdrant.NewValueString(record["검사항목"]),
				"검사금액": qdrant.NewValueString(record["검사금액"]),
				"검사주기": qdrant.NewValueString(record["검사주기"]),
				"비고":   qdrant.NewValueString(record["비고"]),
			},
		}
		points = append(points, point)
	}
	return points
}

func processRegionalData(ctx context.Context, records []map[string]string, embeddingClient *goopenai.Client) []*qdrant.PointStruct {
	var points []*qdrant.PointStruct
	for i, record := range records {
		fullText := fmt.Sprintf("%s %s", record["시도명"], record["시군구명"])
		embedding, err := getEmbedding(ctx, embeddingClient, fullText)
		if err != nil {
			log.Printf("Failed to get embedding for record %d: %v", i, err)
			continue
		}
		point := &qdrant.PointStruct{
			Id:      qdrant.NewIDNum(uint64(i + 1)),
			Vectors: qdrant.NewVectorsDense(embedding),
			Payload: map[string]*qdrant.Value{
				"시도명":  qdrant.NewValueString(record["시도명"]),
				"시군구명": qdrant.NewValueString(record["시군구명"]),
			},
		}
		points = append(points, point)
	}
	return points
}

func getEmbedding(ctx context.Context, client *goopenai.Client, text string) ([]float32, error) {
	processedText := advancedPreprocessText(text)
	resp, err := client.CreateEmbeddings(ctx, goopenai.EmbeddingRequest{
		Input: []string{processedText},
		Model: goopenai.LargeEmbedding3,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}
	return resp.Data[0].Embedding, nil
}

// Add this new function at the end of the file
func searchFAQ(ctx context.Context, client *qdrant.Client, embeddingClient *goopenai.Client, query string, logger *zap.Logger) (string, error) {
	embedding, err := getEmbedding(ctx, embeddingClient, query)
	if err != nil {
		return "", err
	}
	logger.Info("Embedding", zap.Any("embedding", embedding))
	searchResult, err := client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: "faq",
		Query:          qdrant.NewQueryDense(embedding),
		// WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil {
		return "", err
	}

	searchResultString := fmt.Sprintf("%v", searchResult)
	return searchResultString, nil
}

// Add these new functions at the end of the file
func testSearchAndCreateCSV(ctx context.Context, client *qdrant.Client, embeddingClient *goopenai.Client, inputCSVPath, outputCSVPath string, logger *zap.Logger) error {
	// Read input CSV
	inputFile, err := os.Open(inputCSVPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	reader := csv.NewReader(inputFile)
	reader.LazyQuotes = true

	// Skip header
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// Create output CSV
	outputFile, err := os.Create(outputCSVPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	writer.Comma = '\t'
	defer writer.Flush()

	// Write header to output CSV
	if err := writer.Write([]string{"Query", "Top Result", "Relevance", "Error"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Add these variables for evaluation
	var totalQueries int
	var relevantResults int

	// Process each row
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error("Failed to read record", zap.Error(err))
			continue
		}

		query := record[4]
		expectedAnswer := record[5]
		result, searchErr := searchFAQTest(ctx, client, embeddingClient, query, logger)

		var errorMsg string
		if searchErr != nil {
			errorMsg = searchErr.Error()
			logger.Error("Failed to search FAQ", zap.Error(searchErr), zap.String("query", query))
		}

		// Evaluate the result
		isRelevant := evaluateResult(result, expectedAnswer)
		if isRelevant {
			relevantResults++
		}
		totalQueries++

		// Write result to output CSV, including relevance and any error message
		if err := writer.Write([]string{query, result, fmt.Sprintf("%v", isRelevant), errorMsg}); err != nil {
			logger.Error("Failed to write result", zap.Error(err))
			// Continue processing other rows even if writing fails
		}
	}

	// Calculate and log the accuracy
	accuracy := float64(relevantResults) / float64(totalQueries) * 100
	logger.Info("Search evaluation completed",
		zap.Int("totalQueries", totalQueries),
		zap.Int("relevantResults", relevantResults),
		zap.Float64("accuracy", accuracy))

	// Write the accuracy to the CSV file
	if err := writer.Write([]string{"Accuracy", fmt.Sprintf("%.2f%%", accuracy)}); err != nil {
		logger.Error("Failed to write accuracy", zap.Error(err))
	}

	return nil
}

func evaluateResult(result, expectedAnswer string) bool {
	// Simple evaluation: check if the result contains the expected answer
	// You might want to implement a more sophisticated evaluation method
	return strings.Contains(strings.ToLower(result), strings.ToLower(expectedAnswer))
}

// Update the searchFAQTest function to return both question and answer
func searchFAQTest(ctx context.Context, client *qdrant.Client, embeddingClient *goopenai.Client, query string, logger *zap.Logger) (string, error) {
	embedding, err := getEmbedding(ctx, embeddingClient, query)
	if err != nil {
		return "", err
	}
	searchResult, err := client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: "faq",
		Query:          qdrant.NewQueryDense(embedding),
		WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil {
		return "", err
	}

	if len(searchResult) == 0 {
		return "No results found", nil
	}

	// Extract the question and answer from the top result
	topResult := searchResult[0]
	question := topResult.Payload["질문"].GetStringValue()
	answer := topResult.Payload["답변"].GetStringValue()

	return fmt.Sprintf("Q: %s\nA: %s", question, answer), nil
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
