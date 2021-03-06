package data

import (
	"minpro_arya/features/admins"

	"minpro_arya/features/admins/bussiness"

	"gorm.io/gorm"
)

type MysqlAdminRepository struct {
	Conn *gorm.DB
}

func NewMysqlAdminRepository(conn *gorm.DB) admins.Repository {
	return &MysqlAdminRepository{
		Conn: conn,
	}
}

func (rep *MysqlAdminRepository) Register(domain *admins.Domain) (admins.Domain, error) {

	admin := fromDomain(*domain)

	result := rep.Conn.Create(&admin)
	if result.Error != nil {
		return admins.Domain{}, result.Error
	}

	return toDomain(admin), nil
}

func (rep *MysqlAdminRepository) Login(username, password string) (admins.Domain, error) {
	var admin Admins
	err := rep.Conn.First(&admin, "username = ?", username).Error

	if err != nil {
		return admins.Domain{}, bussiness.ErrEmailorPass
	}

	return toDomain(admin), nil
}
