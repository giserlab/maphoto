import axios, { type AxiosInstance, type AxiosError } from "axios";
import { ElMessage } from "element-plus";
import { useAuthStore } from "@/stores/auth";
import router from "@/router";

const api: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "./api/v1",
  timeout: 30000,
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore();
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`;
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  },
);

// Response interceptor
api.interceptors.response.use(
  (response) => {
    const { data } = response;
    if (!data.status) {
      ElMessage.error(data.message || "请求失败");
      return Promise.reject(new Error(data.message));
    }
    return data;
  },
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore();
      authStore.logout();
      router.push("/login");
      ElMessage.error("登录已过期，请重新登录");
    } else {
      const message = error.response?.data?.message || error.message || "网络错误";
      ElMessage.error(message);
    }
    return Promise.reject(error);
  },
);

export default api;
