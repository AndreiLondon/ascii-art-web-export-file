package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestBanners(t *testing.T) {
	casesMap := make(map[string]string)

	str, err := readFile("example01.txt")
	if err == nil {
		casesMap["123??"] = str
	}

	str, err = readFile("example00.txt")
	if err == nil {
		casesMap[`{123}
			<Hello> (World)!`] = str
	}

	execute(t, casesMap, "standard")

	casesMap = make(map[string]string)
	str, err = readFile("example02.txt")
	if err == nil {
		casesMap["$% \"="] = str
	}

	execute(t, casesMap, "shadow")

	casesMap = make(map[string]string)
	str, err = readFile("example03.txt")
	if err == nil {
		casesMap["123 T/fs#R"] = str
	}
	execute(t, casesMap, "thinkertoy")
}

func execute(t *testing.T, casesMap map[string]string, banner string) {
	for k, v := range casesMap {

		data := url.Values{}
		data.Set("request", k)
		data.Set("banner", banner)
		reader := strings.NewReader(data.Encode())
		// Create a request to pass to our handler.
		req, err := http.NewRequest("POST", "/ascii-art", reader)
		if err != nil {
			t.Fatal(err)
		}
		// requires to add it as
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(resultHandler)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := v
		if rr.Body.String() != expected {
			// fmt.Println("Expected:")
			// for _, v := range expected {
			// 	fmt.Print(v)
			// 	fmt.Print(" ")
			// }
			// fmt.Println()
			// fmt.Println("Got:")
			// for _, v := range rr.Body.String() {
			// 	fmt.Print(v)
			// 	fmt.Print(" ")
			// }
			// fmt.Println()
			t.Errorf("handler returned unexpected body: got \n%v want \n%v",
				rr.Body.String(), expected)
		}
	}
}

func TestStatus200FormHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(formHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	expected := http.StatusOK
	if status != expected {
		t.Errorf("handler returned unexpected status code: got %v want %v", status, expected)
	}
}

func TestStatus200ResultHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	data := url.Values{}
	data.Set("request", "text")
	data.Set("banner", "standard")
	reader := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", "/ascii-art", reader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(resultHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	expected := http.StatusOK
	if status != expected {
		t.Errorf("handler returned unexpected status code: got %v want %v", status, expected)
	}
}

func TestStatus400ResultHandlerWrongMethod(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/ascii-art", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(resultHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	expected := http.StatusBadRequest
	if status != expected {
		t.Errorf("handler returned unexpected status code: got %v want %v", status, expected)
	}
}

func TestStatus400FormHandlerInvalidEndPoint(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/1234", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(formHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	expected := http.StatusBadRequest
	if status != expected {
		t.Errorf("handler returned unexpected status code: got %v want %v", status, expected)
	}
}

func TestStatus400ResultHandlerDifferentMethod(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("PUT", "/ascii-art", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(resultHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	expected := http.StatusBadRequest
	if status != expected {
		t.Errorf("handler returned unexpected status code: got %v want %v", status, expected)
	}
}

func TestStatus404FormHandlerTemplateNotFound(t *testing.T) {
	templatePath = "templates/1234.html"

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(formHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	expected := http.StatusNotFound
	if status != expected {
		t.Errorf("handler returned unexpected status code: got %v want %v", status, expected)
	}
}

func TestStatus404ResultHandlerBannerNotFound(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/ascii-art", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(resultHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	expected := http.StatusNotFound
	if status != expected {
		t.Errorf("handler returned unexpected status code: got %v want %v", status, expected)
	}
}
