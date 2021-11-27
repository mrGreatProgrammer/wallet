// page 26 лк 14
package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/mrGreatProgrammer/wallet/pkg/types"
)

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone:   "+992000000001",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	// регисрируем там пользователя
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("cant register account, error = %v", err)
	}

	// пополняем счёт
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("cant deposity account, error = %v", err)
	}

	// выполняем платежи
	// можем создать слайс сразу нужной длины, поскольку знаем размер
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		// тогда здесь работаем просто через index, а не через append
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("cant make, error = %v", err)
		}
	}

	return account, payments, nil
}

func TestService_FindAccountByID_success(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")

	expected := &types.Account{
		ID:      1,
		Phone:   "+992000000001",
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
		t.Errorf("Repeat(): invalid data in repeated, payment = &amp;%v", rep)
		return
	}

	if !reflect.DeepEqual(payment, rep) {
		t.Errorf("Repeat(): invalid data in repeated = %v, payment = %v ", rep, payment)
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

func TestService_FavoriteFromPayment_success(t *testing.T) {
	// создаём сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	addFavorite, err := s.FavoritePayment(payment.ID, "megafon")
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Errorf("FavoritePaymetn(): invalid data %v, err = %v", addFavorite, err)
		return
	}

}

func TestService_FavoriteFromPayment_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FavoritePayment("no", "not")
	if err == nil {
		t.Errorf("FavoritePayment(): must return error, returned nil = %v", err)
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FavoritePayment(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	// создаём сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	addFavorite, err := s.FavoritePayment(payment.ID, "megafon")
	if err != nil {
		t.Error(err)
		return
	}

	payFavorite, err := s.PayFromFavorite(addFavorite.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): invalid data in favorite, payment = %v", payFavorite)
		return
	}


	if !reflect.DeepEqual(payment, payFavorite) {
		t.Errorf("PayFromFavorite(): invalid data in favorite = %v, payment = %v ", payFavorite, payment)
		return
	}
}

func TestService_PayFromFavorite_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.PayFromFavorite("hi")
	if err == nil {
		t.Errorf("PayFromFavorite(): must return error, returned nil = %v", err)
		return
	}

	if err != ErrFavoriteNotFound {
		t.Errorf("PayFromFavorite(): must return ErrFavoriteNotFound, returned = %v", err)
		return
	}
}
