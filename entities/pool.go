package entities

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"gorm.io/gorm"
)

type Pool struct {
	Id            int            `gorm:"primaryKey"`
	BlockNumber   int            `gorm:"type:integer"`
	Factory       string         `gorm:"type:varchar(42)"`
	Address       string         `gorm:"type:varchar(42)"`
	Token0        string         `gorm:"type:varchar(42)"`
	Token1        string         `gorm:"type:varchar(42)"`
	FactoryCommon common.Address `gorm:"-"`
	AddressCommon common.Address `gorm:"-"`
	Token0Common  common.Address `gorm:"-"`
	Token1Common  common.Address `gorm:"-"`
}

type Reserves struct {
	Reserve0  *uint256.Int
	Reserve1  *uint256.Int
	LastBlock uint64
	LastIndex uint
}

func (e *Pool) AfterFind(_ *gorm.DB) (err error) {
	e.FactoryCommon = common.HexToAddress(e.Factory)
	e.AddressCommon = common.HexToAddress(e.Address)
	e.Token0Common = common.HexToAddress(e.Token0)
	e.Token1Common = common.HexToAddress(e.Token1)

	return
}
