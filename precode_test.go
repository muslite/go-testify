package main

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
//    "github.com/stretchr/testify/require"
)

func TestMainHandlerServerRespondsCorrectly(t *testing.T) {

    req := httptest.NewRequest("GET", "/cafe?count=4&city=moscow", nil)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)

    body := responseRecorder.Body.String()
    status := responseRecorder.Code

    assert.Equalf(t, http.StatusOK, status, "expected status code: %d, got %d", http.StatusOK, status)
    assert.NotEmptyf(t, body, "expected not empty responce body")

}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
    totalCount := 4
    req := httptest.NewRequest("GET", "/cafe?count=5&city=moscow", nil)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)

    // здесь нужно добавить необходимые проверки
    // добавим вспомогательных переменных
    body := responseRecorder.Body.String()
    bodyList := strings.Split(body, ",")
    status := responseRecorder.Code

    assert.Equalf(t, http.StatusOK, status, "expected status code: %d, got %d", http.StatusOK, status)
    assert.Lenf(t, bodyList, totalCount,
        "expected caffe count: %d, got %d", totalCount, len(bodyList))

}

// Проверяем ответ на не поддерживаемый город (отсутствующий в базе)
func TestMainHandlerWhenUnknownCountry(t *testing.T) {
    wrongCity := "wrong city value"
    req := httptest.NewRequest("GET", "/cafe?count=1&city=tula", nil)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)

    // добавим вспомогательных переменных
    body := responseRecorder.Body.String()
    status := responseRecorder.Code

    assert.Equalf(t, http.StatusBadRequest, status, "expected status code: %d, got %d", http.StatusBadRequest, status)
    assert.Equalf(t, wrongCity, body, "expected server's message: '%s', got '%s'", wrongCity, body)

}
