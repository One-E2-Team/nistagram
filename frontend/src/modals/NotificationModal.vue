<template>
  <v-container>
    <v-row >
        <v-col cols="12" sm="3">
           <v-menu bottom left>
                <template v-slot:activator="{ on, attrs }">
                        <v-btn dark icon v-bind="attrs" v-on="on">
                        <v-icon  color="blue">mdi-email</v-icon>
                        <v-badge v-if="notifications.length > 0"
                            color="green"
                            :content='notifications.length'
                            >
                            </v-badge>
                    </v-btn>
                    
                </template>
               
                <v-list >
                    
                        <v-list-item
                            v-for="(item, i) in notifications"
                            :key="i"
                        >
                        
                        <v-list-item-content>
                        <v-list-item-title @click="goToMessaging(item.username)" >{{ item.text }}</v-list-item-title>
                        
                        </v-list-item-content>
                        
                        </v-list-item>
                </v-list>
           </v-menu>
        </v-col>
    </v-row>
  </v-container>
</template>

<script>
import axios from 'axios'
import * as comm from '../configuration/communication.js'
export default {
    name: "NotificationModal",
    data() {return {
        notifications: [],
        usernames: []
    }},
    created() {
    },
    mounted(){
        this.getNotifications();

        this.$root.$on('messageSeen', (senderId) => {
            axios({
                    method: "put",
                    url: comm.protocol +'://' + comm.server + '/api/messaging/seen-message/' + senderId,
                    headers: comm.getHeader(),
                }).then(response => {
                if(response.status==200){
                   this.getNotifications();
                }
                });
        });

        this.$root.$on('newNotification', (data) => {
            console.log('hello');
            console.log(data);
            this.getNotifications();
        })
    },
    methods: {
        getNotifications(){
            this.notifications = [];
                axios({
                    method: "get",
                    url: comm.protocol +'://' + comm.server + '/api/messaging/notification',
                    headers: comm.getHeader(),
                }).then(response => {
                if(response.status==200){
                    this.usernames = response.data.usernames;
                    for (let u of this.usernames){
                        let notif = {
                            username : u,
                            text: "You have a new message from " + u + " .",
                        } 
                        this.notifications.push(notif);
                    }
                }
                });
        },
        goToMessaging(username){
             this.$router.push({name: 'Messaging' ,  params: {username: username}});
        }
    }
}
</script>

<style>

</style>