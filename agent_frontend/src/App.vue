<template>
  <v-app id="app">
    <nav-bar/>
    <router-view/>
  </v-app>
</template>

<script>
import NavBar from './components/NavBar.vue'
export default {
  components: { NavBar },
}
</script>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #0f1110;
    }
  }
}
</style>
