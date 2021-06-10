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
                :append-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
                :rules="[rules.required, rules.min]"
                :type="show ? 'text' : 'password'"
                label="Password"
                hint="At least 8 characters"
                counter
                @click:append="show = !show"
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
    data: () => ({
      show: false,
      valid: true,    
      email: '',
      password: '',

      rules: validator.rules
    }),
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
                "password" : this.password
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
                //axios.defaults.headers.common['Authorization'] = 'Bearer ' + comm.getJWTToken().token;
              }
            }) //TODO: redirect
        }
      },
      requestRecovery() {
        let mail = this.email;
        if(mail === ''){
          alert('You must enter email!');
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