<template>
    <v-container>
        <v-row align="left" >
            <v-col cols="12" sm="12" >
                <personal-data v-on:loaded-user='profileLoaded($event)' style="height:200px"   />
                <v-btn v-if="!isMyProfile"
                color="warning"
                elevation="8"
                @click="follow"
                >
                Follow
                </v-btn>
                <follow-requests v-if="isMyProfile"/>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
import PersonalData from '../components/PersonalData.vue'
import FollowRequests from '../components/FollowRequests.vue'
import axios from 'axios'
import * as comm from '../configuration/communication.js'
export default {
    components: {
        PersonalData,
        FollowRequests,
    },
    data: () => ({
      isMyProfile: false,
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
        },
        profileLoaded(loadedProfileID){
            this.profileId = loadedProfileID;
            this.isMyProfile = comm.getLoggedUserID() == this.profileId;
        },
        getProfileRequests(){

        }
    }
}
</script>

<style>

</style>