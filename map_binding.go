package httpserver

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
)

type mapBinding struct{}

func (mapBinding) Name() string {
	return "map"
}

func (mapBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		return errors.Wrapf(err, "read request body err")
	}

	if len(data) == 0 {
		values := req.URL.Query()
		query := make(map[string]string)
		for k, v := range values {
			query[k] = v[0]
		}

		data, err = json.Marshal(query)
		if err != nil {
			return errors.Wrapf(err, "marshal query err")
		}
	}

	if len(data) == 0 {
		return nil
	}

	return binding.JSON.BindBody(data, obj)
}
