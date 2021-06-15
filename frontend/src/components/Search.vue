<template>
  <v-form v-model="valid" lazy-validation ref="form">
    <v-container>
      <v-row align="center" justify="center">
        <v-col cols="12" sm="9" md="5" >
          <v-text-field
            v-model="searchParams"
            label="Search .."
            :rules="[rules.required]"
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
                <v-btn @click="searchType='hashtags'">
                    <v-icon>mdi-pound</v-icon>
                </v-btn>
                </v-btn-toggle>
            </v-card-text>
        </v-card>
        </v-col>

        <v-col cols="12" sm="6" md="1" >
            <v-btn  @click="search()">
                <v-icon>mdi-magnify</v-icon> Search
            </v-btn>
        </v-col>
      </v-row>
    </v-container>
  </v-form>
</template>

<script>
import axios from 'axios'
import * as comm from '../configuration/communication.js'
import * as validator from '../plugins/validator.js'
export default {
    name: "Search",
    data() {
      return {
        valid: true,
        rules: validator.rules,
        searchParams: '',
        searchType: 'accounts' // possible values (accounts|locations|hashtags)
      }
    },
    methods:{
        search(){
          if (this.$refs.form.validate()){
            if (this.searchType == 'locations'){
              this.searchLocation();
            } else if (this.searchType == 'hashtags'){
              this.searchHashTags();
            } else if (this.searchType == 'accounts'){
              this.searchAccounts();
            }
          }
        },
        searchLocation() {
          axios({
            method: "get",
            url: comm.protocol + '://' + comm.server + '/api/post/public/location/' + this.searchParams,
            }).then(response => {
              if (response.status==200) {
                let data = {
                  collection : response.data.collection,
                  type : "posts"
                }
                this.$emit('searched-result', data);
              }
          })
        },
        searchHashTags() {
          axios({
            method: "get",
            url: comm.protocol + '://' + comm.server + '/api/post/public/hashtag/' + this.searchParams,
            }).then(response => {
              if (response.status==200) {
                let data = {
                  collection : response.data.collection,
                  type : "posts"
                }
                this.$emit('searched-result', data);
              }
            });
        },
      searchAccounts(){
        axios({
          method: "get",
          url: comm.protocol + '://' + comm.server + '/api/profile/search/' + this.searchParams,
        }).then(response => {
          if(response.status==200){
            let data = {
              collection : response.data.collection,
              type : "accounts"
            }
            this.$emit('searched-result', data);
          }
        })
      },
    }
}
</script>

<style>

</style>