package service

import (
	"context"
	"testing"

	"github.com/ralleh-ai/ralleh-mcp/internal/brand/model"
	"github.com/ralleh-ai/ralleh-mcp/internal/brand/store"
)

func TestValidateContentFlagsForbiddenTerms(t *testing.T) {
	s, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	ctx := context.Background()
	_, _, err = s.UpsertVoice(ctx, model.BrandVoice{OrgID: "org", BrandID: "brand", ForbiddenTerms: []string{"magical"}, PreferredPhrases: []string{"verification before done"}}, "test", "brand.update_voice", "test")
	if err != nil {
		t.Fatal(err)
	}
	res, err := (Service{Store: s}).ValidateContent(ctx, model.ValidationRequest{OrgID: "org", BrandID: "brand", Content: "Our magical AI is fully autonomous with guaranteed ROI", Rewrite: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.BrandComplianceScore >= 100 || len(res.Violations) == 0 {
		t.Fatalf("expected violations, got %+v", res)
	}
	if res.RewrittenVersion == "" {
		t.Fatal("expected rewrite")
	}
}
