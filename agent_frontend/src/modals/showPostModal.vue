<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog transition="dialog-bottom-transition" width="900">
        <template v-slot:activator="{ on, attrs }" v-if="post.postType==2">
            <span v-bind="attrs"  v-on="on" >
                <post-media :width="width" :height="height" :post="post"/>
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
                           <post-media :width="width" :height="height" :post="post"/>
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
import PostMedia from '../components/Post/PostMedia.vue'
import PostReactionsModal from './PostReactionsModal.vue'
import * as comm from '../configuration/communication.js'
export default {
  components: { PostMedia, PostReactionsModal },
  name: 'ShowPostModal',
  props: ['width','height','post'],
  data(){
      return{
          isUserLogged: comm.isUserLogged(),
          comment: '',
          newReaction: this.reaction,
          searchedTaggedUsers : [],
          cursorStart: -1,
          cursorEnd: -1,
      }
  },
  methods:{
  },
}
</script>