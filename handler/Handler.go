package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/util/log"

	"github.com/lecex/core/client"
	healthPB "github.com/lecex/init/proto/health"
	cli "github.com/micro/go-micro/v2/client"

	"github.com/lecex/init/config"
)

var Conf = config.Conf

// Register 注册
func Register(srv micro.Service) {
	server := srv.Server()
	healthPB.RegisterHealthHandler(server, &Health{})

	go Sync() // 同步前端权限
}

// Sync 同步
func Sync() {
	time.Sleep(5 * time.Second)
	result := readJson("./permissions.json")
	data := []map[string]interface{}{}
	err := json.Unmarshal([]byte(result), &data)
	if err != nil {
		log.Log(err)
	}
	for _, v := range data {
		SyncFrontPermits(v)
	}
}

func SyncFrontPermits(req map[string]interface{}) {
	res := make(map[string]interface{})
	err := client.Call(context.TODO(), Conf.Service["user"], "FrontPermits.UpdateOrCreate", &req, &res, cli.WithContentType("application/json"))
	if err != nil {
		log.Log(err)
		time.Sleep(5 * time.Second)
		SyncFrontPermits(req)
	}
}
func readJson(filePath string) (result string) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	buf := bufio.NewReader(file)
	for {
		s, err := buf.ReadString('\n')
		result += s
		if err != nil {
			if err == io.EOF {
				fmt.Println("Read is ok")
				break
			} else {
				fmt.Println("ERROR:", err)
				return
			}
		}
	}
	return result
}
