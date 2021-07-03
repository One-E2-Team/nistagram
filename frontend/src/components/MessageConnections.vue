<template>
    <v-container>
        <v-row>
                <v-col cols="12" sm="3">
                <v-menu
                    bottom
                    left
                >
                    <template v-slot:activator="{ on, attrs }">
                    <v-btn
                        @click="getMessageConnections()"
                        dark
                        icon
                        v-bind="attrs"
                        v-on="on"
                    >
                        <v-icon color="blue">mdi-cog-outline</v-icon>
                    </v-btn>
                    </template>

                    <v-list>
                    <v-list-item
                        v-for="(item, i) in users"
                        :key="i"
                    >
                        <v-list-item-title>{{ item.username }}</v-list-item-title>
                    </v-list-item>
                    <v-list-item>
                        <v-divider/>
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
    name : "MessageConnections",

    data(){
        return{
            users: {}
        }
    },

    methods:{
       getMessageConnections(){
            axios({
            method: "get",
            url: comm.protocol + '://' + comm.server + '/api/messaging/connections',
            headers: comm.getHeader(),
        }).then(response => {
            if(response.status==200) {
                console.log(response.data);
                this.users = response.data.collection;
            }
        }).catch(reason => {
            console.log(reason);
        });
       },
    },


    mounted(){
        
    }
}
</script>