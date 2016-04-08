package encoding_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/ayiga/go-kit-middlewarer/encoding"
)

func TestDecodeErrorJSON(t *testing.T) {
	buf := new(bytes.Buffer)
	rw := createResponseWriter(buf)

	// server error...
	rw.WriteHeader(500)
	err := encoding.JSON(0).EncodeResponse()(rw, http.ErrContentLength)
	if err != nil {
		t.Fatalf("Unable to Encode Response: %s", err)
	}

	t.Logf("Body Content: %s", buf.String())

	ro := new(http.Response)
	ro.StatusCode = rw.statusCode
	ro.Body = ioutil.NopCloser(buf)
	ro.Header = make(http.Header)
	ro.Header.Set("Content-Type", "application/json")

	resp := new(request)
	r, err := encoding.Default().DecodeResponse(resp)(ro)
	if err != nil {
		t.Fatalf("Unable to Decode Response: %s", err)
	}

	if got, want := reflect.TypeOf(r), reflect.TypeOf(encoding.WrapperError{}); got != want {
		t.Fatalf("Type Of:\ngot:\n%s\nwant:\n%s", got, want)
	}

	err, ok := r.(error)
	if !ok {
		t.Fatal("Unable to cast returned response into an error")
	}

	if got, want := err.Error(), http.ErrContentLength.Error(); got != want {
		t.Errorf(".Error():\ngot:\n\t%s\nwant:\n\t%s", got, want)
	}

	if got, want := r == http.ErrMissingContentLength, false; got != want {
		t.Errorf(".Error():\ngot:\n\t%t\nwant:\n\t%t", got, want)
	}
}

func TestDecodeErrorXML(t *testing.T) {
	buf := new(bytes.Buffer)
	rw := createResponseWriter(buf)

	// server error...
	rw.WriteHeader(500)
	err := encoding.XML(0).EncodeResponse()(rw, http.ErrContentLength)
	if err != nil {
		t.Fatalf("Unable to Encode Response: %s", err)
	}

	t.Logf("Body Content: %s", buf.String())

	ro := new(http.Response)
	ro.StatusCode = rw.statusCode
	ro.Body = ioutil.NopCloser(buf)
	ro.Header = make(http.Header)
	ro.Header.Set("Content-Type", "application/xml")

	resp := new(request)
	r, err := encoding.Default().DecodeResponse(resp)(ro)
	if err != nil {
		t.Fatalf("Unable to Decode Response: %s", err)
	}

	if got, want := reflect.TypeOf(r), reflect.TypeOf(encoding.WrapperError{}); got != want {
		t.Fatalf("Type Of:\ngot:\n%s\nwant:\n%s", got, want)
	}

	err, ok := r.(error)
	if !ok {
		t.Fatal("Unable to cast returned response into an error")
	}

	if got, want := err.Error(), http.ErrContentLength.Error(); got != want {
		t.Errorf(".Error():\ngot:\n\t%s\nwant:\n\t%s", got, want)
	}

	if got, want := r == http.ErrMissingContentLength, false; got != want {
		t.Errorf(".Error():\ngot:\n\t%t\nwant:\n\t%t", got, want)
	}
}

func TestDecodeErrorGob(t *testing.T) {
	buf := new(bytes.Buffer)
	rw := createResponseWriter(buf)

	// server error...
	rw.WriteHeader(500)
	err := encoding.Gob(0).EncodeResponse()(rw, http.ErrContentLength)
	if err != nil {
		t.Fatalf("Unable to Encode Response: %s", err)
	}

	t.Logf("Body Content: %s", buf.String())

	ro := new(http.Response)
	ro.StatusCode = rw.statusCode
	ro.Body = ioutil.NopCloser(buf)
	ro.Header = make(http.Header)
	ro.Header.Set("Content-Type", "application/gob")

	resp := new(request)
	r, err := encoding.Default().DecodeResponse(resp)(ro)
	if err != nil {
		t.Fatalf("Unable to Decode Response: %s", err)
	}

	if got, want := reflect.TypeOf(r), reflect.TypeOf(encoding.WrapperError{}); got != want {
		t.Fatalf("Type Of:\ngot:\n%s\nwant:\n%s", got, want)
	}

	err, ok := r.(error)
	if !ok {
		t.Fatal("Unable to cast returned response into an error")
	}

	if got, want := err.Error(), http.ErrContentLength.Error(); got != want {
		t.Errorf(".Error():\ngot:\n\t%s\nwant:\n\t%s", got, want)
	}

	if got, want := r == http.ErrMissingContentLength, false; got != want {
		t.Errorf(".Error():\ngot:\n\t%t\nwant:\n\t%t", got, want)
	}
}

type CustomDecodableError struct {
	Code   int    `json:"code" xml:"code"`
	Reason string `json:"reason" xml:"reason"`
}

func (cde CustomDecodableError) Error() string {
	return fmt.Sprintf("Code: %d, Reason: %s", cde.Code, cde.Reason)
}

func init() {
	encoding.RegisterError(CustomDecodableError{})
	gob.Register(CustomDecodableError{})
}

func TestDecodeCustomDecodableErrorJSON(t *testing.T) {
	buf := new(bytes.Buffer)
	rw := createResponseWriter(buf)

	testErr := CustomDecodableError{
		Code:   50,
		Reason: "Halp",
	}

	// server error...
	rw.WriteHeader(500)
	err := encoding.JSON(0).EncodeResponse()(rw, &testErr)
	if err != nil {
		t.Fatalf("Unable to Encode Response: %s", err)
	}

	t.Logf("Body Content: %s", buf.String())

	ro := new(http.Response)
	ro.StatusCode = rw.statusCode
	ro.Body = ioutil.NopCloser(buf)
	ro.Header = make(http.Header)
	ro.Header.Set("Content-Type", "application/json")

	resp := new(request)
	r, err := encoding.Default().DecodeResponse(resp)(ro)
	if err != nil {
		t.Fatalf("Unable to Decode Response: %s", err)
	}

	t.Logf("Decode Result: %#v", r)
	if got, want := reflect.TypeOf(r), reflect.TypeOf(testErr); got != want {
		t.Fatalf("Type Of:\ngot:\n%s\nwant:\n%s", got, want)
	}

	castErr, ok := r.(CustomDecodableError)
	if !ok {
		t.Fatal("Unable to cast returned response into an error")
	}

	if got, want := castErr.Error(), testErr.Error(); got != want {
		t.Errorf(".Error():\ngot:\n\t%s\nwant:\n\t%s", got, want)
	}

	if got, want := castErr.Code, testErr.Code; got != want {
		t.Errorf("castErr.Code:\ngot:\n\t%d\nwant:\n\t%d", got, want)
	}

	if got, want := castErr.Reason, testErr.Reason; got != want {
		t.Errorf("castErr.Reason:\ngot:\n\t%s\nwant:\n\t%s", got, want)
	}
}

func TestDecodeCustomDecodableErrorXML(t *testing.T) {
	buf := new(bytes.Buffer)
	rw := createResponseWriter(buf)

	testErr := CustomDecodableError{
		Code:   50,
		Reason: "Halp",
	}

	// server error...
	rw.WriteHeader(500)
	err := encoding.XML(0).EncodeResponse()(rw, &testErr)
	if err != nil {
		t.Fatalf("Unable to Encode Response: %s", err)
	}

	t.Logf("Body Content: %s", buf.String())

	ro := new(http.Response)
	ro.StatusCode = rw.statusCode
	ro.Body = ioutil.NopCloser(buf)
	ro.Header = make(http.Header)
	ro.Header.Set("Content-Type", "application/xml")

	resp := new(request)
	r, err := encoding.Default().DecodeResponse(resp)(ro)
	if err != nil {
		t.Fatalf("Unable to Decode Response: %s", err)
	}

	t.Logf("Decode Result: %#v", r)
	if got, want := reflect.TypeOf(r), reflect.TypeOf(testErr); got != want {
		t.Fatalf("Type Of:\ngot:\n%s\nwant:\n%s", got, want)
	}

	castErr, ok := r.(CustomDecodableError)
	if !ok {
		t.Fatal("Unable to cast returned response into an error")
	}

	if got, want := castErr.Error(), testErr.Error(); got != want {
		t.Errorf(".Error():\ngot:\n\t%s\nwant:\n\t%s", got, want)
	}

	if got, want := castErr.Code, testErr.Code; got != want {
		t.Errorf("castErr.Code:\ngot:\n\t%d\nwant:\n\t%d", got, want)
	}

	if got, want := castErr.Reason, testErr.Reason; got != want {
		t.Errorf("castErr.Reason:\ngot:\n\t%s\nwant:\n\t%s", got, want)
	}
}

func TestDecodeCustomDecodableErrorGob(t *testing.T) {
	buf := new(bytes.Buffer)
	rw := createResponseWriter(buf)

	testErr := CustomDecodableError{
		Code:   50,
		Reason: "Halp",
	}

	// server error...
	rw.WriteHeader(500)
	err := encoding.Gob(0).EncodeResponse()(rw, &testErr)
	if err != nil {
		t.Fatalf("Unable to Encode Response: %s", err)
	}

	t.Logf("Body Content: %s", buf.String())

	ro := new(http.Response)
	ro.StatusCode = rw.statusCode
	ro.Body = ioutil.NopCloser(buf)
	ro.Header = make(http.Header)
	ro.Header.Set("Content-Type", "application/gob")

	resp := new(request)
	r, err := encoding.Default().DecodeResponse(resp)(ro)
	if err != nil {
		t.Fatalf("Unable to Decode Response: %s", err)
	}

	t.Logf("Decode Result: %#v", r)
	if got, want := reflect.TypeOf(r), reflect.TypeOf(testErr); got != want {
		t.Fatalf("Type Of:\ngot:\n%s\nwant:\n%s", got, want)
	}

	castErr, ok := r.(CustomDecodableError)
	if !ok {
		t.Fatal("Unable to cast returned response into an error")
	}

	if got, want := castErr.Error(), testErr.Error(); got != want {
		t.Errorf(".Error():\ngot:\n\t%s\nwant:\n\t%s", got, want)
	}

	if got, want := castErr.Code, testErr.Code; got != want {
		t.Errorf("castErr.Code:\ngot:\n\t%d\nwant:\n\t%d", got, want)
	}

	if got, want := castErr.Reason, testErr.Reason; got != want {
		t.Errorf("castErr.Reason:\ngot:\n\t%s\nwant:\n\t%s", got, want)
	}
}
