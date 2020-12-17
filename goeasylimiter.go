package goeasylimiter

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

const (
	// MinimaLimit is the minimal concurrency limit
	MinimaLimit = 2
)

var (
	log *zap.Logger
)

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), zapcore.InfoLevel),
	)
	log = zap.New(core, zap.AddCaller())
}

// Job is an interface for add jobs.
type Job interface {
	Run() (resp string, err error)
}

// EasyLimiter object
type EasyLimiter struct {
	semp chan struct{} // 控制并发的chan

	wg sync.WaitGroup // waitGroup 用于等待协程执行完成, 并关闭通道\清理资源

	jobChan    chan Job         // Job 队列(实现接口即可, 解耦了任务的具体实现)
	ResultChan chan interface{} // job执行结果队列
}

func NewEasyLimiter(taskCount, limit int) *EasyLimiter {
	if limit <= MinimaLimit {
		limit = MinimaLimit
	}

	c := &EasyLimiter{
		semp:       make(chan struct{}, limit),
		ResultChan: make(chan interface{}, taskCount),
		jobChan:    make(chan Job, taskCount),
	}

	// 创建后马上就监听job队列
	// job队列中有数据且semp队列未满 (满了会阻塞,以此来实现并发控制), 则取出job对象, 交给单独协程处理
	go func() {
		for job := range c.jobChan {
			//c.semp <- struct{}{}

			select {
			case c.semp <- struct{}{}:
			case <-time.After(time.Millisecond * 200):
				//log.Info("goroutine pool full, wait for 200 mis  ", zap.Int("size", len(c.semp)))
			}

			go func(ajob Job) {
				defer func() {
					c.wg.Done()
					<-c.semp
				}()
				//common.ZLogger.Info("开始执行任务")
				result, err := ajob.Run()
				//common.ZLogger.Info("完成执行任务")
				if err != nil {
					fmt.Printf("err:%v", err)
				}
				c.ResultChan <- result

			}(job)
		}
		log.Info("task队列关闭")
	}()

	return c
}

func (c *EasyLimiter) AddJob(job Job) {
	c.wg.Add(1)
	c.jobChan <- job
	//common.ZLogger.Info("添加任务")
}

func (c *EasyLimiter) Wait() {
	// 关闭job队列 ,此时已不会再添加
	close(c.jobChan)
	c.wg.Wait()
	// 关闭result队列,以保证range方式读取chan 程序会正常向下执行
	close(c.ResultChan)
}
