<template>
  <v-form
    ref="form"
    v-model="valid"
    lazy-validation>
    <v-container fluid>
      <v-row align="center"
      justify="center">
        <v-col
          cols="12"
          sm="4"
        >
          <v-text-field
            v-model ="password1"
            :append-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
            :rules="[rules.password]"
            :type="show ? 'text' : 'password'"
            label="Enter password"
            hint="At least 8 characters"
            class="input-group--focused"
            @click:append="show = !show"
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row align="center"
      justify="center">
        <v-col
          cols="12"
          sm="4"
        >
          <v-text-field
            v-model="password2"
            :rules="[rules.required, rules.passwordMatch]"
            :type="'password'"
            label="Repeat password "
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row align="center"
      justify="center">
        <v-col
          cols="12"
          sm="4"
        >
        <v-btn
                :disabled="!valid"
                color="success"
                class="mr-4"
                @click="resetPassword"
                >
                Confirm
                </v-btn>
            </v-col>
        </v-row>
    </v-container>
  </v-form>
</template>

<script>

import * as validator from '../plugins/validator.js'
import * as comm from '../configuration/communication.js'
import axios from 'axios'
  export default {
    data () {
      return {
        valid: true,
        show: false,
        password1: '',
        password2: '',
        rules: {
          required: validator.rules.required,
          min: validator.rules.min,
          passwordMatch: () => (this.password1 === this.password2) || 'Password must match'
        },
      }
    },
    methods:{
      resetPassword(){
        if (this.$refs.form.validate()){
          let id = comm.getUrlVars()['id'];
          let str = comm.getUrlVars()['str'];
          if (!id || !str){
            alert('Bad url!');
            return;
          }
          let data = {
            id: id,
            uuid: str,
            password: this.password1,
          }
          axios({
            method: "post",
            url: comm.protocol + '://' + comm.server + '/api/auth/recover',
            data: JSON.stringify(data)
          }).then(response => {
            if(response.status==200){
              window.location.href = '#/log-in';
            }
          })
        }
      }
    }
  }
</script>