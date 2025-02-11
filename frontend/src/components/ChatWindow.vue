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
        <div v-if="message.role === 'assistant'" class="avatar" :class="message.role">
          <img :src=getAvatar(message.role) alt="Avatar" />
        </div>
        <div class="bubble">
          <div class="content" v-html="toMarkdown(message.content)"></div>
<!--          <div v-if="message.role === 'assistant'" class="typing-indicator" v-show="isLoading">-->
<!--            <span></span><span></span><span></span>-->
<!--          </div>-->
        </div>
        <div v-if="message.role === 'user'" class="avatar" :class="message.role">
          <img :src=getUserAvatar(message.role) alt="Avatar" />
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
import { Chat, HistoryChat } from "../../wailsjs/go/main/App"; // 引入HistoryChat接口
import { marked } from 'marked'
interface ChatMessage {
  role: 'user' | 'assistant'
  content: string
}

// 响应式数据
const messages = ref<ChatMessage[]>([])
const inputText = ref('')
const isLoading = ref(false)
const messagesEnd = ref<HTMLElement | null>(null)
const toMarkdown = (text: string) => {

  return marked(text);
}


// 自动滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesEnd.value) {
      messagesEnd.value.scrollTop = messagesEnd.value.scrollHeight
    }
  })
}

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
onMounted(async () => {
  try {
    // 添加加载状态消息
    messages.value.push({ role: 'assistant', content: '加载中...' })

    // 假设有一个sessionID，这里使用一个示例值
    const sessionID = "test_session3";
    // 调用HistoryChat接口获取历史聊天记录
    const historyConversation = await HistoryChat(sessionID);

    // 移除加载状态消息
    messages.value.pop()

    messages.value = historyConversation.map(conversation => ({
      role: conversation.Role as 'user' | 'assistant',
      content: conversation.Content
    }));
  } catch (error) {
    console.error('获取历史聊天记录失败:', error);
    // 移除加载状态消息
    messages.value.pop()

    messages.value.push({
      role: 'assistant',
      content: '无法加载历史聊天记录，请稍后再试。'
    });
  }

  // 如果没有历史消息，添加初始消息
  if (messages.value.length === 0) {
    messages.value.push({
      role: 'assistant',
      content: '您好！我是 AI 助手，有什么可以帮您？'
    });
  }
})

// 获取头像路径
function getAvatar(role: 'user' | 'assistant'): string {
  // 确保路径是正确的，这里假设图片在 assets/images 目录下
  return role === 'user' ? '/src/assets/images/AgonySec.png' : '/src/assets/images/360.ico'
}
function getUserAvatar(role: 'user' | 'assistant'): string {
  // 确保路径是正确的，这里假设图片在 assets/images 目录下
  return  role === 'user' ? '/src/assets/images/AgonySec.png' : ''
}
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
  align-items: flex-start;
}

.message.user {
  justify-content: flex-end;
}

.message.assistant {
  justify-content: flex-start;
}

.avatar {
  width: 40px;
  height: 40px;

}

.avatar img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
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

.bubble {
  max-width: 70%;
  padding: 12px 18px;
  border-radius: 18px;
  position: relative;
}

.input-area {
  display: flex;
  padding: 16px;
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