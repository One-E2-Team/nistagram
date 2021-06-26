<template>
  <div id="nav">
    <v-container>
      <v-row align="center" justify="center">
        <v-app-bar
          color="blue"
          dense
        >
        <router-link :to="{ name: 'HomePage'}">Home</router-link>
        <v-spacer></v-spacer>
        <button @click="showCartModal()"><img src="../assets/cart.png" width="30" height="30"></button>
        <!--<v-col cols="12" sm="1" class="float-right">
              <settings v-if="isUserLogged"/>
        </v-col>-->
        </v-app-bar>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import * as comm from '../configuration/communication.js'
//import Settings from '../components/Settings.vue'
import { bus } from '../main'
export default {
    name: "NavBar",
    components: {  },
    data(){
      return {
        isUserLogged: comm.getJWTToken() != null
      }
    },
    mounted(){
      this.$root.$on('loggedUser', () => {
        this.isUserLogged = comm.getLoggedUserID != 0;
      })
    },
    methods: {
     showCartModal(){
       bus.$emit('show-cart', true)
     }
    }
}
</script>