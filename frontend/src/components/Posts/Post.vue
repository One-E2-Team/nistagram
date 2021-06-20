<template>
  <v-card
    class="mx-auto" :width="width"
  >
    <post-modal v-if="showTitle" :visible="showDialog" @close="showDialog=false" v-bind:post="post"/>
    <v-list-item v-if="showTitle">
      <v-list-item-content >
        <v-list-item-title  class="text-h6 d-flex justify-space-between ">
          <router-link :to="{ name: 'Profile', params: { username: post.publisherUsername }}">{{post.publisherUsername}}</router-link>
          <v-btn dark icon @click="showDialog = true" v-if="isUserLogged">
            <v-icon color="blue">mdi-dots-horizontal</v-icon>
          </v-btn>
        </v-list-item-title>
      </v-list-item-content>
    </v-list-item>
    <v-carousel :width="width" :height="height">        
            <v-carousel-item
            v-for="item in post.medias" :key="item.filePath"
            reverse-transition="fade-transition"
            transition="fade-transition">
            <video autoplay loop :width="width" :height="height" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="item.filePath.includes('mp4')">
            Your browser does not support the video tag.
            </video>
            <img :width="width" :height="height" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="!item.filePath.includes('mp4')">

            </v-carousel-item>
    </v-carousel>

    <v-card-text class="text--primary">
       <v-container>
         <v-row>
          <v-col>Location: {{post.location}} </v-col>
         </v-row>
         <v-row>
          <v-col> {{post.description}} </v-col>
         </v-row>
         <v-row>
          <v-col class="d-flex justify-space-around ">
            <v-btn-toggle v-if="isUserLogged" v-model="reaction" color="primary" group dense>
              <v-btn :value="'like'" class="ma-2" text icon @click="react('like')">
                <v-icon>mdi-thumb-up</v-icon>
              </v-btn>
              <v-btn :value="'dislike'" class="ma-2" text icon @click="react('dislike')">
                <v-icon>mdi-thumb-down</v-icon>
              </v-btn>
            </v-btn-toggle>
            <v-item-group v-else color="primary" group dense class="v-btn-toggle">
              <v-btn :value="'like'" class="ma-2" text icon @click="react('like')">
                <v-icon>mdi-thumb-up</v-icon>
              </v-btn>
              <v-btn :value="'dislike'" class="ma-2" text icon @click="react('dislike')">
                <v-icon>mdi-thumb-down</v-icon>
              </v-btn>
            </v-item-group>
          </v-col>
         </v-row>
         <v-row cols="12" md="6">
            <v-col>
              <v-textarea solo placeholder="Enter comment..." rows="4" v-model="comment"></v-textarea>
              <v-btn color="normal" elevation="2" @click="commentPost()">
                Comment
              </v-btn>
            </v-col>
          </v-row>
       </v-container>
    </v-card-text>
  </v-card>
</template>

<script>
import PostModal from '../../modals/PostModal.vue'
import * as comm from '../../configuration/communication.js'
import axios from 'axios'
export default {
  components: { PostModal },
  name: 'Post',
  props: ['post','usage', 'myReaction'],
  data() {
    return {
      showDialog : false,
      width: 300,
      height: 200,
      showTitle: false,
      protocol: comm.protocol,
      server: comm.server,
      reaction: null,
      isUserLogged: comm.isUserLogged(),
      comment: '',
    }
  },
  mounted() {
    this.designView();
    if (this.myReaction == 'none') {
      this.reaction = null;
      return;
    }
    this.reaction = this.myReaction;
  },
  methods: {
    designView() {
      if (this.usage == 'Profile') {
        this.width = 300;
        this.height = 400;
        this.showTitle = false;
      } else if (this.usage == 'Explore') {
        this.width = 300;
        this.height = 400;
        this.showTitle = true;
      } else if(this.usage == 'HomePage') {
        this.width = 600;
        this.height = 700;
        this.showTitle = true;
      } else if(this.usage == 'MyReactions'){
        this.width = 300;
        this.height = 400;
        this.showTitle = true;
      }
    }, 
    react (reactionType) {
      if (this.preventActionIfUnauthorized()) {
        return;
      }
      if (reactionType == this.reaction){
        axios({
          method: 'delete',
          url: comm.protocol + '://' + comm.server + '/api/postreaction/react/' + this.post.id,
          headers: comm.getHeader(),
        }).then(response => {
          console.log(response.data);
        });
      } else {
        let dto = {'postId' : this.post.id, 'reactionType' : reactionType}
        axios({
          method: 'post',
          url: comm.protocol + '://' + comm.server + '/api/postreaction/react',
          data: JSON.stringify(dto),
          headers: comm.getHeader(),
        }).then(response => {
          console.log(response.data);
        });
      }
    },
    commentPost() {
      if (this.preventActionIfUnauthorized()) {
        return;
      }
      let dto = {'postId' : this.post.id, 'content' : this.comment}
      axios({
        method: 'post',
        url: comm.protocol + '://' + comm.server + '/api/postreaction/comment',
        data: JSON.stringify(dto),
        headers: comm.getHeader(),
      }).then(response => {
        console.log(response.data);
      });
    },
    preventActionIfUnauthorized() {
      if(!comm.isUserLogged()){
        alert('You must be logged to react on post');
        this.comment = '';
        if (this.isUserLogged) {
          this.$router.go();
        }
        return true;
      }
      return false;
    },
  },
  watch: {
    usage(){
      this.designView();
    }
  },
}
</script>
