package httpresponse

import (
	"encoding/json"
	"net/http"
)

func WriteMessage(w http.ResponseWriter, status int, msg string) {
	var j struct {
		Msg string `json:"message"`
	}

	j.Msg = msg

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(j); err != nil {
		http.Error(w, err.Error(), status)
		return
	}
}

func WriteData(w http.ResponseWriter, status int, data, meta any) {
	switch meta {
	case nil:
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), status)
			return
		}
		return
	default:
		var j struct {
			Data any `json:"data"`
			Meta any `json:"meta"`
		}

		j.Data = data
		j.Meta = meta

		w.Header().Add("content-type", "application/json")
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(j); err != nil {
			http.Error(w, err.Error(), status)
			return
		}
	}

}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteMessage(w, status, err.Error())
}
