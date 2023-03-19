package mod

import (
	"fmt"
	"net/http"
)

func (s *Server) getReports(w http.ResponseWriter, r *http.Request) error {
	showClosed := r.FormValue("closed") == "true"

	fmt.Println("closed =", showClosed)
	return nil
}
