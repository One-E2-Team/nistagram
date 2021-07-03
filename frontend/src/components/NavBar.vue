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
          <v-col cols="12" sm="2" class="float-right">
            <v-row>
              <v-col><campaign-requests-notification v-if="isUserLogged"/></v-col>
              <v-col><message-requests-modal v-if="isUserLogged"/></v-col>
              <v-col><follow-requests v-if="isUserLogged"/></v-col>
              <v-col><connection-recommendation-modal v-if="isUserLogged"/></v-col>
            </v-row>
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
export default {
    name: "NavBar",
    components: {
      FollowRequests,
      MessageRequestsModal,
      ConnectionRecommendationModal,
      Settings,
      CampaignRequestsNotification
    },
    data(){
      return {
        isUserLogged: comm.getLoggedUserUsername() != null,
        messagingSenderWS: function(request, data) {console.log("sender is not resent for request and data", request, data)}
      }
    },
    mounted(){
      this.$root.$on('loggedUser', () => {
        this.isUserLogged = comm.getLoggedUserUsername() != null;
      })
      if (this.isUserLogged) this.startMessagingWebSocket()
    },
    methods: {
      hasRole(role){
        return comm.hasRole(role);
      },
      startMessagingWebSocket(){
        let handler = function(response, data) {
          switch (response) {
            case "message":
              this.$root.$emit('message', data)
              break;
          
            default:
              break;
          }
        }
        let sender = comm.openWebSocketConn(comm.wsProtocol + '://' + comm.wsNotificationServer + '/messaging', handler)
        this.messagingSenderWS = sender
      }
    }
}
</script>