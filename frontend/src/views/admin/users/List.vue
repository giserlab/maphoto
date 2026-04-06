<template>
  <div class="users-list-page">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
          <el-button type="primary" @click="showCreateDialog = true">
            <el-icon><Plus /></el-icon>
            添加用户
          </el-button>
        </div>
      </template>

      <el-table :data="users" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="row.admin ? 'danger' : 'info'">
              {{ row.admin ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastlogin" label="最后登录" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.lastlogin) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editUser(row)">
              编辑
            </el-button>
            <el-button
              v-if="row.id !== currentUserId"
              type="danger"
              link
              @click="deleteUser(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="users.length === 0" description="暂无用户" />
    </el-card>

    <!-- Create User Dialog -->
    <el-dialog
      v-model="showCreateDialog"
      title="添加用户"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        label-position="top"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="createForm.username" placeholder="输入用户名" />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input
            v-model="createForm.password"
            type="password"
            placeholder="输入密码"
            show-password
          />
        </el-form-item>

        <el-form-item label="角色">
          <el-radio-group v-model="createForm.admin">
            <el-radio :label="false">普通用户</el-radio>
            <el-radio :label="true">管理员</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="handleCreate">
          创建
        </el-button>
      </template>
    </el-dialog>

    <!-- Edit User Dialog -->
    <el-dialog
      v-model="showEditDialog"
      title="编辑用户"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="editFormRef"
        :model="editForm"
        :rules="editRules"
        label-position="top"
      >
        <el-form-item label="用户名">
          <el-input v-model="editForm.username" disabled />
        </el-form-item>

        <el-form-item label="新密码" prop="password">
          <el-input
            v-model="editForm.password"
            type="password"
            placeholder="留空表示不修改"
            show-password
          />
        </el-form-item>

        <el-form-item label="角色">
          <el-radio-group v-model="editForm.admin">
            <el-radio :label="false">普通用户</el-radio>
            <el-radio :label="true">管理员</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" :loading="updating" @click="handleUpdate">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import * as usersApi from '@/api/users'
import { useAuthStore } from '@/stores/auth'
import { formatDateTime } from '@/utils/format'
import type { User } from '@/types'

const authStore = useAuthStore()
const loading = ref(false)
const users = ref<User[]>([])
const currentUserId = computed(() => authStore.user?.id)

// Create dialog
const showCreateDialog = ref(false)
const createFormRef = ref<FormInstance>()
const creating = ref(false)
const createForm = reactive({
  username: '',
  password: '',
  admin: false,
})

const createRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少 6 位', trigger: 'blur' },
  ],
}

// Edit dialog
const showEditDialog = ref(false)
const editFormRef = ref<FormInstance>()
const updating = ref(false)
const editForm = reactive({
  id: 0,
  username: '',
  password: '',
  admin: false,
})

const editRules: FormRules = {
  password: [
    { min: 6, message: '密码长度至少 6 位', trigger: 'blur' },
  ],
}

onMounted(() => {
  fetchUsers()
})

async function fetchUsers() {
  loading.value = true
  try {
    const res = await usersApi.getAll()
    users.value = res.data || []
  } catch {
    // API might not exist yet
    users.value = []
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  if (!createFormRef.value) return

  await createFormRef.value.validate(async (valid) => {
    if (!valid) return

    creating.value = true
    try {
      await usersApi.create(createForm)
      ElMessage.success('创建成功')
      showCreateDialog.value = false
      createFormRef.value?.resetFields()
      fetchUsers()
    } finally {
      creating.value = false
    }
  })
}

function editUser(user: User) {
  editForm.id = user.id
  editForm.username = user.username
  editForm.password = ''
  editForm.admin = user.admin
  showEditDialog.value = true
}

async function handleUpdate() {
  if (!editFormRef.value) return

  await editFormRef.value.validate(async (valid) => {
    if (!valid) return

    updating.value = true
    try {
      const data: Partial<User> = { admin: editForm.admin }
      if (editForm.password) {
        data.password = editForm.password
      }
      await usersApi.update(editForm.id, data)
      ElMessage.success('更新成功')
      showEditDialog.value = false
      fetchUsers()
    } finally {
      updating.value = false
    }
  })
}

async function deleteUser(user: User) {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${user.username}" 吗？`,
      '确认删除',
      { type: 'warning' }
    )
    await usersApi.remove(user.id)
    ElMessage.success('删除成功')
    fetchUsers()
  } catch {
    // cancelled
  }
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
