package api

import "net/http"

func BadRequest(w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"reason": "Неверный формат запроса или его параметры"}`))
	return w
}

func InternalServerError(w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"reason": "Ошибка сервера"}`))
	return w
}
