<template>
    <div>
    <h2>Here is your QR code for two factor authentification </h2> <br/>
    <img width="200" height="200" :src=" protocol + '://' + server + '/static/totp/' + qruuid + '.png'">
    </div>
</template>

<script>
import * as comm from '../configuration/communication.js'
export default {
    props: ['qruuid'],
    data(){
        return{
            protocol : comm.protocol,
            server: comm.server
        }
    },
}
</script>

<style>

</style>