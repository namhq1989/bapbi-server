package httprespond

import (
	"errors"
	"net/http"
	"reflect"

	"golang.org/x/text/language"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"

	"github.com/labstack/echo/v4"
)

func isDataNil(data interface{}) bool {
	if data == nil {
		return true
	}
	v := reflect.ValueOf(data)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

func sendResponse(c echo.Context, httpCode int, err error, data interface{}) error {
	lang := c.Get("lang").(string)
	if lang == "" {
		lang = language.English.String()
	}

	if isDataNil(data) {
		data = echo.Map{}
	}

	code, message := apperrors.GetMessage(lang, err)

	return c.JSON(httpCode, echo.Map{
		"data":    data,
		"code":    code,
		"message": message,
	})
}

// R200 response success
func R200(c echo.Context, data interface{}) error {
	return sendResponse(c, http.StatusOK, apperrors.Common.Success, data)
}

// R400 bad request
func R400(c echo.Context, err error, data interface{}) error {
	// redirect to 403 if error is not_allowed
	if errors.Is(err, apperrors.Auth.NotAllowed) {
		return R403(c, err, data)
	}

	// redirect to 404 if error is not_found
	if errors.Is(err, apperrors.Common.NotFound) {
		return R404(c, err, data)
	}

	if err == nil {
		err = apperrors.Common.BadRequest
	}
	return sendResponse(c, http.StatusBadRequest, err, data)
}

// R401 unauthorized
func R401(c echo.Context, err error, data interface{}) error {
	if err == nil {
		err = apperrors.Common.Unauthorized
	}
	return sendResponse(c, http.StatusUnauthorized, err, data)
}

// R403 forbidden
func R403(c echo.Context, err error, data interface{}) error {
	if err == nil {
		err = apperrors.Common.Forbidden
	}
	return sendResponse(c, http.StatusForbidden, err, data)
}

// R404 not found
func R404(c echo.Context, err error, data interface{}) error {
	if err == nil {
		err = apperrors.Common.NotFound
	}
	return sendResponse(c, http.StatusNotFound, err, data)
}
