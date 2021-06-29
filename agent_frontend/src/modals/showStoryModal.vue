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
            <v-card-text>
                <v-container>
                    <v-row justify="center">
                        <v-col cols="12" sm="8">
                          <post-media :width="width" :height="height" :post="post"/>
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
import PostMedia from '../components/Posts/PostMedia.vue'
export default {
  components:{PostMedia, PostModal},
  props: ['visible', 'post'],
  name: 'ShowPostFullScreenModal',
  data(){
    return {
        protocol: comm.protocol,
        server: comm.server,
        width: 400,
        height: 500,
    }
  },
  methods: {
    isUserLogged() {
      return comm.isUserLogged();
    },
    isMyPost() {
      return comm.getLoggedUserID() == this.post.publisherId;
    }
  }
}
</script>