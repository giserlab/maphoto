import { createRouter, createWebHashHistory } from "vue-router";
import { routes } from "./routes";
import { setupRouterGuards } from "./guards";
import NProgress from "@/utils/nprogress";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes,
});

router.beforeEach(() => {
  NProgress.start();
});

router.afterEach(() => {
  NProgress.done();
});
setupRouterGuards(router);

export default router;
