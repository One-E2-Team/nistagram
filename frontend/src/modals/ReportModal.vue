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
import * as validator from '../plugins/validator.js'
export default {
    props: ['visible','postId'],
    data(){
        return {
            rules: validator,
            loading: false
        }
    },
    methods:{
        report(){
            if(this.$refs.form.validate()){
                alert('Validna forma, unesi axios zahtev')
                this.loading = true
                //TODO: axios za slanje reporta i u then-u zaustaviti loading
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