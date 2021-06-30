<template>
  <v-card class="mx-auto" :width="width+50" elevation="24" outlined >    
      <show-post-modal :width="width" :height="height" :post="post"/>
      <post-media :width="width" :height="height" :post="post"/>
    <v-card-text class="text--primary">
       <v-container>
         <v-row>
          <v-col>Location: {{post.location}} </v-col>
         </v-row>
         <v-row>
          <v-col>{{post.hashTags}} </v-col>
         </v-row>
          <v-row>
            <v-col><make-campaign-modal :postId="post.id"/> </v-col>
         </v-row>
       </v-container>
    </v-card-text>
  </v-card>
</template>

<script>
import ShowPostModal from '../../modals/showPostModal.vue'
import MakeCampaignModal from '../../modals/MakeCampaignModal.vue'
export default {
  components: { ShowPostModal, MakeCampaignModal},
  name: 'Post',
  props: ['post'],
  data() {
    return {
      width: 300,
      height: 400,
    }
  },
  mounted() {
    console.log(this.post);
  },
}
</script>
