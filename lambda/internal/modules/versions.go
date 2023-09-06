package modules

import (
	"context"
	"fmt"
	"github.com/opentffoundation/registry/internal/github"
	"github.com/shurcooL/githubv4"

	"strings"
)

// TODO: doc
func GetVersions(ctx context.Context, ghClient *githubv4.Client, namespace string, name string, system string) ([]Version, error) {
	// the repo name should match the format `terraform-<system>-<name>`
	repoName := fmt.Sprintf("terraform-%s-%s", system, name)

	releases, err := github.FetchReleases(ctx, ghClient, namespace, repoName)
	if err != nil {
		return nil, err
	}

	var versions []Version
	for _, release := range releases {
		// Normalize the version name.
		versionName := release.TagName
		if strings.HasPrefix(versionName, "v") {
			versionName = versionName[1:]
		}

		// Construct the Version struct.
		version := Version{
			Version: versionName,
		}

		versions = append(versions, version)
	}
	return versions, nil
}

func GetVersionDownloadUrl(_ context.Context, namespace string, name string, system string, version string) string {
	// the repo name should match the format `terraform-<system>-<name>`
	repoName := fmt.Sprintf("terraform-%s-%s", system, name)

	// nothing fancy, just return the url for the release
	return fmt.Sprintf("git::https://github.com/%s/%s?ref=v%s", namespace, repoName, version)
}
