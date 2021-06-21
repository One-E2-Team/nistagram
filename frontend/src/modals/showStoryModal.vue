<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog transition="dialog-bottom-transition" width="900">
        <template v-slot:activator="{ on, attrs }" v-if="post.postType==1">
            <v-avatar size="48" v-bind="attrs"  v-on="on" color="blue">
              <img src="../assets/profilepicture.jpg" />
            </v-avatar> 
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-list-item >
                <v-list-item-content >
                    <v-list-item-title  class="text-h5 d-flex justify-space-between ">
                        <router-link :to="{ name: 'Profile', params: { username: post.publisherUsername }}">{{post.publisherUsername}}</router-link>
                    </v-list-item-title>
                </v-list-item-content>
            </v-list-item>
            <v-card-text>
                <v-container>
                    <v-row justify="center">
                        <v-col cols="12" sm="6">
                            <v-carousel v-if="post.medias.length>1">        
                                <v-carousel-item
                                    v-for="item in post.medias" :key="item.filePath"
                                    reverse-transition="fade-transition"
                                    transition="fade-transition">
                                    <video autoplay loop :width="width" :height="height" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="item.filePath.includes('mp4')">
                                        Your browser does not support the video tag.
                                    </video>
                                    
                                    <img :width="width" :height="height"  :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="!item.filePath.includes('mp4')">

                                </v-carousel-item>
                            </v-carousel>
                            <div v-else>
                                <video autoplay loop :width="width" :height="height" :src=" protocol + '://' + server + '/static/data/' + post.medias[0].filePath" v-if="post.medias[0].filePath.includes('mp4')">
                                        Your browser does not support the video tag.
                                </video>
                                <img :width="width" :height="height"  :src=" protocol + '://' + server + '/static/data/' + post.medias[0].filePath" v-if="!post.medias[0].filePath.includes('mp4')">
                            </div>
                    </v-col>
                    <v-col cols="12" sm="6" v-if="post.postType==2">
                       Siso
                    </v-col>
                </v-row>
            </v-container>
            </v-card-text>
            <v-card-actions class="justify-end">
              <v-btn
                text
                @click="dialog.value = false"
              >Close</v-btn>
            </v-card-actions>
          </v-card>
        </template>
      </v-dialog>
    </v-col>
  </v-row>
</template>

<script>
import * as comm from '../configuration/communication.js'
export default {
  props: ['visible', 'post'],
  name: 'ShowPostFullScreenModal',
  data(){
    return {
        protocol: comm.protocol,
        server: comm.server,
        width: 400,
        height: 500
    }
  },
}
</script>