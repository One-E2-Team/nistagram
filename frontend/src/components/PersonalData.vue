<template>
  <v-card
    style="witdh: scale"
    class="mx-auto"
    tile
  >
    <v-img
      height="100%"
      src="../assets/profilebackground.jpg"
    >
      <v-row
        align="center"
        class="fill-height"
      >
        <v-col
          align-self="center"
          class="pa-10"
          cols="3"
          
        >
          <v-avatar
            class="profile "
            color="transparent"
            size="150"
            tile
          >
            <v-img style="border-radius: 45%;" src="../assets/profilepicture.jpg" ></v-img>
          </v-avatar>
        </v-col>
        <v-col class="py-0" cols="4">
          <v-list-item
            color="rgba(0, 0, 0, .4)"
            dark
          >
            <v-list-item-content>
              <v-list-item-title class="text-h6 text-left">
                {{profile.username}}
              </v-list-item-title>
              <v-list-item-subtitle class="text-h6 text-left">Name : {{profile.personalData.name}}</v-list-item-subtitle>
              <v-list-item-subtitle class="text-h6 text-left">Surname : {{profile.personalData.surname}}</v-list-item-subtitle>
              <v-list-item-subtitle class="text-h6 text-left">Birth date: {{profile.personalData.birthDate}}</v-list-item-subtitle>
              <v-list-item-subtitle class="text-h6 text-left">web site: {{profile.webSite}}</v-list-item-subtitle>
              <v-list-item-subtitle class="text-h6 text-left">Biography : {{profile.biography}}</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-col>
      </v-row>
    </v-img>
  </v-card>
</template>

<script>
import axios from 'axios'
import * as comm from '../configuration/communication.js'
export default {
name: "PersonalData",

data: () => ({
    profile: {}
}),
methods: {

},
created(){
    let username = comm.getUrlVars()['username'];
    axios({
        method: "get",
        url: 'http://' + comm.server + '/api/profile/get/' + username,
    }).then(response => {
        if(response.status==200){
            this.profile = response.data;
            this.$emit('loaded-user', this.profile.ID)
        }else{
          this.$router.push({name: 'NotFound'})
        }
    }).catch(reason => {
        console.log(reason);
        this.$router.push({name: 'NotFound'})
    });
}

}
</script>

<style>
</style>