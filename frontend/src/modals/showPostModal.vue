<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog transition="dialog-bottom-transition" width="900">
        <template v-slot:activator="{ on, attrs }" v-if="post.postType==2">
            <span v-bind="attrs"  v-on="on" @click="loadComments()">
                <post-media :width="width" :height="height" :post="post" :campaignData="campaignData"/>
            </span>
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-list-item >
                <v-list-item-content >
                    <v-list-item-title  class="text-h5 d-flex justify-space-between ">
                        <router-link :to="{ name: 'Profile', params: { username: post.publisherUsername }}">{{post.publisherUsername}}</router-link>
                    </v-list-item-title>
                </v-list-item-content>
            </v-list-item>
            <v-card-text>
                <v-container>
                    <v-row justify="center">
                        <v-col cols="12" sm="6">
                           <post-media :width="width" :height="height" :post="post" :campaignData="campaignData"/>
                        </v-col>
                    <v-col cols="12" sm="6">
                        <v-row> <v-col>Location: {{post.location}} </v-col></v-row>
                        <v-row><v-col> {{post.description}} </v-col></v-row>
                        <v-row>
                            <v-col class="d-flex justify-space-around ">
                                <v-btn-toggle v-if="isUserLogged" v-model="newReaction" color="primary" group dense>
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
                        <v-row>
                            <v-col>
                                <post-reactions-modal v-bind:postID="post.id"/>
                            </v-col>
                        </v-row>
                        <v-simple-table fixed-header height="200px" v-if="comments.length>0">
                            <template v-slot:default>
                              <tbody>
                                <tr v-for="(c, index) in comments" :key="index">
                                  <td>{{ c.username }}</td>
                                  <td>{{ c.content }}</td>
                                </tr>
                              </tbody>
                            </template>
                        </v-simple-table>
                        <v-row cols="12" md="6">
                            <v-col>
                                <v-textarea solo placeholder="Enter comment..." rows="4" v-model="comment" @keyup="e => findTag(e)"></v-textarea>
                                <v-list rounded>
                                  <v-list-item-group
                                    color="primary"
                                  >
                                    <v-list-item @click="setTag(item)"
                                      v-for="(item, i) in searchedTaggedUsers"
                                      :key="i"
                                    >
                                      <v-list-item-icon>
                                        <v-icon>mdi-account</v-icon>
                                      </v-list-item-icon>
                                      <v-list-item-content>
                                        <v-list-item-title v-text="item"></v-list-item-title>
                                      </v-list-item-content>
                                    </v-list-item>
                                  </v-list-item-group>
                                </v-list>
                                <v-btn color="normal" elevation="2" @click="commentPost()">
                                        Comment
                                </v-btn>
                            </v-col>
                        </v-row>
                    </v-col>
                </v-row>
            </v-container>
            </v-card-text>
            <v-card-actions class="justify-end">
              <v-btn
                text
                @click="dialog.value = false"
              >Close</v-btn>
            </v-card-actions>
          </v-card>
        </template>
      </v-dialog>
    </v-col>
  </v-row>
</template>

<script>
import PostMedia from '../components/Posts/PostMedia.vue'
import PostReactionsModal from './PostReactionsModal.vue'
import * as comm from '../configuration/communication.js'
import axios from 'axios'
export default {
  components: { PostMedia, PostReactionsModal },
  name: 'ShowPostModal',
  props: ['width','height','post','reaction', 'campaignData'],
  data(){
      return{
          isUserLogged: comm.isUserLogged(),
          comment: '',
          newReaction: this.reaction,
          searchedTaggedUsers : [],
          comments: [],
          cursorStart: -1,
          cursorEnd: -1,
      }
  },
  methods:{
      react(newReaction){
          this.$emit('reactionChanged',newReaction)
          this.newReaction = newReaction
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
      commentPost() {
      if (this.preventActionIfUnauthorized()) {
        return;
      }
      let campaignId = this.campaignData == undefined ? 0 : this.campaignData.campaignId;
      let influencerID = this.campaignData == undefined ? 0 : this.campaignData.influencerId;
      let influencerUsername = this.campaignData == undefined ? '' : this.campaignData.influencerUsername;
      let dto = {'postId' : this.post.id, 'content' : this.comment, 'campaignId': campaignId, 'influencerID': influencerID, 'influencerUsername': influencerUsername};
      axios({
        method: 'post',
        url: comm.protocol + '://' + comm.server + '/api/postreaction/comment',
        data: JSON.stringify(dto),
        headers: comm.getHeader(),
      }).then(response => {
        console.log(response.data);
        alert('Successfully added comment!');
        this.comment = '';
        this.loadComments();
      });
    },
    findTag(e){
        let end = e.target.selectionStart -1;
        if (this.comment[end] == '@')
          return;
        for(let i = end; i >= 0; i--){
          if(this.comment[i] == ' ')
            break;
          if(this.comment[i] == '@'){
              this.cursorStart = i + 1;
              this.cursorEnd = end;
              this.searchUsername(this.comment.slice(i+1, end + 1));
              return;
          }
        }
        this.searchedTaggedUsers = [];
      },
      searchUsername(username){
        axios({
          method: "get",
          url: comm.protocol + '://' + comm.server + '/api/profile/search-for-tag/' + username,
          headers: comm.getHeader()
        }).then(response => {
          if(response.status==200){
            this.searchedTaggedUsers = response.data.collection;
          }
        })
      },
      setTag(item){
        this.comment = this.comment.slice(0, this.cursorStart) +
          item + this.comment.slice(this.cursorEnd + 1, this.comment.length);
        this.searchedTaggedUsers = [];
      },
      loadComments() {
        axios({
          method: "get",
          url: comm.protocol + '://' + comm.server + '/api/postreaction/all-comments/' + this.post.id,
          headers: comm.getHeader()
        }).then(response => {
          if(response.status==200){
            this.comments = response.data.collection;
          }
        })
      }
  },

  watch:{
    reaction: function(){
      this.newReaction = this.reaction
    }
  }
}
</script>