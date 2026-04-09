package constant

// TransactionStatus represents the status of a transaction.
type TransactionStatus string

const (
	TransactionPending    TransactionStatus = "pending"
	TransactionProcessing TransactionStatus = "processing"
	TransactionSuccess    TransactionStatus = "success"
	TransactionFailed     TransactionStatus = "failed"
	TransactionExpired    TransactionStatus = "expired"
	TransactionRefunded   TransactionStatus = "refunded"
)
