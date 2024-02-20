package httpHandler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/EwvwGeN/authService/internal/verification"
	"github.com/gorilla/mux"
)

type httpFunc func(http.ResponseWriter, *http.Request)

func (s *Server) confirmUser() httpFunc {
	log := s.log.With(slog.String("handler", "confirmUser"))
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("start confirm user")
		vars := mux.Vars(r)
		log.Debug("get vars", slog.Any("map", vars))
		code, ok := vars["verificationCode"]
		if !ok {
			http.Error(w, ErrInvalidArgument.Error(), http.StatusBadRequest)
			return
		}
		userId, err := verification.DecryptVerificationCode(code)
		if err != nil {
			http.Error(w, fmt.Errorf("cant decode verification code: %w", err).Error(), http.StatusInternalServerError)
		}
		user, err := s.usrProvider.GetUserById(r.Context(), userId)
		if err != nil {
			http.Error(w, fmt.Errorf("cant finde user: %w", err).Error(), http.StatusNotFound)
			return
		}
		err = s.confirmator.ConfirmUser(r.Context(), user.Id)
		if err != nil {
			http.Error(w, fmt.Errorf("cant confirm user: %w", err).Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, map[string]string{
			"message": "user confirmed",
		})
	}
}