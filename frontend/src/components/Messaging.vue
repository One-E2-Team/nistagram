<template>
  <v-container>
    <v-row justify="space-around">
      <v-col md="2">
        <v-card
          class="mx-auto"
          max-width="300"
          tile
        >
        <v-list dense>
           <v-list-item-group
            color="primary"
            >
            <v-list-item
                v-for="(item, i) in usersToChat"
                :key="i"
              >
             
             <v-list-item-content>
              <v-list-item-title v-if="user.messageApproved != false" style="color:red" @click="getAllMessages(item)" >{{ item.username }}</v-list-item-title>
              <v-list-item-title v-else @click="getAllMessages(item)" >{{ item.username }}</v-list-item-title>
            
            </v-list-item-content>
              
            </v-list-item>
           </v-list-item-group>
        </v-list>
        </v-card>
      </v-col>
      <v-col md="8">
        `<v-card width="1000">
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
                {{user.username}}
              </p>
            </v-card-title>
          </v-img>
          <v-card-text>

            <v-timeline
              align-top
              dense
            >
              <v-timeline-item
                v-for="m in messages"
                :key="m.id"
                :color="getColor(m.senderId)"
                small
              >
              <v-container >
                <v-row justify="left">
                  <div>
                    <div class="font-weight-normal" v-if="loggedUserId == m.senderId">
                      <strong>You: </strong> {{ m.text }}
                    </div>
                    <div class="font-weight-normal" v-else>
                      <strong>{{user.username}}: </strong> {{ m.text }}
                    </div>
                    <div class="font-weight-normal" v-if="m.postId != ''">
                      <show-post-from-message-modal :postId="m.postId"/>
                    </div>
                    <div class="font-weight-normal" v-if="m.mediaPath != '' && m.seen == false">
                      <show-media-from-message :medias="m"/>
                    </div>
                  <!--<div>@{{ message.timestamp }}</div>-->
                  </div>
                </v-row>
              </v-container>
              </v-timeline-item>
            </v-timeline>
          </v-card-text>
          
        </v-card>
        <v-card>

          <v-card-text v-if="this.user.messageApproved == true || this.newUser == true">
            <v-textarea 
              v-model="message.text"
              no-resize
              rows="1"
              name="input-7-4"
              label="Enter message here.."
            ></v-textarea>
          </v-card-text>
          <v-file-input v-if="this.user.messageApproved == true || this.newUser == true"
            v-model="message.file"
            label="Input picture.."
          ></v-file-input>
              <v-btn v-if="this.user.messageApproved == true || this.newUser == true"
            class="ma-2"
            color="secondary"
            @click="sendMessage()"
          >
            Send
          </v-btn>
          <p v-if="this.user.messageApproved == false"> You need to approve message request to answer</p>
          <v-btn class="ma-2" color="secondary" @click="deleteMessages()" v-if="this.user.messageApproved == false">Delete messages</v-btn>
        </v-card>
      </v-col>
      <v-col md="2">
         <v-text-field v-model="searchUsername" label="Search username.."  @keyup.enter.native="search()"></v-text-field>
         <template rounded v-if="searchedUsernames.length != 0">
          <v-card
            class="mx-auto"
            max-width="300"
            tile
          >
            <v-list >
              <v-subheader>Search result</v-subheader>
              <v-list-item-group
                color="primary"
              >
                <v-list-item  @click="getUserByUsername(item)"
                  v-for="(item, i) in searchedUsernames"
                  :key="i"
                >
                  <v-list-item-icon>
                     <v-icon>mdi-account</v-icon>
                  </v-list-item-icon>
                  <v-list-item-content>
                    <v-list-item-title v-text="item"></v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-card>
        </template>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import axios from 'axios'
import * as comm from '../configuration/communication.js'
import ShowPostFromMessageModal from '../modals/showPostFromMessageModal.vue'
import ShowMediaFromMessage from '../modals/showMediaFromMessage.vue'
export default {
  components: { ShowPostFromMessageModal, ShowMediaFromMessage },
    name: "Messaging",
    data() {
      return {
        messages: [],
        loggedUserId: 0,
        user : {},
        post : null,
        file : null,
        usersToChat: [],
        searchUsername: '',
        selectedUser : {},
        searchedUsernames: [],
        messagingSenderWS: null,
        protocol: comm.protocol,
        server: comm.server,
        newUser : false,
        message: {
          senderId : 0,
          receiverId : 0,
          text : '',
          mediaPath : '',
          file : null,
          postId : ''
        }
      }
    },
    mounted() {
      this.loggedUserId = comm.getLoggedUserID();
      this.getMessageConnections();
      this.startMessagingWebSocket();
      
      this.$root.$on('seen', ()=>{
        console.log('seen');
        this.getAllMessages(this.user);
      });

      this.$root.$on('privateContent', () =>{
          alert("Content is private!");
      });

      this.$root.$on('sharePost', (data) =>{
        console.log(data);
          axios({
                  method: "get",
                  url: comm.protocol + '://' + comm.server + '/api/profile/get/' + data.username,
                }).then(response => {
                  if(response.status==200){
                        console.log(response.data);
                        let m = {
                        senderId : this.loggedUserId,
                        receiverId : response.data.ID,
                        text : '',
                        mediaPath : '',
                        postId : data.postId
                      }
                      this.sendWS("SendMessage", m);
                          }}).catch(reason => {
                              console.log(reason);
                          });
                });
          axios({
                  method: "post",
                  url: comm.protocol + "://" + comm.server +"/api/connection/messaging/request/" + this.user.profileId,
                  headers: comm.getHeader(),
              }).then((response) => {
                  if(response.status == 200) {
                      this.$emit('messageRequestSended',response.data)
                  }
              })
              .catch((error) => {
                  console.log(error);
              });
    },
    methods : {
      addMessage(data){
        let d = JSON.parse(data);
        if (d.senderId == this.user.profileId)
            this.messages.push(JSON.parse(data)); 
        },
        handler(response, data) {
          switch (response) {
            case "message":
              this.addMessage(data);
              break;
          
            default:
              break;
          }
        },
      startMessagingWebSocket(){
        let ws = new WebSocket(comm.wsProtocol + '://' + comm.wsNotificationServer + '/messaging' + "?token=" + comm.getJWTToken().token)
        let reload = function(event) {
          console.log(event);
          window.location.reload()
        }
        ws.onerror = reload
        ws.onclose = reload
        let h = this.handler
        ws.onmessage = function(event) {
          let temp = JSON.parse(event.data)
          h(temp.response, temp.data)
        }
        this.messagingSenderWS = ws
      },
      sendWS(request, data){
        let req = {
          jwt: comm.getJWTToken().token, 
          request: request,
          data: JSON.stringify(data)
        }
        this.messagingSenderWS.send(JSON.stringify(req))
      },

        getAllMessages(user){
          this.newUser = false;
          this.user = user;
            axios({
            method: "get",
            url: comm.protocol + '://' + comm.server + '/api/messaging/message/' + user.profileId,
            headers: comm.getHeader(),
        }).then(response => {
            if(response.status==200) {
                console.log(response.data.collection);
                this.messages = response.data.collection;
                console.log(this.messages);
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
        },
         getMessageConnections(){
                axios({
                method: "get",
                url: comm.protocol + '://' + comm.server + '/api/messaging/connections',
                headers: comm.getHeader(),
            }).then(response => {
                if(response.status==200) {
                    console.log(response.data);
                    this.usersToChat = response.data.collection;
                }
            }).catch(reason => {
                console.log(reason);
            });
         },
         sendMessage(){
           if (this.newUser == true){
              axios({
                  method: "post",
                  url: comm.protocol + "://" + comm.server +"/api/connection/messaging/request/" + this.user.profileId,
                  headers: comm.getHeader(),
              }).then((response) => {
                  if(response.status == 200) {
                      this.$emit('messageRequestSended',response.data)
                  }
              })
              .catch((error) => {
                  console.log(error);
              });
            }
            if (this.message.file != null){
                 const data = new FormData();
                 data.append("file", this.message.file);
                 axios.defaults.headers.common['Authorization'] = 'Bearer ' + comm.getJWTToken().token;
                 axios({
                    method: "post",
                    url: comm.protocol + "://" + comm.server + "/api/messaging/file",
                    data: data,
                    config: { headers: {...data.headers}}
                  }).then(response => {
                    console.log(response.data.fileName);
                    this.message.mediaPath = response.data.fileName;
                      let data = {
                      senderId : this.loggedUserId,
                      receiverId : this.user.profileId,
                      text : this.message.text,
                      mediaPath : this.message.mediaPath
                    }
                    console.log(this.messagingSenderWS);
                    this.sendWS("SendMessage", data);
                    this.messages.push(data);
                    this.message = {};
                    delete axios.defaults.headers.common["Authorization"];
                  })
                  .catch(response => {
                    delete axios.defaults.headers.common["Authorization"];
                    console.log(response);
                  });
            }else{
                 let data = {
                  senderId : this.loggedUserId,
                  receiverId : this.user.profileId,
                  text : this.message.text,
                  mediaPath : this.message.mediaPath,
                  postId : this.message.postId
                }
                console.log(this.messagingSenderWS);
                this.sendWS("SendMessage", data);
                this.messages.push(data);
                this.message = {};
            }
         },
        search(){
           axios({
              method: "get",
              url: comm.protocol + '://' + comm.server + '/api/profile/search/' + this.searchUsername,
            }).then(response => {
              if(response.status==200){
                this.searchedUsernames = response.data.collection;
                console.log(response.data.collection);
              }
            });
        },
        getUserByUsername(username){
          axios({
                  method: "get",
                  url: comm.protocol + '://' + comm.server + '/api/profile/get/' + username,
                }).then(response => {
                  if(response.status==200){
                        let selectedUser = {"profileId" : response.data.ID, "username" : response.data.username};
                        this.user = selectedUser;
                        this.newUser = true;
                        axios({
                        method: "get",
                        url: comm.protocol + '://' + comm.server + '/api/messaging/message/' + response.data.ID,
                        headers: comm.getHeader(),
                    }).then(response => {
                        if(response.status==200) {
                            console.log(response.data.collection);
                            this.messages = response.data.collection;
                            console.log(this.messages);
                        }
                    }).catch(reason => {
                        console.log(reason);
                    });
                this.searchedUsernames = [];
              }
            });
        },
        deleteMessages(){
            for(let mess of this.messages){
               axios({
              method: "delete",
              url: comm.protocol + '://' + comm.server + '/api/messaging/message/' + mess.id,
              headers: comm.getHeader(),
            }).then(response => {
              if(response.status==200){
                this.messages = [];
              }
            });
            }
        }
    }
}
</script>