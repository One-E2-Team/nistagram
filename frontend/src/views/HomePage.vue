<template>
    <v-container fluid>
      <v-sheet class="mx-auto" elevation="1" max-width="900" >
        <v-slide-group class="pa-4" >
          <v-slide-item v-for="s in stories" :key="s._id" >
            <div class="mx-3">
              <show-story-modal  :post="s.post"/>
              <h3>{{s.post.publisherUsername}}</h3>
            </div>
          </v-slide-item>
        </v-slide-group>
      </v-sheet>
        <v-row>
          <v-col></v-col>
        </v-row>
        <v-row justify="center" align="center" v-for="(p, index) in posts" :key="index">
            <v-col cols="12" sm="6">
                <post v-bind:usage="'HomePage'" v-bind:post="p.post" v-bind:myReaction="p.reaction" :campaignData="campaignData[index]"/>
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
      campaignData: [],
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
        this.campaignData = [];
        for (let r of response) {
          let data = {
            campaignId: r.campaignId,
            influencerId: r.influencerId,
            influencerUsername: r.influencerUsername,
          }
          this.campaignData.push(data);
        }
      },
    },
  }
</script>