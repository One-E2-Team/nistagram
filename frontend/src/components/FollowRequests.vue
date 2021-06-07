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
            >Opening from the bottom</v-toolbar>
            <v-card-text>
              <div class="text-h2 pa-12">Hello world!</div>
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
    data: () => ({
        requests: [],
    }),
    created() {
    },
    methods: {
        getRequests(){
            axios({
                method: "get",
                url: 'http://' + comm.server + '/api/connection/following/request',
            }).then(response => {
              if(response.status==200){
                  this.requests = response.data.collection;
              }
            })
        }
    }
}
</script>

<style>

</style>