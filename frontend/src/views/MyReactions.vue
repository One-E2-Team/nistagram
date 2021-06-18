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
                    <post v-bind:usage="'MyReactions'" v-bind:post="p" />
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
  import * as comm from '../configuration/communication.js'
  export default {
    comments:{Post},
    data () {
      return {
        reactions: ['like','dislike'],
        reactionIcons: ['mdi-thumb-up','mdi-thumb-down'],
        posts:[],
        tabs: null,}
    },
    created(){
        if( this.isPageAvailabe() )
            this.$router.push({name: 'NotFound'})
        this.showPosts('likes')
    },
    methods:{
        showPosts(reaction){
            console.log(reaction)
            //send axios for reactions and post enter in this.posts
        },
        isPageAvailabe(){
            return !comm.isUserLogged()
        }
    }
  }
</script>