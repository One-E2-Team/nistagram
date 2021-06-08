<template>
    <div id="nav">
      <v-container>
        <v-row align="center" justify="center">
          <v-col cols="12" sm="4">
                <router-link to="/">Home</router-link> |
                <router-link to="/explore">Explore</router-link> |
                <v-btn @click="redirectToProfile()">My profile</v-btn>
          </v-col>
        </v-row>
      </v-container>
    </div>
</template>

<script>
import * as comm from '../configuration/communication.js'
import axios from 'axios'
export default {
    name: "NavBar",

    methods : {
      redirectToProfile(){
        if(comm.getLoggedUserID != 0){
             axios({
                method: "get",
                url: 'http://' + comm.server + '/api/profile/get-by-id/' + comm.getLoggedUserID()
            }).then(response => {
              if(response.status==200){
                comm.setJWTToken(response.data);
                this.$router.push('/profile?username=' + response.data.username);
                
              }
            })
        }
      }
    }
}
</script>