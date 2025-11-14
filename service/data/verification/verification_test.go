package verification

import (
	"testing"
)

func TestIsValidDomain(t *testing.T) {
	tests := []struct {
		name   string
		domain string
		want   bool
	}{
		{"example.com", "example.com", true},
		{"sub.example.com", "sub.example.com", true},
		{"localhost", "localhost", true},
		{"invalid_domain", "invalid_domain", false},
		{"example-.com", "example-.com", false},
		{"-example.com", "-example.com", false},
		{"example.com.", "example.com.", false},
		{"example..com", "example..com", false},
		{"1example.com", "1example.com", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidDomain(tt.domain); got != tt.want {
				t.Errorf("IsValidDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidIPv4(t *testing.T) {

	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"intranet address", "192.168.1.1", true},
		{"maximum address ", "255.255.255.255", true},
		{"invalid address ", "256.256.256.256", false},
		{"loopback address", "127.0.0.1", true},
		{"intranet address", "10.0.0.1", true},
		{"intranet address", "172.16.0.1", true},
		{"intranet address", "192.168.0.256", false},
		{"intranet address", "0.0.0.256", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidIPv4(tt.ip); got != tt.want {
				t.Errorf("IsValidIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidIPv6(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"IPV6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", true},
		{"loopback", "::1", true},
		{"IPV6", "2001:db8::ff00:42:8329", true},
		{"IPV6", "2001:db8:85a3::g123", false},
		{"IPV6", "2001:db8:85a3:0000:0000:8a2e:0370:7334:1234", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidIPv6(tt.ip); got != tt.want {
				t.Errorf("IsValidIPv6() = %v, want %v", got, tt.want)
			}
		})
	}
}
