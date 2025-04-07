package main

import (
	"context"
	"dagger/schema/internal/dagger"
	"fmt"
	"path/filepath"
	"strings"
)

type Release struct {
	// source code directory
	// +defaultPath="/"
	Source *dagger.Directory

	// GitLab token
	// +optional
	Token *dagger.Secret
}

func (m *Release) gitCliffContainer() *dagger.Container {
	return dag.Container().
		From("docker.io/orhunp/git-cliff:2.8.0").
		With(func(r *dagger.Container) *dagger.Container {
			if m.Token != nil {
				// TODO this is specific to ACT3, also does not work yet in git-cliff
				return r.WithSecretVariable("GITHUB_TOKEN", m.Token).
					WithEnvVariable("GITHUB_REPO", gitRepo)
			}
			return r
		}).
		WithMountedDirectory("/app", m.Source)
}

// Generate the change log from conventional commit messages (see cliff.toml)
func (r *Release) Changelog(ctx context.Context) *dagger.File {
	const changelogPath = "/app/CHANGELOG.md"
	return r.gitCliffContainer().
		WithExec([]string{"git-cliff", "--bump", "--strip=footer", "-o", changelogPath}).
		File(changelogPath)
}

// Generate the next version from conventional commit messages (see cliff.toml)
func (r *Release) Version(ctx context.Context) (string, error) {
	version, err := r.gitCliffContainer().
		WithExec([]string{"git-cliff", "--bumped-version"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(version)[1:], err
}

// Generate the initial release notes
func (r *Release) Notes(ctx context.Context) (string, error) {
	notes, err := r.gitCliffContainer().
		WithExec([]string{"git-cliff", "--bump", "--unreleased", "--strip=all"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	return notes, nil
}

// generate the files to prepare for a release
func (s *Release) Prepare(ctx context.Context) (*dagger.Directory, error) {
	changelog := s.Changelog(ctx)
	version, err := s.Version(ctx)
	if err != nil {
		return nil, err
	}

	notes, err := s.Notes(ctx)
	if err != nil {
		return nil, err
	}

	return dag.Directory().
			WithFile("CHANGELOG.md", changelog).
			WithNewFile("VERSION", version+"\n").
			WithNewFile(fmt.Sprintf("releases/v%s.md", version), notes),
		nil
}

// Publish the current release.  This should be tagged.
func (s *Schema) Publish(ctx context.Context,
	// source code directory
	// +defaultPath="/"
	src *dagger.Directory,

	// personal access token
	token *dagger.Secret,
) error {
	version, err := src.File("VERSION").Contents(ctx)
	if err != nil {
		return err
	}
	version = strings.TrimSpace(version)
	vVersion := "v" + version

	title := fmt.Sprintf("Release %s", vVersion)
	notes := src.File(filepath.Join("releases", vVersion+".md"))
	return dag.Gh(
		dagger.GhOpts{
			Token:  token,
			Repo:   gitRepo,
			Source: src,
		}).
		Release().
		Create(ctx, vVersion, title,
			dagger.GhReleaseCreateOpts{
				NotesFile: notes,
			})
}
