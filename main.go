package main

import (
	"fmt"
	"time"
)

// первое решение
// func or(channels ...<-chan interface{}) <-chan interface{} {
// 	mergedChan := make(chan interface{})

// 	go func() {
// 		// для каждого канала создаем горутину, которая считает все элементы из канала
// 		// отправляет сигнал в mergedChan и закрывает mergedChan, если канал закрыли
// 		for i := range channels {
// 			go func(i int) {
// 				// блокируемся пока канал не отправит сигнал done или закроется, потом закрываем
// 				<-channels[i]
// 				close(mergedChan)
// 			}(i)
// 		}
// 	}()

// 	return mergedChan
// }

// второе решение
// если ничего не передано в функцию, то закрываем mergedChan и возвращаем его
// если 1 канал передан, то возвращаем его, иначе переходим в default
// считаем из 2 каналов, приходит значение -> селект завершится и канал закрываем
// если нет,то рекурсивно запускаем эту функцию на оставшиеся каналы, что повторяет эти действия
func or(channels ...<-chan interface{}) <-chan interface{} {
	mergedChan := make(chan interface{})

	switch len(channels) {
	case 0:
		close(mergedChan)
		return mergedChan
	case 1:
		return channels[0]
	default:
		go func() {
			defer close(mergedChan)

			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-or(channels[2:]...):
			}
		}()
		return mergedChan
	}
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))

	start = time.Now()
	<-or(
		sig(time.Millisecond*10),
		sig(time.Millisecond*100),
		sig(time.Second),
	)
	fmt.Printf("done after %v\n", time.Since(start))

	start = time.Now()
	<-or()
	fmt.Printf("done after %v\n", time.Since(start))
}
