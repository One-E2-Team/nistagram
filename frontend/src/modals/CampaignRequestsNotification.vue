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
              <div v-for="n in notifications" :key="n.campaign_id" class="d-flex">
                <h3>{{n.post.publisherUsername}}</h3>
                <show-campaign-request-modal :notification="n"/>
                <v-btn text @click="approve(n.request_id,true)" color="success">Approve</v-btn>
                <v-btn text @click="approve(n.request_id,false)" color="error">Decline</v-btn>
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
import ShowCampaignRequestModal from './ShowCampaignRequestModal.vue'
export default {
  components: { ShowCampaignRequestModal },
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
                url: comm.protocol +'://' + comm.server + '/api/campaign/request/my',
                headers: comm.getHeader(),
            }).then(response => {
              if(response.status==200){
                  this.notifications = response.data.collection;
              }
            })
        },
        approve(requestId,isAccepted){
            let data = {
                accepted: isAccepted
            }
            axios({
                method: "put",
                url: comm.protocol +'://' + comm.server + '/api/campaign/request/' + requestId,
                headers: comm.getHeader(),
                data: JSON.stringify(data)
            }).then(response => {
                if(response.status==200){
                    document.getElementById('close'+this.notification.campaign_id).close();
                    if(isAccepted){
                        alert("notifiaction is accepted")
                    }else{
                        alert("notification is declined")
                    }
                }
            })
            }
    }
}
</script>

<style>

</style>