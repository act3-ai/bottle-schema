package main

import (
	"context"
	"dagger/schema/internal/dagger"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

// Run linters.
func (s *Schema) Lint() *Lint {
	return &Lint{
		Source: s.Source,
	}
}

// Lint organizes linting actions.
type Lint struct {
	Source *dagger.Directory
}

// Run all linters: Yaml, Markdown, Golang, and Shell.
func (l *Lint) All(ctx context.Context,
	// Source code directory
	// +defaultPath="/"
	src *dagger.Directory,
) (string, error) {
	p := pool.NewWithResults[string]().WithContext(ctx)

	p.Go(func(ctx context.Context) (string, error) {
		ctx, span := Tracer().Start(ctx, "yamllint")
		defer span.End()
		return l.Yamllint(ctx, l.Source)
	})

	p.Go(func(ctx context.Context) (string, error) {
		ctx, span := Tracer().Start(ctx, "markdownlint")
		defer span.End()
		return l.Markdownlint(ctx, l.Source)
	})

	p.Go(func(ctx context.Context) (string, error) {
		ctx, span := Tracer().Start(ctx, "golangci-lint")
		defer span.End()
		return l.Go(ctx)
	})

	p.Go(func(ctx context.Context) (string, error) {
		ctx, span := Tracer().Start(ctx, "shellcheck")
		defer span.End()
		return l.Shellcheck(ctx, src.Directory("bin"))
	})

	result, err := p.Wait()
	// TODO maybe we should order the lint result strings
	return strings.Join(result, "\n=====\n"), err
}

// Lint yaml files.
func (l *Lint) Yamllint(ctx context.Context,
	// Source code directory
	// +defaultPath="/"
	src *dagger.Directory,
) (string, error) {
	return dag.Container().
		From("docker.io/cytopia/yamllint:1").
		WithWorkdir("/src").
		WithDirectory("/src", src).
		WithExec([]string{"yamllint", "."}).
		Stdout(ctx)
}

// Lint markdown files.
func (l *Lint) Markdownlint(ctx context.Context,
	// source code directory
	// +defaultPath="/"
	src *dagger.Directory,
) (string, error) {
	return dag.Container().
		From("docker.io/davidanson/markdownlint-cli2:v0.14.0").
		WithWorkdir("/src").
		WithDirectory("/src", src).
		WithExec([]string{"markdownlint-cli2", "."}).
		Stdout(ctx)
}

// Lint shell files.
func (l *Lint) Shellcheck(ctx context.Context,
	// Source code directory
	// +defaultPath="/bin"
	src *dagger.Directory,
) (string, error) {
	filenames, err := src.Entries(ctx)
	if err != nil {
		return "", err
	}

	p := pool.New().WithContext(ctx).WithMaxGoroutines(4)
	for _, filename := range filenames {
		p.Go(func(ctx context.Context) error {
			return dag.Shellcheck().
				Check(src.File(filename)).
				Assert(ctx)
		})
	}

	err = p.Wait()
	if err != nil {
		return err.Error(), err
	}
	return "", nil
}

// Lint golang files.
func (l *Lint) Go(ctx context.Context) (string, error) {
	return dag.GolangciLint().
		Run(l.Source, dagger.GolangciLintRunOpts{
			Timeout: "5m",
		}).
		Stdout(ctx)
}
