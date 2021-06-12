<template>
  <v-form
    ref="form"
    v-model="valid"
    lazy-validation
  >
    <v-container >
        <v-row align="center" justify="center">
            <v-col cols="12" sm="4">
            <v-text-field
                v-model="email"
                :rules="[ rules.email , rules.required] "
                label="Mail:"
                required
                ></v-text-field>
            </v-col>
        </v-row>
        <v-row align="center" justify="center">
            <v-col cols="12" sm="4">
            <v-text-field
                v-model="password"
                :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                :rules="[rules.required]"
                :type="showPassword ? 'text' : 'password'"
                label="Password"
                @click:append="showPassword = !showPassword"
                ></v-text-field>
            </v-col>
        </v-row>
        <v-row align="center" justify="center">
            <v-col cols="12" sm="4">
            <v-text-field
                v-model="passCode"
                :append-icon="showPassCode ? 'mdi-eye' : 'mdi-eye-off'"
                :rules="[rules.required]"
                :type="showPassCode ? 'text' : 'password'"
                label="Pass code"
                @click:append="showPassCode = !showPassCode"
                ></v-text-field>
            </v-col>
        </v-row>
        <v-row align="center" justify="center">
            <v-col cols="12" sm="4" >
                <v-btn
                :disabled="!valid"
                color="success"
                class="mr-4"
                @click="login"
                >
                Log in
                </v-btn>
                <v-btn
                color="warning"
                elevation="8"
                @click="requestRecovery"
                >
                Forgot password
                </v-btn>
            </v-col>
        </v-row>
    </v-container>
  </v-form>
</template>

<script>
    import axios from 'axios'
    import * as comm from '../configuration/communication.js'
    import * as validator from '../plugins/validator.js'
  export default {
    data() {return {
      showPassword: false,
      showPassCode: false,
      valid: true,
      email: '',
      password: '',
      passCode: '',
      rules: validator.rules
    }},
    mounted(){
       if (this.isAvailable()){
          this.$router.push({name: 'NotFound'})
        }
    },
    methods: {
      isAvailable(){
        return comm.isUserLogged()
      },
      login () {
        if (this.$refs.form.validate()){
            let credentials = {
                "email" : this.email,
                "password" : this.password,
                "passCode" : this.passCode
            }
            axios({
                method: "post",
                url: comm.protocol +'://' + comm.server + '/api/auth/login',
                data: JSON.stringify(credentials)
            }).then(response => {
              if(response.status==200){
                comm.setJWTToken(response.data);
                this.$router.push({name: "HomePage"})
                this.$root.$emit('loggedUser')
              }
            })
        }
      },
      requestRecovery() {
        let mail = this.email;
        if (this.rules.email(mail) !== true){
          alert('E-mail must be valid');
          return;
        }
        axios({
          method: "post",
          url: comm.protocol + '://' + comm.server + '/api/auth/request-recovery',
          data: JSON.stringify(mail)
        }).then(response => {
          if(response.status==200){
            alert(response.data);
          }
        })
      },
    },
  }
</script>