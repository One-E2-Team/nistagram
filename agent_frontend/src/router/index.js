import Vue from 'vue'
import VueRouter from 'vue-router'
Vue.use(VueRouter)

const routes = [{
    path: '/',
    name: 'Home',
    component: () =>
      import ('../views/Home.vue')
  },
  {
    path: '/not-found',
    name: 'NotFound',
    component: () =>
      import ('../views/PageNotFound.vue')
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
    path: '/homePage',
    name: 'HomePage',
    component: () =>
      import ('./../components/HomePage.vue')
  },
  {
    path: '/newProduct',
    name: 'NewProduct',
    component: () =>
      import ('./../components/NewProduct.vue')
  },
  {
    path: '/edit-product',
    name: 'EditProduct',
    component: () =>
      import ('./../components/EditProduct.vue')
  },
  {
    path: '/my-posts',
    name: 'MyPosts',
    component: () =>
      import ('../views/MyPosts.vue')
  },
  {
    path: '/my-campaigns',
    name: 'MyCampaigns',
    component: () =>
      import ('../views/MyCampaigns.vue')
  },
  {
    path: '/product/:id',
    name: 'Product',
    component: () =>
      import ('./../components/Product.vue')
  },
]
const router = new VueRouter({
  routes
})

export default router