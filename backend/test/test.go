package main

import (
	"DeepSeekClient/backend/chat"
	"context"
	"log"
)

/**
 *
 * @author Agony
 * @date 2025/2/11 14:46
 * @description test
 */
func main() {
	defer chat.CloseDB()

	// 测试对话
	_, err := chat.ChatDP(context.Background(), "test_session3", "我的上一个问题是什么？")
	if err != nil {
		log.Fatal(err)
	}
}
