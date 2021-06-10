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
                        v-for="(item, i) in options"
                        :key="i"
                    >
                        <v-list-item-title @click="relocate(item.name)">{{ item.title }}</v-list-item-title>
                    </v-list-item>
                    <v-list-item>
                        <v-divider/>
                    </v-list-item>
                    <v-list-item>
                        <v-list-item-title @click="logOut()">Log out</v-list-item-title>
                    </v-list-item>
                    </v-list>
                </v-menu>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
import * as comm from '../configuration/communication.js'
export default {
    name : "Settings",

    data(){
        return{
            options: [
                {title: 'Personal settings',
                name: 'PersonalSettings'},
                {title: 'Profile settings',
                name: 'ProfileSettings'},
                {title: 'My profile',
                name: 'Profile'}]
        }
    },

    methods:{
        relocate(componentName){
            console.log(componentName)
            if(componentName == 'Profile') {
                this.$router.push({name: componentName, params: {username: comm.getLoggedUserUsername()}})
                return; //this return must stay here because method propagate and switch to undefined route on line bellow if
            }
            this.$router.push({name:componentName})
        },
        logOut(){
            comm.logOut()
        }
    }
}
</script>
