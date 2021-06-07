<template>
  <v-form>
    <v-container>
      <v-row align="center" justify="center">
        <v-col cols="12" sm="9" md="5" >
          <v-text-field
            v-model="searchParams"
            label="Search .."
          ></v-text-field>
        </v-col>

        <v-col cols="12" sm="6" md="2">
         <v-card flat class="py-0">
            <v-card-text>    
                <v-btn-toggle
                mandatory
                >
                <v-btn @click="searchType='accounts'">
                    <v-icon>mdi-account</v-icon>
                </v-btn>
                <v-btn @click="searchType='locations'">
                    <v-icon>mdi-map-marker</v-icon>
                </v-btn>
                <v-btn  @click="searchType='hashtags'">
                    <v-icon>mdi-pound</v-icon>
                </v-btn>
                </v-btn-toggle>
            </v-card-text>
        </v-card>
        </v-col>

        <v-col cols="12" sm="6" md="1" >
            <v-btn  @click="search()">
                <v-icon >mdi-magnify</v-icon> Search
            </v-btn>
        </v-col>
      </v-row>
    </v-container>
  </v-form>
</template>

<script>
import axios from 'axios'
import * as comm from '../configuration/communication.js'
export default {
    name: "Search",
    data(){
        return{
            searchParams : '',
            searchType: 'accounts' // possible values (accounts|locations|hashtags)
    }},
    methods:{
        search(){
            if (this.searchType == 'locations'){
                this.searchLocation()
            } else if (this.searchType == 'hashtags'){
                this.searchHashTags()
            } else if (this.searchType == 'accounts'){
                this.searchAccounts();
            }
        },
        searchLocation(){
           axios({
            method: "get",
            url: 'http://' + comm.server + '/api/post/public/location/' + this.searchParams,
            }).then(response => {
            if(response.status==200){
              let res = response.data.collection;
               res.forEach((post) => {
                if(post.medias != null){
                  post.medias.forEach((media) =>{
                   media.filePath = "http://" + comm.server +"/static/data/" + media.filePath;
                  });
                }
            });
              this.$emit('searched-result', res);
          }
        })
      },
      searchHashTags(){
        axios({
            method: "get",
            url: 'http://' + comm.server + '/api/post/public/hashtag/' + this.searchParams,
            }).then(response => {
            if(response.status==200){
              let res = response.data.collection;
               res.forEach((post) => {
                if(post.medias != null){
                  post.medias.forEach((media) =>{
                   media.filePath = "http://" + comm.server +"/static/data/" + media.filePath;
                  });
                }
            });
              this.$emit('searched-result', res);
          }
        })
      },
      searchAccounts(){
        axios({
          method: "get",
          url: 'http://' + comm.server + '/api/profile/search/' + this.searchParams,
        }).then(response => {
          if(response.status==200){
            this.$emit('searched-result', response.data.collection);
          }
        })
      },
    }
}
</script>

<style>

</style>