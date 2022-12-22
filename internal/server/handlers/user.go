package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/AlexeyKluev/user-balance/internal/app"
	"github.com/AlexeyKluev/user-balance/internal/usecase"
)

func NewUserBalanceHandler(resources *app.Resources) http.HandlerFunc {
	type Resp struct {
		Balance string `json:"balance"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
		if err != nil {
			resources.Logger.Info("failed parse id in query string", zap.Error(err))
			http.NotFound(w, r)
			return
		}

		balance, err := resources.UserBalanceUC.Balance(r.Context(), id)
		if err != nil {
			if errors.Is(err, usecase.ErrNotFound) {
				http.NotFound(w, r)
				return
			}
			resources.Logger.Error("failed get user balance", zap.Error(err))
			http.Error(w, http.StatusText(500), 500)
			return
		}

		b := Resp{
			Balance: balance,
		}

		marshal, err := json.Marshal(b)
		if err != nil {
			resources.Logger.Error("failed get user balance", zap.Error(err))
			http.Error(w, http.StatusText(500), 500)
			return
		}

		_, err = w.Write(marshal)
		if err != nil {
			resources.Logger.Error("failed get user balance", zap.Error(err))
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}
}
