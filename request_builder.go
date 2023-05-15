package openai

import (
	"bytes"
	"context"
	"net/http"
)

type requestBuilder interface {
	build(ctx context.Context, method, url string, request any) (*http.Request, error)
}

type httpRequestBuilder struct {
	marshaller marshaller
}

func newRequestBuilder() *httpRequestBuilder {
	return &httpRequestBuilder{
		marshaller: &jsonMarshaller{},
	}
}

func (b *httpRequestBuilder) build(ctx context.Context, method, url string, request any) (*http.Request, error) {
	if request == nil {
		return http.NewRequestWithContext(ctx, method, url, nil)
	}

	var reqBytes []byte
	reqBytes, err := b.marshaller.marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		url,
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return nil, err
	}
	if ctx.Value("header") != nil {
		header := ctx.Value("header").(map[string]string)
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	return req, nil
}
