<template>
    <v-container>
        <v-row>
                <v-col cols="12" sm="3">
                <v-menu
                    bottom
                    left
                >
                    <template v-slot:activator="{ on, attrs }">
                    <v-btn width="300"
                        @click="getMessageConnections()"
                        dark
                        icon
                        v-bind="attrs"
                        v-on="on"
                    >
                        <v-icon color="blue">mdi-cog-outline</v-icon>
                    </v-btn>
                    </template>

                    <v-text-field
                        v-model="searchUsername"
                        label="Search username.."
                    ></v-text-field>

                    <v-list>
                    <v-list-item
                        v-for="(item, i) in users"
                        :key="i"
                        v-model="item.messageApproved"
                    >

                        <v-list-item-title class="" @click="goToMessaging(item)" >{{ item.username }}</v-list-item-title>
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
            users: {},
             searchUsername: ''
        }
    },

    methods:{
      
    },


    mounted(){
        
    }
}
</script>