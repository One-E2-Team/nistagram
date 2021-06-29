<template>
  <v-card class="mx-auto" :width="width+50" elevation="24" outlined >    
      <show-post-modal v-if="post.postType == 1"  :width="width" :height="height" :post="post" :reaction="reaction" v-on:reactionChanged="react($event)"/>
      <!--else => show story modal -->
    <v-card-text class="text--primary">
       <v-container>
         <v-row>
          <v-col>Location: {{post.location}} </v-col>
         </v-row>
         <v-row>
          <v-col>{{post.hashTags}} </v-col>
         </v-row>
       </v-container>
    </v-card-text>
  </v-card>
</template>

<script>
import * as comm from '../../configuration/communication.js'
import ShowPostModal from '../../modals/showPostModal.vue'
import axios from 'axios'
export default {
  components: { ShowPostModal },
  name: 'Post',
  props: ['post','usage', 'myReaction'],
  data() {
    return {
      width: 300,
      height: 200,
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
      if (this.usage == 'MultipleView') {
        this.width = 300;
        this.height = 400;
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
          this.reaction = null;
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
          this.reaction = reactionType;
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
        alert('Successfully added comment!');
        this.comment = '';
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
