package agent

import (
	"testing"

	kit "github.com/mark3labs/kit/pkg/kit"
)

// TestApplyProviderEnvOverrides_IteratrWinsOverKit verifies that
// ITERATR_PROVIDER_URL takes precedence over KIT_PROVIDER_URL.
func TestApplyProviderEnvOverrides_IteratrWinsOverKit(t *testing.T) {
	t.Setenv("ITERATR_PROVIDER_URL", "http://iteratr-host:1234/v1")
	t.Setenv("KIT_PROVIDER_URL", "http://kit-host:5678/v1")
	t.Setenv("ITERATR_PROVIDER_API_KEY", "iteratr-key")
	t.Setenv("KIT_PROVIDER_API_KEY", "kit-key")

	opts := &kit.Options{}
	applyProviderEnvOverrides(opts)

	if opts.ProviderURL != "http://iteratr-host:1234/v1" {
		t.Errorf("expected ITERATR_PROVIDER_URL to win; got ProviderURL=%q", opts.ProviderURL)
	}
	if opts.ProviderAPIKey != "iteratr-key" {
		t.Errorf("expected ITERATR_PROVIDER_API_KEY to win; got ProviderAPIKey=%q", opts.ProviderAPIKey)
	}
}

// TestApplyProviderEnvOverrides_KitFallback verifies that KIT_* env vars
// are used as a fallback when the ITERATR_* versions are unset.
func TestApplyProviderEnvOverrides_KitFallback(t *testing.T) {
	t.Setenv("KIT_PROVIDER_URL", "http://kit-host:5678/v1")
	t.Setenv("KIT_PROVIDER_API_KEY", "kit-key")

	opts := &kit.Options{}
	applyProviderEnvOverrides(opts)

	if opts.ProviderURL != "http://kit-host:5678/v1" {
		t.Errorf("expected KIT_PROVIDER_URL fallback; got ProviderURL=%q", opts.ProviderURL)
	}
	if opts.ProviderAPIKey != "kit-key" {
		t.Errorf("expected KIT_PROVIDER_API_KEY fallback; got ProviderAPIKey=%q", opts.ProviderAPIKey)
	}
}

// TestApplyProviderEnvOverrides_NoneSet verifies the helper is a no-op
// when no env vars are present (preserves any pre-existing values on opts).
func TestApplyProviderEnvOverrides_NoneSet(t *testing.T) {
	// Make sure none of the env vars leak in from the test process.
	t.Setenv("ITERATR_PROVIDER_URL", "")
	t.Setenv("KIT_PROVIDER_URL", "")
	t.Setenv("ITERATR_PROVIDER_API_KEY", "")
	t.Setenv("KIT_PROVIDER_API_KEY", "")

	opts := &kit.Options{
		ProviderURL:    "http://pre-existing:9999/v1",
		ProviderAPIKey: "pre-existing-key",
	}
	applyProviderEnvOverrides(opts)

	if opts.ProviderURL != "http://pre-existing:9999/v1" {
		t.Errorf("expected pre-existing URL preserved; got %q", opts.ProviderURL)
	}
	if opts.ProviderAPIKey != "pre-existing-key" {
		t.Errorf("expected pre-existing key preserved; got %q", opts.ProviderAPIKey)
	}
}

// TestApplyProviderEnvOverrides_OnlyURLSet verifies that URL and key are
// independent — setting only one of them should not affect the other.
func TestApplyProviderEnvOverrides_OnlyURLSet(t *testing.T) {
	t.Setenv("ITERATR_PROVIDER_URL", "http://iteratr-host:1234/v1")
	t.Setenv("KIT_PROVIDER_URL", "")
	t.Setenv("ITERATR_PROVIDER_API_KEY", "")
	t.Setenv("KIT_PROVIDER_API_KEY", "")

	opts := &kit.Options{}
	applyProviderEnvOverrides(opts)

	if opts.ProviderURL != "http://iteratr-host:1234/v1" {
		t.Errorf("expected URL to be set; got %q", opts.ProviderURL)
	}
	if opts.ProviderAPIKey != "" {
		t.Errorf("expected API key to be empty; got %q", opts.ProviderAPIKey)
	}
}
