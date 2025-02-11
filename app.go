package main

import (
	"DeepSeekClient/backend/chat"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"time"
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

type Conversation struct {
	SessionID string
	Role      string
	Content   string
	CreatedAt time.Time
}

func (a *App) HistoryChat(sessionID string) []chat.Conversation {
	//defer chat.CloseDB()
	if err := chat.InitDB("data.db"); err != nil {
		a.Error(err.Error())
	}
	//测试数据
	history, err := chat.GetConversationHistory(a.ctx, sessionID, 10)
	if err != nil {
		a.Error(err.Error())
	}
	return history
}
func (a *App) Chat(userInput string) string {
	defer chat.CloseDB()

	// 测试对话
	assistantMessage, err := chat.ChatDP(a.ctx, "test_session3", userInput)
	if err != nil {
		a.Error(err.Error())
	}
	return assistantMessage
}
