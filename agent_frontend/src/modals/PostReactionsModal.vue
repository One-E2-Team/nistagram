<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog
        transition="dialog-bottom-transition"
        max-width="600"
      >
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            color="primary"
            v-bind="attrs"
            v-on="on"
            @click="getReactions()"
          >Show all reactions</v-btn>
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-toolbar
              color="primary"
              dark
            >All reactions</v-toolbar>
            <v-card-text>
              <post-reactions :reactionValues="reactions"/>
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
import axios from 'axios'
import * as comm from '../configuration/communication.js'
import PostReactions from '../components/Posts/PostReactions.vue'
export default {
    name: "PostReactionsModal",
    props: ['postID'],
    components: {PostReactions},
    data() {
        return {
            reactions : {
                likes: [],
                dislikes: []
            },
        }
    },
    created() {
        
    },
    methods: {
        getReactions() {
            axios({
                method: 'get',
                url: comm.protocol + '://' + comm.server + '/api/postreaction/all-reactions/' + this.postID,
            }).then(response => {
                this.reactions.likes = response.data.likes;
                this.reactions.dislikes = response.data.dislikes;
            });
        }
    },
}
</script>
