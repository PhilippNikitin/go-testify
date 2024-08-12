package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Функция urlForTest возвращает указатель на структуру типа url.URL,
// сформированную на основе полученных аргументов и ошибку, если что-то пойдет
// не так во время выполнения функции.

// Была добавлена, чтобы не повторять один и тот же код в каждой тестовой функции
// в соответствии с принципом DRY

// Параметры
// baseUrl - базовый URL, к которому будет добавлен эндпоинт и параметры запроса
// endpoint - эндпоинт, который мы добавим к базовому URL
// count - количество кафе, которое мы хотим получить
// city - город, для которого мы запрашиваем информацию

func urlForTest(baseUrl, endpoint, count, city string) (*url.URL, error) {
	targetUrl, err := url.Parse(baseUrl)

	// если возникла ошибка, возвращаем nil (нулевое значение для указателя) и ошибку
	if err != nil {
		return nil, err
	}

	// добавляем эндпоинт, который был передан в качестве аргумента
	targetUrl.Path += endpoint

	// Подготавливаем параметры запроса
	params := url.Values{}
	params.Add("count", count)
	params.Add("city", city)

	// добавляем параметры запроса к базовому URL при помощи params.Encode()
	targetUrl.RawQuery = params.Encode()

	return targetUrl, nil
}

func TestMainHandlerWhenOk(t *testing.T) {
	// создаем URL-адрес для запроса при помощи urlForTest
	targetUrl, err := urlForTest("http://localhost:8080/", "cafe", "2", "moscow")
	if err != nil {
		fmt.Println("MalformedURL:", err)
	}
	// создаем новый запрос
	req := httptest.NewRequest("GET", targetUrl.String(), nil)

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
	// создаем URL-адрес для запроса при помощи urlForTest
	targetUrl, err := urlForTest("http://localhost:8080/", "cafe", "2", "london")
	if err != nil {
		fmt.Println("MalformedURL:", err)
	}
	// создаем новый запрос
	req := httptest.NewRequest("GET", targetUrl.String(), nil)

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
	// создаем URL-адрес для запроса при помощи urlForTest
	targetUrl, err := urlForTest("http://localhost:8080/", "cafe", "7", "moscow")
	if err != nil {
		fmt.Println("MalformedURL:", err)
	}
	// создаем новый запрос
	req := httptest.NewRequest("GET", targetUrl.String(), nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	body := responseRecorder.Body.String()
	cafes := strings.Split(body, ",")

	// проверяем, что вернулся код ответа 200
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	// также хотел вместо assert.Equal использовать метод пакета assert HTTPSuccess:
	// assert.HTTPSuccess(t, handler, "GET", "/cafe?count=7&city=moscow", nil)
	// но тест возврашал значение FAIL - получался код ответа 400, а не 200
	assert.Equal(t, totalCount, len(cafes))
	// кроме длина слайса кафе, также можно проверить ожидаемый ответ через
	// сравнение строк:
	assert.Equal(t, strings.Join(cafeList["moscow"], ","), body)
}
