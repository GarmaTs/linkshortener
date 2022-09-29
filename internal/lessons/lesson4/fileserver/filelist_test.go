package fileserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFileListHander(t *testing.T) {
	req, err := http.NewRequest("GET", "/?extension=.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := &FileListHander{Dir: "upload"}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	want := ".txt"
	got := handler.extension
	if got != want {
		t.Errorf("Error defining extension, want %s, got %s", want, got)
	}
}
