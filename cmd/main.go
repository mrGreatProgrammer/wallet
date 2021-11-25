package main

import (
	"fmt"

	"github.com/mrGreatProgrammer/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	account, err := svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(account)
	// err = svc.Deposit(account.ID, 10)
	// if err != nil {
	// 	switch err {
	// 	case wallet.ErrAmountMustBePositive:
	// 		fmt.Println("Сумма должна быть положительной")
	// 	case wallet.ErrAccountNotFound:
	// 		fmt.Println("Аккаунт пользователя не найден")
	// 	}
	// 	return
	// }

	fmt.Println(account)
	// fmt.Println(at)
	
	fmt.Println(svc.FindAccountByID(1))
	fmt.Println(svc.FindAccountByID(2))
	// fmt.Println(svc.Repeat())
}