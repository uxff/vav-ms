package v1

import (
	gkhttp "github.com/giskook/go/http"
	"github.com/giskook/vav-ms/base"
	"net/http"
)

func common_reply(w http.ResponseWriter, http_status int, code string, data interface{}, err error) {
	w.WriteHeader(http_status)
	err_msg := ""
	if err != nil {
		err_msg = err.Error()
	}
	if code != "" {
		err_msg += base.ErrorMap[code]
	}
	gkhttp.EncodeResponse(w, code, data, err_msg)
}
