package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/pkfk-discovery/api/internal/domain"
)

var (
	ErrAIProviderNotFound = errors.New("AI provider not found")
	ErrInvalidDiffFormat   = errors.New("invalid diff format")
)

type AIOptimizationService struct {
	aiProviderRepo     domain.AIProviderRepository
	aiInteractionRepo domain.AIInteractionRepository
}

type OptimizationRequest struct {
	DraftID      string            `json:"draft_id"`
	SQLTemplate  string            `json:"sql_template"`
	Context      map[string]string `json:"context"`
	ProviderID   string            `json:"provider_id"`
}

type OptimizationResponse struct {
	DiffPatch   string `json:"diff_patch"`
	Explanation string `json:"explanation"`
}

func NewAIOptimizationService(aiProviderRepo domain.AIProviderRepository, aiInteractionRepo domain.AIInteractionRepository) *AIOptimizationService {
	return &AIOptimizationService{
		aiProviderRepo:     aiProviderRepo,
		aiInteractionRepo:  aiInteractionRepo,
	}
}

func (s *AIOptimizationService) Optimize(ctx context.Context, req *OptimizationRequest, userID string) (*OptimizationResponse, error) {
	// Load AI provider
	providerID, err := uuid.Parse(req.ProviderID)
	if err != nil {
		return nil, ErrAIProviderNotFound
	}

	provider, err := s.aiProviderRepo.GetByID(providerID)
	if err != nil {
		return nil, ErrAIProviderNotFound
	}

	// Validate inputs (strict guardrails)
	allowedInputs := map[string]bool{
		"sql_template":    true,
		"schema":           true,
		"table":            true,
		"column":           true,
		"validation_error": true,
		"explain_output":  true,
	}

	for key := range req.Context {
		if !allowedInputs[key] {
			return nil, fmt.Errorf("disallowed input key: %s", key)
		}
	}

	// Build prompt (limited inputs only)
	prompt := s.buildPrompt(req)

	// Hash inputs for audit
	promptHash := s.hashString(prompt)

	// Call AI provider (placeholder - would integrate with actual provider)
	response, err := s.callAIProvider(ctx, provider, prompt)
	if err != nil {
		return nil, fmt.Errorf("AI provider call failed: %w", err)
	}

	// Validate response is unified diff format
	if !s.isValidDiff(response) {
		return nil, ErrInvalidDiffFormat
	}

	responseHash := s.hashString(response)

	// Log AI interaction
	interaction := &domain.AIInteraction{
		UserID:        uuid.MustParse(userID),
		Provider:      provider.Name,
		Model:         provider.Name, // TODO: Add Model field to AIProvider domain
		PromptHash:    promptHash,
		ResponseHash:  responseHash,
		AdapterDraftID: nil,
	}

	if draftID, err := uuid.Parse(req.DraftID); err == nil {
		interaction.AdapterDraftID = &draftID
	}

	interaction.ID = uuid.New()
	interaction.CreatedAt = time.Now()
	if err := s.aiInteractionRepo.Create(interaction); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to log AI interaction: %v\n", err)
	}

	return &OptimizationResponse{
		DiffPatch:   response,
		Explanation: "AI-generated optimization patch",
	}, nil
}

func (s *AIOptimizationService) buildPrompt(req *OptimizationRequest) string {
	var parts []string

	parts = append(parts, "You are an expert SQL optimizer. Analyze the provided SQL template and suggest improvements.")
	parts = append(parts, "Output ONLY unified diff format affecting only the SQL template file.")
	parts = append(parts, "")
	parts = append(parts, "SQL Template:")
	parts = append(parts, req.SQLTemplate)
	parts = append(parts, "")

	// Add allowed context
	if validationError, ok := req.Context["validation_error"]; ok {
		parts = append(parts, "Validation Error:")
		parts = append(parts, validationError)
		parts = append(parts, "")
	}

	if explainOutput, ok := req.Context["explain_output"]; ok {
		parts = append(parts, "EXPLAIN Output:")
		parts = append(parts, explainOutput)
		parts = append(parts, "")
	}

	return strings.Join(parts, "\n")
}

func (s *AIOptimizationService) callAIProvider(ctx context.Context, provider *domain.AIProvider, prompt string) (string, error) {
	// TODO: Implement actual AI provider integration
	// This is a placeholder that returns a sample diff

	// For now, return a placeholder diff
	return `--- a/sql/template.sql
+++ b/sql/template.sql
@@ -1,3 +1,3 @@
-SELECT * FROM table_name;
+SELECT id, name FROM table_name WHERE id > 0;
`, nil
}

func (s *AIOptimizationService) isValidDiff(diff string) bool {
	// Basic validation: check for unified diff markers
	lines := strings.Split(diff, "\n")
	foundHeader := false
	for _, line := range lines {
		if strings.HasPrefix(line, "--- ") || strings.HasPrefix(line, "+++ ") {
			foundHeader = true
		}
		if strings.HasPrefix(line, "@@") {
			return foundHeader
		}
	}
	return false
}

func (s *AIOptimizationService) hashString(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}

