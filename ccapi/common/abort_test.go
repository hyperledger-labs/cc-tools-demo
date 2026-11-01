package common

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRespond(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cases := []struct {
		name   string
		res    interface{}
		status int
		err    error
		code   int
		body   string
	}{
		{"default", map[string]interface{}{"a": 1.0}, 200, nil, 200, `{"a":1}`},
		{"created", map[string]interface{}{"a": 1.0}, 201, nil, 201, `{"a":1}`},
		{"no content", nil, 204, nil, 204, ``},
		{"non 2xx success", nil, 302, nil, 200, `null`},
		{"zero status", nil, 0, nil, 200, `null`},
		{"error", nil, 404, errors.New("not found"), 404, `{"error":"not found","response":null,"status":404}`},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		Respond(ctx, c.res, c.status, c.err)

		if w.Code != c.code || w.Body.String() != c.body {
			t.Errorf("%s: expected %d %q, got %d %q", c.name, c.code, c.body, w.Code, w.Body.String())
		}
	}
}

func TestParsePayload(t *testing.T) {
	res, err := ParsePayload(nil)
	if res != nil || err != nil {
		t.Errorf("empty payload: expected nil, got %v %v", res, err)
	}

	res, err = ParsePayload([]byte(`{"a":1}`))
	if err != nil || res.(map[string]interface{})["a"] != 1.0 {
		t.Errorf("unexpected result %v %v", res, err)
	}

	_, err = ParsePayload([]byte(`{`))
	if err == nil {
		t.Error("expected error for invalid payload")
	}
}
