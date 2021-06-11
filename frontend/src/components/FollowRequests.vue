<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog
        transition="dialog-bottom-transition"
        max-width="600"
      >
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            color="warning"
            v-bind="attrs"
            v-on="on"
            @click="getRequests()"
          >Following Requests</v-btn>
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-toolbar
              color="primary"
              dark
            >All requests</v-toolbar>
            <v-card-text>
              <div v-for="req in requests" :key="req.profileID">
                {{req.username}}
                <v-btn
                  text
                  @click="approve(req.profileID)">Approve
                </v-btn>
                <v-btn
                  text
                  @click="decline(req.profileID)">Decline
                </v-btn>
              </div>
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
import axios from 'axios'
import * as comm from '../configuration/communication.js'
export default {
    name: "FollowRequests",
    data() {return {
        requests: [],
    }},
    created() {
    },
    methods: {
        getRequests(){
            axios({
                method: "get",
                url: comm.protocol +'://' + comm.server + '/api/connection/following/request',
                headers: comm.getHeader(),
            }).then(response => {
              if(response.status==200){
                  this.requests = response.data.collection;
              }
            })
        },
        approve(id){
          axios({
                method: "post",
                url: comm.protocol + '://' + comm.server + '/api/connection/following/approve/' + id,
                headers: comm.getHeader(),
            }).then(response => {
              if(response.status==200){
                 alert('Success');
              }
            })
        },
        decline(id){
          axios({
                method: "delete",
                url: comm.protocol + '://' + comm.server + '/api/connection/following/request/' + id,
                headers: comm.getHeader(),
            }).then(response => {
              if(response.status==200){
                  alert('Success');
              }
            })
        }
    }
}
</script>

<style>

</style>