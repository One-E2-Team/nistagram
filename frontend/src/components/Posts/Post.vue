<template>
  <v-card
    class="mx-auto"
    max-width="400"
  >
    <v-list-item one-line>
      <v-list-item-content>
        <v-list-item-title class="text-h5 d-flex justify-space-between ">
          <label>Publisher username: {{publisher.username}} </label>
          <v-btn dark icon @click="showOptionDialog()">
            <v-icon color="blue">mdi-dots-horizontal</v-icon>
          </v-btn>
        </v-list-item-title>
      </v-list-item-content>
    </v-list-item>
    <v-carousel>        
        <v-template v-for="item in post.medias" :key="item.filePath" name="temp">
            <v-carousel-item
            reverse-transition="fade-transition"
            transition="fade-transition">
            <video autoplay loop width="600" height="500" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="item.filePath.includes('mp4')">
            Your browser does not support the video tag.
            </video>
            <img width="600" height="500" :src=" protocol + '://' + server + '/static/data/' + item.filePath" v-if="!item.filePath.includes('mp4')">

            </v-carousel-item>
        </v-template>
    </v-carousel>

    <v-card-text class="text--primary">
      <div>Post description: {{post.description}} </div>
    </v-card-text>
  </v-card>
</template>

<script>

export default {
    name: "Post",
    props: ['post','publisher'],
    methods:{
        showOptionDialog(){

        }
    }
}
</script>
