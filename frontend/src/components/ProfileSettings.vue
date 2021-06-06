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
          v-model="isPrivate"
          :label="'Private account'"
        ></v-checkbox>

        <v-checkbox
          v-model="canReceiveMessageFromUnknown"
          :label="`Can be tagged on posts`"
        ></v-checkbox>

        <v-checkbox
          v-model="canBeTagged"
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
      isPrivate:false,
      canReceiveMessageFromUnknown: false,
      canBeTagged: false
    }
  },

  created(){
    //TODO: prezumi podesavanja trenutnog korisnika
    /*axios.get(comm.server)
      .then(response => {
        if(response.status = 200){
          this.isPrivate = response.data.isPrivate
          this.canReceiveMessageFromUnknown = response.data.canReceiveMessageFromUnknown
          this.canBeTagged = response.data.canBeTagged
        }
      })*/
  },

  methods:{
    updateSettings(){
      axios.put(""+comm.server)
      .then(response => {
        if(response.status == 200){
          this.isPrivate = response.data.isPrivate
          this.canReceiveMessageFromUnknown = response.data.canReceiveMessageFromUnknown
          this.canBeTagged = response.data.canBeTagged
        }
      })
    }
  }
}
</script>