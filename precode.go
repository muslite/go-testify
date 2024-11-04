package main

import (
    "net/http"
    "net/http/httptest"
    "strconv"
    "strings"
    "testing"

    "github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
    "moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
    countStr := req.URL.Query().Get("count")
    if countStr == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("count missing"))
        return
    }

    count, err := strconv.Atoi(countStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("wrong count value"))
        return
    }

    city := req.URL.Query().Get("city")

    cafe, ok := cafeList[city]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("wrong city value"))
        return
    }

    if count > len(cafe) {
        count = len(cafe)
    }

    answer := strings.Join(cafe[:count], ",")

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(answer))
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

    require.Equalf(t, http.StatusOK, status, "expected status code: %d, got %d", http.StatusOK, status)
    require.NotEmptyf(t, body, "expected not empty responce body")
    require.Lenf(t, bodyList, totalCount,
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

    require.Equalf(t, http.StatusBadRequest, status, "expected status code: %d, got %d", http.StatusBadRequest, status)
    // require.NotEmptyf(t, body, "expected not empty responce body")
    require.Equalf(t, wrongCity, body, "expected server's message: '%s', got '%s'", wrongCity, body)
}
