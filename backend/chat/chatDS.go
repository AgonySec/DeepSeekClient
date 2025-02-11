package chat

import (
	"DeepSeekClient/backend/config"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // 使用SQLite数据库
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

/**
 *
 * @author Agony
 * @date 2025/2/11 11:14
 * @description chatDS
 */

const (
	defaultSessionID = "default_session"
	createTableSQL   = `CREATE TABLE IF NOT EXISTS conversations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	createIndexSQL     = "CREATE INDEX IF NOT EXISTS idx_session ON conversations(session_id);"
	initTimeout        = 5 * time.Second
	maxHistoryMessages = 10
)

var (
	dbOnce     sync.Once
	dbInstance *sql.DB
	initErr    error
)

// Conversation 表示单条对话记录
type Conversation struct {
	SessionID string
	Role      string
	Content   string
	CreatedAt time.Time
}

func InitDB(dsn string) error {
	dbOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), initTimeout)
		defer cancel()

		// 打开数据库连接
		dbInstance, initErr = sql.Open("sqlite3", dsn)
		if initErr != nil {
			initErr = fmt.Errorf("打开数据库失败: %w", initErr)
			return
		}

		// 配置连接池参数
		dbInstance.SetMaxOpenConns(25)
		dbInstance.SetMaxIdleConns(5)
		dbInstance.SetConnMaxLifetime(30 * time.Minute)

		// 执行建表语句
		if _, err := dbInstance.ExecContext(ctx, createTableSQL); err != nil {
			initErr = fmt.Errorf("创建表失败: %w", err)
			return
		}

		// 创建索引
		if _, err := dbInstance.ExecContext(ctx, createIndexSQL); err != nil {
			initErr = fmt.Errorf("创建索引失败: %w", err)
			return
		}

		// 验证数据库连接
		if err := dbInstance.PingContext(ctx); err != nil {
			initErr = fmt.Errorf("数据库连接验证失败: %w", err)
			return
		}
	})
	return initErr
}

// Chat 处理对话请求
func ChatDP(ctx context.Context, sessionID, userInput string) (string, error) {
	// 初始化数据库（示例DSN，根据实际情况配置）
	if err := InitDB("data.db"); err != nil {
		return "", fmt.Errorf("数据库初始化失败: %w", err)
	}
	apikey, err := GetApiKey()
	if err != nil {
		return "", fmt.Errorf("获取 API Key 失败: %w", err)
	}
	log.Println("apikey=", apikey)
	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	// 获取对话历史
	history, err := GetConversationHistory(ctx, sessionID, maxHistoryMessages)
	if err != nil {
		return "", fmt.Errorf("获取历史记录失败: %w", err)
	}
	log.Println("history=", history)

	// 构建消息链
	messages := buildMessages(history, userInput)
	//messages := []config.Message{
	//	{Role: "system", Content: "You are a helpful assistant"}}
	//messages = append(messages, config.Message{Role: "user", Content: userInput})

	log.Println("messages=", messages)

	// 调用API（示例实现）
	//assistantMsg, err := mockAPICall(ctx, messages)
	// 构造请求数据
	requestData := config.ChatCompletionRequest{
		Model:    "deepseek-chat",
		Messages: messages,
		Stream:   false,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %w", err)
	}
	// 创建请求对象
	req, err := http.NewRequest(
		"POST",
		"https://api.deepseek.com/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apikey)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 处理非200状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API返回错误状态码: %d\n响应内容: %s", resp.StatusCode, string(body))
	}

	// 解析响应数据
	var response config.ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("JSON解析失败: %w\n响应内容: %s", err, string(body))
	}
	var assistantMessage string
	// 输出结果
	if len(response.Choices) > 0 {
		assistantMessage = response.Choices[0].Message.Content
		log.Println("Assistant:", assistantMessage)
		// 添加助手消息到消息列表
		messages = append(messages, config.Message{Role: "assistant", Content: assistantMessage})
	} else {
		log.Println("未收到有效响应")
	}

	// 保存对话记录
	if err := saveConversations(ctx, sessionID, []string{userInput, assistantMessage}); err != nil {
		log.Printf("保存对话记录失败: %v", err)
	}

	return assistantMessage, nil
}
func GetApiKey() (string, error) {

	var apiKey string
	err := dbInstance.QueryRow("SELECT key FROM api_keys LIMIT 1").Scan(&apiKey)
	if err != nil {
		return "", err
	}
	return apiKey, nil

}

// getConversationHistory 获取指定会话的历史记录
func GetConversationHistory(ctx context.Context, sessionID string, limit int) ([]Conversation, error) {
	query := `
		SELECT session_id, role, content, created_at 
		FROM conversations 
		WHERE session_id = ? 
		LIMIT ?`

	rows, err := dbInstance.QueryContext(ctx, query, sessionID, limit)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %w", err)
	}
	defer rows.Close()

	var history []Conversation
	for rows.Next() {
		var c Conversation
		if err := rows.Scan(&c.SessionID, &c.Role, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("扫描记录失败: %w", err)
		}
		history = append(history, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历记录失败: %w", err)
	}

	return history, nil
}

// saveConversations 批量保存对话记录
func saveConversations(ctx context.Context, sessionID string, contents []string) error {
	tx, err := dbInstance.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("启动事务失败: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO conversations (session_id, role, content) VALUES (?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("准备语句失败: %w", err)
	}
	defer stmt.Close()

	// 保存用户输入
	if _, err := stmt.ExecContext(ctx, sessionID, "user", contents[0]); err != nil {
		return fmt.Errorf("插入用户消息失败: %w", err)
	}
	time.Sleep(1 * time.Millisecond)
	// 保存助手回复
	if _, err := stmt.ExecContext(ctx, sessionID, "assistant", contents[1]); err != nil {
		return fmt.Errorf("插入助手消息失败: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// buildMessages 构建消息链
func buildMessages(history []Conversation, currentInput string) []config.Message {
	messages := []config.Message{
		{Role: "system", Content: "You are a helpful assistant"},
	}

	for _, msg := range history {
		messages = append(messages, config.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	// 添加历史消息（按时间正序）
	//for i := len(history) - 1; i >= 0; i-- {
	//	msg := history[i]
	//	messages = append(messages, config.Message{
	//		Role:    msg.Role,
	//		Content: msg.Content,
	//	})
	//}

	// 添加当前输入
	messages = append(messages, config.Message{
		Role:    "user",
		Content: currentInput,
	})

	return messages
}

// mockAPICall 模拟API调用
func mockAPICall(ctx context.Context, messages []map[string]string) (string, error) {
	// 模拟API响应延迟
	select {
	case <-time.After(100 * time.Millisecond):
	case <-ctx.Done():
		return "", ctx.Err()
	}

	// 构造模拟响应
	response := struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}{
		Choices: []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{
			{
				Message: struct {
					Content string `json:"content"`
				}{
					Content: "这是模拟的助手回复",
				},
			},
		},
	}

	if len(response.Choices) == 0 {
		return "", errors.New("空API响应")
	}

	return response.Choices[0].Message.Content, nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if dbInstance != nil {
		return dbInstance.Close()
	}
	return nil
}
