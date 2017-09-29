package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

//github  接口返回格式
type ResultStruct struct {
	Code int        `json:"code"`
	Data DataStruct `json:"data"`
}
type DataStruct struct {
	Ref       string           `json:"ref"`
	File      FileStruct       `json:"file"`
	IsHead    bool             `json:"isHead"`
	CanEdit   bool             `json:"can_edit"`
	HeadComit HeadCommitStruct `json:"headCommit"`
}
type FileStruct struct {
	Data              string              `json:"data"`
	Lang              string              `json:"lang"`
	Size              int                 `json:"size"`
	Previewed         bool                `json:"previewed"`
	LastCommitMessage string              `json:"lastCommitMessage"`
	LastCommitDate    int64               `json:"lastCommitDate"`
	LastCommitId      string              `json:"lastCommitId"`
	LastCommitter     LastCommitterStruct `json:"lastCommitter"`
	Mode              string              `json:"mode"`
	Path              string              `json:"path"`
	Mame              string              `json:"name"`
}
type HeadCommitStruct struct {
	FullMessage  string          `json:"fullMessage"`
	ShortMessage string          `json:"shortMessage"`
	AllMessage   string          `json:"allMessage"`
	CommitId     string          `json:"commitId"`
	CommitTime   int64           `json:"commitTime"`
	Committer    CommitterStruct `json:"committer"`
	NotesCount   int             `json:"notesCount"`
}
type LastCommitterStruct struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Link   string `json:"link"`
}

type CommitterStruct struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Link   string `json:"link"`
}

//hosts文件位置
var filePath = "/etc/hosts"

//该标识符之后的数据会被删掉，之前的数据保存下来。 可以自定义的hosts写在标识符之前
var explodeString = "###################*******************"

func ReadFile(path string) []string {
	var data = []string{}
	fi, err := os.Open(path)
	if err != nil {
		return data
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		stringA := string(a)
		data = append(data, stringA)
		if stringA == explodeString {
			break
		}
	}
	return data
}

//对filePath和explodeString进行初始化
func InitConfig() {
	filePathTmp := flag.String("F", filePath, "hosts文件位置")
	explodeStringTmp := flag.String("E", explodeString, " 标识符")
	flag.Parse()
	filePath = *filePathTmp
	explodeString = *explodeStringTmp

}
func main() {
	InitConfig()
	Data := ReadFile(filePath)

	//解析接口数据
	ApiArray := &ResultStruct{}
	result := Get("https://coding.net/api/user/scaffrey/project/hosts/git/blob/master%252Fhosts-files%252Fhosts")
	if result == "" {
		panic("读取接口数据失败")
	}
	err := json.Unmarshal([]byte(result), ApiArray)
	if err != nil {
		panic("json 解析失败")
	}
	fd, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	fd_time := time.Now().Format("2006-01-02 15:04:05")
	fd_content := strings.Join(Data, "\n") + "\n" + fd_time + "\n" + ApiArray.Data.File.Data
	buf := []byte(fd_content)
	fd.Write(buf)
	fd.Close()
}

func Get(url string) string {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return ""
	}
	//Add 头协议
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Add("Accept-Language", "ja,zh-CN;q=0.8,zh;q=0.6")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Cookie", "")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	//接收服务端返回给客户端的信息
	response, _ := client.Do(request)
	defer response.Body.Close()

	if response.StatusCode == 200 {
		str, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Fatal error ", err.Error())
			return ""
		}
		return string(str)
	} else {
		return ""
	}
}
