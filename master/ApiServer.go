package master

import (
	"crontab/common"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"
)

//任务的HTTP接口
type ApiServer struct {
	httpserver *http.Server
}

//保存任务接口
//POST job={"name":"job1","command":"echo hello","cronExpr":"*****"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	//解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	//取表单中的job字段
	postJob = req.PostForm.Get("job")
	//反序列化job
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}

	//任务保存到ETCD中
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	//返回正常的应答({"errno":0,"msg":"","data":{...}})
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	//返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

var (
	//单例对象
	G_apiServer *ApiServer
)

//删除任务接口
// POST /job/delete name = job1
func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	var (
		err    error
		name   string
		oldJob *common.Job
		bytes  []byte
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	//删除任务名
	name = req.PostForm.Get("name")
	//删除任务
	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}
	//正常应答
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	//异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

//初始化服务
func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)
	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/del", handleJobDelete)
	//启动tcp监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}
	//创建一个HTTP服务
	httpServer = &http.Server{
		Handler:           mux,
		ReadTimeout:       time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		ReadHeaderTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
	}
	//赋值单例
	G_apiServer = &ApiServer{httpserver: httpServer}
	//启动服务端
	go httpServer.Serve(listener)
	return
}
