package main

import (
	"log"

	_routes "minpro_arya/routes"

	_adminService "minpro_arya/features/admins/bussiness"
	_adminRepo "minpro_arya/features/admins/data"
	_adminController "minpro_arya/features/admins/presentation"

	_companyService "minpro_arya/features/company/bussiness"
	_companyRepo "minpro_arya/features/company/data"
	_companyController "minpro_arya/features/company/presentation"

	_customerService "minpro_arya/features/customer/bussiness"
	_customerRepo "minpro_arya/features/customer/data"
	_customerController "minpro_arya/features/customer/presentation"

	_productService "minpro_arya/features/product/bussiness"
	_productRepo "minpro_arya/features/product/data"
	_productController "minpro_arya/features/product/presentation"

	_transService "minpro_arya/features/transactions/bussiness"
	_transRepo "minpro_arya/features/transactions/data"
	_transController "minpro_arya/features/transactions/presentation"

	_newsRepo "minpro_arya/features/news/data"
	_newsController "minpro_arya/features/news/presentation"

	_dbDriver "minpro_arya/config"

	_driverFactory "minpro_arya/drivers"

	_middleware "minpro_arya/middleware"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func dbMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&_adminRepo.Admins{},
		&_companyRepo.Company{},
		&_customerRepo.Customer{},
		&_productRepo.Product{},
		&_transRepo.Transactions{},
		&_newsRepo.Articles{},
	)
}

func main() {
	configDB := _dbDriver.ConfigDB{
		DB_Username: viper.GetString(`database.user`),
		DB_Password: viper.GetString(`database.pass`),
		DB_Host:     viper.GetString(`database.host`),
		DB_Port:     viper.GetString(`database.port`),
		DB_Database: viper.GetString(`database.name`),
	}
	db := configDB.InitDB()
	dbMigrate(db)

	configJWT := _middleware.ConfigJWT{
		SecretJWT:       viper.GetString(`jwt.secret`),
		ExpiresDuration: int64(viper.GetInt(`jwt.expired`)),
	}

	e := echo.New()

	adminRepo := _driverFactory.NewAdminRepository(db)
	adminService := _adminService.NewServiceAdmin(adminRepo, 10, &configJWT)
	adminCtrl := _adminController.NewHandlerAdmin(adminService)

	companyRepo := _driverFactory.NewCompanyRepository(db)
	companyService := _companyService.NewServiceCompany(companyRepo, 10, &configJWT)
	companyCtrl := _companyController.NewHandlerCompany(companyService)

	customerRepo := _driverFactory.NewCustomerRepository(db)
	customerService := _customerService.NewServiceCustomer(customerRepo, 10, &configJWT)
	customerCtrl := _customerController.NewHandlerCustomer(customerService)

	productRepo := _driverFactory.NewProductRepository(db)
	productService := _productService.NewServiceProduct(productRepo)
	productCtrl := _productController.NewHandlerProduct(productService)

	transRepo := _driverFactory.NewTransRepository(db)
	transService := _transService.NewServiceTrans(transRepo)
	transCtrl := _transController.NewHandlerProduct(transService)

	newsRepo := _newsRepo.NewNewsApi()
	newsCtrl := _newsController.NewNewsHandler(newsRepo)

	routesInit := _routes.RouteList{
		JWTMiddleware:  configJWT.Init(),
		AdminRouter:    *adminCtrl,
		CompanyRouter:  *companyCtrl,
		CustomerRouter: *customerCtrl,
		ProductRouter:  *productCtrl,
		TransRouter:    *transCtrl,
		NewsRouter:     *newsCtrl,
	}

	routesInit.RouteRegister(e)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
