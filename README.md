# Sitekit
![tests](https://github.com/ikorby/sitekit/actions/workflows/test.yml/badge.svg)

> An opinionated server-side rendering framework for Go.

Building traditional web applications in Go often means assembling the same pieces over and over again: routing, templates, rendering, middleware, configuration, static assets, SEO, and error handling. Sitekit provides those pieces with a consistent structure while staying close to the tools Go already gives you.

---

## Features

* **CLI Tool (`sk`)** for instant project scaffolding
* **Type-safe Pages** using Go generics
* **Live Reload** out-of-the-box for rapid development
* Standard library first 
* Server-side rendering with `html/template` 
* Routing built on `http.ServeMux` 
* HTTP and handler middleware 
* Environment-based configuration 
* Static file serving & SEO helpers 
* Structured HTTP errors 

---

## Installation

**1. Install the Sitekit CLI:**
```bash
go install [github.com/ikorby/sitekit/cmd/sk@latest](https://github.com/ikorby/sitekit/cmd/sk@latest)
```

**2. Or install the framework directly into an existing project:**
```bash
go get [github.com/ikorby/sitekit](https://github.com/ikorby/sitekit)
```

---

## Quick Start

The fastest way to start is using the CLI. It generates a ready-to-use project structure:

```bash
# Create a new project
sk new my-awesome-site

# Navigate to the project and run it
cd my-awesome-site
go run ./cmd/newsite
```
Your site will be running at `http://localhost:8080`.
* Live Reload is enabled automatically when `SITEKIT_ENV=development`. Just edit your HTML templates and save - the browser will refresh automatically.*

---

## Type-safe Pages

Sitekit uses generics to ensure type safety between your Go handlers and HTML templates:

```go
// 1. Define your data structure
type HomePageData struct {
    WelcomeMessage string
}

// 2. Pass it to the page securely
r.Get("/", func(c *app.Context) error {
    p := page.New("home.html", HomePageData{
        WelcomeMessage: "Hello, World!",
    })
    return c.Render(http.StatusOK, p)
})
```

---

## Non-goals

Sitekit intentionally does **not** provide :
* an ORM 
* dependency injection 
* a custom template language 
* a frontend framework 
* a replacement for `net/http` 

The goal is to build maintainable server-rendered applications .

---

## Status

Sitekit is under active development. While the project is already usable, APIs may continue to evolve before the first stable release .