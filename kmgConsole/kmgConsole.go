package kmgConsole

import (
	"fmt"
	"os"
	"sort"
	"strings"
	//	"runtime/debug"
	//	"github.com/bronze1man/kmg/kmgDebug"
)

var VERSION = "1.0"

type Command struct {
	Name   string //名称,不区分大小写
	Desc   string //暂时没有用到
	Runner func() //运行这个命令的函数
	Hidden bool   //隐藏这个命令,使其在帮助列表里面不显示
}

var actionMap = map[string]Command{
	"version": Command{
		Name:   "Version",
		Runner: version,
	},
}
var exitActionList = []func(){}

func Main() {
	actionName := ""
	if len(os.Args) >= 2 {
		actionName = os.Args[1]
	}
	lowerActionName := strings.ToLower(actionName)

	action, exist := actionMap[lowerActionName]
	if !exist {
		fmt.Println("command [" + actionName + "] not found.(case insensitive)")
		help()
		return
	}

	os.Args = os.Args[1:]
	//搞这个是为了在panic的时候,仍然可以把输出搞到其他地方去,但是此处复杂度很高,很容易不稳定,具体效果有待确定.
	/*
		defer func(){
			r:=recover()
			if r!=nil{
				fmt.Println("panic:",r)
				os.Stderr.Write(kmgDebug.GetAllStack())
			}
			if len(exitActionList) > 0 {
				for _, action := range exitActionList {
					action()
				}
				exitActionList = nil
			}
		}()
	*/
	action.Runner()
	if len(exitActionList) > 0 {
		WaitForExit()
		for _, action := range exitActionList {
			action()
		}
		exitActionList = nil
	}
}

func AddAction(action Command) {
	name := strings.ToLower(action.Name)
	_, exist := actionMap[name]
	if exist {
		panic("command " + action.Name + " already defined.(case insensitive)")
	}
	actionMap[name] = action
}

//如果你使用了这个命令,主线程不会退出,而是会等到用户或者系统发送退出命令才会退出.
//使用命令调用时,注册退出动作,
//如果你使用了这个东西,请不要再使用WaitForExit了
//TODO 这个地方过于不直观
func AddExitAction(f func()) {
	exitActionList = append(exitActionList, f)
}

//调用这个函数来保证使用AddExitAction方法来注册进程退出请求.
// 使用这个来保证kmgConsole.Main一定会等待进程结束
func UseExitActionRegister() {
	AddExitAction(func() {})
}

//avoid initialization loop
func init() {
	AddAction(Command{
		Name:   "Help",
		Runner: help,
	})
}

func help() {
	fmt.Println("Usage: ")
	actionList := make([]Command, 0, len(actionMap))
	for _, command := range actionMap {
		actionList = append(actionList, command)
	}
	sort.Sort(tActionList(actionList))
	for i := 0; i < len(actionList); i++ {
		if actionList[i].Hidden {
			continue
		}
		fmt.Println("\t", actionList[i].Name)
	}
}

func version() {
	fmt.Println(VERSION)
}

type tActionList []Command

func (t tActionList) Len() int      { return len(t) }
func (t tActionList) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t tActionList) Less(i, j int) bool {
	return strings.ToLower(t[i].Name) < strings.ToLower(t[j].Name)
}
