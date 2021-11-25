// page 26 лк 14
package wallet

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/mrGreatProgrammer/wallet/pkg/types"
)

func TestService_FindAccountByID_success(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")

	expected := &types.Account{
		ID: 1,
		Phone: "+992000000001" ,
		Balance: 0,
	}

	result, _ := svc.FindAccountByID(1)
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, result)
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")

	var nilOfTheStruct *types.Account
	expected := nilOfTheStruct


	result, _ := svc.FindAccountByID(3)
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, result)
	}
}

func TestService_Reject_success(t *testing.T) {
	// создаём сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем отменить платёж
	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): cant find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didnt changed, payment = %v", savedPayment)
		return
	}
	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): cant find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(): balance didnt changed, account = %v", savedAccount)
		return
	}
}

func TestService_FindPaymentByID_seccess(t *testing.T) {
	// создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем найти платёж
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	// сравниваем платежи
	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	// созаём сервис
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем найти несуществуюший платёж
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FindPaymentByID(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}

func TestService_Repeat_success(t *testing.T) {
	// создаём сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем повторить платёж
	payment := payments[0]
	rep, err := s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}
	
	if !reflect.DeepEqual(payment, rep) {
		t.Errorf("invalid result, expected: %v, actual: %v, error", payment, rep)
		return
	}
}

func TestService_Repeat_fail(t *testing.T) {
	// создаём сервис
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем повторить несуществующий платёж
	_, err = s.Repeat("hi")
	if err == nil {
		t.Errorf("Repeat(): must return error, returned nil = %v", err)
		return
	}
	
	if err != ErrPaymentNotFound {
		t.Errorf("Repeat(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}