<template>
  <div id="nav">
    <v-app-bar
      color="blue"
      dense
    >
     <router-link :to="{ name: 'Home'}">Home</router-link>
    </v-app-bar>
  </div>
</template>

<script>
import * as comm from '../configuration/communication.js'
export default {
    name: "NavBar",
    components: {
    },
    data(){
      return {
        isUserLogged: comm.getJWTToken() != null
      }
    },
    mounted(){
    },
    methods: {
     
    }
}
</script>