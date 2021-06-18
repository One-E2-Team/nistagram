<template>
    <v-container fluid>
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
                      <video autoplay loop width="100" height="200" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="item.filePath.includes('mp4')">
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
        <v-row justify="center" align="center" v-for="p in posts" :key="p._id">
            <v-col cols="12" sm="4" v-if="p.post.postType == 2">
                <post v-bind:usage="'HomePage'" v-bind:post="p.post" v-bind:myReaction="p.reaction" />
             </v-col>
        </v-row>
  </v-container>
</template>

<script>
  import axios from 'axios'
  import * as comm from '../configuration/communication.js'
  import Post from './Posts/Post.vue'
  export default {
  components: { Post },

    name: 'HomePage',

    created(){
      axios({
        method: "get",
        url: comm.protocol + "://" + comm.server +"/api/post/homePage",
        headers: comm.getHeader(),
      }).then((response) => {
        this.posts = response.data.collection;
      })
      .catch((error) => {
        console.log(error);
      });
    },

    data() {return {
      posts : null,
      server: comm.server,
      protocol: comm.protocol
    }},

    methods: {
      
    }
  }
</script>