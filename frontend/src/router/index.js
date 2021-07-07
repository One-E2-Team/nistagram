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
    path: '/homePage',
    name: 'HomePage',
    component: () =>
      import ('../views/HomePage.vue')
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
  },
  {
    path: '/reactions',
    name: 'Reactions',
    component: () =>
      import ('../views/MyReactions.vue')
  },
  {
    path: '/reports',
    name: 'Reports',
    component: () =>
      import ('./../components/Reports.vue')
  },
  {
    path: '/agent-requests',
    name: 'AgentRequests',
    component: () =>
      import ('./../components/AgentRequests.vue')
  },
  {
    path: '/register-agent',
    name: 'RegisterAgent',
    component: () =>
      import ('./../components/RegisterAgent.vue')
  },
  {
    path: '/messaging/:username',
    name: 'Messaging',
    props: true,
    component: () =>
      import ('./../components/Messaging.vue')
  },
]
const router = new VueRouter({
  routes
})

export default router