package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EmelinDanila/task-manager-api/routes"
	"github.com/EmelinDanila/task-manager-api/tests/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRoutes_Setup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(db)

	routes.SetupRoutes(router, db.GetDB())

	// Тестовые данные
	validUser := map[string]string{
		"email":    "test@example.com",
		"password": "ValidPass123!",
	}
	invalidUser := map[string]string{
		"email":    "",
		"password": "",
	}

	// Функция для отправки JSON-запроса
	sendJSONRequest := func(method, path string, body interface{}) *httptest.ResponseRecorder {
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}

	// ✅ Проверка регистрации
	resp := sendJSONRequest("POST", "/register", invalidUser)
	assert.Equal(t, http.StatusBadRequest, resp.Code, "Invalid registration should return 400")

	resp = sendJSONRequest("POST", "/register", validUser)
	assert.Equal(t, http.StatusCreated, resp.Code, "Valid registration should return 201")

	resp = sendJSONRequest("POST", "/register", validUser)
	assert.Equal(t, http.StatusConflict, resp.Code, "Duplicate registration should return 409")

	// ✅ Проверка логина
	resp = sendJSONRequest("POST", "/login", invalidUser)
	assert.Equal(t, http.StatusBadRequest, resp.Code, "Invalid login should return 400") // Ожидаем `400`

	resp = sendJSONRequest("POST", "/login", map[string]string{"email": "unknown@example.com", "password": "WrongPass123!"})
	assert.Equal(t, http.StatusUnauthorized, resp.Code, "Login with non-existent user should return 401") // Ожидаем `401`

	resp = sendJSONRequest("POST", "/login", validUser)
	assert.Equal(t, http.StatusOK, resp.Code, "Valid login should return 200")
	assert.Contains(t, resp.Body.String(), "token", "Response should contain a token")

	// ✅ Проверка защищенных маршрутов
	protectedRoutes := []struct {
		method string
		path   string
		status int
	}{
		{"POST", "/tasks", http.StatusUnauthorized},
		{"GET", "/tasks", http.StatusUnauthorized},
		{"GET", "/tasks/1", http.StatusUnauthorized},
	}

	for _, test := range protectedRoutes {
		resp := sendJSONRequest(test.method, test.path, nil)
		assert.Equal(t, test.status, resp.Code, "Unauthorized access should return 401")
	}
}
