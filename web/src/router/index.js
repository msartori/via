import { createRouter, createWebHistory } from 'vue-router'
import ClienteView from '../views/ClientView.vue'
import OperadorView from '../views/OperatorView.vue'
import MonitorView from '../views/MonitorView.vue'

const routes = [
  { path: '/', redirect: '/client' },
  { path: '/client', component: ClienteView },
  { path: '/operator', component: OperadorView },
  { path: '/monitor', component: MonitorView },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
