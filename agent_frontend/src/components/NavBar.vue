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
        <v-col cols="12" sm="4" class="d-flex justify-end">
          <v-btn @click="compareCampaigns()" class="mx-2">Compare campaigns</v-btn>
          <APITokenModal />
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
        </v-col>
        </v-app-bar>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import * as comm from '../configuration/communication.js'
import axios from 'axios'
import APITokenModal from '../modals/APITokenModal.vue'
import { bus } from '../main'
export default {
    name: "NavBar",
    components: {APITokenModal  },
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
     },
     compareCampaigns() {
        axios({
          method: "get",
          url: comm.protocol + "://" + comm.server + "/report/pdf",
          headers: comm.getHeader(),
        }).then(response => {
          console.log(response);
          window.open(comm.protocol + "://" + comm.static_server + "/data/reportComparison.pdf", "_blank")
        });
     },
    }
}
</script>