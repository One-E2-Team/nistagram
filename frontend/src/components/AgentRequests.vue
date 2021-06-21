<template>
    <v-container fluid>
      <v-row justify-md="center">
          <template>
            <v-simple-table fixed-header height="300px">
                <template v-slot:default>
                    <thead>
                        <tr>
                            <th class="text-left">Username</th>
                            <th class="text-left">Email</th>
                            <th class="text-left">Website</th>
                            <th class="text-left"></th>
                            <th class="text-left"></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in requests" :key="item.profileId">
                            <td>{{ item.username }}</td>
                            <td>{{ item.email }}</td>
                            <td>{{ item.website }}</td>
                            <td><v-btn color="success" @click="approve(item.profileId)">Accept</v-btn></td>
                            <td><v-btn color="error" @click="decline(item.profileId)">Decline</v-btn></td>
                        </tr>
                    </tbody>
                </template>
            </v-simple-table>
            </template>
      </v-row>
    </v-container>
</template>

<script>
  import axios from 'axios'
  import * as comm from '../configuration/communication.js'
  export default {
    name: 'AgentRequests',

    mounted(){
       if( !comm.hasRole("ADMIN") ) {
          this.$router.push({name: 'NotFound'});
        }
       else{
        this.getRequests();
       }
    },

    data() {
        return {
            requests: [],
        }
    },

    methods: {
        getRequests() {
            axios({
                method: "get",
                url: comm.protocol + "://" + comm.server +"/api/profile/agent-requests",
                headers: comm.getHeader(),
            }).then((response) => {
                this.requests = response.data.collection;
            });
        },
        approve(profileID) {
            let request = {
                "profileId":'' + profileID,
                "accept":true
            }
            axios({
                method: "put",
                url: comm.protocol + "://" + comm.server +"/api/profile/agent-request",
                headers: comm.getHeader(),
                data: JSON.stringify(request),
            }).then((response) => {
                console.log(response);
                this.deleteRequestFromList(profileID);
            });
        },
        decline(profileID) {
            let request = {
                "profileId":'' + profileID,
                "accept":false
            }
            axios({
                method: "put",
                url: comm.protocol + "://" + comm.server +"/api/profile/agent-request",
                headers: comm.getHeader(),
                data: JSON.stringify(request),
            }).then((response) => {
                console.log(response);
                this.deleteRequestFromList(profileID);
            });
        },
        deleteRequestFromList(profileID) {
            for (let i in this.requests) {
                let r = this.requests[i];
                if (r.profileId == profileID) {
                    this.requests.splice(i, 1);
                    return;
                }
            }
        }
    },
  }
</script>