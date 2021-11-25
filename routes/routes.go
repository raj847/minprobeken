package routes

import (
	business "minpro_arya/features/admins/bussiness"
	admins "minpro_arya/features/admins/presentation"
	controller "minpro_arya/features/admins/presentation/response"
	bussiness "minpro_arya/features/company/bussiness"
	company "minpro_arya/features/company/presentation"
	response "minpro_arya/features/company/presentation/response"
	busssiness "minpro_arya/features/customer/bussiness"
	customer "minpro_arya/features/customer/presentation"
	responsee "minpro_arya/features/customer/presentation/response"
	middlewareApp "minpro_arya/middleware"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RouteList struct {
	JWTMiddleware  middleware.JWTConfig
	AdminRouter    admins.AdminHandler
	CompanyRouter  company.CompanyHandler
	CustomerRouter customer.CustomerHandler
}

func (cl *RouteList) RouteRegister(e *echo.Echo) {
	// Admins
	admins := e.Group("admins")
	admins.POST("/register", cl.AdminRouter.Register)
	admins.POST("/login", cl.AdminRouter.Login)

	// Company
	company := e.Group("company")
	company.POST("/register", cl.CompanyRouter.Register)
	company.POST("/login", cl.CompanyRouter.Login)

	customer := e.Group("customer")
	customer.POST("/register", cl.CustomerRouter.Register)
	customer.POST("/login", cl.CustomerRouter.Login)
	// customer.GET("/my-product", cl.TransHandler.GetAllCustomerTrans, middleware.JWTWithConfig(cl.JWTMiddleware), RoleValidationCustomer())

}

func RoleValidationAdmin() echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := middlewareApp.GetUser(c)

			if claims.Role == "admin" {
				return hf(c)
			} else {
				return controller.NewErrorResponse(c, http.StatusForbidden, business.ErrUnathorized)
			}
		}
	}
}

func RoleValidationCompany() echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := middlewareApp.GetUser(c)

			if claims.Role == "company" || claims.Role == "admin" {
				return hf(c)
			} else {
				return response.NewErrorResponse(c, http.StatusForbidden, bussiness.ErrUnathorized)
			}
		}
	}
}

func RoleValidationCustomer() echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := middlewareApp.GetUser(c)

			if claims.Role == "customer" || claims.Role == "admin" {
				return hf(c)
			} else {
				return responsee.NewErrorResponse(c, http.StatusForbidden, busssiness.ErrUnathorized)
			}
		}
	}
}
