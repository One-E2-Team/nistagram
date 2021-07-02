<template>
    <div>
        <v-carousel :width="width" :height="height" v-if="post.medias.length>1">        
            <v-carousel-item
                v-for="item in post.medias" :key="item.filePath"
                reverse-transition="fade-transition"
                transition="fade-transition">
                <video autoplay loop :width="width" :height="height" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="item.filePath.includes('mp4')">
                Your browser does not support the video tag.
                </video>
                <img :width="width" :height="height" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="!item.filePath.includes('mp4')">
            </v-carousel-item>
        </v-carousel>
        <span v-else>
            <video autoplay loop :width="width" :height="height" :src=" protocol + '://' + server + '/static/data/' + post.medias[0].filePath" v-if="post.medias[0].filePath.includes('mp4')">
                    Your browser does not support the video tag.
            </video>
            <img :width="width" :height="height"  :src=" protocol + '://' + server + '/static/data/' + post.medias[0].filePath" v-if="!post.medias[0].filePath.includes('mp4')">
        </span>
    </div>
</template>

<script>
import * as comm from '../../configuration/communication.js'
export default {
    props:['post','width','height'],
    data(){
        return{
            protocol: comm.protocol,
            server: comm.nistagram_server,
        }
    }
}
</script>

<style>

</style>