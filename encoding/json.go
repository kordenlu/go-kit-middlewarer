package encoding

import (
	"encoding/json"
	"io"

	httptransport "github.com/go-kit/kit/transport/http"
)

func init() {
	arr := []rune{'{', '[', '"', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	Register("text/json", JSON(0), arr)
	Register("application/json", JSON(0), arr)
}

// JSONGenerateDecoder returns a JSON Decoder
func JSONGenerateDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

// JSONGenerateEncoder returns a JSON Encoder
func JSONGenerateEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

// JSON is a simple JSON encoder / decoder that conforms to RequestResponseEncoding
type JSON int

// EncodeRequest implements RequestResponseEncoding
func (JSON) EncodeRequest() httptransport.EncodeRequestFunc {
	return MakeRequestEncoder(JSONGenerateEncoder)
}

// DecodeRequest implements RequestResponseEncoding
func (JSON) DecodeRequest(request interface{}) httptransport.DecodeRequestFunc {
	return MakeRequestDecoder(request, JSONGenerateDecoder)
}

// EncodeResponse implements RequestResponseEncoding
func (JSON) EncodeResponse() httptransport.EncodeResponseFunc {
	return MakeResponseEncoder(JSONGenerateEncoder)
}

// DecodeResponse implements RequestResponseEncoding
func (JSON) DecodeResponse(response interface{}) httptransport.DecodeResponseFunc {
	return MakeResponseDecoder(response, JSONGenerateDecoder)
}
