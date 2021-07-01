<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog transition="dialog-top-transition" max-width="600">
        <template v-slot:activator="{ on, attrs }">
            <v-btn>
                <v-icon v-bind="attrs" v-on="on" large class="mx-2">mdi-shield-key-outline </v-icon>
            </v-btn>
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-toolbar color="primary" dark>Nistagram API Token</v-toolbar>
            <v-card-text>
              <v-form ref="form" v-model="valid" lazy-validation class="text-center">
                  <v-text-field
                    background-color="grey lighten-2"
                    :rules="[rules.required]"
                    color="cyan"
                    label="API Token"
                    v-model="token"/>
              </v-form>
            </v-card-text>
            <v-card-actions class="justify-end">
                <v-btn text @click="confirm()">Confirm</v-btn>
                <v-btn text @click="dialog.value = false">Close</v-btn>
            </v-card-actions>
          </v-card>
        </template>
      </v-dialog>
    </v-col>
  </v-row>
</template>

<script>
import * as validator from '../plugins/validator.js'
import axios from 'axios'
import * as comm from '../configuration/communication.js'
export default {
    data() {return {
        token: '',
        rules:validator.rules,
        valid: true,
    }},
    methods: {
       confirm(){
            if(this.$refs.form.validate()){
                axios({
                    method: 'post',
                    url: comm.protocol + '://' + comm.server + '/api-token',
                    headers: comm.getHeader(),
                    data: JSON.stringify(this.token),
                }).then((response) => {
                    if(response.status == 200){
                        alert('API token successfully added!');
                        this.token = '';
                    }
                });
            }
       }
    }
}
</script>