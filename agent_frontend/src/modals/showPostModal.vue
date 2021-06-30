<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog transition="dialog-bottom-transition" width="900">
        <template v-slot:activator="{ on, attrs }" >
            <span v-bind="attrs"  v-on="on" >
                <post-media :width="width" :height="height" :post="post"/>
            </span>
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-card-text>
                <v-container>
                    <v-row justify="center">
                        <v-col cols="12" sm="6">
                           <post-media :width="width" :height="height" :post="post"/>
                        </v-col>
                    <v-col cols="12" sm="6">
                        <v-row> <v-col>Location: {{post.location}} </v-col></v-row>
                        <v-row><v-col> {{post.description}} </v-col></v-row>
                        <!-- <v-row>
                            <v-col>
                                <post-reactions-modal v-bind:postID="post.id"/>
                            </v-col>
                        </v-row> -->
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
import PostMedia from '../components/Post/PostMedia.vue'
//import PostReactionsModal from './PostReactionsModal.vue'
export default {
  components: { PostMedia, /*PostReactionsModal*/ },
  name: 'ShowPostModal',
  props: ['width','height','post'],
  data(){
    return{
      searchedTaggedUsers : [],
      cursorStart: -1,
      cursorEnd: -1,
    }
  },
  methods:{
  },
}
</script>