<template>
    <v-container fluid>
      <search  v-on:searched-result='loadSearchResult($event)' />

      <v-sheet
        class="mx-auto"
        elevation="1"
        max-width="900"
      >
    <v-slide-group v-if="searchType == 'posts'"
      class="pa-4"
    >
      <v-slide-item
        v-for="p in posts" :key="p._id"
      >
        <v-card v-if="p.postType == 1"
          class="ma-4"
          height="200"
          width="100"
        >
          <v-row
            class="fill-height"
            align="center"
            justify="center"
          >
         <v-carousel>
             <v-template v-for="item in p.medias" :key="item.filePath">
                      <v-carousel-item
                      reverse-transition="fade-transition"
                      transition="fade-transition">
                      <video autoplay  width="100" height="200" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="item.filePath.includes('mp4')">
                        Your browser does not support the video tag.
                      </video>
                      <img width="100" height="200" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="!item.filePath.includes('mp4')">

                      </v-carousel-item>
             </v-template>
          </v-carousel>
          </v-row>
        </v-card>
      </v-slide-item>
    </v-slide-group>
  </v-sheet>
        <v-row>
          <v-col></v-col>
        </v-row>
        <template v-if="searchType == 'posts'">
        <v-row  >
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
  export default {
    name: 'Explore',
    components:{
      Search,
      Post
    },
     created(){
      axios({
        method: 'get',
        url: comm.protocol + '://' + comm.server +'/api/post/public',
        headers: comm.getHeader(),
      }).then(response => {
        if (response.status==200){
          this.posts = response.data.collection;
        }
      }).catch((error) => {
        console.log(error);
      });
    },

    data() {
      return {
      posts : [],
      searchType: "posts", //possible values: {accounts, posts}
      server: comm.server,
      protocol: comm.protocol,
      usernames:[]
    }},
    computed : {
      postsWithTypePost : function (){
        return this.posts.filter(function (item){
          return item.postType == 2;
        })
      }
    },

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
      }
    }
  }
</script>