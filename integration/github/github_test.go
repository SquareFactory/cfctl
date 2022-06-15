package github_test

import (
	"testing"

	"github.com/SquareFactory/cfctl/integration/github"
)

func TestLatestRelease(t *testing.T) {
	_, err := github.LatestRelease(false)

	if err != nil {
		t.Error(err)
	}
}
