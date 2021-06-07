<template>
    <v-container>
        <v-row align="left" >
            <v-col cols="12" sm="12" >
                <personal-data v-on:loaded-user='profileId = $event' style="height:200px"   />
                <v-btn v-if="!isMyProfile"
                color="warning"
                elevation="8"
                @click="follow"
                >
                Follow
                </v-btn>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
import PersonalData from '../components/PersonalData.vue'
import axios from 'axios'
import * as comm from '../configuration/communication.js'
export default {
    components: {
        PersonalData
    },
    data: () => ({
      isMyProfile: true,
      profileId: 1,
    }),
    methods: {
        follow(){
            axios({
                method: "post",
                url: 'http://' + comm.server + '/api/connection/following/request/' + this.profileId,
            }).then(response => {
              if(response.status==200){
                  alert('Success');
              }
            })
        }
    }
}
</script>

<style>

</style>