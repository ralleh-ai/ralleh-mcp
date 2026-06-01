package model

import "time"

type Brand struct {
	OrgID       string    `json:"orgId"`
	BrandID     string    `json:"brandId"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Mission     string    `json:"mission,omitempty"`
	Vision      string    `json:"vision,omitempty"`
	Values      []string  `json:"values,omitempty"`
	Industry    string    `json:"industry,omitempty"`
	Competitors []string  `json:"competitors,omitempty"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type BrandVoice struct {
	OrgID            string   `json:"orgId"`
	BrandID          string   `json:"brandId"`
	Tone             []string `json:"tone,omitempty"`
	WritingStyle     string   `json:"writingStyle,omitempty"`
	ReadingLevel     string   `json:"readingLevel,omitempty"`
	VocabularyRules  []string `json:"vocabularyRules,omitempty"`
	ForbiddenTerms   []string `json:"forbiddenTerms,omitempty"`
	PreferredPhrases []string `json:"preferredPhrases,omitempty"`
	Version          int      `json:"version"`
}

type Persona struct {
	PersonaID      string            `json:"personaId"`
	OrgID          string            `json:"orgId"`
	BrandID        string            `json:"brandId"`
	Name           string            `json:"name"`
	Demographics   map[string]string `json:"demographics,omitempty"`
	Interests      []string          `json:"interests,omitempty"`
	Goals          []string          `json:"goals,omitempty"`
	PainPoints     []string          `json:"painPoints,omitempty"`
	BuyingTriggers []string          `json:"buyingTriggers,omitempty"`
	Version        int               `json:"version"`
}

type MessagingFramework struct {
	OrgID                 string            `json:"orgId"`
	BrandID               string            `json:"brandId"`
	ValuePropositions     []string          `json:"valuePropositions,omitempty"`
	ElevatorPitch         string            `json:"elevatorPitch,omitempty"`
	PositioningStatements []string          `json:"positioningStatements,omitempty"`
	ProductMessaging      map[string]string `json:"productMessaging,omitempty"`
	Version               int               `json:"version"`
}

type StyleGuide struct {
	OrgID             string   `json:"orgId"`
	BrandID           string   `json:"brandId"`
	FormattingRules   []string `json:"formattingRules,omitempty"`
	CTARules          []string `json:"ctaRules,omitempty"`
	ComplianceRules   []string `json:"complianceRules,omitempty"`
	LegalRestrictions []string `json:"legalRestrictions,omitempty"`
	Version           int      `json:"version"`
}

type Campaign struct {
	CampaignID     string            `json:"campaignId"`
	OrgID          string            `json:"orgId"`
	BrandID        string            `json:"brandId"`
	Name           string            `json:"name"`
	Objective      string            `json:"objective,omitempty"`
	Channels       []string          `json:"channels,omitempty"`
	Results        map[string]string `json:"results,omitempty"`
	LessonsLearned []string          `json:"lessonsLearned,omitempty"`
	CreatedAt      time.Time         `json:"createdAt"`
}

type ValidationRequest struct {
	OrgID   string `json:"orgId"`
	BrandID string `json:"brandId"`
	Content string `json:"content"`
	Channel string `json:"channel,omitempty"`
	Rewrite bool   `json:"rewrite"`
}

type Violation struct {
	Rule       string `json:"rule"`
	Severity   string `json:"severity"`
	Text       string `json:"text"`
	Suggestion string `json:"suggestion"`
}

type ValidationResult struct {
	BrandComplianceScore int         `json:"brandComplianceScore"`
	Confidence           float64     `json:"confidence"`
	Violations           []Violation `json:"violations"`
	Suggestions          []string    `json:"suggestions"`
	RewrittenVersion     string      `json:"rewrittenVersion,omitempty"`
}

type AuditEvent struct {
	EventID   string    `json:"eventId"`
	OrgID     string    `json:"orgId"`
	BrandID   string    `json:"brandId"`
	Actor     string    `json:"actor"`
	Tool      string    `json:"tool"`
	Action    string    `json:"action"`
	Entity    string    `json:"entity"`
	EntityID  string    `json:"entityId"`
	Version   int       `json:"version"`
	Hash      string    `json:"hash"`
	Reason    string    `json:"reason,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
