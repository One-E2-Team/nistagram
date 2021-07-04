<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog transition="dialog-bottom-transition" width="900">
        <template v-slot:activator="{ on, attrs }" >
            <v-btn v-on="on" @click="test()" v-bind="attrs">Show media</v-btn>
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
              <v-btn text @click="dialog.value = false">Close</v-btn>
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
  components:{PostMedia},
  props: ['medias'],
  name: 'ShowPostFullScreenModal',
  data(){
    return {
        width: 400,
        height: 500,
        showDialog: false,
        post: {
            medias: []
        }
    }
  },
  methods: {
    isUserLogged() {
      return comm.isUserLogged();
    },
    isMyPost() {
      return comm.getLoggedUserID() == this.post.publisherId;
    },
    test(){
      console.log(this.medias);
      this.post = {
              medias: [
                {
                  filePath : this.medias
                }
              ]
          }
      console.log(this.post);
    }
  },
}
</script>