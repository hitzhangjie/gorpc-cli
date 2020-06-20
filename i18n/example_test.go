package i18n

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func TestI18N_DefaultMessage(t *testing.T) {

	bundle := i18n.NewBundle(language.English)
	loc := i18n.NewLocalizer(bundle, language.English.String())

	messages := &i18n.Message{
		ID:          "Emails",
		Description: "The number of unread emails a user has",
		One:         "{{.Name}} has {{.Count}} email.",
		Other:       "{{.Name}} has {{.Count}} emails.",
	}

	messagesCount := 2
	translation := loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: messages,
		TemplateData: map[string]interface{}{
			"Name":  "Theo",
			"Count": messagesCount,
		},
		PluralCount: 2,
	})

	fmt.Println(translation)
}

func TestI18N_LoadMessageFromConfig(t *testing.T) {

	bundle := i18n.NewBundle(language.English)
	loc := i18n.NewLocalizer(bundle, language.English.String())

	// Unmarshaling from files
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("en.json")
	bundle.MustLoadMessageFile("zh.json")

	loc = i18n.NewLocalizer(bundle, "zh")
	messagesCount := 10
	translation := loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "messages",
		TemplateData: map[string]interface{}{
			"Name":  "Alex",
			"Count": messagesCount,
		},
		PluralCount: messagesCount,
	})

	fmt.Println(translation)
}
