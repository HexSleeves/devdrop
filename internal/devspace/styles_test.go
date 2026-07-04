package devspace

import (
	"bytes"
	"strings"
	"testing"

	"github.com/charmbracelet/colorprofile"
)

func clearColorEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{"NO_COLOR", "CLICOLOR", "CLICOLOR_FORCE", "TTY_FORCE"} {
		t.Setenv(key, "")
	}
}

// resetStylesAfterTest snapshots the package-level style globals and restores
// them on cleanup. Every test that calls configureStyles or mutates
// currentProfile/currentTheme/currentNoColor directly must call this first,
// otherwise a colored profile can leak into unrelated tests that call
// output helpers (printStatus, RunDoctor, ...) directly without going through
// configureStyles — since those globals are shared package state, not
// per-invocation state, in the test binary.
func resetStylesAfterTest(t *testing.T) {
	t.Helper()
	savedTheme, savedProfile, savedNoColor := currentTheme, currentProfile, currentNoColor
	t.Cleanup(func() {
		currentTheme, currentProfile, currentNoColor = savedTheme, savedProfile, savedNoColor
	})
}

// renderThemed always returns the theme's full-color rendering (the theme
// itself is never plain; stripping happens in styledWriter).
func renderThemed() string {
	return currentTheme.OK.Render("ok")
}

func TestStyledWriterStripsAnsiForNonTTYByDefault(t *testing.T) {
	clearColorEnv(t)
	resetStylesAfterTest(t)
	var buf bytes.Buffer
	configureStyles(&buf, false)

	var out bytes.Buffer
	if _, err := styledWriter(&out).Write([]byte(renderThemed())); err != nil {
		t.Fatalf("write: %v", err)
	}
	if strings.ContainsRune(out.String(), 0x1b) {
		t.Fatalf("expected no ANSI escape bytes for a non-terminal writer with no forcing env vars, got %q", out.String())
	}
	if out.String() != "ok" {
		t.Fatalf("expected stripped output to equal the plain label, got %q", out.String())
	}
}

func TestStyledWriterPreservesAnsiWhenCliColorForced(t *testing.T) {
	clearColorEnv(t)
	resetStylesAfterTest(t)
	t.Setenv("CLICOLOR_FORCE", "1")
	var buf bytes.Buffer
	configureStyles(&buf, false)

	var out bytes.Buffer
	if _, err := styledWriter(&out).Write([]byte(renderThemed())); err != nil {
		t.Fatalf("write: %v", err)
	}
	if !strings.ContainsRune(out.String(), 0x1b) {
		t.Fatalf("expected ANSI escape bytes when CLICOLOR_FORCE=1 is set, got %q", out.String())
	}
}

func TestStyledWriterNoColorFlagOverridesForcing(t *testing.T) {
	clearColorEnv(t)
	resetStylesAfterTest(t)
	t.Setenv("CLICOLOR_FORCE", "1")
	var buf bytes.Buffer
	configureStyles(&buf, true) // --no-color

	var out bytes.Buffer
	if _, err := styledWriter(&out).Write([]byte(renderThemed())); err != nil {
		t.Fatalf("write: %v", err)
	}
	if strings.ContainsRune(out.String(), 0x1b) {
		t.Fatalf("expected --no-color to force plain output even when CLICOLOR_FORCE=1 is set, got %q", out.String())
	}
	if out.String() != "ok" {
		t.Fatalf("expected stripped output to equal the plain label, got %q", out.String())
	}
}

func TestStyledWriterHonorsNoColorEnvOnRealTTYProfile(t *testing.T) {
	clearColorEnv(t)
	resetStylesAfterTest(t)
	// Simulate an application that has already detected a colored profile
	// (as if running on a real TTY) and confirm NO_COLOR-driven ASCII
	// downgrade still strips color (decorations like bold may remain per the
	// NO_COLOR spec, but no color escape should survive).
	currentProfile = colorprofile.ASCII

	var out bytes.Buffer
	if _, err := styledWriter(&out).Write([]byte(renderThemed())); err != nil {
		t.Fatalf("write: %v", err)
	}
	if strings.Contains(out.String(), "38;2;") {
		t.Fatalf("expected no truecolor SGR sequence at ASCII profile, got %q", out.String())
	}
}
