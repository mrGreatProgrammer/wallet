package wallet

import (
	"reflect"
	"testing"

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