package dao

import (
	"github.com/jmoiron/sqlx"
)

// sort of a locator for DAOs
type DAOs struct {
	WelCashinEthTransDAO  IWelCashinEthTransDAO
	EthCashoutWelTransDAO IEthCashoutWelTransDAO
	EthCashinWelTransDAO  IEthCashinWelTransDAO
	WelCashoutEthTransDAO IWelCashoutEthTransDAO
	EthSysDAO             *ethSysDAO
	WelSysDAO             *welSysDAO
}

func MkDAOs(db *sqlx.DB) *DAOs {
	return &DAOs{
		WelCashinEthTransDAO:  MkWelCashinEthTransDao(db),
		EthCashoutWelTransDAO: MkEthCashoutWelTransDao(db),
		EthCashinWelTransDAO:  MkEthCashinWelTransDao(db),
		WelCashoutEthTransDAO: MkWelCashoutEthTransDao(db),
		EthSysDAO:             MkEthSysDao(db),
		WelSysDAO:             MkWelSysDao(db)}
}
