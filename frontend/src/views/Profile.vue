<template>
    <v-container>
        <v-row align="left" >
            <v-col cols="12" sm="12" >
                <personal-data style="height:200px"   />
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
import PersonalData from '../components/PersonalData.vue'
export default {
    components: {
        PersonalData
    }
}
</script>

<style>

</style>