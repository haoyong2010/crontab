package worker

import (
	"context"
	"crontab/common"
	"math/rand"
	"os/exec"
	"time"
)

//任务执行器
type Executor struct {
}

var (
	G_executor *Executor
)

//执行一个任务
func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd     *exec.Cmd
			output  []byte
			err     error
			result  *common.JobExecuteResult
			jobLock *JobLock
		)
		//任务结果
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}
		//初始化分布式锁
		jobLock = G_jobMgr.CreateJobLock(info.Job.Name)
		//记录开始时间
		result.StartTime = time.Now()

		//上锁
		//随机睡（0～1秒）
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		err = jobLock.TryLock()
		//释放锁
		defer jobLock.Unlock()
		if err != nil { //上锁失败
			result.Err = err
			result.EndTime = time.Now()
		} else {
			//上锁成功后，重置任务启动时间
			result.StartTime = time.Now()
			//执行shell命令
			cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
			//执行并捕获输出
			output, err = cmd.CombinedOutput()
			//记录任务结束时间
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
		}
		//执行任务完成后，把执行的结果返回给scheduler，scheduler会从executingTable中删除执行记录
		G_scheduler.PushJobResult(result)
	}()
}

//初始化执行器
func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
