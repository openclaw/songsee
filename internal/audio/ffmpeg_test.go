package audio

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestResolveFFmpeg(t *testing.T) {
	ffmpegPath := installFakeFFmpeg(t)
	t.Setenv("PATH", filepath.Dir(ffmpegPath)+string(os.PathListSeparator)+os.Getenv("PATH"))

	path, err := resolveFFmpeg("")
	if err != nil {
		t.Fatalf("resolveFFmpeg: %v", err)
	}
	if path != ffmpegPath {
		t.Fatalf("expected %s, got %s", ffmpegPath, path)
	}
}

func TestResolveFFmpegExplicit(t *testing.T) {
	ffmpegPath := installFakeFFmpeg(t)
	path, err := resolveFFmpeg(ffmpegPath)
	if err != nil {
		t.Fatalf("resolveFFmpeg explicit: %v", err)
	}
	if path != ffmpegPath {
		t.Fatalf("expected %s, got %s", ffmpegPath, path)
	}
}

func TestResolveFFmpegMissing(t *testing.T) {
	t.Setenv("PATH", "")
	if _, err := resolveFFmpeg(""); err == nil {
		t.Fatalf("expected error")
	}
}

func TestDecodeWithFFmpegFile(t *testing.T) {
	ffmpegPath := installFakeFFmpeg(t)
	input := filepath.Join(t.TempDir(), "input.bin")
	if err := os.WriteFile(input, []byte("audio"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	pcm, err := DecodeWithFFmpeg(input, nil, 22050, ffmpegPath, 0)
	if err != nil {
		t.Fatalf("DecodeWithFFmpeg: %v", err)
	}
	if pcm.SampleRate != 22050 {
		t.Fatalf("sample rate = %d", pcm.SampleRate)
	}
	if len(pcm.Samples) == 0 {
		t.Fatalf("empty samples")
	}
}

func TestDecodeWithFFmpegStdin(t *testing.T) {
	ffmpegPath := installFakeFFmpeg(t)
	pcm, err := DecodeWithFFmpeg("", bytes.NewReader([]byte("audio")), 44100, ffmpegPath, 0)
	if err != nil {
		t.Fatalf("DecodeWithFFmpeg stdin: %v", err)
	}
	if len(pcm.Samples) == 0 {
		t.Fatalf("empty samples")
	}
}

func TestDecodeWithFFmpegBadPath(t *testing.T) {
	_, err := DecodeWithFFmpeg("missing.mp3", nil, 0, "/no/such/ffmpeg", 0)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestDecodeFileHonorsOptionsFFmpegTimeout(t *testing.T) {
	ffmpegPath := installHangingFFmpeg(t)
	input := filepath.Join(t.TempDir(), "input.opus")
	if err := os.WriteFile(input, []byte("not-wav-or-mp3"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	start := time.Now()
	_, err := DecodeFile(input, Options{
		FFmpegPath:    ffmpegPath,
		FFmpegTimeout: 200 * time.Millisecond,
	})
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if elapsed > 3*time.Second {
		t.Fatalf("timeout took too long: %v", elapsed)
	}
	msg := strings.ToLower(err.Error())
	if !strings.Contains(msg, "timed out") && !strings.Contains(msg, "deadline") {
		t.Fatalf("want timeout/deadline error, got %v", err)
	}
}

func TestDecodeWithFFmpegRejectsNegativeTimeout(t *testing.T) {
	_, err := DecodeWithFFmpeg("input.bin", nil, 44100, "/bin/true", -time.Second)
	if err == nil {
		t.Fatal("expected error for negative timeout")
	}
	if !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("want timeout validation error, got %v", err)
	}
}

func TestDecodeWithFFmpegTimeout(t *testing.T) {
	ffmpegPath := installHangingFFmpeg(t)
	input := filepath.Join(t.TempDir(), "input.bin")
	if err := os.WriteFile(input, []byte("audio"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	start := time.Now()
	_, err := DecodeWithFFmpeg(input, nil, 44100, ffmpegPath, 200*time.Millisecond)
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if elapsed > 3*time.Second {
		t.Fatalf("timeout took too long: %v", elapsed)
	}
	msg := strings.ToLower(err.Error())
	if !strings.Contains(msg, "timed out") && !strings.Contains(msg, "deadline") {
		t.Fatalf("want timeout/deadline error, got %v", err)
	}
}

func TestDecodeWithFFmpegTimeoutWithInheritedPipes(t *testing.T) {
	dir := t.TempDir()
	ffmpegPath := filepath.Join(dir, "ffmpeg")
	pidPath := filepath.Join(dir, "child.pid")
	script := "#!/bin/sh\n/bin/sleep 30 &\necho $! > \"" + pidPath + "\"\nwait\n"
	if err := os.WriteFile(ffmpegPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		data, err := os.ReadFile(pidPath)
		if err != nil {
			return
		}
		pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err == nil && pid > 0 {
			if process, err := os.FindProcess(pid); err == nil {
				_ = process.Kill()
			}
		}
	})

	// Allow the wrapper to start on a busy host; the child outlives this bound.
	const timeout = 5 * time.Second
	start := time.Now()
	_, err := DecodeWithFFmpeg("input.bin", nil, 44100, ffmpegPath, timeout)
	if elapsed := time.Since(start); elapsed > timeout+3*time.Second {
		t.Fatalf("inherited pipes delayed timeout: %v", elapsed)
	}
	if _, statErr := os.Stat(pidPath); statErr != nil {
		t.Fatalf("wrapper did not start its child: %v", statErr)
	}
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("expected timeout, got %v", err)
	}
}

func installFakeFFmpeg(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "ffmpeg")
	script := "#!/bin/sh\nprintf '\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x3f'\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}

func installHangingFFmpeg(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "ffmpeg")
	script := "#!/bin/sh\nexec /bin/sleep 60\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}
