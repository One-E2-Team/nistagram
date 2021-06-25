<template>
    <v-container fluid>
      <v-row >
          <v-col cols="12" sm="4" v-for="r in requests" :key="r.ID">
            <v-card
            class="ma-4"
            height="400"
            width="300"
            >
            <v-row
                class="fill-height"
                align="center"
                justify="center"
            >
            <img width="300" height="400" :src=" protocol + '://' + server + '/static/data/' + r.imagePath" >
            </v-row>
            <v-card-text>
                <div>{{r.name}} {{r.surname}}, Category: {{r.category.name}}</div>
            </v-card-text>
            <v-row>
                <v-col>
                    <v-btn @click="accept(r.ID)" color="success">Accept</v-btn>
                </v-col>
                <v-col></v-col>
                <v-col>
                    <v-btn @click="reject(r.ID)" color="error">Reject</v-btn>
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

    name: 'VerificationRequests',

    mounted(){
       if( !comm.hasRole("ADMIN") )
          this.$router.push({name: 'NotFound'});
       else{
        this.getRequests();
       }
    },

    data() {return {
      requests: [],
      server: comm.server,
      protocol: comm.protocol
    }},

    methods: {
      accept(id){
          let data = {"verificationId" : id, "status" : true};
          let json = JSON.stringify(data);
          axios({
                method: "put",
                url: this.protocol + "://" + this.server +"/api/profile/verification-request",
                data: json,
                headers: comm.getHeader(),
            }).then((response) => {
            console.log(response.data);
            this.getRequests();
            })
            .catch((error) => {
            console.log(error);
            });
      },
      reject(id){
          let data = {"verificationId" : id, "status" : false};
          let json = JSON.stringify(data);
          axios({
                method: "put",
                url: this.protocol + "://" + this.server +"/api/profile/verification-request",
                data: json,
                headers: comm.getHeader(),
            }).then((response) => {
            console.log(response.data);
            this.getRequests();
            })
            .catch((error) => {
            console.log(error);
            });
      },
       getRequests(){
           axios({
                method: "get",
                url: this.protocol + "://" + this.server +"/api/profile/verification-requests",
                headers: comm.getHeader(),
            }).then((response) => {
            this.requests = response.data.collection;
            })
            .catch((error) => {
            console.log(error);
            });
      }
    }
  }
</script>