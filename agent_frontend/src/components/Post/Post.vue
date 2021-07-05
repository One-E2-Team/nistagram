<template>
  <v-card class="mx-auto" :width="width+50" elevation="24" outlined >    
      <show-post-modal :width="width" :height="height" :post="post"/>
    <v-card-text class="text--primary">
       <v-container>
         <v-row>
          <v-col>Location: {{post.location}} </v-col>
         </v-row>
         <v-row>
          <v-col>{{post.hashTags}} </v-col>
         </v-row>
          <v-row>
            <v-col v-if="campaign == undefined"><make-campaign-modal :postId="post.id"/></v-col>
            <template v-else>
              <v-col><v-btn @click="generatePdf(campaign.ID)">GENERATE PDF</v-btn></v-col>
              <v-col><update-campaign-parameters-modal :campaignId="campaign.ID"/></v-col>
              <v-col><v-btn @click="deleteCampaign(campaign.ID)">DELETE CAMPAIGN</v-btn></v-col>
            </template>
         </v-row>
       </v-container>
    </v-card-text>
  </v-card>
</template>

<script>
import ShowPostModal from '../../modals/showPostModal.vue'
import MakeCampaignModal from '../../modals/MakeCampaignModal.vue'
import UpdateCampaignParametersModal from '../../modals/UpdateCampaignParametersModal.vue'
import axios from 'axios'
import * as comm from '../../configuration/communication.js'
export default {
  components: { ShowPostModal, MakeCampaignModal, UpdateCampaignParametersModal},
  name: 'Post',
  props: ['post', 'campaign'],
  data() {
    return {
      width: 300,
      height: 400,
    }
  },
  methods: {
    deleteCampaign(id) {
      axios({
        method: 'delete',
        url: comm.protocol + '://' + comm.server +'/campaign/' + id,
        headers: comm.getHeader(),
      }).then((response) => {
        if(response.status == 200){
          alert("Successfully deleted campaign!");
        }
      });
    },
    generatePdf(id){
      axios({
        method: 'post',
        url: comm.protocol + '://' + comm.server +'/report/campaign/' + id,
        headers: comm.getHeader(),
      }).then((response) => {
        if(response.status == 200){
          console.log(response.data);
            axios({
            method: 'get',
            url: comm.protocol + '://' + comm.server +'/report/pdf/' + id,
            headers: comm.getHeader(),
          }).then((response) => {
            if(response.status == 200){
              console.log('success');
              console.log(response.data);
            }
          });
        }
      });
    }
  }
}
</script>
