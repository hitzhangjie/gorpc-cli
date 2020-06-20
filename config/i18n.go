package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var loc *i18n.Localizer

func initializeI18NMessages(dir string) {

	bundle := i18n.NewBundle(language.English)

	// Unmarshaling from files
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	en := filepath.Join(dir, "active.en.json")
	zh := filepath.Join(dir, "active.zh.json")

	bundle.MustLoadMessageFile(en)
	bundle.MustLoadMessageFile(zh)

	locale, err := GetLocale()
	if err != nil {
		panic(err)
	}
	loc = i18n.NewLocalizer(bundle, locale)
}

// LoadTranslation 加载对应当前locale的消息
func LoadTranslation(messageID string, data map[string]interface{}) string {

	translation := loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	return translation
}

func GetLocale() (string, error) {

	var (
		locale string
		err    error
	)
	switch runtime.GOOS {
	case "linux", "darwin":
		// Check the LANG environment variable, common on UNIX.
		// XXX: we can easily override as a nice feature/bug.
		envlang, ok := os.LookupEnv("LANG")
		if !ok {
			//err = errors.New("env LANG not set")
			//return "", err
			return "en", nil
		}
		locale = strings.Split(envlang, ".")[0]
	case "windows":
		// Exec powershell Get-Culture on Windows.
		cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		locale = strings.Trim(string(output), "\r\n")
	default:
		err = fmt.Errorf("cannot determine locale")
		return "", err
	}

	return locale, nil
}
