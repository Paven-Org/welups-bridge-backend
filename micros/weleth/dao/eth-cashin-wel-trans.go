package dao

import (
	"bridge/micros/weleth/model"
	"bridge/service-managers/logger"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type IEthCashinWelTransDAO interface {
	CreateEthCashinWelTrans(t *model.EthCashinWelTrans) (int64, error)
	GetUnconfirmedTx2Treasury(from, treasury, token, amount string) (*model.TxToTreasury, error)

	CreateTx2Treasury(t *model.TxToTreasury) error

	UpdateEthCashinWelTx(t *model.EthCashinWelTrans) error

	SelectTransByDepositTxHash(txHash string) (*model.EthCashinWelTrans, error)
	SelectTransByIssueTxHash(txHash string) ([]*model.EthCashinWelTrans, error)
	SelectTransById(id string) (*model.EthCashinWelTrans, error)
}

// sort of a locator for DAOs
type ethCashinWelTransDAO struct {
	db *sqlx.DB
}

func (w *ethCashinWelTransDAO) CreateTx2Treasury(t *model.TxToTreasury) error {
	db := w.db
	log := logger.Get()

	q := db.Rebind(`INSERT INTO tx_to_treasury(
									tx_id,
									from_address,
									treasury_address,
									token_address,
									amount,
									tx_fee,
									status) VALUES (?,?,?,?,?,?,?)`)
	_, err := db.Exec(q, t.TxID, t.FromAddress, t.TreasuryAddr, t.TokenAddr, t.Amount, t.TxFee, t.Status)

	if err != nil {
		log.Err(err).Msgf("Error while inserting tx to treasury %s", t.TxID)
		return err
	}
	return nil
}

func (w *ethCashinWelTransDAO) GetUnconfirmedTx2Treasury(from, treasury, token, amount string) (*model.TxToTreasury, error) {
	db := w.db
	log := logger.Get()

	var res model.TxToTreasury
	q := db.Rebind(
		`SELECT * FROM tx_to_treasury
			WHERE from_address = ?,
						treasury_address = ?,
						token_address = ?,
						amount = ?
						status = 'unconfirmed'
			ORDER BY created_at DESC`)

	err := db.Get(&res, q, from, treasury, token, amount)
	if err == sql.ErrNoRows {
		log.Info().Msg("[GetUnconfirmedTx2Treasury] no tx found")
		return nil, nil
	}
	if err != nil {
		log.Err(err).Msg("[GetUnconfirmedTx2Treasury] error while querying DB")
		return nil, err
	}

	return &res, nil
}

func (w *ethCashinWelTransDAO) CreateEthCashinWelTrans(t *model.EthCashinWelTrans) (int64, error) {
	log := logger.Get()
	tx, err := w.db.Beginx()
	if err != nil {
		log.Err(err).Msg("Can't start transaction")
		return -1, err
	}

	qCreate := tx.
		Rebind(
			`INSERT INTO eth_cashin_wel_trans(
				eth_tx_hash,
				wel_issue_tx_hash,
				eth_token_addr,
				wel_token_addr,
				network_id,
				eth_wallet_addr,
				wel_wallet_addr,
				amount,
				commnission_fee,
				status) VALUES (?,?,?,?,?,?,?,?,?,?) RETURNING id`)
	var id int64
	err = tx.
		Get(&id,
			qCreate,
			t.EthTxHash,
			t.WelIssueTxHash,
			t.EthTokenAddr,
			t.WelTokenAddr,
			t.NetworkID,
			t.EthWalletAddr,
			t.WelWalletAddr,
			t.Amount,
			t.CommissionFee,
			t.Status)

	if err != nil {
		log.Err(err).Msgf("Error while inserting EthCashinWel tx with eth tx hash %s", t.EthTxHash)
		tx.Rollback()
		return id, err
	}

	qUpdateTx2Treasury := tx.Rebind(`UPDATE tx_to_treasury SET status='isCashin' WHERE tx_id = ?`)
	_, err = tx.Exec(qUpdateTx2Treasury, t.EthTxHash)
	if err != nil {
		log.Err(err).Msgf("Error while inserting EthCashinWel tx with eth tx hash %s", t.EthTxHash)
		tx.Rollback()
		return id, err
	}
	tx.Commit()

	return id, nil
}

func (w *ethCashinWelTransDAO) UpdateEthCashinWelTx(t *model.EthCashinWelTrans) error {
	db := w.db
	log := logger.Get()

	q := db.
		Rebind(
			`UPDATE eth_cashin_wel_trans SET
		    eth_tx_hash = ?,
		    wel_issue_tx_hash = ?,
		    eth_token_addr = ?,
		    wel_token_addr = ?,
		    network_id = ?,
		    eth_wallet_addr = ?,
		    wel_wallet_addr = ?,
		    amount = ?,
		    commnission_fee = ?,
		    status = ? 
		    WHERE id = ?`)
	_, err := db.
		Exec(q,
			t.EthTxHash,
			t.WelIssueTxHash,
			t.EthTokenAddr,
			t.WelTokenAddr,
			t.NetworkID,
			t.EthWalletAddr,
			t.WelWalletAddr,
			t.Amount,
			t.CommissionFee,
			t.Status,
			t.ID)

	if err != nil {
		log.Err(err).Msgf("Error while inserting EthCashinWel tx with eth tx hash %s", t.EthTxHash)
		return err
	}
	return nil
}

func (w *ethCashinWelTransDAO) SelectTransByDepositTxHash(txHash string) (*model.EthCashinWelTrans, error) {
	var t = &model.EthCashinWelTrans{}
	err := w.db.Get(t, "SELECT * FROM wel_cashin_eth_trans WHERE eth_tx_hash = $1", txHash)
	return t, err
}

func (w *ethCashinWelTransDAO) SelectTransByIssueTxHash(txHash string) ([]*model.EthCashinWelTrans, error) {
	var txs = []*model.EthCashinWelTrans{}
	err := w.db.Select(txs, "SELECT * FROM wel_cashin_eth_trans WHERE wel_issue_tx_hash = $1", txHash)
	return txs, err
}

func (w *ethCashinWelTransDAO) SelectTransById(id string) (*model.EthCashinWelTrans, error) {
	var t = &model.EthCashinWelTrans{}
	err := w.db.Get(t, "SELECT * FROM wel_cashin_eth_trans WHERE id = $1", id)
	return t, err
}

func MkEthCashinWelTransDao(db *sqlx.DB) *ethCashinWelTransDAO {
	return &ethCashinWelTransDAO{
		db: db,
	}
}