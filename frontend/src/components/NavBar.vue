<template>
    <div id="nav">
      <v-container>
        <v-row align="center" justify="center">
          <v-col cols="12" sm="4">
                <router-link to="/">Home</router-link> |
                <router-link to="/explore">Explore</router-link> |
                <v-btn @click="redirectToProfile()" v-if="username">My profile</v-btn>
          </v-col>
        </v-row>
      </v-container>
    </div>
</template>

<script>
import * as comm from '../configuration/communication.js'
export default {
    name: "NavBar",
    data: () => ({
      username: comm.getLoggedUserUsername()
    }),
    methods : {
      redirectToProfile(){
        this.username = comm.getLoggedUserUsername();
        if(this.username){
          this.$router.push('/profile?username=' + this.username);
        }
      }
    }
}
</script>