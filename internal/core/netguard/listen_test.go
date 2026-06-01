package netguard

import "testing"

func TestValidateLocalListen(t *testing.T) {
	for _, addr := range []string{"127.0.0.1:8621", "[::1]:8621", "localhost:8621"} {
		if err := ValidateLocalListen(addr, false); err != nil {
			t.Fatalf("expected %s to be accepted: %v", addr, err)
		}
	}
	for _, addr := range []string{"0.0.0.0:8621", "192.168.1.10:8621", ":8621"} {
		if err := ValidateLocalListen(addr, false); err == nil {
			t.Fatalf("expected %s to be rejected", addr)
		}
	}
	if err := ValidateLocalListen("0.0.0.0:8621", true); err != nil {
		t.Fatalf("explicit override should allow non-loopback: %v", err)
	}
}
