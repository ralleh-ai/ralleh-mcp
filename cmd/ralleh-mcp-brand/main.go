package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ralleh-ai/ralleh-mcp/internal/brand/model"
	"github.com/ralleh-ai/ralleh-mcp/internal/brand/service"
	"github.com/ralleh-ai/ralleh-mcp/internal/brand/store"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/health"
	"github.com/ralleh-ai/ralleh-mcp/internal/core/source"
)

func capabilities() map[string]bool {
	return map[string]bool{"canStoreBrandMemory": true, "canValidateContent": true, "canAuditWrites": true, "canManageMultipleBrands": true}
}

func main() {
	dbPath := flag.String("db", "/tmp/ralleh-mcp-brand.db", "SQLite database path")
	healthOnly := flag.Bool("health", false, "print health JSON and exit")
	healthServer := flag.Bool("health-server", false, "serve local-only HTTP health endpoints")
	healthListen := flag.String("health-listen", "127.0.0.1:8625", "health server listen address")
	allowNonLoopback := flag.Bool("allow-non-loopback-health", false, "allow health server to bind outside loopback; requires external firewall/auth controls")
	createBrand := flag.Bool("create-brand", false, "create/update a brand profile")
	updateVoice := flag.Bool("update-voice", false, "create/update brand voice")
	getProfile := flag.Bool("get-profile", false, "get brand profile")
	getVoice := flag.Bool("get-voice", false, "get brand voice")
	validate := flag.Bool("validate-content", false, "validate content against brand")
	audit := flag.Bool("audit-log", false, "print brand audit log")
	orgID := flag.String("org", "org_default", "organization ID")
	brandID := flag.String("brand", "brand_default", "brand ID")
	name := flag.String("name", "", "brand name")
	description := flag.String("description", "", "brand description")
	mission := flag.String("mission", "", "brand mission")
	voiceTone := flag.String("tone", "", "comma-separated tone values")
	forbidden := flag.String("forbidden", "", "comma-separated forbidden terms")
	preferred := flag.String("preferred", "", "comma-separated preferred phrases")
	content := flag.String("content", "", "content to validate")
	rewrite := flag.Bool("rewrite", false, "return rewritten content when validating")
	flag.Parse()

	st, err := store.Open(*dbPath)
	if err != nil {
		fatal(err)
	}
	defer st.Close()
	healthStatus := health.Evaluate("ralleh-mcp-brand", source.Registry{Sources: map[string]source.Source{"sqlite": {ID: "sqlite", Name: "SQLite Brand Store"}}, Collections: map[string]source.Collection{"brand_memory": {ID: "brand_memory", DefaultSources: []string{"sqlite"}}}}, capabilities())
	if *healthOnly {
		writeJSON(healthStatus)
		if !healthStatus.Ready {
			os.Exit(1)
		}
		return
	}
	if *healthServer {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		if err := health.ServeLocal(ctx, *healthListen, *allowNonLoopback, healthStatus); err != nil && ctx.Err() == nil {
			fatal(err)
		}
		return
	}

	ctx := context.Background()
	if *createBrand {
		brand, evt, err := st.UpsertBrand(ctx, model.Brand{OrgID: *orgID, BrandID: *brandID, Name: *name, Description: *description, Mission: *mission}, "cli", "brand.create_brand", "cli upsert")
		if err != nil {
			fatal(err)
		}
		writeJSON(map[string]any{"brand": brand, "auditEvent": evt})
		return
	}
	if *updateVoice {
		voice, evt, err := st.UpsertVoice(ctx, model.BrandVoice{OrgID: *orgID, BrandID: *brandID, Tone: splitCSV(*voiceTone), ForbiddenTerms: splitCSV(*forbidden), PreferredPhrases: splitCSV(*preferred)}, "cli", "brand.update_voice", "cli upsert")
		if err != nil {
			fatal(err)
		}
		writeJSON(map[string]any{"voice": voice, "auditEvent": evt})
		return
	}
	if *getProfile {
		b, err := st.GetBrand(ctx, *orgID, *brandID)
		if err != nil {
			fatal(err)
		}
		writeJSON(b)
		return
	}
	if *getVoice {
		v, err := st.GetVoice(ctx, *orgID, *brandID)
		if err != nil {
			fatal(err)
		}
		writeJSON(v)
		return
	}
	if *validate {
		res, err := (service.Service{Store: st}).ValidateContent(ctx, model.ValidationRequest{OrgID: *orgID, BrandID: *brandID, Content: *content, Rewrite: *rewrite})
		if err != nil {
			fatal(err)
		}
		writeJSON(res)
		return
	}
	if *audit {
		events, err := st.AuditLog(ctx, *orgID, *brandID)
		if err != nil {
			fatal(err)
		}
		writeJSON(events)
		return
	}
	writeJSON(map[string]any{"service": "ralleh-mcp-brand", "status": "scaffold_ready", "capabilities": capabilities()})
}

func splitCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
func writeJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fatal(err)
	}
}
func fatal(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
