package router

import (
	"Gophermarket/go-musthave-diploma-tpl/internal/handlers/ordersHandler"
	"Gophermarket/go-musthave-diploma-tpl/internal/handlers/userHandler"
	"Gophermarket/go-musthave-diploma-tpl/internal/handlers/withdrawHandler"
	middleware2 "Gophermarket/go-musthave-diploma-tpl/internal/router/middleware"

	"github.com/go-chi/chi/v5"
)

const (
	createUserPath          = "/api/user/register"
	createOrderPath         = "/api/user/orders"
	createWithdrawPath      = "/api/user/balance/withdraw"
	loginUserPath           = "/api/user/login"
	getUserBalancePath      = "/api/user/balance"
	getAllUserWithdrawsPath = "/api/user/withdrawals"
)

func NewRouter(h userHandler.UserHandler, o ordersHandler.OrderHandler,
	w withdrawHandler.WithdrawHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware2.GzipDecompressHandle, middleware2.GzipCompressHandle)
	r.Post(createUserPath, h.Create)
	r.Post(createOrderPath, o.Create)
	r.Post(createWithdrawPath, w.Create)
	r.Post(loginUserPath, h.Login)
	r.Get(createOrderPath, o.GetAll)
	r.Get(getUserBalancePath, w.GetUserBalance)
	r.Get(getAllUserWithdrawsPath, w.GetAll)
	return r
}
