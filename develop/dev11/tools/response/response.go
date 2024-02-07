package response

import (
	"dev11/core"
	"encoding/json"
	"log"
	"net/http"
)

func Resp(w http.ResponseWriter, res interface{}, errResp interface{}, status ...int) {
	var st int
	for _, s := range status {
		st = s
		break
	}

	if st < 100 {
		st = 200
	}
	var respSuccess core.SuccessResponse
	var respError core.ErrorResponse
	if errResp == nil {
		respSuccess = core.SuccessResponse{
			Result: res,
		}
	} else {
		respError = core.ErrorResponse{
			Error: errResp,
		}
	}

	if respSuccess.Result != nil {
		respSucc, err := json.Marshal(respSuccess)

		if err != nil {
			log.Printf("response: can not marshal response: %s\n", err.Error())
			return
		}
		w.WriteHeader(st)
		if _, err := w.Write(respSucc); err != nil {
			log.Printf("response: can not write response: %s\n", err.Error())
		}
		return
	}

	if respError.Error != nil {
		respErr, err := json.Marshal(respSuccess)

		if err != nil {
			log.Printf("response: can not marshal response: %s\n", err.Error())
			return
		}

		http.Error(w, string(respErr), st)
	}
}
