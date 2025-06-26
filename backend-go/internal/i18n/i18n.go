package i18n

import (
	"encoding/json"
	"fmt"
	"os"
)

// Messages holds the loaded translation key-value pairs
var Messages map[string]string

// LoadMessages loads the translations from a JSON file for the given locale.
// The JSON file should be located at "assets/i18n/{locale}.json"
func LoadMessages(locale string) error {
	path := fmt.Sprintf("assets/i18n/%s.json", locale)

	file, err := os.Open(path)
	if err != nil {
		// Return error with translated message key for file open failure
		return fmt.Errorf("%s: %w", T("error.i18n.load"), err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&Messages); err != nil {
		// Return error with translated message key for JSON decode failure
		return fmt.Errorf("%s: %w", T("error.i18n_json_decode_failed"), err)
	}

	return nil
}

// T returns the translation for the given key or the key itself if missing.
// Used throughout the backend to retrieve localized strings.
func T(key string) string {
	if val, ok := Messages[key]; ok {
		return val
	}
	return key
}
