# goeasylimiter

lightweight goroutine concurrency limiter 

  

* Installation
  > go get github.com/zhangxiaodel/goeasylimiter

* Usage
    > 编写运行任务 (需实现Run接口 )
    ``` go 
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
    ```
    > sample 
   ```go 
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

    
   ``

* comment
    > 可以看到同时有5个goroutine 在并发跑, 处理结果也在同时读取了
  ```
    2020-12-17T09:52:47.169+0800	info	goeasylimiter/goeasylimiter_test.go:36	begin
    2020-12-17T09:52:49.174+0800	info	goeasylimiter/goeasylimiter_test.go:31	job run: 3/60, test
    2020-12-17T09:52:49.174+0800	info	goeasylimiter/goeasylimiter_test.go:31	job run: 0/60, test
    2020-12-17T09:52:49.174+0800	info	goeasylimiter/goeasylimiter_test.go:31	job run: 1/60, test
    2020-12-17T09:52:49.174+0800	info	goeasylimiter/goeasylimiter_test.go:31	job run: 4/60, test
    2020-12-17T09:52:49.174+0800	info	goeasylimiter/goeasylimiter_test.go:31	job run: 2/60, test
    2020-12-17T09:52:49.174+0800	info	goeasylimiter/goeasylimiter_test.go:57	read from result chan	{"result": "job run: 3/60, test"}
    2020-12-17T09:52:49.174+0800	info	goeasylimiter/goeasylimiter_test.go:57	read from result chan	{"result": "job run: 0/60, test"}
    2020-12-17T09:52:49.174+0800	info	goeasylimiter/goeasylimiter_test.go:57	read from result chan	{"result": "job run: 1/60, test"}
    2020-12-17T09:52:49.174+0800	info	goeasylimiter/goeasylimiter_test.go:57	read from result chan	{"result": "job run: 4/60, test"}
    2020-12-17T09:52:49.174+0800	info	goeasylimiter/goeasylimiter_test.go:57	read from result chan	{"result": "job run: 2/60, test"}
    2020-12-17T09:52:49.375+0800	info	goeasylimiter/goeasylimiter_test.go:31	job run: 5/60, test
    2020-12-17T09:52:49.375+0800	info	goeasylimiter/goeasylimiter_test.go:57	read from result chan	{"result": "job run: 5/60, test"}
    2020-12-17T09:52:49.574+0800	info	goeasylimiter/goeasylimiter_test.go:31	job run: 6/60, test
    2020-12-17T09:52:49.575+0800	info	goeasylimiter/goeasylimiter_test.go:57	read from result chan	{"result": "job run: 6/60, test"}
    2020-12-17T09:52:49.775+0800	info	goeasylimiter/goeasylimiter_test.go:31	job run: 7/60, test
    2020-12-17T09:52:49.775+0800	info	goeasylimiter/goeasylimiter_test.go:57	read from result chan	{"result": "job run: 7/60, test"}
    2020-12-17T09:52:49.980+0800	info	goeasylimiter/goeasylimiter_test.go:31	job run: 8/60, test
    2020-12-17T09:52:49.980+0800	info	goeasylimiter/goeasylimiter_test.go:57	read from result chan	{"result": "job run: 8/60, test"}
    2020-12-17T09:52:50.182+0800	info	goeasylimiter/goeasylimiter_test.go:31	job run: 9/60, test
    2020-12-17T09:52:50.182+0800	info	goeasylimiter/goeasylimiter_test.go:57	read from result chan	{"result": "job run: 9/60, test"}
    2020-12-17T09:52:50.381+0800	info	goeasylimiter/goeasylimiter_test.go:31	job run: 10/60, test
  ```