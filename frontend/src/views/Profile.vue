<template>
    <v-container>
        <v-row align="left" >
            <v-col cols="12" sm="11">
                <personal-data v-on:loaded-user='profileLoaded($event)' style="height:200px" v-bind:username="username"/>
                <template v-if="isUserLoggedIn()">
                    <v-btn v-if="!isMyProfile && followType==followTypeValues[1]" color="warning" elevation="8" @click="follow">
                        Follow
                    </v-btn>
                    <v-btn v-if="!isMyProfile && followType!=followTypeValues[1]" color="normal" elevation="8" @click="unfollow">
                        {{unfollowButtonText}}
                    </v-btn>
                    <follow-requests v-if="isMyProfile"/>
                    <v-btn v-if="isMyProfile" color="normal" elevation="8" @click="redirectToCreatePost()">
                    Create post
                    </v-btn>
                    <profile-options-drop-menu v-if="!isMyProfile" 
                        v-bind:profileId="profileId" v-bind:conn="connection" v-bind:blocked="isBlocked" 
                        v-on:connectionChanged='connection=$event' v-on:blockChanged='isBlocked=$event' class="mx-2">
                            <v-icon>mdi-menu-down</v-icon>
                    </profile-options-drop-menu>
                </template>
            </v-col>
        </v-row>
        <v-row v-if="followType == followTypeValues[0] || !isPrivateProfile || isMyProfile">
            <v-col cols="12" sm="4" v-for="p in posts" :key="p._id">
               <post v-bind:usage="'Profile'" v-bind:post="p.post" v-bind:myReaction="p.reaction" />
            </v-col>
        </v-row>
        <v-row v-else-if="followType != followTypeValues[0] && isPrivateProfile">
            <v-col cols="12" sm="4">
               <h3 class="display-2 font-weight-bold mb-3"> This profile is private !</h3>
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
            //isFollowed: false,
            followTypeValues: ['following', 'not_following', 'sent_request'],
            followType: '',
            isPrivateProfile: true,
            unfollowButtonText: '',
            isBlocked: false,
            connection: null,
        }
    },
    methods: {
        profileLoaded(loadedProfile){
            this.profileId = loadedProfile.ID;
            this.isPrivateProfile = loadedProfile.profileSettings.isPrivate
            this.isMyProfile = comm.getLoggedUserID() == loadedProfile.ID;
            if (this.isMyProfile){
                this.loadMyPosts();
            } else {
                this.loadPostsFromUsername();
                if (comm.isUserLogged()) {
                    this.checkIsUserBlocked()
                    this.loadUsersConnection();
                }
            }
        },
        follow(){
            axios({
                method: "post",
                url: comm.protocol + '://' + comm.server + '/api/connection/following/request/' + this.profileId,
                headers: comm.getHeader(),
            }).then(response => {
                if (response.status==200) {
                    this.prepareFollowButtons(response.data);
                    if (!response.data.approved) {
                        this.connection = null;
                    } else {
                        this.connection = response.data;
                    }
                }
            })
        },
        unfollow(){
            axios({
                method: "put",
                url: comm.protocol + '://' + comm.server + '/api/connection/unfollow/' + this.profileId,
                headers: comm.getHeader(),
            }).then(response => {
              if (response.status==200 && response.data.status == 'ok'){
                  this.followType = this.followTypeValues[1];
                  this.connection = null;
                  this.posts = [];
              }
            })
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
                headers: comm.getHeader(),
            }).then(response => {
                if (response.status==200)
                    this.posts = response.data.collection;  
            });
        },
        isUserLoggedIn(){
            return comm.isUserLogged();
        },
        loadPostsFromUsername(){
            axios({ method: "get",
                    url: comm.protocol + '://' + comm.server + '/api/post/profile/' + this.username, 
                    headers: comm.getHeader(),
                    }).then(response => {
                        if(response.status==200)
                            this.posts = response.data.collection;
                    });
        },
        loadUsersConnection(){
            axios({ method: "get",
                url: comm.protocol + '://' + comm.server + '/api/connection/following/my-properties/' + this.profileId, 
                headers: comm.getHeader(),
                }).then(response => {
                    if(response.status==200) {
                        this.prepareFollowButtons(response.data);
                        if (!response.data.approved) {
                            this.connection = null;
                        } else {
                            this.connection = response.data;
                        }
                    }
                });
        },
        prepareFollowButtons(responseData) {
            if (responseData == null) { //TODO: improve response
                this.followType = this.followTypeValues[1];
                return;
            }
            if (!responseData.approved && responseData.connectionRequest) {
                this.unfollowButtonText = 'Cancel follow request';
                this.followType = this.followTypeValues[2];
            } else if (!responseData.approved) {
                this.followType = this.followTypeValues[1];
            } else {
                this.unfollowButtonText = 'Unfollow';
                this.followType = this.followTypeValues[0];
            }
        },
        checkIsUserBlocked(){
            axios({
                    method: "get",
                    url: comm.protocol + "://" + comm.server +"/api/connection/block/" + this.profileId,
                    headers: comm.getHeader(),
                }).then((response) => {
                console.log(response.data);
                if(response.status == 200)
                this.isBlocked = response.data.blocked
                }).catch((error) => {
                    console.log(error);
                });
        }
    },
}
</script>

<style>

</style>