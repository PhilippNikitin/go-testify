package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenOk(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// получаем объединенную строку с названиями кафе из ответа,
	// и далее - слайс с отдельными названиями
	body := responseRecorder.Body.String()
	cafes := strings.Split(body, ",")

	// проверяем, если вернулся код ответа 400
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	// вместо assert.Equal сначала хотел использовать assert.HTTPSuccess:
	// assert.HTTPSuccess(t, handler, "GET", "/cafe?count=2&city=moscow", nil)
	// но почему-то эта строка выше не проходит проверку:
	// возвращается код 400, не понятно, почему

	// проверяем, если тело ответа не пустое, через require - если оно пустое,
	// дальнейшие проверки не имеют смысла
	require.NotEmpty(t, responseRecorder.Body.String())
	// если тело ответа не пустое, проверяем, что в ответе указано 2 кафе:
	assert.Equal(t, 2, len(cafes))
	// проверяем - что вернулись именно первые названия, а не случайные
	// сравнивать слайсы мы технически не можем, поэтому сравниваем строку
	// "Мир кофе,Сладкоежка" и responseRecorder.Body.String() (переменная body)
	assert.Equal(t, "Мир кофе,Сладкоежка", body)
}

func TestMainHandlerWhenWrongCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=london", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// проверяем, что вернулся код ответа 400
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	// проверяем, что в теле ответа вернулась ошибка: "wrong city value"
	assert.HTTPBodyContains(t, handler, "GET", "/cafe?count=2&city=london", nil, "wrong city value")
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4 // общее количество кафе
	req := httptest.NewRequest("GET", "/cafe?count=5&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	body := responseRecorder.Body.String()
	cafes := strings.Split(body, ",")

	// проверяем, что вернулся код ответа 200
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	// также хотел вместо assert.Equal использовать метод пакета assert HTTPSuccess:
	// assert.HTTPSuccess(t, handler, "GET", "/cafe?count=5&city=moscow", nil)
	// но тест возврашал значение FAIL - получался код ответа 400, а не 200
	assert.Equal(t, totalCount, len(cafes))
	// кроме длина слайса кафе, также можно проверить ожидаемый ответ через
	// сравнение строк:
	assert.Equal(t, strings.Join(cafeList["moscow"], ","), body)
}
