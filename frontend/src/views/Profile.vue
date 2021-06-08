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
        <v-row align="center" justify="center">
            <v-col cols="12" sm="4" v-for="p in posts" :key="p._id">
                <v-card justify="center" align="center"
                    outlined
                    width="600"
                >
                <v-carousel>
                
                <v-template v-for="item in p.medias" :key="item.filePath">
                      <v-carousel-item
                      reverse-transition="fade-transition"
                      transition="fade-transition">
                      <video autoplay loop style="object-fit:contain;" :src="'http://' + server + '/static/data/' + item.filePath" v-if="item.filePath.includes('mp4')">
                        Your browser does not support the video tag.
                      </video>
                      <img style="object-fit:contain;" :src="'http://' + server + '/static/data/' + item.filePath" v-if="!item.filePath.includes('mp4')">

                      </v-carousel-item>
             </v-template>
             </v-carousel>
                </v-card>
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
      posts: [],
      username: '',
      server: comm.server
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
            this.username = comm.getUrlVars()['username'];
            if(this.isMyProfile){
                axios({
                method: "get",
                url: 'http://' + comm.server + '/api/post/my'
            }).then(response => {
              if(response.status==200){
                  this.posts = response.data.collection;
              }
            })

            }else{
                 axios({
                method: "get",
                url: 'http://' + comm.server + '/api/post/profile/' + this.username,
            }).then(response => {
              if(response.status==200){
                  this.posts = response.data.collection;
              }
            })
            }
        },
        getProfileRequests(){

        }
    },
    mounted(){
       
    }
}
</script>

<style>

</style>