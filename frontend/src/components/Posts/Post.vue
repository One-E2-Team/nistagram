<template>
  <v-card
    class="mx-auto"
    :width="width"
    :height="height"
  >
    <post-modal v-if="showTitle" :visible="showDialog" @close="showDialog=false" v-bind:post="post"/>
    <v-list-item v-if="showTitle">
      <v-list-item-content >
        <v-list-item-title  class="text-h6 d-flex justify-space-between ">
          <label> {{post.publisherUsername}} </label>
          <v-btn dark icon @click="showOptionDialog()">
            <v-icon color="blue">mdi-dots-horizontal</v-icon>
          </v-btn>
        </v-list-item-title>
      </v-list-item-content>
    </v-list-item>
    <v-carousel :width="width" :height="height">        
        <v-template v-for="item in post.medias" :key="item.filePath" name="temp">
            <v-carousel-item
            reverse-transition="fade-transition"
            transition="fade-transition">
            <video autoplay loop :width="width" :height="height" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="item.filePath.includes('mp4')">
            Your browser does not support the video tag.
            </video>
            <img :width="width" :height="height" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="!item.filePath.includes('mp4')">

            </v-carousel-item>
        </v-template>
    </v-carousel>

    <v-card-text class="text--primary">
      <div> {{post.description}} </div>
    </v-card-text>
  </v-card>
</template>

<script>
import PostModal from '../../modals/PostModal.vue'
import * as comm from '../../configuration/communication.js'
export default {
    components: { PostModal },
    name: "Post",
    props: ['post','usage'],
    data() {
        return {
            showDialog : false,
            width: 300,
            height: 200,
            showTitle: false,
            protocol: comm.protocol,
            server: comm.server,
        }
    },
    mounted() {
      this.designView();
    },
    methods: {
        showOptionDialog(){
            this.showDialog = true;
        },
        designView() {
          if (this.usage == 'Profile') {
            this.width = 300;
            this.height = 400;
            this.showTitle = false;
          } else if (this.usage == 'Explore') {
            this.width = 300;
            this.height = 400;
            this.showTitle = true;
          } else if(this.usage == 'HomePage') {
            this.width = 600;
            this.height = 700;
            this.showTitle = true;
          }
        }
    },
    watch: {
      usage(){
        this.designView();
      }
    }
}
</script>
