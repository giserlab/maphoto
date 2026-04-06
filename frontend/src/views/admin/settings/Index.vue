<template>
  <div class="settings-page">
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-position="top"
      v-loading="loading"
      class="settings-form"
    >
      <el-card class="mb-20">
        <template #header><span>地图配置</span></template>

        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item label="地图标题" prop="title">
              <el-input v-model="form.title" placeholder="输入地图标题" />
            </el-form-item>
          </el-col>

          <el-col :xs="24" :sm="12">
            <el-form-item label="外部链接" prop="link">
              <el-input v-model="form.link" placeholder="输入外部链接" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :xs="24" :sm="8">
            <el-form-item label="默认经度">
              <el-input-number v-model="form.lon" :precision="4" style="width: 100%" />
            </el-form-item>
          </el-col>

          <el-col :xs="24" :sm="8">
            <el-form-item label="默认纬度">
              <el-input-number v-model="form.lat" :precision="4" style="width: 100%" />
            </el-form-item>
          </el-col>

          <el-col :xs="24" :sm="8">
            <el-form-item label="默认缩放级别">
              <el-input-number v-model="form.zoom" :min="1" :max="20" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :xs="24" :sm="8">
            <el-form-item label="最小缩放">
              <el-input-number v-model="form.minZoom" :min="1" :max="20" style="width: 100%" />
            </el-form-item>
          </el-col>

          <el-col :xs="24" :sm="8">
            <el-form-item label="最大缩放">
              <el-input-number v-model="form.maxZoom" :min="1" :max="20" style="width: 100%" />
            </el-form-item>
          </el-col>

          <el-col :xs="24" :sm="8">
            <el-form-item label="图标大小">
              <el-input-number v-model="form.iconSize" :min="10" :max="100" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="自动居中">
          <el-switch v-model="form.autoCenter" />
        </el-form-item>
      </el-card>

      <el-card class="mb-20">
        <template #header><span>其他设置</span></template>

        <el-form-item label="备注">
          <el-input v-model="form.note" type="textarea" :rows="4" placeholder="输入备注信息" />
        </el-form-item>
      </el-card>

      <el-form-item>
        <el-button type="primary" :loading="saving" @click="handleSubmit"> 保存设置 </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ElMessage } from "element-plus";
import type { FormInstance, FormRules } from "element-plus";
import { useUserStore } from "@/stores/user";

const userStore = useUserStore();
const formRef = ref<FormInstance>();
const loading = ref(false);
const saving = ref(false);

const form = reactive({
  title: "",
  link: "",
  iconSize: 30,
  lon: 116.4074,
  lat: 39.9042,
  zoom: 10,
  maxZoom: 18,
  minZoom: 3,
  tolorance: 0.5,
  autoCenter: false,
  note: "",
});

const rules: FormRules = {
  title: [{ required: true, message: "请输入地图标题", trigger: "blur" }],
};

onMounted(async () => {
  loading.value = true;
  try {
    await userStore.fetchConfig();
    if (userStore.config) {
      Object.assign(form, userStore.config);
    }
  } finally {
    loading.value = false;
  }
});

async function handleSubmit() {
  if (!formRef.value) return;

  await formRef.value.validate(async (valid) => {
    if (!valid) return;

    saving.value = true;
    try {
      await userStore.updateConfig(form);
      ElMessage.success("保存成功");
    } finally {
      saving.value = false;
    }
  });
}
</script>

<style scoped>
.settings-form {
  width: 100%;
}

.mb-20 {
  margin-bottom: 20px;
}
</style>
