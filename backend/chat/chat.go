package chat

import (
	"DeepSeekClient/backend/config"
	"DeepSeekClient/backend/db"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

/**
 *
 * @author Agony
 * @date 2025/2/7 21:48
 * @description chat
 */

func Chat(userInput string, conversationID string) string {
	// 初始化数据库连接
	if err := db.InitDB(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	apiKey, err := db.GetAPIKey()
	if err != nil {
		log.Fatalf("获取 API Key 失败: %v", err)
	}
	log.Println("apikey=", apiKey)

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 读取之前的对话内容
	//conversationID := "default_conversation" // 这里可以使用更复杂的逻辑来生成或获取conversationID
	messages, err := db.GetMessagesWithRole(conversationID)
	if err != nil {
		log.Fatalf("获取对话内容失败: %v", err)
	}
	if len(messages) == 0 {
		messages = []config.Message{
			{Role: "system", Content: "You are a helpful assistant"}}
	}
	// 读取用户输入
	log.Printf("User: %s", userInput)
	// 添加用户消息到消息列表
	messages = append(messages, config.Message{Role: "user", Content: userInput})

	// 构造请求数据
	requestData := config.ChatCompletionRequest{
		Model:    "deepseek-chat",
		Messages: messages,
		Stream:   false,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		panic(fmt.Sprintf("JSON编码失败: %v", err))
	}

	// 创建请求对象
	req, err := http.NewRequest(
		"POST",
		"https://api.deepseek.com/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		panic(fmt.Sprintf("创建请求失败: %v", err))
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("请求失败: %v", err))
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("读取响应失败: %v", err))
	}

	// 处理非200状态码
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("API返回错误状态码: %d\n响应内容: %s",
			resp.StatusCode, string(body)))
	}

	// 解析响应数据
	var response config.ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		panic(fmt.Sprintf("JSON解析失败: %v\n响应内容: %s",
			err, string(body)))
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

	// 保存对话内容到数据库
	if err := db.SaveMessagesWithRole(messages, conversationID); err != nil {
		log.Fatalf("保存对话内容失败: %v", err)
	}

	return assistantMessage
}
