<template>
  <v-card>
    <v-toolbar flat>
      <template v-slot:extension>
        <v-tabs v-model="tabs" fixed-tabs>
          <v-tabs-slider></v-tabs-slider>
           <v-tab v-for="(type, index) in postTypes" :key="index" :href="'#' + type" class="primary--text">
            <v-icon>{{tabNames[index]}}</v-icon>
          </v-tab>
        </v-tabs>
      </template>
    </v-toolbar>

    <v-tabs-items v-model="tabs">
      <v-tab-item v-for="(type, index) in postTypes" :key="index" :value="type" >
        <v-card flat>
          <v-card-text>
            <v-row>
                <v-col cols="12" sm="4" v-for="p in filteredPosts(type)" :key="p.post.id">
                    <post v-bind:post="p.post" :campaign="p.campaign"/>
                </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-tab-item>
    </v-tabs-items>
  </v-card>
</template>

<script>
  import Post from '../components/Post/Post.vue'
  export default {
    components:{Post},
    props: ['posts', ],
    data () {
      return {
        postTypes: ['2','1'],
        tabNames: ['Posts','Stories'],
        tabs: 2,
        myPosts: [],
      }
    },
    /*created(){
        if(!this.isPageAvailable()) {
          this.$router.push({name: 'NotFound'})
          return;
        }
    },*/

    methods: {
      filteredPosts(type) {
        return this.posts.filter(function (item) {
          return item.post.postType == type;
        });
      },
    },

    // watch: {
    //   posts: function() {
    //     this.myPosts = this.posts;
    //   }
    // },

    /*methods:{
        isPageAvailable(){
            return comm.isUserLogged()
        }
    }*/
  }
</script>