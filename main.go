package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

// мой баланс
type myAmmount struct {
	amMytex sync.Mutex
	amount  int
}

// возможные операции над балансом
type BankClient interface {
	Deposit(amount int)          //зачисление
	Withdrawal(amount int) error //списание
	Balance() int                //баланс
}

// вывод баланса
func (m *myAmmount) Balance() *int {
	return &m.amount
}

// снятие суммы
func (m *myAmmount) Withdrawal(w int) *int {
	i := 0
	if w == 0 { //если не передали сколько списать - списываем рандомно
		i = rand.Intn(11) //а давайте раз уж 5 потоков, то списывать больше чем в задании, чтобы слишком быстро не росла сумма
	} else {
		i = w
	}
	if m.amount >= i {
		m.amMytex.Lock()
		m.amount = m.amount - i
		m.amMytex.Unlock()
	} else {
		fmt.Println("Недостаточно денег на счете")
	}
	return &m.amount
}

// зачисление суммы
func (m *myAmmount) Deposit(d int) *int {
	i := 0
	if d == 0 { //если не передали сколько начислить - начисляем рандомно
		i = rand.Intn(7) //а давайте раз уж 10 потоков, то зачислять меньше чем в задании, чтобы слишком быстро не росла сумма
	} else {
		i = d
	}
	m.amMytex.Lock()
	m.amount = m.amount + i
	m.amMytex.Unlock()
	return &m.amount
}

func main() {
	var m myAmmount
	for i := 0; i < 10; i++ {
		go func() {
			for {
				m.Deposit(0)
				sleep()
			}
		}()
	}

	for i := 0; i < 5; i++ {
		go func() {
			for {
				m.Withdrawal(0)
				sleep()
			}
		}()
	}

	var command string
	var amm int
	for { //читаем команды из консоли пока не отправят на выход
		_, err := fmt.Scanln(&command)
		checkerr(err)

		switch command {
		case "balance":
			fmt.Println(*m.Balance())
		case "deposit":
			fmt.Println("введите сумму:")
			_, err := fmt.Scanln(&amm)
			checkerr(err)
			fmt.Println("Баланс ", *m.Deposit(amm))
		case "withdrawal":
			fmt.Println("введите сумму:")
			_, err := fmt.Scanln(&amm)
			checkerr(err)
			fmt.Println("Баланс ", *m.Withdrawal(amm))
		case "exit":
			os.Exit(1)
		default:
			fmt.Println("Unsupported command. You can use commands: balance, deposit, withdrawal, exit")
		}
	}

}

// ждем от 0.5 до 1 сек
func sleep() {
	duration := rand.Float32()/2 + 0.5
	time.Sleep(time.Second * time.Duration(duration))
}

// обработка ошибок
func checkerr(err error) {
	if err != nil {
		fmt.Println("FATAL: ", err.Error())
		os.Exit(1)
	}
}
