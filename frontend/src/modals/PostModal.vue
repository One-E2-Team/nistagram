<template>
  <div class="text-center">
    <v-dialog
      v-model="show"
      max-width="250px"
    >

      <v-card>
        <v-card-text>
            <v-btn v-if="isBlocked"
              label="Report"
              class="my-2"
              style="color:red"
              @click="showReportModal = true"
              width="200"
            >Report</v-btn><br/>
            <v-btn
              label="Block"
              width="200"
              @click="block()"
            >Block</v-btn><br/>
        </v-card-text>
      </v-card>
    </v-dialog>
    <report-modal :visible="showReportModal" @close="showReportModal=false" v-bind:postId="post.id"/>
  </div>
</template>

<script>
import ReportModal from './ReportModal.vue';
import axios from 'axios'
import * as comm from '../configuration/communication.js'
export default {
  components: { ReportModal },
  props: ['visible', 'post'],
  data(){
    return {
      showReportModal: false,
      isBlocked: true,
    }
  },
  created(){
    //this.checkIfBlocked()
  },
  methods: {
      block(){
        axios({
                method: "put",
                url: comm.protocol + "://" + comm.server +"/api/connection/block/" + this.post.publisherId,
                headers: comm.getHeader(),
            }).then((response) => {
            console.log(response.data);
          })
          .catch((error) => {
            console.log(error);
          });
      },
      // checkIfBlocked(){
      //   axios({
      //           method: "get",
      //           url: comm.protocol + "://" + comm.server +"/api/connection/block/" + this.post.publisherId,
      //           headers: comm.getHeader(),
      //       }).then((response) => {
      //       console.log(response.data);
      //       if(response.status == 200)
      //         this.isBlocked = response.data == 'true'
      //       }).catch((error) => {
      //           console.log(error);
      //       });
      // }
  },
  computed: {
    show: {
      get () {
        return this.visible;
      },
      set (value) {
        if (!value) {
          this.$emit('close');
        }
      }
    }
  }
}
</script>