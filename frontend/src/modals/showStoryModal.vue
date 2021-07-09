<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <post-modal :visible="showDialog" @close="showDialog=false" v-bind:post="post"/>
      <v-dialog transition="dialog-bottom-transition" width="900">
        <template v-slot:activator="{ on, attrs }" v-if="post.postType==1">
            <v-avatar size="48" v-bind="attrs"  v-on="on" color="blue">
              <img src="../assets/profilepicture.jpg" />
            </v-avatar> 
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-list-item >
                <v-list-item-content >
                    <v-list-item-title  class="text-h6 d-flex justify-space-between">
                        <router-link v-if="campaignData == undefined || campaignData.influencerUsername == ''" :to="{ name: 'Profile', params: { username: post.publisherUsername }}">{{post.publisherUsername}}</router-link>
                        <router-link v-else-if="campaignData.influencerUsername != ''" :to="{ name: 'Profile', params: { username: campaignData.influencerUsername }}">{{campaignData.influencerUsername}}</router-link>
                        <v-btn dark icon @click="showDialog = true" v-if="isUserLogged() && !isMyPost()">
                          <v-icon color="blue">mdi-dots-horizontal</v-icon>
                        </v-btn>
                    </v-list-item-title>
                </v-list-item-content>
            </v-list-item>
            <v-card-text>
                <v-container>
                    <v-row justify="center">
                        <v-col cols="12" sm="8">
                             <post-media :width="width" :height="height" :post="post" :campaignData="campaignData"/>
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
import * as comm from '../configuration/communication.js'
import PostMedia from '../components/Posts/PostMedia.vue'
import PostModal from './PostModal.vue'
export default {
  components:{PostMedia, PostModal},
  props: ['visible', 'post', 'campaignData'],
  name: 'ShowPostFullScreenModal',
  data(){
    return {
        protocol: comm.protocol,
        server: comm.server,
        width: 400,
        height: 500,
        showDialog: false,
    }
  },
  methods: {
    isUserLogged() {
      return comm.isUserLogged();
    },
    isMyPost() {
      return comm.getLoggedUserID() == this.post.publisherId;
    }
  }
}
</script>