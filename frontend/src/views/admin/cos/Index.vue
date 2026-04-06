<template>
  <div class="cos-manager">
    <el-card class="config-card">
      <template #header>
        <div class="card-header">
          <span>腾讯云COS配置</span>
          <el-button type="primary" @click="saveConfig">
            <el-icon><Check /></el-icon>
            保存配置
          </el-button>
        </div>
      </template>

      <el-collapse v-model="activeNames">
        <el-collapse-item title="展开/收起配置" name="config">
          <el-form :model="config" label-width="120px">
            <el-form-item label="SecretId">
              <el-input
                v-model="config.secretId"
                placeholder="请输入SecretId"
                show-password
              />
            </el-form-item>
            <el-form-item label="SecretKey">
              <el-input
                v-model="config.secretKey"
                placeholder="请输入SecretKey"
                show-password
              />
            </el-form-item>
            <el-form-item label="Bucket">
              <el-input v-model="config.bucket" placeholder="请输入Bucket名称，如：mybucket-1250000000" />
            </el-form-item>
            <el-form-item label="Region">
              <el-input v-model="config.region" placeholder="请输入Region，如：ap-beijing" />
            </el-form-item>
            <el-form-item label="访问域名">
              <el-input v-model="config.domain" placeholder="可选：自定义CDN加速域名，如：https://cdn.example.com" />
            </el-form-item>
            <el-divider content-position="left">文件夹路径配置</el-divider>
            <el-form-item label="封面文件夹">
              <el-input v-model="config.thumbsPrefix" placeholder="默认为：thumbs" />
              <div class="form-tip">封面图片存储的文件夹路径，如：images/thumbs</div>
            </el-form-item>
            <el-form-item label="照片文件夹">
              <el-input v-model="config.photosPrefix" placeholder="默认为：photos" />
              <div class="form-tip">照片存储的文件夹路径，如：images/photos</div>
            </el-form-item>
          </el-form>

          <el-alert
            v-if="!isConfigValid"
            title="请填写完整的COS配置信息后才能使用文件管理功能"
            type="warning"
            :closable="false"
            show-icon
            style="margin-top: 16px;"
          />
        </el-collapse-item>
      </el-collapse>
    </el-card>

    <el-card v-if="isConfigValid" class="file-manager">
      <template #header>
        <div class="card-header">
          <span>文件管理</span>
          <div class="header-actions">
            <el-radio-group v-model="currentFolder" size="small">
              <el-radio-button label="thumbs">封面 ({{ config.thumbsPrefix || 'thumbs' }})</el-radio-button>
              <el-radio-button label="photos">照片 ({{ config.photosPrefix || 'photos' }})</el-radio-button>
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

      <el-table
        v-loading="loading"
        :data="fileList"
        stripe
        row-key="Key"
        style="width: 100%"
        empty-text="暂无文件"
      >
        <el-table-column label="预览" width="100">
          <template #default="{ row }">
            <el-image
              :src="getFileUrl(row)"
              :preview-src-list="[getFileUrl(row)]"
              fit="cover"
              class="preview-image"
            />
          </template>
        </el-table-column>
        <el-table-column prop="Key" label="文件名" min-width="200">
          <template #default="{ row }">
            <div class="filename">{{ getFileName(row.Key) }}</div>
            <div class="file-path text-secondary">{{ row.Key }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="Size" label="大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.Size) }}
          </template>
        </el-table-column>
        <el-table-column prop="LastModified" label="修改时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.LastModified) }}
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
          <el-tag>{{ currentPrefix }}/</el-tag>
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
import { Check, Upload, Refresh } from '@element-plus/icons-vue'
import COS from 'cos-js-sdk-v5'

interface CosConfig {
  secretId: string
  secretKey: string
  bucket: string
  region: string
  domain: string
  thumbsPrefix: string
  photosPrefix: string
}

interface CosFile {
  Key: string
  Size: number
  LastModified: string
  ETag: string
}

const STORAGE_KEY = 'cos_config'

// 折叠面板状态，默认折叠（空数组表示折叠）
const activeNames = ref<string[]>([])

// COS配置
const config = ref<CosConfig>({
  secretId: '',
  secretKey: '',
  bucket: '',
  region: '',
  domain: '',
  thumbsPrefix: 'thumbs',
  photosPrefix: 'photos',
})

// 当前文件夹
const currentFolder = ref<'thumbs' | 'photos'>('thumbs')

// 计算属性：获取当前文件夹前缀
const currentPrefix = computed(() => {
  const prefix = currentFolder.value === 'thumbs' ? config.value.thumbsPrefix : config.value.photosPrefix
  return prefix || currentFolder.value
})

// 文件列表
const fileList = ref<CosFile[]>([])
const loading = ref(false)
const allFiles = ref<CosFile[]>([]) // 存储所有文件用于前端分页

// 分页
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0,
})

// COS实例
let cosInstance: COS | null = null

// 上传相关
const showUploadDialog = ref(false)
const uploadRef = ref()
const uploadFileList = ref<File[]>([])
const uploading = ref(false)

// 重命名相关
const showRenameDialog = ref(false)
const renaming = ref(false)
const renameForm = ref({
  oldKey: '',
  newName: '',
})

// 计算属性：配置是否有效
const isConfigValid = computed(() => {
  return config.value.secretId &&
    config.value.secretKey &&
    config.value.bucket &&
    config.value.region
})

// 初始化COS实例
function initCos() {
  if (!isConfigValid.value) return

  cosInstance = new COS({
    SecretId: config.value.secretId,
    SecretKey: config.value.secretKey,
  })
}

// 从localStorage加载配置
function loadConfig() {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved) {
    try {
      config.value = JSON.parse(saved)
      initCos()
    } catch {
      // ignore parse error
    }
  }
}

// 保存配置到localStorage
function saveConfig() {
  if (!isConfigValid.value) {
    ElMessage.warning('请填写完整的配置信息')
    return
  }
  localStorage.setItem(STORAGE_KEY, JSON.stringify(config.value))
  initCos()
  ElMessage.success('配置已保存')
  refreshList()
}

// 获取文件列表（获取全部，前端分页和排序）
async function fetchFileList() {
  if (!cosInstance || !isConfigValid.value) return

  loading.value = true
  try {
    // 获取所有文件（COS单次最多1000个，一般用户不会超过这个数量）
    const result = await cosInstance.getBucket({
      Bucket: config.value.bucket,
      Region: config.value.region,
      Prefix: `${currentPrefix.value}/`,
      MaxKeys: 1000,
    })

    // 过滤掉目录本身，并按修改时间从新到旧排序
    allFiles.value = (result.Contents || [])
      .filter((item: CosFile) => item.Key !== `${currentPrefix.value}/`)
      .sort((a: CosFile, b: CosFile) => {
        return new Date(b.LastModified).getTime() - new Date(a.LastModified).getTime()
      })

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

// 获取文件URL
function getFileUrl(file: CosFile): string {
  if (config.value.domain) {
    return `${config.value.domain}/${file.Key}`
  }
  return `https://${config.value.bucket}.cos.${config.value.region}.myqcloud.com/${file.Key}`
}

// 获取文件名（去掉路径前缀）
function getFileName(key: string): string {
  return key.replace(`${currentPrefix.value}/`, '')
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
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

// 复制URL
function copyUrl(file: CosFile) {
  const url = getFileUrl(file)
  navigator.clipboard.writeText(url).then(() => {
    ElMessage.success('链接已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

// 删除文件
async function handleDelete(file: CosFile) {
  try {
    await ElMessageBox.confirm(
      `确定要删除文件 "${getFileName(file.Key)}" 吗？`,
      '确认删除',
      { type: 'warning' }
    )

    if (!cosInstance) return

    await cosInstance.deleteObject({
      Bucket: config.value.bucket,
      Region: config.value.region,
      Key: file.Key,
    })

    ElMessage.success('删除成功')
    refreshList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(`删除失败: ${error.message || error}`)
    }
  }
}

// 重命名
function handleRename(file: CosFile) {
  renameForm.value.oldKey = file.Key
  renameForm.value.newName = getFileName(file.Key)
  showRenameDialog.value = true
}

// 确认重命名
async function confirmRename() {
  if (!renameForm.value.newName.trim()) {
    ElMessage.warning('请输入新文件名')
    return
  }

  renaming.value = true
  try {
    if (!cosInstance) return

    const newKey = `${currentPrefix.value}/${renameForm.value.newName.trim()}`

    // COS重命名是通过复制+删除实现的
    await cosInstance.putObjectCopy({
      Bucket: config.value.bucket,
      Region: config.value.region,
      Key: newKey,
      CopySource: `${config.value.bucket}.cos.${config.value.region}.myqcloud.com/${renameForm.value.oldKey}`,
    })

    // 删除原文件
    await cosInstance.deleteObject({
      Bucket: config.value.bucket,
      Region: config.value.region,
      Key: renameForm.value.oldKey,
    })

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
function handleFileChange(file: any) {
  uploadFileList.value.push(file.raw)
}

// 文件移除
function handleFileRemove(file: any) {
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
  let successCount = 0

  for (const file of uploadFileList.value) {
    try {
      if (!cosInstance) continue

      const key = `${currentPrefix.value}/${file.name}`

      await cosInstance.putObject({
        Bucket: config.value.bucket,
        Region: config.value.region,
        Key: key,
        Body: file,
        ContentType: file.type,
      })

      successCount++
    } catch (error: any) {
      ElMessage.error(`上传 ${file.name} 失败: ${error.message || error}`)
    }
  }

  uploading.value = false
  showUploadDialog.value = false
  uploadFileList.value = []
  uploadRef.value?.clearFiles()

  ElMessage.success(`成功上传 ${successCount} 个文件`)
  refreshList()
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
  loadConfig()
  if (isConfigValid.value) {
    fetchFileList()
  }
})
</script>

<style scoped>
.cos-manager {
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

.config-card {
  margin-bottom: 0;
}

.file-manager {
  flex: 1;
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

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.4;
}

:deep(.el-collapse) {
  border: none;
}

:deep(.el-collapse-item__header) {
  font-size: 14px;
  color: #606266;
  padding-left: 0;
}

:deep(.el-collapse-item__content) {
  padding-bottom: 0;
}

:deep(.el-collapse-item__wrap) {
  border-bottom: none;
}
</style>
