<template>
    <v-menu bottom left >
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
export default {
    name : "Settings",
    props: ['profileId'],
    name: 'ProfileOptions',
    methods:{
        mute(){            
        },
        block(){
        }
    }
}
</script>
