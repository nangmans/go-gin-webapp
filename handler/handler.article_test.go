package handler_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nangmans14/gin-web/handler"
	common "github.com/nangmans14/gin-web/test"
)

func TestShowIndexPageUnauthenticated(t *testing.T) {
	r := common.GetRouter(true)

	r.GET("/", handler.ShowIndexPage)

	req, _ := http.NewRequest("GET", "/", nil)

	common.TestHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Home Page</title>") > 0

		return statusOK && pageOK

	})
}
