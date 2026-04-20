package ratelimit

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Cooldown != 5*time.Minute {
		t.Errorf("expected cooldown 5m, got %v", cfg.Cooldown)
	}
	if cfg.MaxPerHour != 12 {
		t.Errorf("expected MaxPerHour 12, got %d", cfg.MaxPerHour)
	}
}

func TestValidate_NegativeCooldown_UsesDefault(t *testing.T) {
	cfg := Config{Cooldown: -1 * time.Second, MaxPerHour: 5}
	_ = cfg.Validate()
	if cfg.Cooldown != DefaultConfig().Cooldown {
		t.Errorf("expected default cooldown after negative input, got %v", cfg.Cooldown)
	}
}

func TestValidate_NegativeMaxPerHour_SetsZero(t *testing.T) {
	cfg := Config{Cooldown: time.Minute, MaxPerHour: -3}
	_ = cfg.Validate()
	if cfg.MaxPerHour != 0 {
		t.Errorf("expected MaxPerHour 0 after negative input, got %d", cfg.MaxPerHour)
	}
}

func TestValidate_ValidConfig_Unchanged(t *testing.T) {
	cfg := Config{Cooldown: 2 * time.Minute, MaxPerHour: 6}
	_ = cfg.Validate()
	if cfg.Cooldown != 2*time.Minute {
		t.Errorf("expected cooldown unchanged, got %v", cfg.Cooldown)
	}
	if cfg.MaxPerHour != 6 {
		t.Errorf("expected MaxPerHour unchanged, got %d", cfg.MaxPerHour)
	}
}

func TestNewFromConfig_CreatesFunctionalLimiter(t *testing.T) {
	cfg := Config{Cooldown: 10 * time.Minute, MaxPerHour: 6}
	rl := NewFromConfig(cfg)
	if rl == nil {
		t.Fatal("expected non-nil RateLimiter")
	}
	if !rl.Allow("host1") {
		t.Error("expected first call to Allow to return true")
	}
	if rl.Allow("host1") {
		t.Error("expected second call within cooldown to return false")
	}
}

func TestNewFromConfig_ZeroCooldown_UsesDefault(t *testing.T) {
	cfg := Config{Cooldown: 0, MaxPerHour: 0}
	rl := NewFromConfig(cfg)
	if rl == nil {
		t.Fatal("expected non-nil RateLimiter")
	}
	if !rl.Allow("host2") {
		t.Error("expected first call to Allow to return true")
	}
}
