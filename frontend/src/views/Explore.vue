<template>
    <v-container fluid>
      <search  v-on:searched-result='loadSearchResult($event)' />

      <v-sheet
        class="mx-auto"
        elevation="1"
        max-width="900"
      >
    <v-slide-group v-if="searchType == 'posts'" class="pa-4" >
      <v-slide-item v-for="s in stories" :key="s._id" >
        <div class="mx-3">
          <show-post-full-screen-modal :visible="showPostFullScreenModal" @close="showPostFullScreenModal=false" v-bind:post="s.post"/>
          <h3>{{s.post.publisherUsername}}</h3>
        </div>
      </v-slide-item>
    </v-slide-group>
  </v-sheet>
        <v-row>
          <v-col></v-col>
        </v-row>
        <template v-if="searchType == 'posts'">
        <v-row>
            <v-col cols="12" sm="3" v-for="p in posts" :key="p._id" >
               <post v-bind:usage="'Explore'" v-bind:post="p.post" v-bind:myReaction="p.reaction"/>
             </v-col>
        </v-row>
        </template>

        <template v-if="searchType=='accounts'">
          <v-card
            class="mx-auto"
            max-width="300"
            tile
          >
            <v-list rounded>
              <v-subheader>Search result</v-subheader>
              <v-list-item-group
                color="primary"
              >
                <v-list-item @click="redirect(item)"
                  v-for="(item, i) in usernames"
                  :key="i"
                >
                  <v-list-item-icon>
                     <v-icon>mdi-account</v-icon>
                  </v-list-item-icon>
                  <v-list-item-content>
                    <v-list-item-title v-text="item"></v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-card>
        </template>
  </v-container>
</template>

<script>
import axios from 'axios'
import * as comm from '../configuration/communication.js'
import Search from '../components/Search.vue'
import Post from '../components/Posts/Post.vue'
import ShowPostFullScreenModal from '../modals/showStoryModal.vue'
  export default {
    name: 'Explore',
    components:{
      Search,
      Post,
      ShowPostFullScreenModal
    },
     created(){
      axios({
        method: 'get',
        url: comm.protocol + '://' + comm.server +'/api/post/public',
        headers: comm.getHeader(),
      }).then(response => {
        if (response.status==200){
          this.setPostAndStories(response.data.collection);
        }
      }).catch((error) => {
        console.log(error);
      });
    },

    data() {
      return {
      posts : [],
      stories: [],
      searchType: "posts", //possible values: {accounts, posts}
      server: comm.server,
      protocol: comm.protocol,
      usernames:[],
      showPostFullScreenModal : false
    }},

    methods: {
      loadSearchResult(searchResult){
        this.searchType = searchResult.type;
        if(this.searchType == "posts"){
          this.posts = searchResult.collection;
          this.usernames = [];
        }else if(this.searchType == "accounts"){
          this.usernames = searchResult.collection;
          this.posts = [];
        }
      },
      redirect(username){
        this.$router.push({name: 'Profile', params: {username: username}})
      },
      openStory(story){
        console.log(story)
        this.storiesForUsername
      },
      setPostAndStories(posts){
            for (let p of posts){
                if(p.post.postType == 1){
                    this.stories.push(p)
                }
                else if(p.post.postType == 2){ 
                    this.posts.push(p)
                }
            }
            console.log(this.posts)
            console.log(this.stories)
        }
    },
  }
</script>