<template>
  <v-container>
    <v-row justify="space-around">
      <v-card width="1000">
        <v-img
          height="120px"
          src="https://cdn.pixabay.com/photo/2020/07/12/07/47/bee-5396362_1280.jpg"
        >
          <v-card-title class="white--text mt-8">
            <v-avatar size="56">
              <img
                alt="user"
                src="https://cdn.pixabay.com/photo/2020/06/24/19/12/cabbage-5337431_1280.jpg"
              >
            </v-avatar>
            <p class="ml-3">
              Username
            </p>
          </v-card-title>
        </v-img>
        <v-card-text>

          <v-timeline
            align-top
            dense
          >
            <v-timeline-item
              v-for="message in messages"
              :key="message.id"
              :color="getColor(message.senderId)"
              small
            >
            <v-container >
              <v-row justify="left">
                <div>
                  <div class="font-weight-normal" v-if="loggedUserId == message.senderId">
                    <strong>You: </strong> <!--@{{ message.timestamp }}-->
                  </div>
                  <div class="font-weight-normal" v-else>
                    <strong>Username: </strong> {{ message.text }}
                  </div>
                <!--<div>@{{ message.timestamp }}</div>-->
                </div>
              </v-row>
            </v-container>
            </v-timeline-item>
          </v-timeline>
        </v-card-text>
         
      </v-card>
       <v-card width="1000">

        <v-card-text>
          <v-textarea
            v-model="text"
             no-resize
            rows="1"
            name="input-7-4"
            label="Enter message here.."
          ></v-textarea>
        </v-card-text>
            <v-btn
          class="ma-2"
          color="secondary"
          @click="sendMessage()"
        >
          Send
        </v-btn>
      </v-card>
    </v-row>
  </v-container>
</template>

<script>
import axios from 'axios'
import * as comm from '../configuration/communication.js'
export default {
    name: "Messaging",
    data() {
      return {
        messages: [],
        loggedUserId: 0,
        text : '',
        post : null,
        file : null
      }
    },
    mounted() {
       this.getAllMessages();
       this.loggedUserId = comm.getLoggedUserID();
    },
    methods : {
        getAllMessages(){
            axios({
            method: "get",
            url: comm.protocol + '://' + comm.server + '/api/messaging/message/' + 2,
            headers: comm.getHeader(),
        }).then(response => {
            if(response.status==200) {
                console.log(response.data);
                this.messages = response.data.collection;
            }
        }).catch(reason => {
            console.log(reason);
        });
        },
        getColor(id){
          if (id == this.loggedUserId){
            return 'deep-purple lighten-1';
          }else{
            return 'green';
          }
        }
    }
}
</script>