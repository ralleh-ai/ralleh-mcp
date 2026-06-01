package health

import (
	"context"
	"net/http"

	"github.com/ralleh-ai/ralleh-mcp/internal/core/netguard"
)

// ServeLocal starts the health server on a loopback address by default.
func ServeLocal(ctx context.Context, addr string, allowNonLoopback bool, status Status) error {
	if err := netguard.ValidateLocalListen(addr, allowNonLoopback); err != nil {
		return err
	}
	server := &http.Server{Addr: addr, Handler: Handler(status)}
	errCh := make(chan error, 1)
	go func() { errCh <- server.ListenAndServe() }()
	select {
	case <-ctx.Done():
		_ = server.Shutdown(context.Background())
		return ctx.Err()
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}
