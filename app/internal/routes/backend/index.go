package backend

import (
	"net/http"

	"github.com/Owicca/tr/internal/infra"
)

func init() {
	adminRouter := infra.S.Router.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/", Index).Methods(http.MethodGet).Name("back_index")
}

func Index(w http.ResponseWriter, r *http.Request) {
	infra.S.HTML(w, r, http.StatusOK, "back/index", data)
}
