package server

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"
)

func returnSucceeded(response *restful.Response) {
	response.AddHeader("Content-Type", "text/json")
	//response.WriteHeader(http.StatusOK)
	response.WriteAsJson(&EdgeOtaResponse{
		Code:   http.StatusOK,
		Status: SUCCEEDED,
	})
}

func returnFaileure(response *restful.Response, err error) {
	response.AddHeader("Content-Type", "text/json")
	response.WriteHeader(http.StatusInternalServerError)
	response.WriteAsJson(&EdgeOtaResponse{
		Code:    http.StatusInternalServerError,
		Status:  FAILURE,
		Message: fmt.Sprintf("read entity error: [+%v]", err),
	})

}
