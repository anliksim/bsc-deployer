package api

import (
	"github.com/anliksim/bsc-deployer/model"
	"github.com/anliksim/bsc-deployer/util"
	"github.com/gorilla/mux"
	"github.com/nvellon/hal"
	"net/http"
)

var baseUrl string

func Register(r *mux.Router, base string) {
	baseUrl = base
	r.HandleFunc(Base, getBase)
}

func getBase(w http.ResponseWriter, r *http.Request) {
	res := hal.NewResource(&model.None{}, Url(baseUrl, ""))
	res.AddNewLink("v1", Url(baseUrl, "v1"))
	util.RespondJson(w, res)
}
