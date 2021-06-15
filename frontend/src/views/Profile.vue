<template>
    <v-container>
        <post-modal :visible="showFollowOption" @close="showFollowOption=false"/>
        <v-row align="left" >
            <v-col cols="12" sm="11" >
                <personal-data v-on:loaded-user='profileLoaded($event)' style="height:200px" v-bind:username="username"/>
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
                <v-btn v-if="!isMyProfile" class="mx-2" fab dark small color="cyan" @click="showFollowOptionDialog()">
                    <v-icon>mdi-menu-down</v-icon>
                </v-btn>
            </v-col>
        </v-row>
        <v-row align="center" justify="center">
            <v-col cols="12" sm="4" v-for="p in posts" :key="p._id">
               <post v-bind:usage="'Explore'" v-bind:post="p" />
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
import PersonalData from '../components/PersonalData.vue'
import FollowRequests from '../components/FollowRequests.vue'
import Post from '../components/Posts/Post.vue'
import axios from 'axios'
import * as comm from '../configuration/communication.js'
import Post from '../components/Posts/Post.vue'

export default {
    components: {
        PersonalData,
        FollowRequests,
        Post},
    props: ['username'],
    data() {return {
      isMyProfile: false,
      profileId: 1,
      posts: [],
      server: comm.server,
      protocol: comm.protocol,
      showFollowOption: false
    }},
    methods: {
        follow(){
            axios({
                method: "post",
                url: comm.protocol + '://' + comm.server + '/api/connection/following/request/' + this.profileId,
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
                url: comm.protocol + '://' + comm.server + '/api/post/my',
                headers: comm.getHeader(),
            }).then(response => {
              if(response.status==200){
                  this.posts = response.data.collection;
              }
            })

            }else{
                axios({
                method: "get",
                url: comm.protocol + '://' + comm.server + '/api/post/profile/' + this.username,
            }).then(response => {
              if(response.status==200){
                  this.posts = response.data.collection;
              }
            })
            }
        },
        showFollowOptionDialog(){
            this.showFollowOption = true
        },
        redirectToCreatePost(){
            this.$router.push({name:'Post'});
        }
    },
}
</script>

<style>

</style>