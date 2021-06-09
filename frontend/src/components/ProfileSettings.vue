<template>
<v-container>
  <v-row align="center" justify="center">
    <v-col cols="12" sm="6" >
      <h1 class="display-2 font-weight-bold mb-3">
         Profile settings
      </h1>
      <v-form
        ref="form"
      >
      <v-checkbox
          v-model="settings.isPrivate"
          :label="'Private account'"
        ></v-checkbox>

        <v-checkbox
          v-model="settings.canReceiveMessageFromUnknown"
          :label="`Can be tagged on posts`"
        ></v-checkbox>

        <v-checkbox
          v-model="settings.canBeTagged"
          :label="`Can recieve message from unknown profiles`"
        ></v-checkbox>

        <v-btn
          color="success"
          class="mr-4"
          @click="updateSettings"
        >
          Confirm
        </v-btn>

      </v-form>
    </v-col>
  </v-row>
</v-container>
</template>

<script>
import * as comm from "../configuration/communication.js"
import axios from 'axios'
export default {
  data(){
    return{
      settings: {
        isPrivate:false,
        canReceiveMessageFromUnknown: false,
        canBeTagged: false
      },
    }
  },

  created(){
    axios({
      method: 'get',
      url: "http://" + comm.server + "/api/profile/my-profile-settings",
      headers: comm.getHeader(),
    }).then(response => {
        if(response.status == 200){
          this.settings.isPrivate = response.data.isPrivate
          this.settings.canReceiveMessageFromUnknown = response.data.canReceiveMessageFromUnknown
          this.settings.canBeTagged = response.data.canBeTagged
        }
      });
  },

  methods:{
    updateSettings(){
      axios({
      method: 'put',
      url: "http://" + comm.server + "/api/profile/my-profile-settings",
      headers: comm.getHeader(),
      data: JSON.stringify(this.settings)
      }).then(response => {
        if(response.status == 200){
          console.log(response.data)
        }
      });
    }
  }
}
</script>