package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // 使用SQLite数据库
)

// 假设已经有一个数据库连接
var dbs *sql.DB

func InitDB() error {
	var err error
	dbs, err = sql.Open("sqlite3", "./data.db") // 打开或创建一个SQLite数据库文件
	if err != nil {
		return fmt.Errorf("无法打开数据库: %v", err)
	}
	defer dbs.Close()

	// 测试数据库连接是否成功
	err = dbs.Ping()
	if err != nil {
		return fmt.Errorf("无法连接到数据库: %v", err)
	}

	// 创建表（如果不存在）
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS api_keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		conversation_id TEXT NOT NULL,
		message TEXT NOT NULL
	);
	`
	_, err = dbs.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}

	return nil
}

// 保存apikey到数据库
func saveAPIKey(apiKey string) error {
	_, err := dbs.Exec("INSERT INTO api_keys (key) VALUES (?)", apiKey)
	return err
}

// 从数据库获取apikey
func GetAPIKey() (string, error) {
	var err error
	dbs, err = sql.Open("sqlite3", "./data.db") // 打开或创建一个SQLite数据库文件
	if err != nil {
		return "无法打开数据库", fmt.Errorf("无法打开数据库: %v", err)
	}
	defer dbs.Close()

	var apiKey string
	err = dbs.QueryRow("SELECT key FROM api_keys LIMIT 1").Scan(&apiKey)
	if err != nil {
		return "", err
	}
	return apiKey, nil
}

// 保存messages到数据库
func saveMessages(messages []string, conversationID string) error {
	for _, message := range messages {
		_, err := dbs.Exec("INSERT INTO messages (conversation_id, message) VALUES (?, ?)", conversationID, message)
		if err != nil {
			return err
		}
	}
	return nil
}

// 从数据库获取messages
func getMessages(conversationID string) ([]string, error) {
	rows, err := dbs.Query("SELECT message FROM messages WHERE conversation_id = ?", conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var message string
		if err := rows.Scan(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// 处理连续对话
func handleConversation(newMessages []string, conversationID string) error {
	// 获取现有消息
	existingMessages, err := getMessages(conversationID)
	if err != nil {
		return err
	}

	// 合并新消息
	allMessages := append(existingMessages, newMessages...)

	// 保存所有消息
	return saveMessages(allMessages, conversationID)
}
