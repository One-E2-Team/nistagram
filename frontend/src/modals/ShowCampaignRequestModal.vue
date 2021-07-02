<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog transition="dialog-bottom-transition" width="900">
        <template v-slot:activator="{ on, attrs }">
            <span v-bind="attrs"  v-on="on" >
                <v-btn text>Show</v-btn>
            </span>
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-list-item >
                <v-list-item-content >
                    <v-list-item-title  class="text-h5 d-flex justify-space-between">
                        <router-link :to="{ name: 'Profile', params: { username: notification.post.publisherUsername }}">{{notification.post.publisherUsername}}</router-link>
                    </v-list-item-title>
                </v-list-item-content>
            </v-list-item>
            <v-card-text>
                <v-container>
                    <v-row justify="center">
                        <v-col cols="12" sm="6">
                            <post-media :width="width" :height="height" :post="notification.post"/>
                        </v-col>
                    <v-col cols="12" sm="6">
                        <v-row> <v-col>Location: {{notification.post.location}} </v-col></v-row>
                        <v-row> <v-col>Hash tags: {{notification.post.hashTags}} </v-col></v-row>
                        <v-row><v-col>Description: {{notification.post.description}} </v-col></v-row>
                        <v-row justify="center">
                            <v-col>    
                                <p>Date: {{notification.start}} - {{notification.end}}</p>
                            </v-col>
                         </v-row>
                        <v-row justify="space-around">
                            <v-col cols="12" sm="8" md="8">
                                <v-sheet elevation="17"  height="50" >
                                    <v-chip-group mandatory class="primary--text">
                                        <v-chip v-for="time in notification.timestamps" :key="time">{{ time }}</v-chip>
                                    </v-chip-group>
                                </v-sheet>
                            </v-col>
                        </v-row>
                    </v-col>
                </v-row>
            </v-container>
            </v-card-text>
            <v-card-actions class="justify-end">
              <v-btn text @click="accept(true)" color="success">Accept</v-btn>
              <v-btn text @click="accept(false)" color="error">Decline</v-btn>
              <v-btn text @click="dialog.value = false" :id="'close'+notification.campaign_id">Close</v-btn>
            </v-card-actions>
          </v-card>
        </template>
      </v-dialog>
    </v-col>
  </v-row>
</template>

<script>
import PostMedia from '../components/Posts/PostMedia.vue'
import * as comm from '../configuration/communication.js'
import axios from 'axios'
export default {
   components: { PostMedia },
  name: 'ShowCampaignRequestModal',
  props: ['notification'],
  data(){
      return{
          isUserLogged: comm.isUserLogged(),
          comment: '',
          newReaction: this.reaction,
          searchedTaggedUsers : [],
          cursorStart: -1,
          cursorEnd: -1,
          width: 300,
          height: 400,
      }
  },
  methods: {
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
    accept(isAccepted){
        let data = {
            accepted: isAccepted
        }
        axios({
            method: "put",
            url: comm.protocol +'://' + comm.server + '/api/campaign/request/' + this.notification.request_id,
            headers: comm.getHeader(),
            data: JSON.stringify(data)
        }).then(response => {
            if(response.status==200){
                alert('Notification successfully ' + (isAccepted ? 'accepted' : 'declined') + '!');
                document.getElementById('close'+this.notification.campaign_id).click();
            }
        })
    }
      
  },
}
</script>