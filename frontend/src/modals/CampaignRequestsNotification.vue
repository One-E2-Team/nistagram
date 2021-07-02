<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog transition="dialog-bottom-transition" max-width="600">
        <template v-slot:activator="{ on, attrs }">
          <v-icon v-bind="attrs" v-on="on" color="blue" @click="getNotifications()">mdi-bell</v-icon>
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-toolbar color="primary" dark>Campaign requests</v-toolbar>
            <v-card-text>
              <div v-for="n in notifications" :key="n.campaign_id">
                <p>{{n.post.publisherUsername}}</p>
               
                <v-btn text @click="approve(n.request_id)">Approve</v-btn>
                <v-btn text @click="decline(n.request_id)">Decline</v-btn>
              </div>
            </v-card-text>
            <v-card-actions class="justify-end">
              <v-btn text @click="dialog.value = false">Close</v-btn>
            </v-card-actions>
          </v-card>
        </template>
      </v-dialog>
    </v-col>
  </v-row>
</template>

<script>
import axios from 'axios'
import * as comm from '../configuration/communication.js'
export default {
    name: "CampaignRequestNotifications",
    data() {
        return {
            notifications: [],
        }
    },
    methods: {
        getNotifications(){
            axios({
                method: "get",
                url: comm.protocol +'://' + comm.server + '/campaign/request/my',
                headers: comm.getHeader(),
            }).then(response => {
              if(response.status==200){
                  this.notifications = response.data.collection;
              }
            })
        },
    }
}
</script>

<style>

</style>