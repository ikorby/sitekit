package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if err := run(os.Args, "."); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func run(args []string, baseDir string) error {
	if len(args) < 3 || args[1] != "new" {
		return fmt.Errorf("usage: sk new <project-name>")
	}

	projectName := args[2]
	targetDir := filepath.Join(baseDir, projectName)

	dirs := []string{
		filepath.Join(targetDir, "cmd", "newsite"),
		filepath.Join(targetDir, "templates", "layouts"),
		filepath.Join(targetDir, "templates", "pages"),
		filepath.Join(targetDir, "static", "css"),
		filepath.Join(targetDir, "config"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	files := map[string]string{
		filepath.Join(targetDir, "cmd", "newsite", "main.go"):         mainGoTemplate(),
		filepath.Join(targetDir, "templates", "layouts", "base.html"): baseHTMLTemplate(),
		filepath.Join(targetDir, "templates", "pages", "home.html"):   homeHTMLTemplate(),
		filepath.Join(targetDir, ".env"):                              "SITEKIT_ENV=development\nSITEKIT_PORT=8080",
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", path, err)
		}
	}

	// 1. Инициализируем go.mod
	cmdInit := exec.Command("go", "mod", "init", projectName)
	cmdInit.Dir = targetDir
	if err := cmdInit.Run(); err != nil {
		return fmt.Errorf("failed to init go module: %w", err)
	}

	// 2. Скачиваем сам фреймворк
	cmdGet := exec.Command("go", "get", "github.com/ikorby/sitekit@latest")
	cmdGet.Dir = targetDir
	if err := cmdGet.Run(); err != nil {
		return fmt.Errorf("failed to get sitekit dependency: %w", err)
	}

	// 3. Подчищаем зависимости
	cmdTidy := exec.Command("go", "mod", "tidy")
	cmdTidy.Dir = targetDir
	if err := cmdTidy.Run(); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	fmt.Printf("✅ Sitekit project '%s' successfully created!\n", projectName)
	fmt.Printf("cd %s && go run ./cmd/newsite\n", projectName)
	return nil
}

func mainGoTemplate() string {
	return `package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/ikorby/sitekit/app"
	"github.com/ikorby/sitekit/config"
	apperrors "github.com/ikorby/sitekit/errors"
	"github.com/ikorby/sitekit/middleware"
	"github.com/ikorby/sitekit/page"
	"github.com/ikorby/sitekit/render"
	"github.com/ikorby/sitekit/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger := slog.Default()
	renderer := render.New(os.DirFS("templates"), cfg.IsDevelopment())
	
	// Теперь мы подключаем логгер, обработчик ошибок и все middleware фреймворка
	application := app.New(cfg, 
		app.WithRenderer(renderer),
		app.WithLogger(logger),
		app.WithErrorHandler(apperrors.Handler(logger)),
		app.WithMiddleware(
			middleware.Recovery(logger),
			middleware.Logger(logger),
			middleware.Security(cfg),
		),
	)
	
	r := router.New(application)
	r.Get("/", func(c *app.Context) error {
		p := page.New("home.html", map[string]string{"Title": "Welcome to Sitekit"})
		return c.Render(http.StatusOK, p)
	})

	if err := application.Run(); err != nil {
		slog.Error("server failed", "error", err)
	}
}
`
}

func baseHTMLTemplate() string {
	return `{{define "layout"}}<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Meta.Title}}</title>
</head>
<body>
    {{template "content" .}}
</body>
</html>{{end}}`
}

func homeHTMLTemplate() string {
	return `{{define "content"}}
    <h1>{{.Data.Title}}</h1>
    <p>Your site is ready.</p>
{{end}}`
}
