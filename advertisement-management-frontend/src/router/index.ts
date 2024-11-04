import { createRouter, createWebHashHistory } from 'vue-router';
import Ads from '../views/Ads.vue'; // 导入您的视图组件

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Ads // 设置默认路由
    },
    // 可以在这里添加更多路由
];

const router = createRouter({
    history: createWebHashHistory(),
    routes,
});

export default router;