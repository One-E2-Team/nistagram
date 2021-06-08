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
                      <v-carousel-item
                      v-for="(item,i) in p.medias"
                      :key="i"
                      :src= "item.filePath"
                      reverse-transition="fade-transition"
                      transition="fade-transition"
                      ></v-carousel-item>
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
                    <v-carousel-item
                    v-for="(item,i) in p.medias"
                    :key="i"
                    :src= "item.filePath"
                    reverse-transition="fade-transition"
                    transition="fade-transition"
                    ></v-carousel-item>
                </v-carousel>
                <v-card-text>{{p.description}}</v-card-text>
                <v-card-text>{{p.publishDate}}</v-card-text>
                </v-card>
             </v-col>
        </v-row>
  </v-container>
</template>

<script>
  import axios from 'axios'
  import * as comm from '../configuration/communication.js'
  export default {

    name: 'CreatePost',

    mounted(){
        axios({
                method: "post",
                url: "http://" + comm.server +"/api/post/homePage",
                headers: comm.getHeader(),
            }).then((response) => {
            let res = response.data.collection;
            res.forEach((post) => {
                if(post.medias != null){
                  post.medias.forEach((media) =>{
                    media.filePath = "http://" + comm.server +"/static/data/" + media.filePath;
                  });

                }
            });
            this.posts = res;
    })
    .catch((error) => {
      console.log(error);
    });
    },

    data: () => ({
      posts : null
    }),

    methods: {
      
    }
  }
</script>