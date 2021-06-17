<template>
    <v-dialog
        v-model="show"
        max-width="500px"
      >
        <v-card :loading="loading">
          <v-card-title>
            <span>Report</span>            
          </v-card-title>
          <v-card-text>
              <v-form ref="form" v-model="valid" lazy-validation class="text-center">
                  <v-textarea
                  background-color="grey lighten-2"
                  :rules="[rules.required, rules.max255]"
                  color="cyan"
                  label="Reason"
                  :counter="255"
                  v-model="reason"
                ></v-textarea>
              </v-form>
          </v-card-text>
          <v-card-actions>
            <v-btn
              color="primary"
              text
              @click="report()"
            >
              Send
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
</template>

<script>
import axios from 'axios'
import * as comm from '../configuration/communication.js'
import * as validator from '../plugins/validator.js'
export default {
    props: ['visible','postId'],
    data(){
        return {
            rules: validator.rules,
            loading: false,
            valid: false,
            reason: '',
        }
    },
    methods:{
        report(){
            if(this.$refs.form.validate()){
                this.loading = true
                let data = {postId : this.postId, reason : this.reason}
                axios({
                method: "post",
                url: comm.protocol + "://" + comm.server + "/api/postreaction/report",
                data : JSON.stringify(data),
                headers: comm.getHeader(),
              }).then((response) => {
              console.log(response.data);
              this.loading = false;
              })
              .catch((error) => {
                console.log(error);
                this.loading = false;
              });
            }
        }
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