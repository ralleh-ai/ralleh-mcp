package mcpstdio

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const ProtocolVersion = "2024-11-05"

type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
	Handler     Handler        `json:"-"`
}

type Handler func(context.Context, map[string]any) (any, error)

type request struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params,omitempty"`
}

type response struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Result  any              `json:"result,omitempty"`
	Error   *rpcError        `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Serve(ctx context.Context, in io.Reader, out io.Writer, serverName, version string, tools []Tool) error {
	toolMap := make(map[string]Tool, len(tools))
	for _, tool := range tools {
		if strings.TrimSpace(tool.Name) == "" || tool.Handler == nil {
			return fmt.Errorf("invalid MCP tool %q", tool.Name)
		}
		toolMap[tool.Name] = tool
	}

	enc := json.NewEncoder(out)
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var req request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			continue
		}
		if req.ID == nil || strings.HasPrefix(req.Method, "notifications/") {
			continue
		}
		res := handle(ctx, req, serverName, version, tools, toolMap)
		if err := enc.Encode(res); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func handle(ctx context.Context, req request, serverName, version string, tools []Tool, toolMap map[string]Tool) response {
	base := response{JSONRPC: "2.0", ID: req.ID}
	switch req.Method {
	case "initialize":
		base.Result = map[string]any{
			"protocolVersion": ProtocolVersion,
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]any{"name": serverName, "version": version},
		}
	case "ping":
		base.Result = map[string]any{}
	case "tools/list":
		list := make([]Tool, 0, len(tools))
		for _, tool := range tools {
			list = append(list, Tool{Name: tool.Name, Description: tool.Description, InputSchema: tool.InputSchema})
		}
		base.Result = map[string]any{"tools": list}
	case "tools/call":
		var params struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			base.Error = &rpcError{Code: -32602, Message: "invalid tools/call params"}
			return base
		}
		tool, ok := toolMap[params.Name]
		if !ok {
			base.Error = &rpcError{Code: -32602, Message: "unknown tool: " + params.Name}
			return base
		}
		payload, err := tool.Handler(ctx, params.Arguments)
		if err != nil {
			base.Result = map[string]any{"isError": true, "content": []map[string]any{{"type": "text", "text": err.Error()}}}
			return base
		}
		text, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			base.Error = &rpcError{Code: -32603, Message: err.Error()}
			return base
		}
		base.Result = map[string]any{"content": []map[string]any{{"type": "text", "text": string(text)}}}
	default:
		base.Error = &rpcError{Code: -32601, Message: "method not found: " + req.Method}
	}
	return base
}
