<template>
    <v-container>
        <v-row>
            <v-col cols="12" sm="11">
                <personal-data v-on:loaded-user='profileLoaded($event)' style="height:200px" v-bind:username="username"/>
                <template v-if="isUserLoggedIn()">
                    <template v-if="!isBlocked">
                        <v-btn v-if="!isMyProfile() && followTypeValue==followType.NOT_FOLLOW" color="warning" elevation="8" @click="follow">
                            Follow
                        </v-btn>
                        <v-btn v-if="!isMyProfile() && followTypeValue!=followType.NOT_FOLLOW" color="normal" elevation="8" @click="unfollow">
                            <p v-if="followTypeValue==followType.REQUEST_SENDED" class="my-3">Cancel follow request</p>
                            <p v-else-if="followTypeValue==followType.FOLLOW" class="my-3">Unfollow</p>
                        </v-btn>
                        <v-btn v-if="isMyProfile()" color="normal" elevation="8" @click="redirectToCreatePost()">
                            Create post
                        </v-btn>
                    </template>
                    <profile-options-drop-menu v-if="!isMyProfile()" 
                        :profileId="profile.ID" :conn="connection" :blocked="isBlocked" :msgConn="messageConnection"
                        v-on:connectionChanged='connectionChanged($event)' v-on:blockChanged='isBlocked=$event' 
                        v-on:messageRequestSended='messageConnection=$event'
                        class="mx-2">
                            <v-icon>mdi-menu-down</v-icon>
                    </profile-options-drop-menu>
                </template>
            </v-col>
        </v-row>
        <v-row>
            <v-col cols="12" sm="4">
                <v-slide-group class="pa-4" >
                    <v-slide-item v-for="s in stories" :key="s._id" >
                        <div class="mx-3">
                            <show-story-modal  :post="s.post"/>
                            <h3>{{s.post.publisherUsername}}</h3>
                        </div>
                    </v-slide-item>
                </v-slide-group>
            </v-col>
        </v-row>
        <v-row v-if="followTypeValue == followType.FOLLOW || !isProfilePrivate || isMyProfile()">
            <v-col cols="12" sm="4" v-for="p in posts" :key="p._id">
               <post v-bind:usage="'Profile'" v-bind:post="p.post" v-bind:myReaction="p.reaction" />
            </v-col>
        </v-row>
        <v-row v-else-if="followTypeValue != followType.FOLLOW && isProfilePrivate">
            <v-col cols="12" sm="4">
               <h3 class="display-2 font-weight-bold mb-3"> This profile is private !</h3>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
import PersonalData from '../components/PersonalData.vue'
import Post from '../components/Posts/Post.vue'
import axios from 'axios'
import * as comm from '../configuration/communication.js'
import ProfileOptionsDropMenu from '../components/DropMenu/ProfileOptionsDropMenu.vue'
import ShowStoryModal from '../modals/showStoryModal.vue'
export default {
    components: {
        PersonalData,
        Post,
        ProfileOptionsDropMenu,
        ShowStoryModal},
    props: ['username'],
    data() {
        return {
            profile: {},
            isProfilePrivate: null,
            posts: [],
            stories: [],
            server: comm.server,
            protocol: comm.protocol,
            showFollowOption: false,
            followTypeValue: -1,
            isBlocked: false,
            connection: null,
            messageConnection: null,
            followType : { FOLLOW : 0, NOT_FOLLOW : 1, REQUEST_SENDED : 2 }
        }
    },
    methods: {
        profileLoaded(loadedProfile){
            this.posts = [];
            this.stories = [];
            this.profile = loadedProfile
            this.isProfilePrivate = loadedProfile.profileSettings.isPrivate
            
            if (this.isMyProfile()){
                this.loadMyPosts();
                return
            }
    
            this.loadConnectionAndPosts();
            
            if (comm.isUserLogged()) {
                this.checkIsUserBlocked();
                this.checkMessageConnection();
            }
            
        },
        loadConnectionAndPosts(){
             axios({ method: "get",
                    url: comm.protocol + '://' + comm.server + '/api/connection/following/my-properties/' + this.profile.ID, 
                    headers: comm.getHeader(),
                }).then(response => {
                    if(response.status==200) {
                        this.prepareFollowButtons(response.data);
                        if (response.data && !response.data.approved) {
                            this.connection = null;
                        } else {
                            this.connection = response.data;
                        }
                        if (!this.isProfilePrivate || (this.connection != null && this.connection.approved))
                            this.loadPostsFromUsername()
                    }
                });
        },
        isMyProfile(){
            return comm.getLoggedUserID() == this.profile.ID;
        },
        follow(){
            axios({
                method: "post",
                url: comm.protocol + '://' + comm.server + '/api/connection/following/request/' + this.profile.ID,
                headers: comm.getHeader(),
            }).then(response => {
                if (response.status==200) {
                    this.prepareFollowButtons(response.data);
                    if (response.data && !response.data.approved) {
                        this.connection = null;
                    } else {
                        this.connection = response.data;
                        this.messageConnection = {};
                        this.messageConnection.notifyMessage = true;
                    }
                }
            })
        },
        unfollow(){
            axios({
                method: "put",
                url: comm.protocol + '://' + comm.server + '/api/connection/unfollow/' + this.profile.ID,
                headers: comm.getHeader(),
            }).then(response => {
              if (response.status==200 && response.data.status == 'ok'){
                  this.followTypeValue = this.followType.NOT_FOLLOW;
                  this.connection = null;
                  if (this.isProfilePrivate) {
                    this.posts = [];
                    this.stories = [];
                  }
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
                    this.setPostAndStories(response.data.collection)
            });
        },
        isUserLoggedIn(){
            return comm.isUserLogged();
        },
        loadPostsFromUsername(){
            axios({ method: "get",
                    url: comm.protocol + '://' + comm.server + '/api/post/profile/' + this.profile.username, 
                    headers: comm.getHeader(),
                    }).then(response => {
                        if(response.status==200){   
                            this.setPostAndStories(response.data.collection)
                        }
                    });
        },
        prepareFollowButtons(responseData) {
            if (responseData == null) {
                this.followTypeValue = this.followType.NOT_FOLLOW;
                
            }
            else if (!responseData.approved && responseData.connectionRequest) {
                this.followTypeValue = this.followType.REQUEST_SENDED
            } else if (!responseData.approved) {
                this.followTypeValue = this.followType.NOT_FOLLOW;
            } else {
                this.followTypeValue = this.followType.FOLLOW;
            }
        },
        checkIsUserBlocked(){
            axios({
                    method: "get",
                    url: comm.protocol + "://" + comm.server +"/api/connection/block/" +this.profile.ID,
                    headers: comm.getHeader(),
                }).then((response) => {
                if(response.status == 200)
                this.isBlocked = response.data.blocked
                }).catch((error) => {
                    console.log(error);
                });
        },
        checkMessageConnection(){
             axios({
                    method: "get",
                    url: comm.protocol + "://" + comm.server +"/api/connection/messaging/my-properties/" + this.profile.ID,
                    headers: comm.getHeader(),
                }).then((response) => {
                if(response.status == 200){
                   this.messageConnection = response.data
                }
                }).catch((error) => {
                    console.log(error);
                });
        },

        setPostAndStories(posts){
            for (let p of posts){
                if(p.post.postType == 1){
                    this.stories.push(p);
                }
                else if(p.post.postType == 2){ 
                    this.posts.push(p);
                }
            }
        },
        connectionChanged(newConnection) {
            if (newConnection == null) {
                this.followTypeValue = this.followType.NOT_FOLLOW;
            }
            this.connection = newConnection;
        }
    },
}
</script>

<style>

</style>