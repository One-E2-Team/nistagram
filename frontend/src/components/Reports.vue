<template>
    <v-container fluid>
      <v-row >
          <v-col cols="12" sm="4" v-for="r in reports" :key="r.reportId">
            <v-card class="ma-4" :height="height" :width="width">
            <v-carousel :width="width" :height="height">        
                <v-template v-for="item in r.medias" :key="item.filePath" name="temp">
                    <v-carousel-item
                    reverse-transition="fade-transition"
                    transition="fade-transition">
                    <video autoplay loop :width="width" :height="height" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="item.filePath.includes('mp4')">
                    Your browser does not support the video tag.
                    </video>
                    <img :width="width" :height="height" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="!item.filePath.includes('mp4')">

                    </v-carousel-item>
                </v-template>
            </v-carousel>
            <v-card-text>
                <div>Username: {{r.publisherUsername}}, Description: {{r.description}}, Reason: {{r.reason}}</div>
            </v-card-text>
            <v-row>
                <v-col>
                    <v-btn @click="deletePost(r.postId)">Delete post</v-btn>
                </v-col>
                <v-col></v-col>
                <v-col>
                    <v-btn @click="deleteProfile(r.publisherId)">Delete profile</v-btn>
                </v-col>
            </v-row>
            </v-card>
          </v-col>
      </v-row>
    </v-container>
</template>

<script>
  import axios from 'axios'
  import * as comm from '../configuration/communication.js'
  export default {

    name: 'Reports',

    mounted(){
       if( !comm.hasRole("ADMIN") )
          this.$router.push({name: 'NotFound'});
       else{
        this.getReports();
       }
    },

    data() {return {
      reports: [],
      server: comm.server,
      protocol: comm.protocol,
      width : 400,
      height : 500
    }},

    methods: {
      deletePost(postId){
          axios({
                method: "delete",
                url: this.protocol + "://" + this.server +"/api/post/" + postId,
                headers: comm.getHeader(),
            }).then((response) => {
            console.log(response.data);
            this.getReports();
            })
            .catch((error) => {
            console.log(error);
            });
      },
      deleteProfile(profileId){
          axios({
                method: "delete",
                url: this.protocol + "://" + this.server +"/api/profile/" + profileId,
                headers: comm.getHeader(),
            }).then((response) => {
            console.log(response.data);
            this.getReports();
            })
            .catch((error) => {
            console.log(error);
            });
      },
       getReports(){
           this.reports = [];
           axios({
                method: "get",
                url: this.protocol + "://" + this.server +"/api/postreaction/report",
                headers: comm.getHeader(),
            }).then((response) => {
            this.reports = response.data.collection;
            })
            .catch((error) => {
            console.log(error);
            });
      }
    }
  }
</script>