<template>
    <v-menu bottom right >
        <template v-slot:activator="{ on, attrs }">
            <v-btn dark icon v-bind="attrs" v-on="on" class="mx-2" fab small color="cyan">
                <v-icon>mdi-menu-down</v-icon>
            </v-btn>
        </template>

        <v-list>
            <v-list-item v-if="isBlocked">
                <v-list-item-title @click="toggleBlock()">Unblock</v-list-item-title>
            </v-list-item>
            <v-list-item v-else>
                <v-list-item-title @click="toggleBlock()">Block</v-list-item-title>
            </v-list-item>
            <template v-if="messageConnection != null">
                <v-list-item v-if="messageConnection.notifyMessage">
                    <v-list-item-title @click="toggle('notify/message')">Don't notify on message</v-list-item-title>
                </v-list-item>
                <v-list-item v-else>
                    <v-list-item-title @click="toggle('notify/message')">Notify on message</v-list-item-title>
                </v-list-item>
            </template>
            <v-list-item v-else>
                    <v-list-item-title @click="sendMessageRequest()">Send message request</v-list-item-title>
            </v-list-item>
            <template v-if="connection != null">
                <v-list-item v-if="connection.muted">
                    <v-list-item-title @click="toggle('mute')">Unmute</v-list-item-title>
                </v-list-item>
                <v-list-item v-else>
                    <v-list-item-title @click="toggle('mute')">Mute</v-list-item-title>
                </v-list-item>
                <v-list-item v-if="connection.notifyPost">
                    <v-list-item-title @click="toggle('notify/post')">Don't notify on post</v-list-item-title>
                </v-list-item>
                <v-list-item v-else>
                    <v-list-item-title @click="toggle('notify/post')">Notify on post</v-list-item-title>
                </v-list-item>
                <v-list-item v-if="connection.notifyStory">
                    <v-list-item-title @click="toggle('notify/story')">Don't notify on story</v-list-item-title>
                </v-list-item>
                <v-list-item v-else>
                    <v-list-item-title @click="toggle('notify/story')">Notify on story</v-list-item-title>
                </v-list-item>
                <v-list-item v-if="connection.notifyComment">
                    <v-list-item-title @click="toggle('notify/comment')">Don't notify on comment</v-list-item-title>
                </v-list-item>
                <v-list-item v-else>
                    <v-list-item-title @click="toggle('notify/comment')">Notify on comment</v-list-item-title>
                </v-list-item>
                <v-list-item v-if="connection.closeFriend">
                    <v-list-item-title @click="toggle('closeFriend')">Make regular friend</v-list-item-title>
                </v-list-item>
                <v-list-item v-else>
                    <v-list-item-title @click="toggle('closeFriend')">Make close friend</v-list-item-title>
                </v-list-item>
            </template>
        </v-list>
    </v-menu>

</template>

<script>
import axios from 'axios'
import * as comm from '../../configuration/communication.js'
export default {
    props: ['profileId', 'conn','blocked','msgConn'],
    name: 'ProfileOptions',
    data(){
        return{
            isBlocked: true,
            connection: null,
            messageConnection: null,
        }
    },
    created(){
        this.isBlocked = this.blocked
    },
    methods:{
        toggle(point){
            axios({
                method: "put",
                url: comm.protocol + "://" + comm.server +"/api/connection/"+ point + '/' + this.profileId,
                headers: comm.getHeader(),
            }).then((response) => {
                if(response.status == 200){
                    if (point == 'notify/message')
                        this.$emit('messageRequestSended',response.data)
                    else
                        this.$emit('connectionChanged', response.data)
                }
            })
            .catch((error) => {
                console.log(error);
            });
        },
        toggleBlock(){
            axios({
                method: "put",
                url: comm.protocol + "://" + comm.server +"/api/connection/block/" + this.profileId,
                headers: comm.getHeader(),
            }).then((response) => {
                if(response.status == 200) {
                    this.isBlocked = !this.isBlocked
                    this.$emit('blockChanged', this.isBlocked);
                    if(!this.isBlocked){
                        this.connection = null
                        this.$emit('connectionChanged', this.connection)
                        this.messageConnection = null,
                        this.$emit('messageRequestSended',this.messageConnection)
                    }
                }
            })
            .catch((error) => {
                console.log(error);
            });
        },
        sendMessageRequest(){
             axios({
                method: "post",
                url: comm.protocol + "://" + comm.server +"/api/connection/messaging/request/" + this.profileId,
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
    },
    watch:{
        blocked: function() { 
          this.isBlocked = this.blocked
        },
        conn: function() {
            this.connection = this.conn
        },
        msgConn: function(){
            this.messageConnection = this.msgConn
        }
    }
}
</script>
