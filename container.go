package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type ContainerINfo struct {
	Pid         string `json:"pid"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	CreatedTime string `json:"createdTime"`
	Status      string `json:"status"`
}

var (
	RUNNING             string = "running"
	STOP                string = "stop"
	Exit                string = "exit"
	DefaultInfoLocation string = "/var/run/jiandocker/%s/"
	ConfigName          string = "config.json"
)

// id生成
func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// 记录容器信息
func recordContainerINfo(containerPID int, command string, containerName string) (string, error) {
	// 准备json信息
	id := randStringBytes(10)
	createTime := time.Now().Format("2006-01-02 15:04:05")
	if containerName == "" {
		containerName = id
	}
	containerINfo := &ContainerINfo{
		Id:          id,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Status:      RUNNING,
		Name:        containerName,
	}

	// 生成json
	jsonBytes, err := json.Marshal(containerINfo)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	// 创建目录
	dir := fmt.Sprintf(DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dir, 0622); err != nil {
		log.Errorf("mkdir error %s error %v", dir, err)
		return "", err
	}

	// 创建文件
	fileName := dir + "/" + ConfigName
	file, err := os.Create(fileName)
	if err != nil {
		log.Errorf("create file error %s error %v", fileName, err)
		return "", err
	}

	// 将json写入文件
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("file write json error %v", err)
		return "", err
	}

	return containerName, nil
}

// 删除容器信息
func deleteContainerinfo(containerName string) {
	dir := fmt.Sprintf(DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dir); err != nil {
		log.Errorf("delete dir %s error %v", dir, err)
	}
}
