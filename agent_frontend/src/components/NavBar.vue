<template>
  <div id="nav">
    <v-container>
      <v-row align="center" justify="center">
        <v-app-bar
          color="blue"
          dense
        >
        <v-col cols="12" sm="4">
          <v-btn @click="goToHomePage()" class="mx-2">
            <v-icon large>
              mdi-home-outline
            </v-icon>
          </v-btn>
          <v-spacer></v-spacer>
        </v-col>
        <v-col cols="12" sm="4">
          <router-link :to="{ name: 'MyPosts'}">Posts | </router-link> 
          <router-link :to="{ name: 'MyCampaigns'}">Campaigns</router-link> 
        </v-col>
        <v-col cols="12" sm="4">
          <v-btn class="mx-2" >
            <v-icon large>
              mdi-shield-key-outline
            </v-icon>
          </v-btn>
          <v-btn @click="goToNewProduct()" class="mx-2">
            <v-icon large>
              mdi-plus-circle-outline
            </v-icon>
          </v-btn>
          <v-btn @click="showCartModal()" class="mx-2">
            <v-icon large>
              mdi-cart-variant
            </v-icon>
          </v-btn>
          <!--<v-col cols="12" sm="1" class="float-right">
                <settings v-if="isUserLogged"/>
          </v-col>-->
        </v-col>
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