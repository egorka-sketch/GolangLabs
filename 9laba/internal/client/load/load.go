package load

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"9laba/internal/client/api"
)

func Run(baseURL, login, password string, workers, iterations int) {
	var ok, failed atomic.Int64
	var firstErr sync.Once
	fail := func(worker int, err error) {
		failed.Add(1)
		firstErr.Do(func() {
			fmt.Printf("Ошибка: клиент %d: %v\n", worker, err)
		})
	}

	fmt.Printf("Запуск: %d пользователей по %d итераций\n", workers, iterations)
	start := time.Now()

	var wg sync.WaitGroup
	for i := range workers {
		wg.Go(func() {
			c := api.New(baseURL)
			ctx := context.Background()
			if err := c.Login(ctx, login, password); err != nil {
				fail(i, err)
				return
			}
			for j := range iterations {
				in := api.UserInput{Name: fmt.Sprintf("Нагрузка %d-%d", i, j), BirthDay: "2000-01-01"}
				u, err := c.CreateUser(ctx, in)
				if err != nil {
					fail(i, err)
					continue
				}
				ok.Add(1)
				if _, err := c.GetUser(ctx, u.ID); err != nil {
					fail(i, err)
				} else {
					ok.Add(1)
				}
				if err := c.DeleteUser(ctx, u.ID); err != nil {
					fail(i, err)
				} else {
					ok.Add(1)
				}
			}
			if err := c.Logout(ctx); err != nil {
				fail(i, err)
			}
		})
	}
	wg.Wait()

	elapsed := time.Since(start)
	total := ok.Load() + failed.Load()
	fmt.Printf("Время выполнения %v\n", elapsed.Round(time.Millisecond))
	fmt.Printf("Операций: %d ошибок: %d\n", total, failed.Load())
}
