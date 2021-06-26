<template>
  <div id="nav">
    <v-container>
      <v-row align="center" justify="center">
        <v-app-bar
          color="blue"
          dense
        >
        <button @click="goToHomePage()"><img src="../assets/home.png" width="65" height="65"></button>
        <v-spacer></v-spacer>
        <button @click="goToNewProduct()"><img src="../assets/add.png" width="40" height="40"></button>
        <button @click="showCartModal()"><img src="../assets/cart.png" width="40" height="40"></button>
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
     },
     goToHomePage(){
       this.$router.push({name: "HomePage"})
     },
     goToNewProduct(){
       this.$router.push({name: "NewProduct"})
     }
    }
}
</script>