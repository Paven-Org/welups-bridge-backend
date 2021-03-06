package msweleth

import (
	welethService "bridge/micros/weleth/temporal"
)

const (
	TaskQueue = welethService.WelethServiceQueue

	GetWelToEthCashinByTxHash  = "GetWelToEthCashinByTxHashWF"
	GetEthToWelCashoutByTxHash = "GetEthToWelCashoutByTxHashWF"
	GetWelToEthCashin          = "GetWelToEthCashinWF"
	GetEthToWelCashout         = "GetEthToWelCashoutWF"

	GetWelToEthCashinClaimRequest  = "GetWelToEthCashinClaimRequestWF"
	GetEthToWelCashoutClaimRequest = "GetEthToWelCashoutClaimRequestWF"

	GetEthToWelCashinByTxHash  = "GetEthToWelCashinByTxHashWF"
	GetWelToEthCashoutByTxHash = "GetWelToEthCashoutByTxHashWF"
	GetEthToWelCashin          = "GetEthToWelCashinWF"
	GetWelToEthCashout         = "GetWelToEthCashoutWF"

	GetTx2TreasuryBySender = "GetTx2TreasuryBySenderWF"

	CreateW2ECashinClaimRequestWF  = "CreateW2ECashinClaimRequestWF"
	CreateE2WCashoutClaimRequestWF = "CreateE2WCashoutClaimRequestWF"

	WaitForPendingW2ECashinClaimRequestWF  = "WaitForPendingW2ECashinClaimRequestWF"
	WaitForPendingE2WCashoutClaimRequestWF = "WaitForPendingE2WCashoutClaimRequestWF"
)
