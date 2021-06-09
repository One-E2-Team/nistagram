<template>
    <v-container>
        <v-row align="left" >
            <v-col cols="12" sm="11" >
                <personal-data v-on:loaded-user='profileLoaded($event)' style="height:200px"   />
                <v-btn v-if="!isMyProfile"
                color="warning"
                elevation="8"
                @click="follow"
                >
                Follow
                </v-btn>
                <follow-requests v-if="isMyProfile"/>
                <v-btn v-if="isMyProfile"
                color="normal"
                elevation="8"
                @click="redirectToCreatePost()"
                >
                Create post
                </v-btn>
            </v-col>
        </v-row>
        <v-row align="center" justify="center">
            <v-col cols="12" sm="4" v-for="p in posts" :key="p._id">
                <v-card justify="center" align="center"
                    outlined
                    width="600"
                >
                <v-carousel>
                
                <v-template v-for="item in p.medias" :key="item.filePath" name="temp">
                      <v-carousel-item
                      reverse-transition="fade-transition"
                      transition="fade-transition">
                      <video autoplay loop width="600" height="500" :src="'http://' + server + '/static/data/' + item.filePath" v-if="item.filePath.includes('mp4')">
                        Your browser does not support the video tag.
                      </video>
                      <img width="600" height="500" :src="'http://' + server + '/static/data/' + item.filePath" v-if="!item.filePath.includes('mp4')">

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
    props: ['username'],
    data: () => ({
      isMyProfile: false,
      profileId: 1,
      posts: [],
      server: comm.server
    }),
    methods: {
        follow(){
            axios({
                method: "post",
                url: 'http://' + comm.server + '/api/connection/following/request/' + this.profileId,
                headers: comm.getHeader(),
            }).then(response => {
              if(response.status==200){
                  alert('Success');
              }
            })
        },
        profileLoaded(loadedProfileID){
            this.profileId = loadedProfileID;
            this.isMyProfile = comm.getLoggedUserID() == loadedProfileID;
            if(this.isMyProfile){
                axios({
                method: "get",
                url: 'http://' + comm.server + '/api/post/my',
                headers: comm.getHeader(),
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

        },
        redirectToCreatePost(){
            this.$router.push({name:'Post'});
        }
    },
    mounted(){
       
    }
}
</script>

<style>

</style>