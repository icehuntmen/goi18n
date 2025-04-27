package goi18n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/template"

	"github.com/bwmarrin/discordgo"
)

// Vars represents a map of variables to be injected into translations.
type Vars map[string]interface{}

const (
	defaultLocale   = discordgo.EnglishUS
	leftDelim       = "{{"
	rightDelim      = "}}"
	keyDelim        = "."
	executionPolicy = "missingkey=error"
)

// LanguageI18N is the main struct for handling translations.
type LanguageI18N struct {
	defaultLocale discordgo.Locale
	translations  map[discordgo.Locale]bundle
	loadedBundles map[string]bundle
}

// bundle represents a map of translations for a specific locale.
type bundle map[string][]string

// NewLanguageI18N creates a new instance of LanguageI18N.
func NewLanguageI18N() *LanguageI18N {
	return &LanguageI18N{
		defaultLocale: defaultLocale,
		translations:  make(map[discordgo.Locale]bundle),
		loadedBundles: make(map[string]bundle),
	}
}

// SetDefault sets the locale used as a fallback.
func (l *LanguageI18N) SetDefault(language discordgo.Locale) {
	l.defaultLocale = language
}

// LoadBundle loads a translation file corresponding to a specified locale.
func (l *LanguageI18N) LoadBundle(locale discordgo.Locale, path string) error {
	loadedBundle, found := l.loadedBundles[path]
	if !found {
		buf, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Failed to read file '%s': %v\n", path, err)
			return err
		}

		var jsonContent map[string]interface{}
		err = json.Unmarshal(buf, &jsonContent)
		if err != nil {
			fmt.Printf("Failed to unmarshal JSON content from '%s': %v\n", path, err)
			return err
		}

		newBundle := l.mapBundleStructure(jsonContent)

		fmt.Printf("Bundle '%s' loaded with '%s' content\n", locale, path)
		l.loadedBundles[path] = newBundle
		l.translations[locale] = newBundle
	} else {
		fmt.Printf("Bundle '%s' loaded with '%s' content (already loaded for other locales)\n", locale, path)
		l.translations[locale] = loadedBundle
	}

	return nil
}

// Get gets a translation corresponding to a locale and a key.
func (l *LanguageI18N) Get(locale discordgo.Locale, key string, variables Vars) string {
	bundles, found := l.translations[locale]
	if !found {
		if locale != l.defaultLocale {
			fmt.Printf("Bundle '%s' is not loaded, trying to translate key '%s' in '%s'\n", locale, key, l.defaultLocale)
			return l.GetDefault(key, variables)
		}

		fmt.Printf("Bundle '%s' is not loaded, cannot translate '%s', key returned\n", locale, key)
		return key
	}

	raws, found := bundles[key]
	if !found || len(raws) == 0 {
		if locale != l.defaultLocale {
			fmt.Printf("No label found for key '%s' in '%s', trying to translate it in %s\n", key, locale, l.defaultLocale)
			return l.GetDefault(key, variables)
		}

		fmt.Printf("No label found for key '%s' in '%s', key returned\n", locale, key)
		return key
	}

	raw := raws[rand.Intn(len(raws))]

	if variables != nil && strings.Contains(raw, leftDelim) {
		t, err := template.New("").Delims(leftDelim, rightDelim).Option(executionPolicy).Parse(raw)
		if err != nil {
			fmt.Printf("Cannot parse raw corresponding to key '%s' in '%s': %v\n", locale, key, err)
			return raw
		}

		var buf bytes.Buffer
		err = t.Execute(&buf, variables)
		if err != nil {
			fmt.Printf("Cannot inject variables in raw corresponding to key '%s' in '%s': %v\n", locale, key, err)
			return raw
		}
		return buf.String()
	}

	return raw
}

// GetDefault gets a translation corresponding to the default locale and a key.
func (l *LanguageI18N) GetDefault(key string, variables Vars) string {
	return l.Get(l.defaultLocale, key, variables)
}

// GetLocalizations retrieves translations from every loaded bundle.
func (l *LanguageI18N) GetLocalizations(key string, variables Vars) *map[discordgo.Locale]string {
	localizations := make(map[discordgo.Locale]string)

	for locale := range l.translations {
		localizations[locale] = l.Get(locale, key, variables)
	}

	return &localizations
}

// mapBundleStructure maps the JSON structure to a bundle.
func (l *LanguageI18N) mapBundleStructure(jsonContent map[string]interface{}) bundle {
	bundle := make(map[string][]string)
	for key, content := range jsonContent {
		switch v := content.(type) {
		case string:
			bundle[key] = []string{v}
		case []interface{}:
			values := make([]string, 0)
			for _, value := range v {
				values = append(values, fmt.Sprintf("%v", value))
			}
			bundle[key] = values
		case map[string]interface{}:
			subValues := l.mapBundleStructure(v)
			for subKey, subValue := range subValues {
				bundle[fmt.Sprintf("%s%s%s", key, keyDelim, subKey)] = subValue
			}
		default:
			bundle[key] = []string{fmt.Sprintf("%v", v)}
		}
	}

	return bundle
}
