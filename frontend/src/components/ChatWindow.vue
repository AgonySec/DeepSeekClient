<template>
  <div class="chat-container">
    <!-- 消息展示区域 -->
    <div ref="messagesEnd" class="messages">
      <div
          v-for="(message, index) in messages"
          :key="index"
          class="message"
          :class="message.role"
      >
        <div class="bubble">
          <div class="content">{{ message.content }}</div>
          <div v-if="message.role === 'assistant'" class="typing-indicator" v-show="isLoading">
            <span></span><span></span><span></span>
          </div>
        </div>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="input-area">
      <textarea
          v-model="inputText"
          @keydown.enter.exact.prevent="sendMessage"
          placeholder="输入消息..."
          :disabled="isLoading"
      ></textarea>
      <button @click="sendMessage" :disabled="isLoading || !inputText.trim()">
        {{ isLoading ? '发送中...' : '发送' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import {Chat} from "../../wailsjs/go/main/App";

interface ChatMessage {
  role: 'user' | 'assistant'
  content: string
}

// 响应式数据
const messages = ref<ChatMessage[]>([])
const inputText = ref('')
const isLoading = ref(false)
const messagesEnd = ref<HTMLElement | null>(null)

// 自动滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesEnd.value) {
      messagesEnd.value.scrollTop = messagesEnd.value.scrollHeight
    }
  })
}

// 模拟 API 调用（替换为实际的后端接口）
const fetchAIResponse = async (prompt: string): Promise<string> => {
  // 这里应该替换为实际的后端接口调用
  // 示例使用模拟延迟
  await new Promise(resolve => setTimeout(resolve, 1000))
  return `这是 AI 的回复，针对："${prompt}"`
}
const aiResponse = ref()
// 发送消息处理
const sendMessage = async () => {
  const content = inputText.value.trim()
  if (!content || isLoading.value) return

  // 添加用户消息
  messages.value.push({ role: 'user', content })
  inputText.value = ''
  scrollToBottom()

  try {
    isLoading.value = true
    // 添加临时 AI 消息
    messages.value.push({ role: 'assistant', content: '' })

    // 调用 API
    const result = await Chat(content)
    console.log(result)
    // 更新最后一条消息
    messages.value[messages.value.length - 1].content = result
  } catch (error) {
    console.error('API 调用失败:', error)
    messages.value.push({
      role: 'assistant',
      content: '抱歉，请求处理失败，请稍后再试。'
    })
  } finally {
    isLoading.value = false
    scrollToBottom()
  }
}

// 初始化示例对话
onMounted(() => {
  messages.value.push({
    role: 'assistant',
    content: '您好！我是 AI 助手，有什么可以帮您？'
  })
})
</script>

<style scoped>
.chat-container {
  height: 100%;
  background: #f8f9fa;
}

.messages {
  height: 80vh;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.message {
  display: flex;
}

.message.user {
  justify-content: flex-end;
}

.message.assistant {
  justify-content: flex-start;
}

.bubble {
  max-width: 70%;
  padding: 12px 18px;
  border-radius: 18px;
  position: relative;
}

.user .bubble {
  background: #007bff;
  color: white;
  border-bottom-right-radius: 4px;
}

.assistant .bubble {
  background: #ffffff;
  color: #333;
  border: 1px solid #e0e0e0;
  border-bottom-left-radius: 4px;
}

.input-area {
  display: flex;
  padding: 20px;
  background: white;
  border-top: 1px solid #eee;
  gap: 10px;
}

textarea {
  flex: 1;
  padding: 12px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  resize: none;
  min-height: 48px;
  max-height: 120px;
}

button {
  padding: 0 20px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: opacity 0.2s;
}

button:disabled {
  background: #ccc;
  cursor: not-allowed;
  opacity: 0.7;
}

/* 打字动画 */
.typing-indicator {
  display: inline-flex;
  margin-left: 8px;
}

.typing-indicator span {
  width: 6px;
  height: 6px;
  margin: 0 2px;
  background: #666;
  border-radius: 50%;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}
</style>