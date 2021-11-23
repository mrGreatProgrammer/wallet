package types

// Money представляет собой денежную сумму в минимальных единицах (центы, копейки, дирамы и т.д.)
type Money int64

// PaymentCategory представляет собой категорию, в которой был совершён платёж (авто, аптеки, рестораны и т.д.)
type PaymentCategory string

// PaymentStatus представляет собой статус пдатежа.
type PaymentStatus string

// Предопределённые статусы платежей
const (
	PaymentStatusOk PaymentStatus = "OK"
	PaymentStatusFail PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

// Payment представляет информацию о платеже.
type Payment struct {
	ID string
	Amount Money
	Category PaymentCategory
	Status PaymentStatus
}

type Phone string

// Account представляет информаию о счёте пользователя
type Account struct {
	ID int64
	Phone Phone
	Balance Money
}