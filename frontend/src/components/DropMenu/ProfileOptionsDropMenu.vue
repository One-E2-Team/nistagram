<template>
    <v-menu bottom right >
        <template v-slot:activator="{ on, attrs }">
            <v-btn dark icon v-bind="attrs" v-on="on" class="mx-2" fab small color="cyan">
                <v-icon>mdi-menu-down</v-icon>
            </v-btn>
        </template>

        <v-list>
            <v-list-item>
                <v-list-item-title @click="mute()">Mute</v-list-item-title>
            </v-list-item>
            <v-list-item>
                <v-list-item-title @click="block()">Block</v-list-item-title>
            </v-list-item>
        </v-list>
    </v-menu>

</template>

<script>
import axios from 'axios'
import * as comm from '../../configuration/communication.js'
export default {
    props: ['profileId'],
    name: 'ProfileOptions',
    methods:{
        mute(){   
            axios({
                method: "put",
                url: comm.protocol + "://" + comm.server +"/api/connection/mute/" + this.profileId,
                headers: comm.getHeader(),
            }).then((response) => {
            console.log(response.data);
            })
            .catch((error) => {
            console.log(error);
            });
    },
        block(){
            axios({
                method: "put",
                url: comm.protocol + "://" + comm.server +"/api/connection/block/" + this.profileId,
                headers: comm.getHeader(),
            }).then((response) => {
            console.log(response.data);
            })
            .catch((error) => {
            console.log(error);
            });
        }
    }
}
</script>
