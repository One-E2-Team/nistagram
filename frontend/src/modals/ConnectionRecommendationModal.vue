<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog transition="dialog-bottom-transition" max-width="600">
        <template v-slot:activator="{ on, attrs }">
          <v-icon v-bind="attrs"
            v-on="on"
            color="blue"
            @click="getRecommendations()">mdi-account-group </v-icon>
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-toolbar color="primary" dark>Recommended Connections</v-toolbar>
            <v-card-text>
              <div v-for="r in recommendations" :key="r.profileID">
                {{r.username}} with confidence {{r.confidence}}
                <v-btn
                  text
                  @click="followRequest(r.profileID)">Follow
                </v-btn>
                <v-btn
                  text
                  @click="dismiss(r.profileID)">Dismiss
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
    name: "ConnectionRecommendationModal",
    data() {return {
        recommendations: [],
    }},
    created() {
    },
    methods: {
        getRecommendations(){
            axios({
                method: "get",
                url: comm.protocol +'://' + comm.server + '/api/connection/recommendation',
                headers: comm.getHeader(),
            }).then(response => {
              if(response.status==200){
                  this.recommendations = response.data.collection;
              }
            })
        },
        followRequest(id){
          axios({
                method: "post",
                url: comm.protocol + '://' + comm.server + '/api/connection/following/request/' + id,
                headers: comm.getHeader(),
            }).then(response => {
                if (response.status==200) {
                    this.dismiss(id)
                }
            })
        },
        dismiss(id){
            this.recommendations = this.recommendations.filter(value => value.profileID != id)
        }
    }
}
</script>

<style>

</style>