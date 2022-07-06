package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	type poolToken = struct{}
	lPoolChan := make(chan poolToken, int(pool))
	defer close(lPoolChan)

	lUsers := make([]user, 0, n)
	var lUserLock sync.Mutex

	var lWaitGroup sync.WaitGroup
	var i int64
	for i = 0; i < n; i++ {
		lWaitGroup.Add(1)

		go func(aId int64) {
			defer lWaitGroup.Done()

			lPoolChan <- poolToken{} // block if pool limit exceeded
			defer func() {
				<-lPoolChan // release pool lock
			}()

			lUser := getOne(aId)

			lUserLock.Lock()
			defer lUserLock.Unlock()
			lUsers = append(lUsers, lUser)
		}(i)
	}
	lWaitGroup.Wait()

	return lUsers
}
