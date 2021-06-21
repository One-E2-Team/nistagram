<template>
  <v-card>
    <v-toolbar flat>
      <template v-slot:extension>
        <v-tabs v-model="tabs" fixed-tabs>
          <v-tabs-slider></v-tabs-slider>
           <v-tab v-for="(reaction, index) in reactions" :key="index" :href="'#' + reaction" class="primary--text" @click="showPosts(reaction)">
            <v-icon>{{reactionIcons[index]}}</v-icon>
          </v-tab>
        </v-tabs>
      </template>
    </v-toolbar>

    <v-tabs-items v-model="tabs">
      <v-tab-item v-for="(reaction, index) in reactions" :key="index" :value="reaction" >
        <v-card flat>
          <v-card-text>
            <v-row v-for="(username, i) in getSelectedReactions(reaction)" :key="i">
                <v-col>
                    <router-link :to="{name: 'Profile', params: {username: username}}">{{username}}</router-link> 
                </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-tab-item>
    </v-tabs-items>
  </v-card>
</template>

<script>
  export default {
    name: 'PostReactions',
    components:{},
    props: ['reactionValues'],
    data () {
      return {
            reactions: ['like','dislike'],
            reactionIcons: ['mdi-thumb-up','mdi-thumb-down'],
            tabs: null,
        }
    },
    created(){
    },
    methods: {
        getSelectedReactions(reactionType) {
            if (reactionType == 'like') {
                return this.reactionValues.likes;
            } else if (reactionType == 'dislike') {
                return this.reactionValues.dislikes;
            }
            return [];
        }
    },
  }
</script>
