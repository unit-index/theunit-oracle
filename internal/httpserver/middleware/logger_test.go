""

package middleware

import (
"net/http"
"net/http/httptest"
"strings"
"testing"

"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"

"github.com/toknowwhy/theunit-oracle/pkg/log"
"github.com/toknowwhy/theunit-oracle/pkg/log/callback"
)

func TestLogger_InfoLevel(t *testing.T) {
var recordedLogMsgs []string
var recordedLogFields []log.Fields
l := callback.New(log.Info, func (level log.Level, fields log.Fields, msg string) {
if level != log.Info {
return
}
recordedLogMsgs = append(recordedLogMsgs, msg)
recordedLogFields = append(recordedLogFields, fields)
})

r := httptest.NewRequest("GET", "/", nil)
w := httptest.NewRecorder()
h := (&Logger{Log: l}).Handle(http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {}))
h.ServeHTTP(w, r)

require.Len(t, recordedLogMsgs, 1)
assert.Equal(t, httpRequestLog, recordedLogMsgs[0])
assert.Equal(t, "GET", recordedLogFields[0]["method"])
assert.Equal(t, "/", recordedLogFields[0]["url"])
assert.NotEmpty(t, recordedLogFields[0]["duration"])
assert.NotEmpty(t, recordedLogFields[0]["remoteAddr"])
}

func TestLogger_DebugLevel(t *testing.T) {
var recordedLogMsgs []string
var recordedLogFields []log.Fields
l := callback.New(log.Debug, func (level log.Level, fields log.Fields, msg string) {
if level != log.Debug {
return
}
recordedLogMsgs = append(recordedLogMsgs, msg)
recordedLogFields = append(recordedLogFields, fields)
})

r := httptest.NewRequest("GET", "/", strings.NewReader("request"))
w := httptest.NewRecorder()
h := (&Logger{Log: l}).Handle(http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
writer.Write([]byte("response"))
writer.WriteHeader(http.StatusNotFound)
}))
h.ServeHTTP(w, r)

require.Len(t, recordedLogMsgs, 1)
assert.Equal(t, httpRequestLog, recordedLogMsgs[0])
assert.Equal(t, "GET", recordedLogFields[0]["method"])
assert.Equal(t, "/", recordedLogFields[0]["url"])
assert.Equal(t, http.StatusNotFound, recordedLogFields[0]["status"])
assert.Equal(t, "request", recordedLogFields[0]["request"])
assert.Equal(t, "response", recordedLogFields[0]["response"])
assert.NotEmpty(t, recordedLogFields[0]["duration"])
assert.NotEmpty(t, recordedLogFields[0]["remoteAddr"])
}
