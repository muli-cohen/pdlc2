package subnet_test

import (
	"testing"

	"github.com/muli-cohen/pdlc2/subnet"
)

func TestParse(t *testing.T) {
	t.Run("typical /24", func(t *testing.T) {
		n, err := subnet.Parse("192.168.1.10/24")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Address != "192.168.1.10" {
			t.Errorf("Address: got %q, want %q", n.Address, "192.168.1.10")
		}
		if n.NetworkAddress != "192.168.1.0" {
			t.Errorf("NetworkAddress: got %q, want %q", n.NetworkAddress, "192.168.1.0")
		}
		if n.Broadcast != "192.168.1.255" {
			t.Errorf("Broadcast: got %q, want %q", n.Broadcast, "192.168.1.255")
		}
		if n.Mask != "255.255.255.0" {
			t.Errorf("Mask: got %q, want %q", n.Mask, "255.255.255.0")
		}
		if n.FirstHost != "192.168.1.1" {
			t.Errorf("FirstHost: got %q, want %q", n.FirstHost, "192.168.1.1")
		}
		if n.LastHost != "192.168.1.254" {
			t.Errorf("LastHost: got %q, want %q", n.LastHost, "192.168.1.254")
		}
		if n.PrefixLength != 24 {
			t.Errorf("PrefixLength: got %d, want 24", n.PrefixLength)
		}
		if n.UsableHosts != 254 {
			t.Errorf("UsableHosts: got %d, want 254", n.UsableHosts)
		}
	})

	t.Run("/0 network", func(t *testing.T) {
		n, err := subnet.Parse("0.0.0.0/0")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.NetworkAddress != "0.0.0.0" {
			t.Errorf("NetworkAddress: got %q, want %q", n.NetworkAddress, "0.0.0.0")
		}
		if n.Broadcast != "255.255.255.255" {
			t.Errorf("Broadcast: got %q, want %q", n.Broadcast, "255.255.255.255")
		}
		if n.UsableHosts != 4294967294 {
			t.Errorf("UsableHosts: got %d, want 4294967294", n.UsableHosts)
		}
	})

	t.Run("/31 boundary", func(t *testing.T) {
		n, err := subnet.Parse("10.0.0.1/31")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Broadcast != "" {
			t.Errorf("Broadcast: got %q, want empty", n.Broadcast)
		}
		if n.FirstHost != "" {
			t.Errorf("FirstHost: got %q, want empty", n.FirstHost)
		}
		if n.LastHost != "" {
			t.Errorf("LastHost: got %q, want empty", n.LastHost)
		}
		if n.UsableHosts != 0 {
			t.Errorf("UsableHosts: got %d, want 0", n.UsableHosts)
		}
	})

	t.Run("/32 boundary", func(t *testing.T) {
		n, err := subnet.Parse("10.0.0.1/32")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Broadcast != "" {
			t.Errorf("Broadcast: got %q, want empty", n.Broadcast)
		}
		if n.FirstHost != "" {
			t.Errorf("FirstHost: got %q, want empty", n.FirstHost)
		}
		if n.LastHost != "" {
			t.Errorf("LastHost: got %q, want empty", n.LastHost)
		}
		if n.UsableHosts != 0 {
			t.Errorf("UsableHosts: got %d, want 0", n.UsableHosts)
		}
	})

	t.Run("host-bit masking", func(t *testing.T) {
		n1, err := subnet.Parse("192.168.1.10/24")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		n2, err := subnet.Parse("192.168.1.200/24")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n1.NetworkAddress != n2.NetworkAddress {
			t.Errorf("NetworkAddress mismatch: %q vs %q", n1.NetworkAddress, n2.NetworkAddress)
		}
	})

	t.Run("/8 usable hosts", func(t *testing.T) {
		n, err := subnet.Parse("10.0.0.0/8")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.UsableHosts != 16777214 {
			t.Errorf("UsableHosts: got %d, want 16777214", n.UsableHosts)
		}
	})

	t.Run("rejection: missing prefix", func(t *testing.T) {
		_, err := subnet.Parse("192.168.1.1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("rejection: prefix /33", func(t *testing.T) {
		_, err := subnet.Parse("192.168.1.1/33")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("rejection: octet 300", func(t *testing.T) {
		_, err := subnet.Parse("300.1.1.1/24")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("rejection: wrong octet count", func(t *testing.T) {
		_, err := subnet.Parse("192.168.1/24")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("rejection: non-numeric octet", func(t *testing.T) {
		_, err := subnet.Parse("192.168.abc.1/24")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
