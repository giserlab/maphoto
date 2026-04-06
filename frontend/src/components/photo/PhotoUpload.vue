<template>
  <div class="photo-url-input">
    <el-form @submit.prevent="handleSubmit">
      <el-form-item label="图片地址">
        <el-input
          v-model="url"
          placeholder="请输入图片URL地址"
          clearable
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSubmit" :disabled="!url.trim()">
          添加
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const emit = defineEmits<{
  success: [url: string]
}>()

const url = ref('')

function handleSubmit() {
  const trimmedUrl = url.value.trim()
  if (!trimmedUrl) {
    ElMessage.error('请输入图片地址')
    return
  }
  // 简单的URL格式验证
  if (!trimmedUrl.startsWith('http://') && !trimmedUrl.startsWith('https://')) {
    ElMessage.error('请输入有效的HTTP/HTTPS地址')
    return
  }
  emit('success', trimmedUrl)
  url.value = ''
}
</script>

<style scoped>
.photo-url-input {
  padding: 16px 0;
}
</style>
