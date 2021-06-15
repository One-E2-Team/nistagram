import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import HomePage from '../components/HomePage'
Vue.use(VueRouter)

const routes = [{
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/homePage',
    name: 'HomePage',
    component: HomePage
  },
  {
    path: '/explore',
    name: 'Explore',
    component: () =>
      import ('../views/Explore.vue')
  },
  {
    path: '/not-found',
    name: 'NotFound',
    component: () =>
      import ('../views/PageNotFound.vue')
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
  },
  {
    path: '/profile/settings',
    name: 'ProfileSettings',
    component: () =>
      import ('./../components/ProfileSettings.vue')
  },
  {
    path: '/profile/settings/personal',
    name: 'PersonalSettings',
    component: () =>
      import ('./../components/PersonalSettings.vue')
  },
  {
    path: '/user/:username',
    name: 'Profile',
    props: true,
    component: () =>
      import ('./../views/Profile.vue')
  },
  {
    path: '/2fa-totp/:qruuid',
    name: 'TwoFactorAuth',
    props: true,
    component: () =>
      import ('./../components/TwoFactorAuth.vue')
  },

  {
    path: '/createVerificationRequest',
    name: 'CreateVerificationRequest',
    component: () =>
      import ('./../components/CreateVerificationRequest.vue')
  },

  {
    path: '/verificationRequests',
    name: 'VerificationRequests',
    component: () =>
      import ('./../components/VerificationRequests.vue')
  }

const router = new VueRouter({
  routes
})

export default router