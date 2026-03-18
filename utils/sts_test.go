package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

func TestGenerateTokenUsesChinaTimezone(t *testing.T) {
	local := time.FixedZone("UTC-7", -7*3600)
	now := time.Date(2026, 3, 18, 20, 0, 0, 0, local)

	expected := generateTokenAt("ak", "sk", now.In(chinaTimezone))
	got := generateTokenAt("ak", "sk", now)

	if got != expected {
		t.Fatalf("expected token %s, got %s", expected, got)
	}
}

func TestGenerateTokenConvertsTimezoneBeforeFormatting(t *testing.T) {
	local := time.FixedZone("UTC-7", -7*3600)
	now := time.Date(2026, 3, 18, 20, 0, 0, 0, local)

	localTokenText := fmt.Sprintf(
		"%s_%s_%s_%d",
		now.Format("2006-1-2"),
		"ak",
		"sk",
		now.Hour(),
	)
	localSum := md5.Sum([]byte(localTokenText))
	localToken := hex.EncodeToString(localSum[:])

	chinaToken := generateTokenAt("ak", "sk", now)

	if chinaToken == localToken {
		t.Fatalf("expected token to be derived from UTC+8 time, but token matched local time")
	}
}
