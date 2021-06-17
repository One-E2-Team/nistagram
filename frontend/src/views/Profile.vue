<template>
    <v-container>
        <v-row align="left" >
            <v-col cols="12" sm="11" >
                <personal-data v-on:loaded-user='profileLoaded($event)' style="height:200px" v-bind:username="username"/>
                <v-btn v-if="!isMyProfile && !isFollowed" color="warning" elevation="8" @click="follow">
                    Follow
                </v-btn>
                <v-btn v-if="!isMyProfile && isFollowed" color="normal" elevation="8" @click="unfollow">
                    Unfollow
                </v-btn>
                <follow-requests v-if="isMyProfile"/>
                <v-btn v-if="isMyProfile"
                color="normal"
                elevation="8"
                @click="redirectToCreatePost()"
                >
                Create post
                </v-btn>
                <profile-options-drop-menu v-if="!isMyProfile" v-bind:profileId="profileId" class="mx-2">
                    <v-icon>mdi-menu-down</v-icon>
                </profile-options-drop-menu>
            </v-col>
        </v-row>
        <v-row>
            <v-col cols="12" sm="4" v-for="p in posts" :key="p._id">
               <post v-bind:usage="'Profile'" v-bind:post="p" />
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
import ProfileOptionsDropMenu from '../components/DropMenu/ProfileOptionsDropMenu.vue'
export default {
    components: {
        PersonalData,
        FollowRequests,
        Post,
        ProfileOptionsDropMenu},
    props: ['username'],
    data() {
        return {
            isMyProfile: false,
            profileId: 0,
            posts: [],
            server: comm.server,
            protocol: comm.protocol,
            showFollowOption: false,
            isFollowed: false,
        }
    },
    methods: {
        follow(){
            axios({
                method: "post",
                url: comm.protocol + '://' + comm.server + '/api/connection/following/request/' + this.profileId,
                headers: comm.getHeader(),
            }).then(response => {
              if (response.status==200){
                  alert('Success');
              }
            })
        },
        unfollow(){
            alert('I need endpoint to do that you know igor fy_master -_- ')
            //TODO: send axios for unfollow and when response status is 200 then this.isFollowed set on false
        },
        profileLoaded(loadedProfileID){
            this.profileId = loadedProfileID;
            this.isMyProfile = comm.getLoggedUserID() == loadedProfileID;
            if (this.isMyProfile){
                this.loadMyPosts();
            } else {
                this.loadPostsFromUsername();
                this.loadIfUserAreFollowed();
            }
        },
        showFollowOptionDialog() {
            this.showFollowOption = true;
        },
        redirectToCreatePost() {
            this.$router.push({name:'Post'});
        },
        loadMyPosts(){
            axios({ method: "get", 
                    url: comm.protocol + '://' + comm.server + '/api/post/my',
                    headers: comm.getHeader()
                    }).then(response => {
                        if (response.status==200)
                            this.posts = response.data.collection;  
                        });
        },
        loadPostsFromUsername(){
            axios({ method: "get",
                    url: comm.protocol + '://' + comm.server + '/api/post/profile/' + this.username, 
                    headers: comm.getHeader(),
                    }).then(response => {
                        if(response.status==200)
                            this.posts = response.data;
                        });
        },
        loadIfUserAreFollowed(){
            axios({ method: "get",
                    url: comm.protocol + '://' + comm.server + '/api/connection/following/my-properties/' + this.profileId, 
                    headers: comm.getHeader(),
                    }).then(response => {
                        if(response.status==200)
                            this.isFollowed = response.data;
                        });
        }
    },
}
</script>

<style>

</style>