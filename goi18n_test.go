package goi18n

import (
	"github.com/bwmarrin/discordgo"
	"testing"
)

func TestLoadMultipleBundles(t *testing.T) {
	i18n := NewLanguageI18N()

	// Загрузка тестовых файлов для разных языков
	languages := map[discordgo.Locale]string{
		discordgo.Russian:   "locales/ru.json",
		discordgo.EnglishUS: "locales/en.json",
		discordgo.German:    "locales/de.json",
	}

	for locale, path := range languages {
		err := i18n.LoadBundle(locale, path)
		if err != nil {
			t.Fatalf("Failed to load bundle for %s: %v", locale, err)
		}

		// Проверка наличия загруженного пакета
		bundle, exists := i18n.translations[locale]
		if !exists {
			t.Fatalf("%s bundle was not loaded", locale)
		}

		// Проверка наличия ключей в пакете
		if _, ok := bundle["hello.world"]; !ok {
			t.Errorf("Key 'hello.world' is missing in %s bundle", locale)
		}
		if _, ok := bundle["goodbye"]; !ok {
			t.Errorf("Key 'goodbye' is missing in %s bundle", locale)
		}
	}
}

func TestGetTranslationsForMultipleLocales(t *testing.T) {
	i18n := NewLanguageI18N()

	// Загрузка тестовых файлов для разных языков
	err := i18n.LoadBundle(discordgo.Russian, "locales/ru.json")
	if err != nil {
		t.Fatalf("Failed to load Russian bundle: %v", err)
	}
	err = i18n.LoadBundle(discordgo.EnglishUS, "locales/en.json")
	if err != nil {
		t.Fatalf("Failed to load English bundle: %v", err)
	}
	err = i18n.LoadBundle(discordgo.German, "locales/de.json")
	if err != nil {
		t.Fatalf("Failed to load German bundle: %v", err)
	}

	// Установка языка по умолчанию
	i18n.SetDefault(discordgo.EnglishUS)

	tests := []struct {
		locale  discordgo.Locale
		key     string
		vars    Vars
		wantAny []string // Ожидаемые варианты перевода
	}{
		{discordgo.Russian, "hello.world", Vars{"name": "World"}, []string{"Привет, World!"}},
		{discordgo.Russian, "goodbye", nil, []string{"Пока!", "До свидания!"}},
		{discordgo.EnglishUS, "hello.world", Vars{"name": "World"}, []string{"Hello, World!"}},
		{discordgo.EnglishUS, "goodbye", nil, []string{"See you!", "Goodbye!"}},
		{discordgo.German, "hello.world", Vars{"name": "Welt"}, []string{"Hallo, Welt!"}},
		{discordgo.German, "goodbye", nil, []string{"Tschüss!", "Auf Wiedersehen!"}},
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

func TestFallbackBetweenLocales(t *testing.T) {
	i18n := NewLanguageI18N()

	// Загрузка тестовых файлов для разных языков
	err := i18n.LoadBundle(discordgo.Russian, "locales/ru.json")
	if err != nil {
		t.Fatalf("Failed to load Russian bundle: %v", err)
	}
	err = i18n.LoadBundle(discordgo.EnglishUS, "locales/en.json")
	if err != nil {
		t.Fatalf("Failed to load English bundle: %v", err)
	}

	// Установка языка по умолчанию
	i18n.SetDefault(discordgo.EnglishUS)

	// Тест fallback'a для немецкого языка (отсутствующий пакет)
	translation := i18n.Get(discordgo.German, "hello.world", Vars{"name": "Welt"})
	expected := "Hello, Welt!" // Fallback к английскому языку
	if translation != expected {
		t.Errorf("Expected fallback to default locale, but got '%s'", translation)
	}
}

func TestMissingKeyInSomeLocales(t *testing.T) {
	i18n := NewLanguageI18N()

	// Загрузка тестовых файлов для разных языков
	err := i18n.LoadBundle(discordgo.Russian, "locales/ru.json")
	if err != nil {
		t.Fatalf("Failed to load Russian bundle: %v", err)
	}
	err = i18n.LoadBundle(discordgo.EnglishUS, "locales/en.json")
	if err != nil {
		t.Fatalf("Failed to load English bundle: %v", err)
	}

	// Установка языка по умолчанию
	i18n.SetDefault(discordgo.EnglishUS)

	// Добавляем новый ключ только в русский пакет
	russianBundle := i18n.translations[discordgo.Russian]
	russianBundle["new.key"] = []string{"Новый ключ"}

	// Попытка получить этот ключ для других языков
	translation := i18n.Get(discordgo.EnglishUS, "new.key", nil)
	if translation != "new.key" {
		t.Errorf("Expected key 'new.key' to be returned for English, but got '%s'", translation)
	}

	translation = i18n.Get(discordgo.Russian, "new.key", nil)
	if translation != "Новый ключ" {
		t.Errorf("Expected 'Новый ключ' for Russian, but got '%s'", translation)
	}
}
