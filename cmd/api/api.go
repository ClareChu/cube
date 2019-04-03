package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/inconshreveable/go-update"
	"github.com/manifoldco/promptui"
	xwebsocket "golang.org/x/net/websocket"
	"gopkg.in/src-d/go-git.v4"
	"hidevops.io/hiboot/pkg/model"
	hiboot_io "hidevops.io/hiboot/pkg/utils/io"

	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func GetInput(label string) (userInput string) {

	validate := func(input string) error {
		if len(input) < 8 {
			return errors.New("Password must have more than 8 characters")
		}
		return nil
	}

	checkName := func(input string) error {
		if label == USERNAME {
			if input == "" {
				return errors.New("Please Input username!")
			}
		}
		return nil
	}

	if label == PASSWORD {
		u := promptui.Prompt{
			Label:    label,
			Mask:     '*',
			Validate: validate,
		}
		userInput, _ = u.Run()
	} else {
		u := promptui.Prompt{
			Label:    label,
			Validate: checkName,
		}
		userInput, _ = u.Run()
	}
	return userInput
}

//定义用户用以HTTP登陆的JSON对象
type LoginAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type BaseResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

//启动PipelineStart
func PipelineStart(user *User, start *PipelineStarts) error {

	if _, err := StartInit(user, start); err != nil {
		return err
	}

	if err := Start(*start, GetPipelineStartApi(user.Server), user.Token); err != nil {
		return err
	}
	return nil
}

//初始化Start数据信息
func StartInit(user *User, start *PipelineStarts) (*PipelineStarts, error) {

	if user.Token == "" {
		err := fmt.Errorf("token can not be empty,login to get")
		return nil, err
	}
	if start.Name == "" && start.Namespace == "" {
		name, namespace, err := GetProjectInfoByCurrPath()
		if err != nil {
			return nil, err
		}

		start.Namespace = namespace
		start.Name = name
	}

	if start.Name == "" && start.Namespace != "" {
		err := fmt.Errorf("app can not be empty")
		return nil, err
	}
	if start.Name != "" && start.Namespace == "" {
		err := fmt.Errorf("project can not be empty")
		return nil, err
	}
	fmt.Println("app: ", start.Name)
	fmt.Println("project: ", start.Namespace)
	if start.Context != "" {
		fmt.Println("context: ", start.Context)
	}

	//如果 SourceCode 为空，则从发送http请求获取 SourceCode 信息，
	//如果SourceCode获取失败。则尝试本地推断SourceCode信息
	if start.SourceCode == "" {
		codeType, err := GetSourceCodeType(GetSourceCodeTypeApi(user.Server, start.Name, start.Namespace), user.Token)
		if err != nil {

			//如果发送http请求获取不到代码类型，则做本地补偿检测
			codeTypeStr, errs := sourceCodeSpot()
			if errs != nil {
				fmt.Println("[ERROR] source code get failed")
				fmt.Println("[ERROR] ", err)
				os.Exit(0)
			}
			codeType = codeTypeStr
		}
		start.SourceCode = codeType

		fmt.Println("source code: ", start.SourceCode)
	}

	return start, nil
}

//通过HTTP请求,启动 pipeline
func Start(start PipelineStarts, url, token string) error {

	jsonByte, err := json.Marshal(start)
	//fmt.Println("PipelineStart", string(jsonByte))

	if err != nil {
		fmt.Println("Login Failed ", err)
		return err
	}
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)
	if err != nil {
		//隐藏登陆完整URL信息
		errs := strings.Split(err.Error(), ":")
		err = errors.New(errs[len(errs)-1])
		fmt.Println("[ERROR] ", err)
		fmt.Println("Startup failed,Please check if the server is correct")
		os.Exit(0)
		//return errors.New("Startup failed,Please check if the server is correct")
	}
	defer resp.Body.Close()
	byteResp, _ := ioutil.ReadAll(resp.Body)
	resData := model.BaseResponse{}

	if err := json.Unmarshal(byteResp, &resData); err != nil {
		return err
	}

	if resData.Code != 200 {
		fmt.Println("resp", string(byteResp))
		return errors.New("pipeline start filed")
	}
	return nil
}

//通过HTTP登陆，返回Token
func Login(url, username, password string) (token string, err error) {
	if username == "" || password == "" {
		err := errors.New("username or password cannot be empty")
		return "", err
	}
	myAuth := LoginAuth{Username: username, Password: password}
	jsonByte, err := json.Marshal(myAuth)
	if err != nil {
		fmt.Println("Login Failed ", err)
		return token, err
	}

	myToken := model.BaseResponse{}
	//resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonByte))

	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err == nil {
		defer resp.Body.Close()
		byteResp, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(byteResp, &myToken)
		if err == nil {
			if myToken.Code == 200 {
				token = myToken.Data.(map[string]interface{})["token"].(string)
				err = errors.New(myToken.Message)
			} else if myToken.Code == 203 {
				err = errors.New("username or password is incorrectly entered")
				return "", err
			}
		}
	} else {
		//隐藏登陆完整URL信息
		errs := strings.Split(err.Error(), ":")
		err = errors.New(errs[len(errs)-1])
		fmt.Println("ERROR ", err)
		fmt.Println("Startup failed,Please check if the server is correct")
		os.Exit(0)
		//return "",errors.New("login request failed,please check if the server is correct")
	}
	if token == "" {
		return token, errors.New("token get failed")
	}
	err = nil
	return token, err
}

//获取用户HOME目录
func GetHomeDir() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	if "windows" == runtime.GOOS {
		fmt.Println("windows")
		return homeWindows()
	}
	return homeUnix()
}

//获取*unix系统家目录，不对外提供服务。给GetHomeDir调用
func homeUnix() (string, error) {
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}
	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}
	return result, nil
}

//获取Windows系统家目录，不对外提供服务。给GetHomeDir调用
func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}
	return home, nil
}

//检查指定目录或者文件是否存在
func pathExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//向指定文件路径写入数据
func WriteTextByPath(filePath, text string) error {
	//temporaryFile := fmt.Sprintf("./script-%d.sh", time.Now().Unix())
	fileObj, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fileObj.Close()
	if _, err := fileObj.WriteString(text); err == nil {
		return err
	}
	fileObj.Sync()
	return nil
}

//从指定文件路径读数据
func ReadTextByPath(filePath string) (string, error) {
	fileObj, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	//if fileObj,err := os.OpenFile(name,os.O_RDONLY,0644); err == nil {
	defer fileObj.Close()
	contents, err := ioutil.ReadAll(fileObj)
	if err != nil {
		return "", err
	}
	result := strings.Replace(string(contents), "\n", "", 1)
	return result, nil
}

//废弃
func WsLogsOut(server, path, query string) {

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: server, Path: path, RawQuery: query}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

//连接ws服务端并且循环接受数据并打印到控制台
func ClientLoop(url string, f func(string) error) error {
	WS, err := xwebsocket.Dial(url, "", "http://localserver/")
	if err != nil {
		return errors.New("failed to connect websocket")
	}
	defer func() {
		if WS != nil {
			WS.Close()
		}
	}()

	WS.Write([]byte("hello service"))

	var msg = make([]byte, 2048)
	for {
		if n, err := WS.Read(msg); err != nil {
			return err
		} else {
			message := string(msg[:n])

			if strings.Contains(message, "in namespace") || strings.Contains(message, "Server closed") {
				continue
			}
			if strings.Contains(message, "Information acquisition failed") {
				fmt.Println("[ERROR] ", message)
				os.Exit(0)
			}

			if err := f(message); err != nil {
				return nil
			}

		}
	}
}

//读取pipeline config信息
func ReadConfig() (user *User, filePath string, err error) {
	homeDir, err := GetHomeDir()
	if err != nil {
		log.Println("Error", err)
		return nil, "", err
	}

	filePath = filepath.Join(homeDir, CLI_DIR, CONFIG_FILE)

	if !pathExists(filePath) {
		if err := os.MkdirAll(filepath.Join(homeDir, CLI_DIR), 0777); err != nil {
			log.Println("Error", err)
			return nil, "", err
		}
		u := User{Server: DEFAULT_SERVER}
		if err := WriteConfig(&u, filePath); err != nil {
			fmt.Println("E", err)
			return nil, "", err
		}

		return &u, filePath, nil
	}

	text, err := ReadTextByPath(filePath)
	if err != nil {
		log.Println("Error", err)
		return nil, "", err
	}

	if text != "" {
		if err := json.Unmarshal([]byte(text), &user); err != nil {
			fmt.Println("Error", err)
			return nil, "", err
		}
	}
	return user, filePath, nil
}

//把user结构转化成字json字符串写在配置文件中
func WriteConfig(user *User, filePath string) error {
	userByte, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error", err)
		return err
	}
	if err := WriteTextByPath(filePath, string(userByte)); err != nil {
		fmt.Println("Error", err)
		return err
	}
	return nil
}

//根据项目的name和namespace获取代码类型
func GetSourceCodeType(url, token string) (string, error) {
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", url, nil)
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)
	if err != nil {
		//隐藏登陆完整URL信息
		errs := strings.Split(err.Error(), ":")
		err = errors.New(errs[len(errs)-1])
		return "", err
	}

	defer resp.Body.Close()
	byteResp, _ := ioutil.ReadAll(resp.Body)
	resData := model.BaseResponse{}

	if err := json.Unmarshal(byteResp, &resData); err != nil {
		return "", err
	}

	if resData.Code != 200 {
		err := errors.New("failed to get the code type")
		return "", err
	}
	return resData.Data.(string), nil
}

//根据本地文件属性 推断代码类型
func sourceCodeSpot() (string, error) {
	if b := hiboot_io.EnsureWorkDir(1, ".git"); !b {
		return "", fmt.Errorf("did not find git project")
	}

	sourceCode := ""
	if pathExists("pom.xml") {
		sourceCode = "java"
	} else if pathExists("package.json") {
		sourceCode = "nodejs"
	} else {
		err := fmt.Errorf("no code type found")
		return "", err
	}

	return sourceCode, nil
}

//获取当前项目的name和namespace
func GetProjectInfoByCurrPath() (string, string, error) {

	if b := hiboot_io.EnsureWorkDir(1, ".git"); !b {
		return "", "", fmt.Errorf("did not find git project")
	}

	r, err := git.PlainOpen(".")
	if err != nil {
		fmt.Println("Error", err)
		return "", "", err
	}

	remotes, err := r.Remotes()
	if err != nil {
		fmt.Println("Error", err)
		return "", "", err
	}

	var gitPath string
	for _, r := range remotes {
		reg := regexp.MustCompile(`(?i:http|git).*.git`)
		regStr := reg.FindAllString(r.String(), -1)
		if len(regStr) != 0 {
			gitPath = regStr[0]
		}
	}
	if gitPath == "" {
		return "", "", fmt.Errorf("failed to get git remotes")
	}

	gitPath = strings.Replace(gitPath, ".git", "", -1)
	strList := strings.Split(gitPath, ":")
	strList = strings.Split(strList[len(strList)-1], "/")
	if len(strList) < 2 {
		return "", "", fmt.Errorf("project info gets failed")
	}
	return strList[len(strList)-1], strList[len(strList)-2], nil
}

//mkdir cube/client-darwin-386 cube/client-darwin-amd64 cube/client-linux-386 cube/client-linux-amd64 cube/client-windows-386 cube/client-windows-amd64
//根据升级地址，更新客户端
func _update(url string, options update.Options) error {

	// request the new file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, options)
	if err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			fmt.Println(fmt.Sprintf("Failed to rollback from bad update: %v", rerr))
		}
	}
	return err
}

//向服务端发请求获取升级相关信息，并执行升级
func DoUpdate(url string) error {

	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", url, nil)
	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		//隐藏登陆完整URL信息
		errs := strings.Split(err.Error(), ":")
		err = errors.New(errs[len(errs)-1])
		fmt.Println("[ERROR] Version verification failed,please check if the server is correct")
		return errors.New("login request failed,please check if the server is correct")
	}

	resJson, _ := json.Marshal(resp)

	defer resp.Body.Close()
	byteResp, _ := ioutil.ReadAll(resp.Body)
	resData := model.BaseResponse{}
	if err := json.Unmarshal(byteResp, &resData); err != nil {
		return errors.New("abnormal information acquisition")
	}
	if resData.Code != 200 {
		return errors.New("abnormal information acquisition")
	}

	updateInfo := resData.Data.(map[string]interface{})
	if updateInfo["enable"].(bool) {

		if updateInfo["url"].(string) == "" {
			fmt.Println("[INFO] ", resJson)
			fmt.Println("[ERROR] abnormal information acquisition.")
			return errors.New("abnormal information acquisition")
		}
		isUpdate := GetInput("discover new version，Whether to upgrade. y/n ")
		if isUpdate == "y" || isUpdate == "Y" {

			fmt.Println("upgrading...")

			if err := _update(updateInfo["url"].(string), update.Options{}); err != nil {
				fmt.Println("[ERROR]", err)
			}

			fmt.Println("update successed.")
			os.Exit(0)
		}
	}

	return nil
}

//作为参数来控制日志输出
func LogOut(message string) error {
	fmt.Print(message)
	return nil
}

//作为参数来控制日志输出
func BuildLogOut(message string) error {
	fmt.Print(message)
	if END_STR != "" {
		if strings.Contains(message, END_STR) {
			err := fmt.Errorf("end for build")
			return err
		}
	}
	return nil
}
