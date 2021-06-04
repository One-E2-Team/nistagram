import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'

Vue.use(VueRouter)

const routes = [{
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/about',
    name: 'About',
    component: () =>
      import ('../views/About.vue')
  },
  {
    path: '/post',
    name: 'Post',
    component: () =>
      import ('../components/CreatePost.vue')
  },
  {
    path: '/log-in',
    name: 'Login',
    component: () =>
      import ('./../components/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () =>
      import ('./../components/Register.vue')
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () =>
      import ('./../components/ResetPassword.vue')
  }
]

const router = new VueRouter({
  routes
})

export default router