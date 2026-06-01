package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"github.com/ralleh-ai/ralleh-mcp/internal/brand/model"
)

type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	if err := s.Migrate(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) Migrate(ctx context.Context) error {
	stmts := []string{
		`PRAGMA journal_mode=WAL;`,
		`CREATE TABLE IF NOT EXISTS brands (org_id TEXT NOT NULL, brand_id TEXT NOT NULL, data_json TEXT NOT NULL, version INTEGER NOT NULL, created_at TEXT NOT NULL, updated_at TEXT NOT NULL, PRIMARY KEY(org_id, brand_id));`,
		`CREATE TABLE IF NOT EXISTS brand_voice (org_id TEXT NOT NULL, brand_id TEXT NOT NULL, data_json TEXT NOT NULL, version INTEGER NOT NULL, updated_at TEXT NOT NULL, PRIMARY KEY(org_id, brand_id));`,
		`CREATE TABLE IF NOT EXISTS personas (org_id TEXT NOT NULL, brand_id TEXT NOT NULL, persona_id TEXT NOT NULL, data_json TEXT NOT NULL, version INTEGER NOT NULL, updated_at TEXT NOT NULL, PRIMARY KEY(org_id, brand_id, persona_id));`,
		`CREATE TABLE IF NOT EXISTS campaigns (org_id TEXT NOT NULL, brand_id TEXT NOT NULL, campaign_id TEXT NOT NULL, data_json TEXT NOT NULL, created_at TEXT NOT NULL, PRIMARY KEY(org_id, brand_id, campaign_id));`,
		`CREATE TABLE IF NOT EXISTS versions (org_id TEXT NOT NULL, brand_id TEXT NOT NULL, entity TEXT NOT NULL, entity_id TEXT NOT NULL, version INTEGER NOT NULL, snapshot_json TEXT NOT NULL, hash TEXT NOT NULL, created_at TEXT NOT NULL, PRIMARY KEY(org_id, brand_id, entity, entity_id, version));`,
		`CREATE TABLE IF NOT EXISTS audit_events (event_id TEXT PRIMARY KEY, org_id TEXT NOT NULL, brand_id TEXT NOT NULL, actor TEXT NOT NULL, tool TEXT NOT NULL, action TEXT NOT NULL, entity TEXT NOT NULL, entity_id TEXT NOT NULL, version INTEGER NOT NULL, hash TEXT NOT NULL, reason TEXT, created_at TEXT NOT NULL);`,
	}
	for _, stmt := range stmts {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) UpsertBrand(ctx context.Context, brand model.Brand, actor, tool, reason string) (model.Brand, model.AuditEvent, error) {
	now := time.Now().UTC()
	current, err := s.GetBrand(ctx, brand.OrgID, brand.BrandID)
	if err == nil {
		brand.Version = current.Version + 1
		brand.CreatedAt = current.CreatedAt
	} else {
		brand.Version = 1
		brand.CreatedAt = now
	}
	brand.UpdatedAt = now
	data, hash, err := marshalHash(brand)
	if err != nil {
		return brand, model.AuditEvent{}, err
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO brands(org_id, brand_id, data_json, version, created_at, updated_at) VALUES(?,?,?,?,?,?) ON CONFLICT(org_id, brand_id) DO UPDATE SET data_json=excluded.data_json, version=excluded.version, updated_at=excluded.updated_at`, brand.OrgID, brand.BrandID, string(data), brand.Version, brand.CreatedAt.Format(time.RFC3339), brand.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return brand, model.AuditEvent{}, err
	}
	if err := s.insertVersion(ctx, brand.OrgID, brand.BrandID, "brand", brand.BrandID, brand.Version, data, hash, now); err != nil {
		return brand, model.AuditEvent{}, err
	}
	evt, err := s.insertAudit(ctx, brand.OrgID, brand.BrandID, actor, tool, "upsert", "brand", brand.BrandID, brand.Version, hash, reason, now)
	return brand, evt, err
}

func (s *Store) GetBrand(ctx context.Context, orgID, brandID string) (model.Brand, error) {
	var raw string
	if err := s.db.QueryRowContext(ctx, `SELECT data_json FROM brands WHERE org_id=? AND brand_id=?`, orgID, brandID).Scan(&raw); err != nil {
		return model.Brand{}, err
	}
	var b model.Brand
	return b, json.Unmarshal([]byte(raw), &b)
}

func (s *Store) UpsertVoice(ctx context.Context, voice model.BrandVoice, actor, tool, reason string) (model.BrandVoice, model.AuditEvent, error) {
	now := time.Now().UTC()
	current, err := s.GetVoice(ctx, voice.OrgID, voice.BrandID)
	if err == nil {
		voice.Version = current.Version + 1
	} else {
		voice.Version = 1
	}
	data, hash, err := marshalHash(voice)
	if err != nil {
		return voice, model.AuditEvent{}, err
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO brand_voice(org_id, brand_id, data_json, version, updated_at) VALUES(?,?,?,?,?) ON CONFLICT(org_id, brand_id) DO UPDATE SET data_json=excluded.data_json, version=excluded.version, updated_at=excluded.updated_at`, voice.OrgID, voice.BrandID, string(data), voice.Version, now.Format(time.RFC3339))
	if err != nil {
		return voice, model.AuditEvent{}, err
	}
	if err := s.insertVersion(ctx, voice.OrgID, voice.BrandID, "voice", voice.BrandID, voice.Version, data, hash, now); err != nil {
		return voice, model.AuditEvent{}, err
	}
	evt, err := s.insertAudit(ctx, voice.OrgID, voice.BrandID, actor, tool, "upsert", "voice", voice.BrandID, voice.Version, hash, reason, now)
	return voice, evt, err
}

func (s *Store) GetVoice(ctx context.Context, orgID, brandID string) (model.BrandVoice, error) {
	var raw string
	if err := s.db.QueryRowContext(ctx, `SELECT data_json FROM brand_voice WHERE org_id=? AND brand_id=?`, orgID, brandID).Scan(&raw); err != nil {
		return model.BrandVoice{}, err
	}
	var v model.BrandVoice
	return v, json.Unmarshal([]byte(raw), &v)
}

func (s *Store) StoreCampaign(ctx context.Context, c model.Campaign, actor, tool, reason string) (model.Campaign, model.AuditEvent, error) {
	now := time.Now().UTC()
	if c.CampaignID == "" {
		c.CampaignID = fmt.Sprintf("camp_%d", now.UnixNano())
	}
	c.CreatedAt = now
	data, hash, err := marshalHash(c)
	if err != nil {
		return c, model.AuditEvent{}, err
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO campaigns(org_id, brand_id, campaign_id, data_json, created_at) VALUES(?,?,?,?,?)`, c.OrgID, c.BrandID, c.CampaignID, string(data), now.Format(time.RFC3339))
	if err != nil {
		return c, model.AuditEvent{}, err
	}
	evt, err := s.insertAudit(ctx, c.OrgID, c.BrandID, actor, tool, "insert", "campaign", c.CampaignID, 1, hash, reason, now)
	return c, evt, err
}

func (s *Store) Campaigns(ctx context.Context, orgID, brandID string) ([]model.Campaign, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT data_json FROM campaigns WHERE org_id=? AND brand_id=? ORDER BY created_at DESC`, orgID, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Campaign{}
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		var c model.Campaign
		if err := json.Unmarshal([]byte(raw), &c); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Store) AuditLog(ctx context.Context, orgID, brandID string) ([]model.AuditEvent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT event_id, org_id, brand_id, actor, tool, action, entity, entity_id, version, hash, COALESCE(reason,''), created_at FROM audit_events WHERE org_id=? AND brand_id=? ORDER BY created_at DESC`, orgID, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.AuditEvent{}
	for rows.Next() {
		var e model.AuditEvent
		var created string
		if err := rows.Scan(&e.EventID, &e.OrgID, &e.BrandID, &e.Actor, &e.Tool, &e.Action, &e.Entity, &e.EntityID, &e.Version, &e.Hash, &e.Reason, &created); err != nil {
			return nil, err
		}
		e.CreatedAt, _ = time.Parse(time.RFC3339, created)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) insertVersion(ctx context.Context, orgID, brandID, entity, entityID string, version int, data []byte, hash string, now time.Time) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO versions(org_id, brand_id, entity, entity_id, version, snapshot_json, hash, created_at) VALUES(?,?,?,?,?,?,?,?)`, orgID, brandID, entity, entityID, version, string(data), hash, now.Format(time.RFC3339))
	return err
}

func (s *Store) insertAudit(ctx context.Context, orgID, brandID, actor, tool, action, entity, entityID string, version int, hash, reason string, now time.Time) (model.AuditEvent, error) {
	e := model.AuditEvent{EventID: fmt.Sprintf("evt_%d", now.UnixNano()), OrgID: orgID, BrandID: brandID, Actor: actor, Tool: tool, Action: action, Entity: entity, EntityID: entityID, Version: version, Hash: hash, Reason: reason, CreatedAt: now}
	_, err := s.db.ExecContext(ctx, `INSERT INTO audit_events(event_id, org_id, brand_id, actor, tool, action, entity, entity_id, version, hash, reason, created_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`, e.EventID, e.OrgID, e.BrandID, e.Actor, e.Tool, e.Action, e.Entity, e.EntityID, e.Version, e.Hash, e.Reason, e.CreatedAt.Format(time.RFC3339))
	return e, err
}

func marshalHash(v any) ([]byte, string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, "", err
	}
	sum := sha256.Sum256(data)
	return data, hex.EncodeToString(sum[:]), nil
}
