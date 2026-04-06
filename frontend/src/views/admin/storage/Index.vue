<template>
  <div class="storage-manager">
    <el-card class="file-manager">
      <template #header>
        <div class="card-header">
          <span>本地图片存储管理</span>
          <div class="header-actions">
            <el-radio-group v-model="currentFolder" size="small">
              <el-radio-button label="thumbs">封面 (thumbs)</el-radio-button>
              <el-radio-button label="photos">照片 (photos)</el-radio-button>
            </el-radio-group>
            <el-button type="primary" size="small" @click="showUploadDialog = true">
              <el-icon><Upload /></el-icon>
              上传文件
            </el-button>
            <el-button size="small" @click="refreshList">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <el-alert
        type="info"
        :closable="false"
        class="storage-info"
      >
        <template #title>
          <div class="storage-tips">
            <span>文件存储在服务器本地，路径: /uploads/</span>
            <span>支持格式: JPG、PNG、GIF、WebP，单个文件最大 20MB</span>
          </div>
        </template>
      </el-alert>

      <el-table
        v-loading="loading"
        :data="fileList"
        stripe
        row-key="path"
        style="width: 100%"
        empty-text="暂无文件"
      >
        <el-table-column label="预览" width="100">
          <template #default="{ row }">
            <el-image
              :src="row.url"
              :preview-src-list="[row.url]"
              fit="cover"
              class="preview-image"
            />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="文件名" min-width="200">
          <template #default="{ row }">
            <div class="filename">{{ row.name }}</div>
            <div class="file-path text-secondary">{{ row.path }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column prop="lastModified" label="修改时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.lastModified) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="copyUrl(row)">
              复制链接
            </el-button>
            <el-button link type="primary" size="small" @click="handleRename(row)">
              重命名
            </el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 上传对话框 -->
    <el-dialog
      v-model="showUploadDialog"
      title="上传文件"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form label-width="80px">
        <el-form-item label="目标目录">
          <el-tag>{{ currentFolder }}/</el-tag>
        </el-form-item>
        <el-form-item label="选择文件">
          <el-upload
            ref="uploadRef"
            action="#"
            :auto-upload="false"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
            accept="image/*"
            :file-list="uploadFileList"
            multiple
            :limit="10"
          >
            <el-button type="primary">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip">
                支持 JPG、PNG、GIF、WebP 格式，单个文件不超过 20MB
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showUploadDialog = false">取消</el-button>
        <el-button type="primary" :loading="uploading" @click="handleUpload">
          开始上传
        </el-button>
      </template>
    </el-dialog>

    <!-- 重命名对话框 -->
    <el-dialog
      v-model="showRenameDialog"
      title="重命名文件"
      width="400px"
      :close-on-click-modal="false"
    >
      <el-form label-width="80px">
        <el-form-item label="新文件名">
          <el-input v-model="renameForm.newName" placeholder="请输入新文件名" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRenameDialog = false">取消</el-button>
        <el-button type="primary" :loading="renaming" @click="confirmRename">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, Refresh } from '@element-plus/icons-vue'
import type { UploadFile } from 'element-plus'
import * as storageApi from '@/api/storage'
import type { StorageFile } from '@/api/storage'

const currentFolder = ref<'thumbs' | 'photos'>('thumbs')

// 文件列表
const fileList = ref<StorageFile[]>([])
const loading = ref(false)
const allFiles = ref<StorageFile[]>([])

// 分页
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0,
})

// 上传相关
const showUploadDialog = ref(false)
const uploadRef = ref()
const uploadFileList = ref<File[]>([])
const uploading = ref(false)

// 重命名相关
const showRenameDialog = ref(false)
const renaming = ref(false)
const renameForm = ref({
  oldFile: null as StorageFile | null,
  newName: '',
})

// 获取文件列表
async function fetchFileList() {
  loading.value = true
  try {
    const res = await storageApi.getFiles(currentFolder.value)
    allFiles.value = res.data || []
    pagination.value.total = allFiles.value.length
    applyPagination()
  } catch (error: any) {
    ElMessage.error(`获取文件列表失败: ${error.message || error}`)
  } finally {
    loading.value = false
  }
}

// 应用前端分页
function applyPagination() {
  const start = (pagination.value.page - 1) * pagination.value.pageSize
  const end = start + pagination.value.pageSize
  fileList.value = allFiles.value.slice(start, end)
}

// 刷新列表
function refreshList() {
  pagination.value.page = 1
  fetchFileList()
}

// 获取文件名
function getFileName(key: string): string {
  return key.split('/').pop() || key
}

// 格式化文件大小
function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化日期
function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleString('zh-CN')
}

// 复制URL
function copyUrl(file: StorageFile) {
  navigator.clipboard.writeText(file.url).then(() => {
    ElMessage.success('链接已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

// 删除文件
async function handleDelete(file: StorageFile) {
  try {
    await ElMessageBox.confirm(
      `确定要删除文件 "${file.name}" 吗？`,
      '确认删除',
      { type: 'warning' }
    )

    const [folder, filename] = file.path.split('/')
    await storageApi.deleteFile(folder, filename)
    ElMessage.success('删除成功')
    refreshList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(`删除失败: ${error.message || error}`)
    }
  }
}

// 重命名
function handleRename(file: StorageFile) {
  renameForm.value.oldFile = file
  renameForm.value.newName = file.name
  showRenameDialog.value = true
}

// 确认重命名
async function confirmRename() {
  if (!renameForm.value.newName.trim()) {
    ElMessage.warning('请输入新文件名')
    return
  }

  if (!renameForm.value.oldFile) return

  renaming.value = true
  try {
    const [folder, filename] = renameForm.value.oldFile.path.split('/')
    await storageApi.renameFile(folder, filename, renameForm.value.newName.trim())
    ElMessage.success('重命名成功')
    showRenameDialog.value = false
    refreshList()
  } catch (error: any) {
    ElMessage.error(`重命名失败: ${error.message || error}`)
  } finally {
    renaming.value = false
  }
}

// 文件选择变化
function handleFileChange(file: UploadFile) {
  if (file.raw) {
    uploadFileList.value.push(file.raw)
  }
}

// 文件移除
function handleFileRemove(file: UploadFile) {
  const index = uploadFileList.value.findIndex(f => f.name === file.name)
  if (index > -1) {
    uploadFileList.value.splice(index, 1)
  }
}

// 上传文件
async function handleUpload() {
  if (uploadFileList.value.length === 0) {
    ElMessage.warning('请选择要上传的文件')
    return
  }

  uploading.value = true
  try {
    const res = await storageApi.uploadFiles(currentFolder.value, uploadFileList.value)
    const result = res.data

    if (result.failed && result.failed.length > 0) {
      ElMessage.warning(`${result.failed.length} 个文件上传失败`)
    }

    if (result.count > 0) {
      ElMessage.success(`成功上传 ${result.count} 个文件`)
    }

    showUploadDialog.value = false
    uploadFileList.value = []
    uploadRef.value?.clearFiles()
    refreshList()
  } catch (error: any) {
    ElMessage.error(`上传失败: ${error.message || error}`)
  } finally {
    uploading.value = false
  }
}

// 分页变化
function handleSizeChange(size: number) {
  pagination.value.pageSize = size
  pagination.value.page = 1
  applyPagination()
}

function handlePageChange(page: number) {
  pagination.value.page = page
  applyPagination()
}

// 监听文件夹变化
watch(currentFolder, () => {
  refreshList()
})

onMounted(() => {
  fetchFileList()
})
</script>

<style scoped>
.storage-manager {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.storage-info {
  margin-bottom: 16px;
}

.storage-tips {
  display: flex;
  gap: 20px;
}

.preview-image {
  width: 60px;
  height: 60px;
  border-radius: 4px;
  object-fit: cover;
}

.filename {
  font-weight: 500;
  color: #303133;
}

.file-path {
  font-size: 12px;
  margin-top: 4px;
}

.text-secondary {
  color: #909399;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

:deep(.el-upload__tip) {
  color: #909399;
}
</style>
