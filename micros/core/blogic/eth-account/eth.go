package ethLogic

import (
	"bridge/libs"
	"bridge/micros/core/model"
	ethService "bridge/micros/core/service/eth"
	"context"

	"github.com/ethereum/go-ethereum/crypto"
	"go.temporal.io/sdk/client"
)

//AddEthAccount(address string, status string)
func AddEthAccount(address, status string) error {
	log.Info().Msgf("[ethAccount logic internal] Creating ethAccount %s...", address)

	if !verifyAddress(address) {
		err := model.ErrEthInvalidAddress
		log.Err(err).Msgf("[ethAccount logic internal] Address %s invalid", address)
		return err
	}

	err := ethDAO.AddEthAccount(address, status)
	if err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] Failed to create ethAccount %s", address)
		return err
	}
	return nil
}

//RemoveEthAccount(address string)
func RemoveEthAccount(address string) error {
	log.Info().Msgf("[ethAccount logic internal] Start removing ethAccount %s...", address)
	_, err := ethDAO.GetEthAccount(address)
	if err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] Failed to retrieve ethAccount %s", address)
		return err
	}

	log.Info().Msgf("[ethAccount logic internal] Removing ethAccount %s...", address)
	if err := ethDAO.RemoveEthAccount(address); err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] Failed to remove ethAccount %s", address)
		return err
	}

	return nil
}

//GetEthAccount(address string)
func GetEthAccount(address string) (*model.EthAccount, error) {
	log.Info().Msgf("[EthAccount logic internal] Getting EthAccount %s", address)
	ethAccount, err := ethDAO.GetEthAccount(address) // should eventually get by ID instead, but this is more convenient for now
	if err != nil {
		log.Err(err).Msgf("[EthAccount logic internal] Failed to retrieve EthAccount %s's info", address)
		return nil, err
	}

	return ethAccount, nil
}

//GetAllEthAccounts(offset uint, size uint)
func GetAllEthAccounts(offset uint, size uint) ([]model.EthAccount, error) {
	log.Info().Msgf("[ethAccount logic internal] Getting ethAccounts...")
	ethAccounts, err := ethDAO.GetAllEthAccounts(offset, size)
	if err != nil {
		log.Err(err).Msg("[ethAccount logic internal] Failed to retrieve ethAccounts")
		return nil, err
	}

	return ethAccounts, nil
}

//GetAllRoles()
func GetAllRoles() ([]string, error) {
	log.Info().Msgf("[ethAccount logic internal] Getting all roles...")
	roles, err := ethDAO.GetAllRoles()
	if err != nil {
		log.Err(err).Msg("[ethAccount logic internal] Failed to retrieve roles")
		return nil, err
	}

	return roles, nil
}

//GetEthAccountRoles(address string)
func GetEthAccountRoles(address string) ([]string, error) {
	log.Info().Msgf("[EthAccount logic] Getting EthAccount %s's roles...", address)
	roles, err := ethDAO.GetEthAccountRoles(address)
	if err != nil {
		log.Err(err).Msgf("[EthAccount logic] Failed to retrieve EthAccount %s's roles", address)
		return nil, err
	}

	return roles, nil
}

//GetEthAccountsWithRole(role string, offset uint, size uint)
func GetEthAccountsWithRole(role string, offset uint, size uint) ([]model.EthAccount, error) {
	log.Info().Msgf("[ethAccount logic internal] Getting ethAccounts with role %s...", role)
	ethAccounts, err := ethDAO.GetEthAccountsWithRole(role, offset, size)
	if err != nil {
		log.Err(err).Msg("[ethAccount logic internal] Failed to retrieve ethAccounts with role " + role)
		return nil, err
	}

	return ethAccounts, nil
}

//GrantRole(address string, role string)
func GrantRole(address, role string, callerkey string) (string, error) {
	log.Info().Msgf("[ethAccount logic internal] Start granting role %s to ethAccount %s...", role, address)

	key, err := crypto.HexToECDSA(callerkey)
	if err != nil {
		log.Err(err).Msg("[ethAccount logic internal] Invalid private key")
		return "", err // invalid key
	}

	callerAddress := crypto.PubkeyToAddress(key.PublicKey).Hex()
	log.Info().Msgf("[ethAccount logic internal] caller address: %s", callerAddress)

	acc, err := ethDAO.GetEthAccount(callerAddress)
	if err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] Unable to get ethereum account %s from DB", callerAddress)
		return "", err // invalid key
	}
	if acc.Status != "ok" {
		log.Info().Msgf("[ethAccount logic internal] account %s locked and cannot do administrative tasks", callerAddress)
		return "", model.ErrEthAccountLocked
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", callerkey)
	// call workflow
	log.Info().Msgf("[ethAccount logic internal] Calling GovContractService workflow...")
	wo := client.StartWorkflowOptions{
		TaskQueue: ethService.GovContractQueue,
	}

	we, err := tempcli.ExecuteWorkflow(ctx, wo, ethService.GrantRoleWorkflow, address, role)
	if err != nil {
		log.Err(err).Msg("[ethAccount logic internal] Unable to call GrantRoleWorkflow")
		return "", err
	}
	log.Info().Str("Workflow", we.GetID()).Str("runID=", we.GetRunID()).Msg("dispatched")

	var txhash string
	if err := we.Get(ctx, &txhash); err != nil {
		log.Err(err).Msg("[ethAccount logic internal] GrantRoleWorkflow failed")
		return txhash, err
	}
	return txhash, nil
}

//RevokeRole(address string, role string)
func RevokeRole(address, role string, callerkey string) (string, error) {
	log.Info().Msgf("[ethAccount logic internal] Start revoking role %s to ethAccount %s...", role, address)

	key, err := crypto.HexToECDSA(callerkey)
	if err != nil {
		log.Err(err).Msg("[ethAccount logic internal] Invalid private key")
		return "", err // invalid key
	}

	callerAddress := crypto.PubkeyToAddress(key.PublicKey).Hex()
	log.Info().Msgf("[ethAccount logic internal] caller address: %s", callerAddress)

	acc, err := ethDAO.GetEthAccount(callerAddress)
	if err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] Unable to get ethereum account %s from DB", callerAddress)
		return "", err // invalid key
	}
	if acc.Status != "ok" {
		log.Info().Msgf("[ethAccount logic internal] account %s locked and cannot do administrative tasks", callerAddress)
		return "", model.ErrEthAccountLocked
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "callerkey", callerkey)
	// call workflow
	log.Info().Msgf("[ethAccount logic internal] Calling GovContractService workflow...")
	wo := client.StartWorkflowOptions{
		TaskQueue: ethService.GovContractQueue,
	}
	we, err := tempcli.ExecuteWorkflow(ctx, wo, ethService.RevokeRoleWorkflow, address, role)
	if err != nil {
		log.Err(err).Msg("[ethAccount logic internal] Unable to call RevokeRoleWorkflow")
		return "", err
	}
	log.Info().Str("Workflow", we.GetID()).Str("runID=", we.GetRunID()).Msg("dispatched")

	var txhash string
	if err := we.Get(ctx, &txhash); err != nil {
		log.Err(err).Msg("[ethAccount logic internal] RevokeRoleWorkflow failed")
		return txhash, err
	}
	return txhash, nil
}

//SetEthAccountStatus(address string, status string)
func SetEthAccountStatus(address, status string) error {
	log.Info().Msgf("[ethAccount logic internal] Start setting status %s to ethAccount %s...", status, address)
	//ethAccount, err := ethDAO.GetEthAccount(address)
	//if err != nil {
	//	log.Err(err).Msgf("[ethAccount logic internal] Failed to retrieve ethAccount %s", address)
	//	return err
	//}

	log.Info().Msgf("[ethAccount logic internal] setting status %s to ethAccount %s...", status, address)
	if err := ethDAO.SetEthAccountStatus(address, status); err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] Failed to set status %s to ethAccount %s", status, address)
		return err
	}

	return nil
}

// system keys

func SetCurrentAuthenticator(prikey string) error {
	sysAccounts.Lock()
	defer sysAccounts.Unlock()

	key, err := crypto.HexToECDSA(prikey)
	if err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] invalid private key")
		return err
	}

	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	log.Info().Msgf("[ethAccount logic internal] Set current authenticator to %s", address)

	accs, err := ethDAO.GetEthAccountsWithRole(model.EthAccountRoleAuthenticator, 0, 1000) // should've made the DAO to branch out queries instead, but deadline
	if err != nil {
		log.Err(err).Msgf("[ethAccount logic internal] couldn't retrieve authenticator accounts")
		return err
	}
	match := libs.DropWhile(func(a model.EthAccount) bool { return a.Address != address }, accs)
	if len(match) < 1 {
		err = model.ErrEthAccountNotFound
		log.Err(err).Msgf("[ethAccount logic internal] authenticator %s not found", address)
		return err
	}

	if match[0].Status != "ok" {
		err = model.ErrEthAccountLocked
		log.Err(err).Msgf("[ethAccount logic internal] authenticator %s is locked", address)
		return err
	}

	sysAccounts.authenticator.Address = address
	sysAccounts.authenticator.Prikey = prikey
	sysAccounts.authenticator.Status = match[0].Status
	return nil
}

//GetEthPrikeyIfExists(address string)

//SetPriKey(address string, key string)

//UnsetPrikey(address string)
