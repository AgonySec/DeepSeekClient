<script lang="ts" setup>
import ChatWindow from "./components/ChatWindow.vue";
import {CreateSession, GetAPI, GetSessionList, GetTitle, SetAPI} from "../wailsjs/go/main/App";
import {reactive, ref} from "vue";
import { ElNotification } from "element-plus";
import { Delete, Edit, Search, Share, Upload,Setting } from '@element-plus/icons-vue'

const sessionList = ref<{ id: string; title: string }[]>([]);
const dialogFormVisible = ref(false)
const formLabelWidth = '100px'
const form = reactive({
  api: '',

})
let selectSessionID = ref('')

function getSessionList() {
  GetAPI().then((res) => {
    if(res.code !==200){
      ElNotification({
        title: "API配置",
        message: "请配置APIKey",
        type: "error",
      })
      return
    }
    form.api = res.data
  })
  GetSessionList().then((res) => {
    if(res.code !== 200){
      ElNotification({
        title: "获取session列表",
        message: res.msg,
        type: "error",
      })
      return
    }
    const fetchTitles = (res.data as string[]).map(sessionID =>
        GetTitle(sessionID).then(titleRes => {
          // if (titleRes.code !== 200) {
          //   throw new Error(titleRes.msg);
          // }
          return { id: sessionID, title: titleRes.data };
        })
    );
    Promise.all(fetchTitles).then(titles => {
      sessionList.value = titles;
    }).catch((err) => {
      ElNotification({
        title: 'Error',
        message: `Failed to retrieve session list: ${err}`,
        type: 'error',
      });
    });
  }).catch((err) => {
    ElNotification({
      title: 'Error',
      message: `Failed to retrieve session list: ${err}`,
      type: 'error',
    });
  });
}

function handleSessionClick(sessionId: string) {
  selectSessionID.value = sessionId
  ElNotification({
    title: 'sessionId',
    message: sessionId,
    type: 'success',
  });
}

function createSession() {
  CreateSession().then((res) => {
    if(res.code !==200){
      ElNotification({
        title: "New Session",
        message: res.msg,
        type: "error",
      })
      return
    }
    selectSessionID.value = res.data
    ElNotification({
      title: '创建会话成功',
      message: res.msg,
      type: 'success',
    })
    getSessionList(); // 刷新会话列表
  })
}
function defaultSession() {
  CreateSession().then((res) => {
    if(res.code !==200){
      ElNotification({
        title: "默认会话",
        message: res.msg,
        type: "error",
      })
      return
    }
    selectSessionID.value = res.data
    ElNotification({
      title: '默认会话',
      message: res.msg,
      type: 'success',
    });
  })

}
function deleteSession(sessionId: string) {
  // 假设有一个删除会话的API，这里调用它
  // DeleteSession(sessionId).then(() => {
    ElNotification({
      title: '删除会话',
      message: `会话 ${sessionId} 已删除`,
      type: 'success',
    });
    getSessionList(); // 刷新会话列表
  // })
}
function setConfig() {
  SetAPI(form.api).then((res) => {
    if(res.code !==200){
      ElNotification({
        title: "API配置",
        message: res.msg,
        type: "error",
      })
      return
    }
    ElNotification({
      title: "API配置",
      message: res.msg,
      type: "success",
    })
  })
}

getSessionList();
defaultSession()
</script>

<template>
  <div class="common-layout">
    <el-container>
      <el-aside width="200px" class="sidebar">
        <el-button class="newSession" type="primary" plain @click="createSession">开启新对话</el-button>
        <div class="item" v-for="(session, index) in sessionList" :key="index" @click="handleSessionClick(session.id)">
          {{ session.title }}
        </div>
        <el-button class="setting" type="info" :icon="Setting" @click="dialogFormVisible= true" circle />

        <el-dialog v-model="dialogFormVisible" title="请输入Api" width="500">
          <el-form :model="form">
            <el-form-item label="API地址：" :label-width="formLabelWidth">
              <el-input type="password"  show-password v-model="form.api" autocomplete="off"  />
            </el-form-item>
          </el-form>
          <template #footer>
            <div class="dialog-footer">
              <el-button @click="dialogFormVisible = false">取消</el-button>
              <el-button type="primary" @click="setConfig">
                确定
              </el-button>
            </div>
          </template>
        </el-dialog>

      </el-aside>
      <el-container>
        <el-main>
          <ChatWindow :sessionID="selectSessionID"/>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<style>
.setting {
  margin-top: 20px;
  margin-bottom: 20px;

}
.newSession {
  margin-top: 20px;
  margin-bottom: 20px;
  width: 50%;
}
.item {
  border-radius: 4px;
  display: flex;
  justify-content: space-between;
  width: 80%;
  padding: 8px;
  margin: 0 auto;
  border-bottom: 1px solid #8d8a8a;
  cursor: pointer;
}

.sidebar {
  background-color: #f8f8fa;
  color: #333;
  height: 100vh;
  overflow-y: auto;
}

.common-layout {
  flex-direction: column;
  justify-content: space-between;
  background-color: #ffffff;
}
</style>