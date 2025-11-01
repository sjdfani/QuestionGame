package httpmsg

import (
	"QuestionGame/pkg/errmsg"
	"QuestionGame/pkg/richerror"
	"net/http"
)

func Error(err error) (message string, code int) {
	switch err := err.(type) {
	case richerror.RichError:
		msg := err.GetMessage()
		code := mapKindToStatusCode(err.GetKind())
		if code > 500 {
			msg = errmsg.ErrorMsgSomethingWentWrong
		}
		return msg, code

	default:
		return err.Error(), http.StatusBadRequest
	}
}

func mapKindToStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}
