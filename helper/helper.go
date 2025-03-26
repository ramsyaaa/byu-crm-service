package helper

import "github.com/go-playground/validator/v10"

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func ErrorValidationFormat(err error, validationMessages map[string]string) map[string]string {
	errors := make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		// Buat key yang sesuai dengan field dan tag error
		key := e.Field() + "." + e.Tag()
		if message, exists := validationMessages[key]; exists {
			errors[e.Field()] = message
		} else {
			errors[e.Field()] = e.Error()
		}
	}

	return errors
}
