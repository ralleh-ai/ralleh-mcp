package service

import (
	"context"
	"database/sql"
	"strings"

	"github.com/ralleh-ai/ralleh-mcp/internal/brand/model"
	"github.com/ralleh-ai/ralleh-mcp/internal/brand/store"
)

type Service struct{ Store *store.Store }

func (s Service) ValidateContent(ctx context.Context, req model.ValidationRequest) (model.ValidationResult, error) {
	voice, err := s.Store.GetVoice(ctx, req.OrgID, req.BrandID)
	if err != nil && err != sql.ErrNoRows {
		return model.ValidationResult{}, err
	}
	violations := []model.Violation{}
	contentLower := strings.ToLower(req.Content)
	for _, term := range voice.ForbiddenTerms {
		if term == "" {
			continue
		}
		if strings.Contains(contentLower, strings.ToLower(term)) {
			violations = append(violations, model.Violation{Rule: "forbidden_term", Severity: "medium", Text: term, Suggestion: "Remove or replace with approved brand language."})
		}
	}
	for _, restriction := range []string{"guaranteed roi", "replace your whole team", "fully autonomous", "no oversight"} {
		if strings.Contains(contentLower, restriction) {
			violations = append(violations, model.Violation{Rule: "legal_restriction", Severity: "high", Text: restriction, Suggestion: "Avoid unsupported autonomy, replacement, or ROI claims."})
		}
	}
	score := 100 - len(violations)*18
	if score < 0 {
		score = 0
	}
	suggestions := []string{}
	if len(violations) > 0 {
		suggestions = append(suggestions, "Revise flagged terms and legal/compliance claims.")
	}
	if len(voice.PreferredPhrases) > 0 {
		suggestions = append(suggestions, "Consider using approved phrases where natural: "+strings.Join(voice.PreferredPhrases, ", "))
	}
	res := model.ValidationResult{BrandComplianceScore: score, Confidence: 0.82, Violations: violations, Suggestions: suggestions}
	if req.Rewrite && len(violations) > 0 {
		res.RewrittenVersion = rewrite(req.Content, violations, voice)
	}
	return res, nil
}

func rewrite(content string, violations []model.Violation, voice model.BrandVoice) string {
	out := content
	for _, v := range violations {
		out = strings.ReplaceAll(out, v.Text, "")
	}
	out = strings.ReplaceAll(out, "  ", " ")
	if len(voice.PreferredPhrases) > 0 {
		out = strings.TrimSpace(out) + " " + voice.PreferredPhrases[0] + "."
	}
	return strings.TrimSpace(out)
}
