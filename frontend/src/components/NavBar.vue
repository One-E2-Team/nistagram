<template>
    <div id="nav">
      <v-container>
        <v-row align="center" justify="center">
          <v-col cols="12" sm="4"><v-spacer/></v-col>
          <v-col cols="12" sm="4">
                <router-link to="/">Home</router-link> |
                <router-link to="/explore">Explore</router-link> |
          </v-col>
          <v-col cols="12" sm="4">
              <v-spacer />
              <settings v-if="isUserLogged"></settings>
          </v-col>
        </v-row>
      </v-container>
    </div>
</template>

<script>
import * as comm from '../configuration/communication.js'
import Settings from '../components/Settings.vue'
export default {
    name: "NavBar",
    components: {
      Settings
    },
    data(){
      return {
        isUserLogged: comm.getLoggedUserUsername() != null
      }
    },
    mounted(){
      this.$root.$on('loggedUser', () => {
        this.isUserLogged = comm.getLoggedUserUsername() != null;
      })
    }
}
</script>