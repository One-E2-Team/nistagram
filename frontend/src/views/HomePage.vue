<template>
    <v-container fluid>
      <v-sheet class="mx-auto" elevation="1" max-width="900" >
        <v-slide-group class="pa-4" >
          <v-slide-item v-for="(s, index) in stories" :key="index">
            <div class="mx-3">
              <show-story-modal  :post="s.post" :campaignData="campaignDataStory[index]"/>
              <h3 v-if="campaignDataStory == undefined || campaignDataStory[index].influencerUsername == ''">{{s.post.publisherUsername}}</h3>
              <h3 v-else-if="campaignDataStory[index].influencerUsername != ''">{{campaignDataStory[index].influencerUsername}}</h3>
            </div>
          </v-slide-item>
        </v-slide-group>
      </v-sheet>
        <v-row>
          <v-col></v-col>
        </v-row>
        <v-row justify="center" align="center" v-for="(p, index) in posts" :key="index">
            <v-col cols="12" sm="6">
                <post v-bind:usage="'HomePage'" v-bind:post="p.post" v-bind:myReaction="p.reaction" :campaignData="campaignDataPost[index]"/>
             </v-col>
        </v-row>
  </v-container>
</template>

<script>
  import axios from 'axios'
  import * as comm from '../configuration/communication.js'
  import Post from '../components/Posts/Post.vue'
  import ShowStoryModal from '../modals/showStoryModal.vue'
  export default {
    components: { Post, ShowStoryModal },
    name: 'HomePage',
    created(){
      axios({
        method: "get",
        url: comm.protocol + "://" + comm.server +"/api/post/homePage",
        headers: comm.getHeader(),
      }).then((response) => {
        if(response.status == 200){
          this.setPostAndStories(response.data.collection)
          this.setCampaignData(response.data.collection);
        }
      })
      .catch((error) => {
        console.log(error);
      });
    },

    data() {return {
      posts : [],
      stories: [],
      server: comm.server,
      protocol: comm.protocol,
      campaignDataPost: [],
      campaignDataStory: [],
    }},

    methods: {
      setPostAndStories(posts){
        for (let p of posts) {
          if(p.post.postType == 1) {
            this.stories.push(p);
          }
          else if(p.post.postType == 2) { 
            this.posts.push(p);
          }
        }
      },
      setCampaignData(response) {
        this.campaignDataPost = [];
        this.campaignDataStory = [];
        for (let r of response) {
          let data = {
            campaignId: r.campaignId,
            influencerId: r.influencerId,
            influencerUsername: r.influencerUsername,
          }
          if (r.post.postType == 1) {
            this.campaignDataStory.push(data);
          } else if (r.post.postType == 2) {
            this.campaignDataPost.push(data);
          }
        }
      },
    },
  }
</script>