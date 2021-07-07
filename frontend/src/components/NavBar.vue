<template>
    <div id="nav">
      <v-container>
        <v-row align="center" justify="center">
          <v-col cols="12" sm="4"><v-spacer/></v-col>
          <v-col cols="12" sm="4">
            
                <router-link v-if="isUserLogged" :to="{ name: 'HomePage'}">Home page</router-link> 
                <router-link v-else :to="{ name: 'Home'}">Home</router-link> |
                <router-link :to="{ name: 'Explore'}">Explore </router-link>
                <template v-if="isUserLogged"> | <router-link :to="{ name: 'Reactions'}" >Reactions</router-link> </template>
                <template v-if="isUserLogged"> | <router-link :to="{ name: 'Messaging'}">Chat </router-link> </template>
                <template v-if="hasRole('ADMIN')">
                  <v-menu offset-y>
                    <template v-slot:activator="{ on, attrs }">
                      <v-btn icon v-bind="attrs" v-on="on">
                        <v-icon>mdi-dots-vertical</v-icon>
                      </v-btn>
                    </template>
                    <v-list>
                      <v-list-item>
                        <v-list-item-title><router-link :to="{ name: 'VerificationRequests'}">Verification Requests</router-link></v-list-item-title>
                      </v-list-item>
                      <v-list-item>
                        <v-list-item-title><router-link :to="{ name: 'Reports'}">Reports</router-link></v-list-item-title>
                      </v-list-item>
                      <v-list-item>
                        <v-list-item-title><router-link :to="{ name: 'AgentRequests'}">Agent Requests</router-link></v-list-item-title>
                      </v-list-item>
                       <v-list-item>
                        <v-list-item-title><router-link :to="{ name: 'RegisterAgent'}">Create new agent</router-link></v-list-item-title>
                      </v-list-item>
                    </v-list>
                </v-menu>
                </template>   
          </v-col>
          <v-col cols="12" sm="1" class="float-right">
              <v-spacer/>
          </v-col>
          <v-col cols="12" sm="1" class="float-right">
            <v-row>
              <v-col><campaign-requests-notification v-if="isUserLogged"/></v-col>
              <v-col><message-requests-modal v-if="isUserLogged"/></v-col>
              <v-col><follow-requests v-if="isUserLogged"/></v-col>
              <v-col><connection-recommendation-modal v-if="isUserLogged"/></v-col>
            </v-row>
          </v-col>
          <v-col cols="12" sm="1" class="float-right">
              <notification-modal v-if="isUserLogged"/>
          </v-col>
          <v-col cols="12" sm="1" class="float-right">
              <settings v-if="isUserLogged"/>
          </v-col>
        </v-row>
      </v-container>
    </div>
</template>

<script>
import * as comm from '../configuration/communication.js'
import Settings from '../components/Settings.vue'
import FollowRequests from '../components/FollowRequests.vue'
import MessageRequestsModal from '../modals/MessageRequestsModal.vue'
import ConnectionRecommendationModal from '../modals/ConnectionRecommendationModal.vue'
import CampaignRequestsNotification from '../modals/CampaignRequestsNotification.vue'
import NotificationModal from '../modals/NotificationModal.vue'
export default {
    name: "NavBar",
    components: {
      FollowRequests,
      MessageRequestsModal,
      ConnectionRecommendationModal,
      Settings,
      CampaignRequestsNotification,
      NotificationModal
    },
    data(){
      return {
        isUserLogged: comm.getLoggedUserUsername() != null,
      }
    },
    mounted(){
      this.$root.$on('loggedUser', () => {
        this.isUserLogged = comm.getLoggedUserUsername() != null;
      })
      this.startMessagingWebSocket();
    },
    methods: {
      hasRole(role){
        return comm.hasRole(role);
      },
      addNotification(data){
        console.log(data);
        this.$root.$emit('newNotification');
      },
       handler(response, data) {
          switch (response) {
            case "notify":
              this.addNotification(data);
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
      }
    }
}
</script>