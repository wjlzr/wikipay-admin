package main

import (
	"wikipay-admin/database"
	"wikipay-admin/server/job"

	"github.com/jasonlvhit/gocron"
)

//
func main() {
	//数据库初始化
	database.Setup()
	//redis connect
	//redis.Connect()

	job.MonitorJobsStart()
	<-gocron.Start()
}
