package apperrors

import (
	"os"
	"path/filepath"

	"google.golang.org/grpc/status"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

const (
	langEn = "en"
	langVi = "vi"
)

var localizers = map[string]*i18n.Localizer{
	langEn: nil,
	langVi: nil,
}

func Init() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	dir, _ := os.Getwd()
	path := filepath.Join(dir, "internal", "utils", "error", "i18n")
	bundle.MustLoadMessageFile(filepath.Join(path, "active.en.toml"))
	bundle.MustLoadMessageFile(filepath.Join(path, "active.vi.toml"))

	localizers[langEn] = i18n.NewLocalizer(bundle, language.English.String(), langEn)
	localizers[langVi] = i18n.NewLocalizer(bundle, language.Vietnamese.String(), langVi)
}

func GetMessage(lang string, err error) (code, msg string) {
	key := err.Error()

	if grpcErr, ok := status.FromError(err); ok {
		key = grpcErr.Message()
	}

	if localizers[lang] == nil {
		return "_", key
	}

	msg, localizeErr := localizers[lang].Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	if localizeErr != nil {
		return "-", key
	}

	return key, msg
}
