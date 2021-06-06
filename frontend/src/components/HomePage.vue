<template>
    <v-container fluid>
        <v-row justify="center" align="center" v-for="p in posts" :key="p._id">
            <v-col cols="12" sm="4">
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
        axios.get("http://" + comm.server +"api/post/").then((response) => {
            let res = response.data.collection;
            console.log(res);
            for(let p of res){
                for(let media of p.medias){
                    media.filePath = "http://" + comm.server +"/static/data/" + media.filePath;
                }
            }
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