# goi18n for DiscorgGo

[![Go](https://img.shields.io/badge/Go-1.20+-blue.svg)](https://golang.org/) [![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

`goi18n` is a lightweight Go package for internationalization (i18n) that simplifies the process of managing translations in your applications. It supports loading translations from JSON files, injecting variables into translation strings, fallback to a default language, and random selection from multiple translation options.

---

## Table of Contents

- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
    - [Creating an Instance](#creating-an-instance)
    - [Loading Translation Files](#loading-translation-files)
    - [Getting Translations](#getting-translations)
    - [Handling Nested Keys](#handling-nested-keys)
    - [Fallback to Default Language](#fallback-to-default-language)
- [API Reference](#api-reference)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

---

## Introduction

The `goi18n` package is designed to help developers manage translations for multi-language applications. It provides a simple and efficient way to load translations from JSON files, inject variables into translation strings, and handle fallbacks when a specific translation is not available.

Key features:
- Supports nested keys in JSON files.
- Allows variable injection using Go's `text/template` syntax.
- Provides fallback to a default language if a translation is missing.
- Randomly selects from multiple translation options for a single key.

---

## Installation

To install the `goi18n` package, use the following command:

```bash
go get github.com/icehuntmen/goi18n
```

Replace `github.com/icehuntmen/goi18n` with the actual repository URL if it differs.

---

## Usage

### Creating an Instance

Start by creating an instance of `LanguageI18N`:

```go
package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/icehuntmen/goi18n"
)

func main() {
	i18n := goi18n.NewLanguageI18N()

	// Set the default language
	i18n.SetDefault(discordgo.EnglishUS)

	// Load translation files
	err := i18n.LoadBundle(discordgo.Russian, "locales/ru.json")
	if err != nil {
		panic(err)
	}

	// Get a translation
	translation := i18n.Get(discordgo.Russian, "hello.world", goi18n.Vars{"name": "World"})
	fmt.Println(translation) // Output: Привет, World!
}
```

---

### Loading Translation Files

Translation files are stored in JSON format. Here's an example of a Russian translation file (`locales/ru.json`):

```json
{
  "hello": {
    "world": "Привет, {{.name}}!"
  },
  "goodbye": ["Пока!", "До свидания!"]
}
```

You can load this file using the `LoadBundle` method:

```go
err := i18n.LoadBundle(discordgo.Russian, "locales/ru.json")
if err != nil {
    panic(err)
}
```

---

### Getting Translations

To retrieve a translation for a specific locale and key, use the `Get` method:

```go
translation := i18n.Get(discordgo.Russian, "hello.world", goi18n.Vars{"name": "World"})
fmt.Println(translation) // Output: Привет, World!
```

If the key contains variables (e.g., `{{.name}}`), they will be replaced with the values provided in the `Vars` map.

---

### Handling Nested Keys

The package supports nested keys in JSON files. For example, if your JSON file looks like this:

```json
{
  "hello": {
    "world": "Hello, {{.name}}!"
  }
}
```

You can access the nested key using dot notation:

```go
translation := i18n.Get(discordgo.EnglishUS, "hello.world", goi18n.Vars{"name": "World"})
fmt.Println(translation) // Output: Hello, World!
```

---

### Fallback to Default Language

If a translation is not available for the requested locale, the package will fall back to the default language:

```go
defaultTranslation := i18n.GetDefault("hello.world", goi18n.Vars{"name": "World"})
fmt.Println(defaultTranslation) // Output: Hello, World!
```

You can set the default language using the `SetDefault` method:

```go
i18n.SetDefault(discordgo.EnglishUS)
```

---

## API Reference

### `NewLanguageI18N()`

Creates a new instance of `LanguageI18N`.

```go
i18n := goi18n.NewLanguageI18N()
```

### `SetDefault(locale discordgo.Locale)`

Sets the default locale used as a fallback.

```go
i18n.SetDefault(discordgo.EnglishUS)
```

### `LoadBundle(locale discordgo.Locale, path string) error`

Loads a translation file for the specified locale.

```go
err := i18n.LoadBundle(discordgo.Russian, "locales/ru.json")
```

### `Get(locale discordgo.Locale, key string, vars Vars) string`

Retrieves a translation for the specified locale and key. Replaces variables in the translation string.

```go
translation := i18n.Get(discordgo.Russian, "hello.world", goi18n.Vars{"name": "World"})
```

### `GetDefault(key string, vars Vars) string`

Retrieves a translation for the default locale.

```go
defaultTranslation := i18n.GetDefault("hello.world", goi18n.Vars{"name": "World"})
```

### `GetLocalizations(key string, vars Vars) *map[discordgo.Locale]string`

Retrieves all loaded translations for the specified key.

```go
localizations := i18n.GetLocalizations("hello.world", goi18n.Vars{"name": "World"})
for locale, translation := range *localizations {
    fmt.Printf("%s: %s\n", locale, translation)
}
```

---

## Testing

To test the package, create JSON files for different locales (e.g., `locales/ru.json`, `locales/en.json`, `locales/de.json`) and run the tests using Go's testing framework.

Example test:

```go
package goi18n

import (
	"github.com/bwmarrin/discordgo"
	"testing"
)

func TestGetTranslationsForMultipleLocales(t *testing.T) {
	i18n := NewLanguageI18N()

	// Load translation files
	err := i18n.LoadBundle(discordgo.Russian, "locales/ru.json")
	if err != nil {
		t.Fatalf("Failed to load Russian bundle: %v", err)
	}
	err = i18n.LoadBundle(discordgo.EnglishUS, "locales/en.json")
	if err != nil {
		t.Fatalf("Failed to load English bundle: %v", err)
	}

	// Set default language
	i18n.SetDefault(discordgo.EnglishUS)

	// Test translations
	tests := []struct {
		locale  discordgo.Locale
		key     string
		vars    goi18n.Vars
		wantAny []string
	}{
		{discordgo.Russian, "hello.world", goi18n.Vars{"name": "World"}, []string{"Привет, World!"}},
		{discordgo.EnglishUS, "hello.world", goi18n.Vars{"name": "World"}, []string{"Hello, World!"}},
	}

	for _, tt := range tests {
		got := i18n.Get(tt.locale, tt.key, tt.vars)
		found := false
		for _, expected := range tt.wantAny {
			if got == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Get(%s, %s, %v): expected one of '%v', but got '%s'", tt.locale, tt.key, tt.vars, tt.wantAny, got)
		}
	}
}
```

Run the tests using:

```bash
go test ./...
```

---

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you'd like to contribute improvements or fixes.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.