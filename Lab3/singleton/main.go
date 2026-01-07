package main

import (
	"fmt"
	"sync"
)

type single struct {
	data string
}

var (
	singleInstance *single
	once           sync.Once
)

func GetInstance() *single {
	once.Do(func() {
		fmt.Println("Creating Singleton instance...")
		singleInstance = &single{"I am the one!"}
	})
	return singleInstance
}

func main() {
	for i := 0; i < 3; i++ {
		go func() {
			item := GetInstance()
			fmt.Printf("Адрес: %p, Данные: %s\n", item, item.data)
		}()
	}

	fmt.Scanln()
}
