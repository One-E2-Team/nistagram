<template>
    <div>
        <v-carousel :width="width" :height="height" v-if="post.medias.length>1" v-model="postIndex">        
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
            <p @click="redirect()">{{post.medias[0].webSite}}</p>
        </span>
        <p v-if="post.medias.length>1" @click="redirect()">{{currentWebsite}}</p>
    </div>
</template>

<script>
import * as comm from '../../configuration/communication.js'
import axios from 'axios'
export default {
    props:['post','width','height', 'campaignData'],
    name: 'PostMedia',
    data(){
        return{
            protocol: comm.protocol,
            server: comm.server,
            postIndex: 0,
        }
    },
    computed: {
        currentWebsite() {
            return this.post.medias[this.postIndex].webSite;
        },
        currentMediaId() {
            return this.post.medias[this.postIndex].id;
        },
    },
    methods: {
        redirect() {
            let campaignData = this.returnCampaignData();
            axios({
                method: 'get',
                url: comm.protocol + '://' + comm.server + '/redirect/'+campaignData.campaignId+'/'+campaignData.influencerId+'/'+this.currentMediaId,
                headers: comm.getHeader(),
            });
            window.open(this.currentWebsite, "_blank");
        },
        returnCampaignData(){
            let ret = {};
            ret.campaignId = this.campaignData == undefined ? 0 : this.campaignData.campaignId;
            ret.influencerId = this.campaignData == undefined ? 0 : this.campaignData.influencerId;
            return ret;
        },
    }
}
</script>

<style>

</style>