import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import './style.css'

import JumpList from './views/JumpList.vue'
import JumpForm from './views/JumpForm.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/log' },
    { path: '/log', component: JumpList, meta: { title: 'Logbook' } },
    { path: '/log/new', component: JumpForm, meta: { title: 'New Jump' } },
    { path: '/log/:id/edit', component: JumpForm, meta: { title: 'Edit Jump' } },
  ],
})

const app = createApp(App)
app.use(router)
app.mount('#app')
