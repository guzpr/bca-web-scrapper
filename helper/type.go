package helper

type Transaction struct {
	Date                   string            `json:"date"`
	Method                 TransactionMethod `json:"method"`
	Description            string            `json:"description"`
	DestinationAccountName string            `json:"destination_account_name"`
	Amount                 string            `json:"amount"`
	Type                   string            `json:"type"`
}

type TransactionMethod string

const (
	TransactionMethodTransferEBanking = TransactionMethod("TRSF E-BANKING DB")
	TransactionMethodTarikanATM       = TransactionMethod("TARIKAN ATM")
	TransactionMethodSwitchingCR      = TransactionMethod("SWITCHING CR")
)

type TransactionMetadata string

const (
	TransactionMetadataFTFVA = TransactionMetadata("FTFVA")
	TransactionMetadataFTSCY = TransactionMetadata("FTSCY")
)
