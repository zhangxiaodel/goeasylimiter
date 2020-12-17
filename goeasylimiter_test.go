package goeasylimiter

import (
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"testing"
	"time"
)

func RandomInt(n int) int {
	rand.Seed(time.Now().UnixNano())
	v := rand.Intn(n)
	if v < 1 {
		v = 1
	}
	return v
}

type SampleJob struct {
	total int
	idx   int
	key   string
}

func (job SampleJob) Run() (resp string, err error) {
	v := fmt.Sprintf("job run: %d/%d, %v", job.idx, job.total, job.key)

	//time.Sleep(time.Second * time.Duration(RandomInt(5)))
	time.Sleep(time.Second * 2)
	log.Info(v)
	return v, nil
}

func TestNewEasyLimiter(t *testing.T) {
	log.Info("begin")
	total:=60
	limiter:=NewEasyLimiter(total, 5)

	for i := 0; i < total; i++ {
		//添加任务
		limiter.AddJob(&SampleJob{
			total: total,
			idx:   i,
			key:   "test",
		})
	}

	// 控制完成并退出的信号
	chanSign:=make(chan interface{})

	// new goroutine to read
	go func() {
		for {
			select {
			case result, ok:= <- limiter.resultChan:
				log.Info("read from result chan", zap.Any("result", result))
				if !ok{
					log.Info("result chan 关闭/读取完成, 退出当前goroutine")
					chanSign <-1
					return
				}

			}
		}
	}()


	//等待任务完成
	limiter.Wait()

	<- chanSign

	log.Info("done")

}
