package models

import "gorm.io/gorm"

type Tenant struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
}

type Seller struct {
	gorm.Model
	TenantId uint   `json:"tenantid" gorm:"not null"`
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"not null;unique"`
	Phone    string `json:"phone" gorm:"not null;unique"`
	Address  string `json:"address" gorm:"not null"`
	Tenant   Tenant `gorm:"foreignKey:TenantId;references:ID"`
}

type Hub struct {
	gorm.Model
	TenantId uint   `json:"tenantid" gorm:"not null"`
	Name     string `json:"name" gorm:"not null"`
	Address  string `json:"address" gorm:"not null"`
	Tenant   Tenant `gorm:"foreignKey:TenantId;references:ID"`
}

type Sku struct {
	gorm.Model
	TenantId uint   `json:"prodid" gorm:"not null"`
	Name     string `json:"name" gorm:"not null"`
	SkuCode  string `json:"skucode" gorm:"not null"`
	Desc     string `json:"desc" gorm:"not null"`
	Price    uint   `json:"price" gorm:"not null"`
	Tenant   Tenant `gorm:"foreignKey:TenantId;references:ID"`
}

type Inventory struct {
	gorm.Model
	HubId uint `json:"hubid" gorm:"not null"`
	SkuId uint `json:"skuid" gorm:"not null"`
	Qty   uint `json:"qty" gorm:"not null"`
	Hub   Hub  `gorm:"foreignKey:HubId;references:ID"`
	Sku   Sku  `gorm:"foreignKey:SkuId;references:ID"`
}
