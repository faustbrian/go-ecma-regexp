package ecmascript_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Node and Deno are two command names, but only one independent engine family.
func TestReleaseDifferentialRejectsTwoV8Runtimes(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node is unavailable for the dual-V8 regression fixture")
	}
	node, err = filepath.Abs(node)
	if err != nil {
		t.Fatal(err)
	}
	tools := t.TempDir()
	if err := os.Symlink(node, filepath.Join(tools, "node")); err != nil {
		t.Fatal(err)
	}
	// Preserve Deno's eval invocation while executing the same V8 implementation.
	wrapper := "#!/bin/sh\nshift\nexec '" + strings.ReplaceAll(node, "'", "'\"'\"'") + "' -e \"$@\"\n"
	if err := os.WriteFile(filepath.Join(tools, "deno"), []byte(wrapper), 0700); err != nil {
		t.Fatal(err)
	}
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, binary, "-test.run=^TestDifferentialMatchingAgainstJavaScriptEngines$", "-test.count=1")
	// The parent owns the mutation-profile lock; the child only tests the
	// release-family guard and must not reacquire that same lock.
	for _, variable := range os.Environ() {
		if !strings.HasPrefix(variable, "GOLIB_GREMLINS_COVERAGE_PROFILE=") {
			command.Env = append(command.Env, variable)
		}
	}
	command.Env = append(command.Env, "PATH="+tools, "ECMA_RELEASE_DIFFERENTIAL=1")
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatal("release differential accepted Node and Deno without an independent JavaScriptCore engine")
	}
	if !strings.Contains(string(output), "requires V8 and JavaScriptCore") {
		t.Fatalf("release differential failed for the wrong reason: %v\n%s", err, output)
	}
}
