<template>
  <v-card>
    <v-toolbar flat>
      <template v-slot:extension>
        <v-tabs v-model="tabs" fixed-tabs>
          <v-tabs-slider></v-tabs-slider>
           <v-tab v-for="(reaction, index) in reactions" :key="index" :href="'#' + reaction" class="primary--text" @click="showPosts(reaction)">
            <v-icon>{{reactionIcons[index]}}</v-icon>
          </v-tab>
        </v-tabs>
      </template>
    </v-toolbar>

    <v-tabs-items v-model="tabs">
      <v-tab-item v-for="(reaction, index) in reactions" :key="index" :value="reaction" >
        <v-card flat>
          <v-card-text>
            <v-row>
                <v-col cols="12" sm="4" v-for="p in posts" :key="p._id">
                    <post v-bind:usage="'MyReactions'" v-bind:post="p" v-bind:myReaction="reaction" />
                </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-tab-item>
    </v-tabs-items>
  </v-card>
</template>

<script>
  import Post from '../components/Posts/Post.vue'
  import axios from 'axios'
  import * as comm from '../configuration/communication.js'
  export default {
    components:{Post},
    data () {
      return {
        reactions: ['like','dislike'],
        reactionIcons: ['mdi-thumb-up','mdi-thumb-down'],
        posts:[],
        tabs: null,}
    },
    created(){
        if(!this.isPageAvailable()) {
          this.$router.push({name: 'NotFound'})
          return;
        }
        this.showPosts('like');
    },
    methods:{
        showPosts(reaction){
          axios({
            method: 'get',
            url: comm.protocol + "://" + comm.server + "/api/postreaction/my-reactions/" + reaction,
            headers: comm.getHeader(),
          }).then(response => {
            if(response.status == 200){
              this.posts = response.data.collection;
            }
          });
        },
        isPageAvailable(){
            return comm.isUserLogged()
        }
    }
  }
</script>