package dir

import (
	"os"
	"testing"
	"time"
)

func Test_verifyPath(t *testing.T) {
	root := buildTestPath(t)
	path := buildSubdir(t, root)

	tests := []struct {
		name        string
		path        string
		trustedRoot string
		want        string
		wantErr     bool
	}{
		{
			name:        "Regular path",
			path:        path,
			trustedRoot: root,
			want:        path,
			wantErr:     false,
		},
		{
			name:        "Traversal denied",
			path:        "/../../../../../../etc/passwd",
			trustedRoot: root,
			want:        doesNotExist,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := verifyPath(tt.path, tt.trustedRoot)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("verifyPath() = %v, want %v", got, tt.want)
			}
		})
	}

	//sleep for 10 seconds to allow the cleanup to happen
	time.Sleep(10 * time.Second)
}

func buildTestPath(t *testing.T) string {
	path, _ := os.MkdirTemp("/tmp", "test*")

	actual, err := absoluteCanonicalPath(path)
	if err != nil {
		t.Fatalf("absoluteCanonicalPath() error = %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(actual)
	})

	return actual
}
func buildSubdir(t *testing.T, root string) string {
	path, _ := os.MkdirTemp(root, "test*")
	return path
}
