<template>
    <v-menu bottom right >
        <template v-slot:activator="{ on, attrs }">
            <v-btn dark icon v-bind="attrs" v-on="on" class="mx-2" fab small color="cyan">
                <v-icon>mdi-menu-down</v-icon>
            </v-btn>
        </template>

        <v-list>
            <v-list-item v-if="isBlocked">
                <v-list-item-title @click="toggle('block')">Unblock</v-list-item-title>
            </v-list-item>
            <v-list-item v-else>
                <v-list-item-title @click="toggle('block')">Block</v-list-item-title>
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
                <v-list-item v-if="connection.notifyMessage">
                    <v-list-item-title @click="toggle('notify/message')">Don't notify on message</v-list-item-title>
                </v-list-item>
                <v-list-item v-else>
                    <v-list-item-title @click="toggle('notify/message')">Notify on message</v-list-item-title>
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
    props: ['profileId'],
    name: 'ProfileOptions',
    data(){
        return{
            isBlocked: true,
            connection: {},
        }
    },

    created(){
        this.checkConnection()
    },
    methods:{
        checkConnection(){
            axios({
                method: "get",
                url: comm.protocol + "://" + comm.server +"/api/connection/following/my-properties/" + this.profileId,
                headers: comm.getHeader(),
            }).then((response) => {
                if(response.status == 200)
                    this.connection = response.data
            })
            .catch((error) => {
                console.log(error);
            });
        },
        toggle(point){
            axios({
                method: "put",
                url: comm.protocol + "://" + comm.server +"/api/connection/"+ point + '/' + this.profileId,
                headers: comm.getHeader(),
            }).then((response) => {
                if(response.status == 200)
                    this.connection = response.data
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
                if(response.status == 200)
                    this.isBlocked = !this.isBlocked

            })
            .catch((error) => {
                console.log(error);
            });
        },
        checkIfBlocked(){
            axios({
                    method: "get",
                    url: comm.protocol + "://" + comm.server +"/api/connection/block/" + this.post.publisherId,
                    headers: comm.getHeader(),
                }).then((response) => {
                console.log(response.data);
                if(response.status == 200)
                this.isBlocked = response.data == 'true'
                }).catch((error) => {
                    console.log(error);
                });
        }
    }
}
</script>
