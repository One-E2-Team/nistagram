<template>
    <v-container fluid>
      <search  v-on:searched-result='posts = $event' />

      <v-sheet
        class="mx-auto"
        elevation="1"
        max-width="900"
      >
    <v-slide-group
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
                      <video autoplay  width="600" height="500" :src="item.filePath" v-if="item.filePath.includes('mp4')">
                        Your browser does not support the video tag.
                      </video>
                      <img width="600" height="500" :src="item.filePath" v-if="!item.filePath.includes('mp4')">

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
        <v-row justify="center" align="center" v-for="p in posts" :key="p._id">
            <v-col cols="12" sm="4" v-if="p.postType == 2">
                <v-card justify="center" align="center"
                    outlined
                    width="600"
                >
                <v-card-title>{{p.publisherUsername}}</v-card-title>
                <v-carousel>
             <v-template v-for="item in p.medias" :key="item.filePath">
                      <v-carousel-item
                      reverse-transition="fade-transition"
                      transition="fade-transition">
                      <video autoplay  width="600" height="500" :src="item.filePath" v-if="item.filePath.includes('mp4')">
                        Your browser does not support the video tag.
                      </video>
                      <img width="600" height="500" :src="item.filePath" v-if="!item.filePath.includes('mp4')">

                      </v-carousel-item>
             </v-template>
          </v-carousel>
                <v-card-text>{{p.description}} {{p.hashTags}} {{p.location}}</v-card-text>
                <v-card-text>{{p.publishDate}}</v-card-text>
                </v-card>
             </v-col>
        </v-row>
  </v-container>
</template>

<script>
import axios from 'axios'
import * as comm from '../configuration/communication.js'
import Search from './Search.vue'
  export default {
    name: 'Explore',
    components:{
      Search
    },
     mounted(){
        axios.get("http://" + comm.server +"/api/post/public").then((response) => {
            let res = response.data.collection;
            res.forEach((post) => {
                if(post.medias != null){
                  post.medias.forEach((media) =>{
                    console.log(media)
                   media.filePath = "http://" + comm.server +"/static/data/" + media.filePath;
                  });
                }
            });
            this.posts = res;
            this.allPosts = res;
    })
    .catch((error) => {
      console.log(error);
    });
    },

    data: () => ({
      posts : null,
      allPosts: null
    }),

    methods: {
      
      
    }
  }
</script>