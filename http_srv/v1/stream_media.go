package v1

import (
	"encoding/json"
	gkbase "github.com/giskook/go/base"
	gkhttp "github.com/giskook/go/http"
	rc "github.com/giskook/vav-common/redis_cli"
	"github.com/giskook/vav-ms/base"
	"github.com/giskook/vav-ms/redis_cli"
	"net/http"
)

func stream_media_get(w http.ResponseWriter, r *http.Request) (int, string, interface{}, error) {
	stream_medias, err := rc.GetInstance().GetStreamMedia(redis_cli.AV_STREAM_MEDIA, "0", "-1")
	if err != nil {
		return http.StatusInternalServerError, base.HTTP_INTERNAL_SERVER_ERROR_GET_STREAM_MEDIA, nil, err
	}

	return http.StatusOK, base.HTTP_OK, stream_medias, nil
}

func stream_media_post(w http.ResponseWriter, r *http.Request) (int, string, error) {
	r.ParseForm()
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var stream_medias base.StreamMedias
	err := decoder.Decode(&stream_medias)
	if err != nil {
		gkbase.ErrorCheck(err)
		return http.StatusBadRequest, base.HTTP_BAD_REQUEST_DECODE, err
	}
	if stream_medias.StreamMedias == nil ||
		len(stream_medias.StreamMedias) == 0 {
		gkbase.ErrorCheck(base.ERROR_BAD_REQUEST_MISSING)
		return http.StatusBadRequest, base.HTTP_BAD_REQUEST_MISSING, base.ERROR_BAD_REQUEST_MISSING
	}

	for _, v := range stream_medias.StreamMedias {
		if v.AccessUUID == "" || v.DomainInner == "" || v.DomainOuter == "" {
			gkbase.ErrorCheck(base.ERROR_BAD_REQUEST_MISSING)
			return http.StatusBadRequest, base.HTTP_BAD_REQUEST_MISSING, base.ERROR_BAD_REQUEST_MISSING
		}
	}

	err = rc.GetInstance().SetStreamMedia(redis_cli.AV_STREAM_MEDIA, stream_medias.StreamMedias)
	if err != nil {
		gkbase.ErrorCheck(err)
		return http.StatusInternalServerError, base.HTTP_INTERNAL_SERVER_ERROR_ADD_STREAM_MEDIA, err
	}

	return http.StatusCreated, base.HTTP_OK, nil
}

func StreamMedia(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			gkbase.ErrorPrintStack()
			w.WriteHeader(http.StatusInternalServerError)
			gkhttp.EncodeResponse(w, base.HTTP_INTERNAL_SERVER_ERROR, nil, "")
		}
	}()
	gkhttp.RecordReq(r)

	var http_status int
	var internal_status string
	var err error
	var data interface{}

	switch r.Method {
	case http.MethodPost:
		http_status, internal_status, err = stream_media_post(w, r)
	case http.MethodGet:
		http_status, internal_status, data, err = stream_media_get(w, r)
	}

	common_reply(w, http_status, internal_status, data, err)
}