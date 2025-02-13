package main

import (
	"DeepSeekClient/backend/chat"
	"context"
	"github.com/labstack/gommon/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
func (a *App) ChatSiliconflow(userInput string) {

}
func (a *App) Debug(msg string) {
	runtime.LogDebug(a.ctx, msg)
}
func (a *App) Error(msg string) {
	runtime.LogError(a.ctx, msg)
}

func (a *App) HistoryChat(sessionID string) interface{} {
	//defer chat.CloseDB()
	if err := chat.InitDB("data.db"); err != nil {
		a.Error(err.Error())
		return map[string]interface{}{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}

	history, err := chat.GetConversationHistory(a.ctx, sessionID, 100)
	if err != nil {
		a.Error(err.Error())
		return map[string]interface{}{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	return map[string]interface{}{
		"code": 200,
		"msg":  "获取历史聊天信息",
		"data": history,
	}
}
func (a *App) Chat(userInput string, sessionID string) interface{} {
	defer chat.CloseDB()
	if err := chat.InitDB("data.db"); err != nil {
		a.Error(err.Error())
		return map[string]interface{}{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}

	assistantMessage, err := chat.ChatDP(a.ctx, sessionID, userInput)
	if err != nil {
		a.Error(err.Error())
		return map[string]interface{}{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	return map[string]interface{}{
		"code": 200,
		"msg":  "chat",
		"data": assistantMessage,
	}
}
func (a *App) GetTitle(sessionId string) interface{} {
	if err := chat.InitDB("data.db"); err != nil {
		a.Error(err.Error())
		return map[string]interface{}{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	title, err := chat.GetSessionTitle(a.ctx, sessionId)
	if title == "" {
		title = "New Session"
	}
	if err != nil {
		a.Error(err.Error())
		return map[string]interface{}{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
			"data": title,
		}
	}

	a.Debug(title)
	return map[string]interface{}{
		"code": 200,
		"msg":  "查询标题",
		"data": title,
	}
}
func (a *App) GetSessionList() interface{} {
	if err := chat.InitDB("data.db"); err != nil {
		a.Error(err.Error())
		return map[string]interface{}{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	sessionList, err := chat.GetSessionList(a.ctx)
	if err != nil {
		a.Error(err.Error())
		return map[string]interface{}{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	log.Debug(sessionList)
	return map[string]interface{}{
		"code": 200,
		"msg":  "获取session列表",
		"data": sessionList,
	}
}
func (a *App) CreateSession() interface{} {
	if err := chat.InitDB("data.db"); err != nil {
		a.Error(err.Error())
		return map[string]interface{}{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	sessionId, err := chat.CreateSession(a.ctx)
	if err != nil {
		a.Error(err.Error())
		return map[string]interface{}{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	return map[string]interface{}{
		"code": 200,
		"msg":  "New Session",
		"data": sessionId,
	}
}
func (a *App) SetAPI(api string) interface{} {
	if err := chat.InitDB("data.db"); err != nil {
		a.Error(err.Error())
	}
	if err := chat.SetAPI(a.ctx, api); err != nil {
		a.Error(err.Error())
		return map[string]interface{}{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	return map[string]interface{}{
		"code": 200,
		"msg":  "设置API完成",
	}
}
func (a *App) GetAPI() interface{} {
	if err := chat.InitDB("data.db"); err != nil {
		a.Error(err.Error())
	}
	apiKey, err := chat.GetApiKey()
	if err != nil {
		return map[string]interface{}{
			"code": -1,
			"msg":  "ERROR:" + err.Error(),
		}
	}
	return map[string]interface{}{
		"code": 200,
		"msg":  "获取APIKey",
		"data": apiKey,
	}
}
