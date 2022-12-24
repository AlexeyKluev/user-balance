package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/AlexeyKluev/user-balance/internal/app"
	"github.com/AlexeyKluev/user-balance/internal/domain/dto"
	"github.com/AlexeyKluev/user-balance/internal/usecase"
)

type userBalanceResp struct {
	Balance string `json:"balance"`
}

// NewUserBalanceHandler godoc
// @Summary      Баланс пользователя
// @Description  Возвращает баланс пользователя по id
// @Tags         Balance
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  userBalanceResp
// @Failure      404
// @Failure      500
// @Router       /users/{id}/balance [get]
func NewUserBalanceHandler(resources *app.Resources) http.HandlerFunc {

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

		b := userBalanceResp{
			Balance: balance,
		}

		marshal, err := json.Marshal(b)
		if err != nil {
			resources.Logger.Error("failed get user balance", zap.Error(err))
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(marshal)
		if err != nil {
			resources.Logger.Error("failed get user balance", zap.Error(err))
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}
}

type accuralFundsReq struct {
	Amount int64 `json:"amount" validate:"required,gte=0"`
}

// NewAccrualFundsHandler godoc
// @Summary      Зачисление средств на баланс
// @Description  Добавляет средства на баланс пользователя
// @Tags         Balance
// @Accept       json
// @Produce      json
// @Param        id      path      int             true  "User ID"
// @Param 		 request body      accuralFundsReq true "123"
// @Success      201
// @Failure      403
// @Failure      404
// @Failure      422
// @Failure      500
// @Router       /users/{id}/accural [post]
func NewAccrualFundsHandler(resources *app.Resources) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
		if err != nil {
			resources.Logger.Info("failed parse id in query string", zap.Error(err))
			http.NotFound(w, r)
			return
		}

		var req accuralFundsReq
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			resources.Logger.Info("failed parse body", zap.Error(err))
			http.Error(w, http.StatusText(422), 422)
			return
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			resources.Logger.Info("validate failed", zap.Error(err))
			http.Error(w, http.StatusText(422), 422)
			return
		}

		err = resources.AccuralFundsUC.Accural(r.Context(), dto.AccuralDTO{
			UserID: id,
			Amount: req.Amount,
		})
		if err != nil {
			if errors.Is(err, usecase.ErrUserIsBanned) {
				http.Error(w, http.StatusText(403), 403)
				return
			}
			if errors.Is(err, usecase.ErrNotFound) {
				http.Error(w, http.StatusText(404), 404)
				return
			}
			resources.Logger.Error("failed accural user balance", zap.Error(err))
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	}
}

type reservationFundsReq struct {
	ServiceID int64 `json:"service_id" validate:"required,gte=0"`
	OrderID   int64 `json:"order_id" validate:"required,gte=0"`
	Amount    int64 `json:"amount" validate:"required,gte=0"`
}

// NewReservationFundsHandler godoc
// @Summary      Резервирование средств
// @Description  Резервирует средства с баланса пользователя
// @Tags         Balance
// @Accept       json
// @Produce      json
// @Param        id      path      int             	   true  "User ID"
// @Param 		 request body      reservationFundsReq true "123"
// @Success      201
// @Failure      400
// @Failure      403
// @Failure      404
// @Failure      422
// @Failure      500
// @Router       /users/{id}/reservation [post]
func NewReservationFundsHandler(resources *app.Resources) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
		if err != nil {
			resources.Logger.Info("failed parse id in query string", zap.Error(err))
			http.NotFound(w, r)
			return
		}

		var req reservationFundsReq
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			resources.Logger.Info("failed parse body", zap.Error(err))
			http.Error(w, http.StatusText(422), 422)
			return
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			resources.Logger.Info("validate failed", zap.Error(err))
			http.Error(w, http.StatusText(422), 422)
			return
		}

		if err := resources.ReservationFundsUC.Reservation(r.Context(), dto.ReservationDTO{
			UserID:    id,
			OrderID:   req.OrderID,
			ServiceID: req.ServiceID,
			Amount:    req.Amount,
		}); err != nil {
			if errors.Is(err, usecase.ErrUserIsBanned) {
				http.Error(w, http.StatusText(403), 403)
				return
			}
			if errors.Is(err, usecase.ErrNotFound) {
				http.Error(w, http.StatusText(404), 404)
				return
			}
			if errors.Is(err, usecase.ErrInsufficientBalance) {
				http.Error(w, http.StatusText(400), 400)
				return
			}
			resources.Logger.Error("failure reservation balance", zap.Error(err))
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	}
}
