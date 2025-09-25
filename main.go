package main

import (
	"fmt"
	"time"
)

// первое решение
// func or(channels ...<-chan interface{}) <-chan interface{} {
// 	mergedChan := make(chan interface{})
// 	var once sync.Once

// 	go func() {
// 		// для каждого канала создаем горутину, которая ждет закрытия канала и закрывает mergedChan
// 		// используем sync.Once чтобы гарантировать закрытие один раз(избежать панику)
// 		for i := range channels {
// 			go func(i int) {
// 				// блокируемся пока канал не отправит сигнал done или закроется, потом закрываем
// 				_, ok := <-channels[i]
// 				if !ok {
// 					once.Do(func() {
// 						close(mergedChan)
// 					})
// 				}

// 			}(i)
// 		}
// 	}()

// 	return mergedChan
// }

// второе решение
// если ничего не передано в функцию, то закрываем mergedChan и возвращаем его
// если 1 канал передан, то возвращаем его, иначе переходим в default
// читаем из 2 каналов, приходит значение -> селект завершится и канал закрываем
// если нет,то рекурсивно запускаем эту функцию на оставшиеся каналы, что повторяет эти действия
func Or(channels ...<-chan interface{}) <-chan interface{} {
	mergedChan := make(chan interface{})

	switch len(channels) {
	case 0:
		close(mergedChan)
		return mergedChan
	case 1:
		return channels[0]
	case 2:
		go func() {
			defer close(mergedChan)

			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		}()
		return mergedChan
	default:
		go func() {
			defer close(mergedChan)

			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-Or(channels[2:]...):
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
