<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog transition="dialog-bottom-transition" width="900">
        <template v-slot:activator="{ on, attrs }">
            <span v-bind="attrs"  v-on="on" >
                <v-btn text>Share</v-btn>
            </span>
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-card-text>
                <v-container>
                    <v-row justify="center">
                        <v-col>
                             <v-text-field 
                                v-model="searchUsername"
                                label="Search .."
                            ></v-text-field>
                        </v-col>
                        <v-col>
                            <v-btn @click="search()">Search</v-btn>
                        </v-col>
                    </v-row>
                    <v-sheet class="mx-auto" elevation="1" max-width="900" >
                    <v-slide-group class="pa-4" >
                        <v-slide-item v-for="u in usernames" :key="u" >
                            <v-btn @click="share(u)">{{u}}</v-btn>
                        </v-slide-item>
                    </v-slide-group>
                </v-sheet>
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
import axios from 'axios'
export default {
  name: 'SharePostModal',
  props: ['postId'],
  data(){
      return{
          usernames : [],
          searchUsername : ''
      }
  },
  methods: {
      search(){
          axios({
          method: "get",
          url: comm.protocol + '://' + comm.server + '/api/profile/search/' + this.searchUsername,
        }).then(response => {
          if(response.status==200){
              this.usernames = response.data.collection;
          }
        })
      },

      share(username){
          let data = {
              username : username,
              postId : this.postId
          }
          this.$router.push({name: "Messaging"})
          this.$root.$emit('sharePost', data);
      }
      
  },
}
</script>