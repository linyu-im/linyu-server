package i18n

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"io/fs"
	"path/filepath"
)

var Bundle *i18n.Bundle

func InitI18n(localesDir string) {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("yml", yaml.Unmarshal)

	var loadedFiles []string

	err := filepath.WalkDir(localesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && (filepath.Ext(path) == ".yml" || filepath.Ext(path) == ".yaml") {
			if _, err := Bundle.LoadMessageFile(path); err != nil {
				panic(fmt.Sprintf("[i18n] Failed to load %s: %v", path, err))
			}
			loadedFiles = append(loadedFiles, path)
		}
		return nil
	})

	if err != nil {
		panic(fmt.Sprintf("[i18n] Failed to scan locales directory: %v", err))
	}

	if len(loadedFiles) == 0 {
		panic("[i18n] No locale files found in ./locales directory")
	}
	fmt.Println("[i18n] All locale files loaded successfully")
}

func T(c *gin.Context, messageID string, data map[string]interface{}) string {
	lang, exists := c.Get("lang")
	if !exists {
		return messageID
	}
	localizer := i18n.NewLocalizer(Bundle, lang.(string))
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		return messageID
	}
	return msg
}
