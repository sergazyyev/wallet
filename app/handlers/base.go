package handlers

import (
	"encoding/json"
	"github.com/sergazyyev/wallet/app/errs"
	"github.com/sergazyyev/wallet/app/validators"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func (s *Server) validateRequest(req *http.Request, data validators.JSONRequestValidator) error {
	v := govalidator.New(govalidator.Options{
		Request:         req,
		Rules:           data.Rules(),
		Messages:        data.Messages(),
		Data:            data,
		RequiredDefault: true,
	})
	res := v.ValidateJSON()

	if len(res) > 0 {
		return validators.NewValidationError(res)
	}

	return nil
}

func (s *Server) errorResponse(w http.ResponseWriter, e error) {
	var err error
	switch t := e.(type) {
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(e.Error()))
	case *validators.ValidationError:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(t.Messages())
	case *errs.CustomError:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(t.Status())
		err = json.NewEncoder(w).Encode(map[string]string{"message": t.Error()})
	}

	//consider cant write to network and panics in this goroutine
	if err != nil {
		panic(err)
	}
}

func (s *Server) response(w http.ResponseWriter, payload interface{}, contentType string) {
	var err error
	switch t := payload.(type) {
	default:
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(payload)
	case []byte:
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		} else {
			w.Header().Set("Content-Type", http.DetectContentType(t))
		}
		_, err = w.Write(t)
	}

	//consider cant write to network and panics in this goroutine
	if err != nil {
		panic(err)
	}
}
