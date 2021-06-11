<template>
<v-container>
  <v-alert
    v-if="alert"
    :value="alertText"
    color="red"
    type="error"
    dismissible
    text
    v-model="alertText"
  >{{alertText}}</v-alert>
  <v-row align="center" justify="center">
    <v-col cols="12" sm="9" >
  <v-stepper v-model="e1">
    
      
    <v-stepper-header >
      <v-stepper-step :complete="e1 > 1" step="1">
        Credentials
      </v-stepper-step>

      <v-divider></v-divider>

      <v-stepper-step :complete="e1 > 2" step="2">
        Personal data
      </v-stepper-step>

      <v-divider></v-divider>

      <v-stepper-step step="3">
        About me
      </v-stepper-step>
    </v-stepper-header>

    <v-stepper-items class="text-center">
      <!-- Step 1  component-->
      <v-stepper-content step="1">
        <v-form ref="form1" v-model="valid1" lazy-validation>
          <v-container >
            <v-row align="center" justify="center">
              <v-col cols="12" sm="6" >
                <v-text-field
                    v-model="credentials.email"
                    :rules="[ rules.email , rules.required] "
                    label="Mail:"
                    required
                    ></v-text-field>
              </v-col>
            </v-row>
            <v-row align="center" justify="center">
              <v-col cols="12" sm="6">
                <v-text-field
                    v-model="credentials.password"
                    :append-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
                    :rules="[ rules.password ]"
                    :type="show ? 'text' : 'password'"
                    label="Password"
                    hint="At least 8 characters, 1 lower, 1 capital letter, 1 number and 1 special character"
                    counter
                    @click:append="show = !show"
                    ></v-text-field>
              </v-col>
            </v-row>
             <v-row align="center" justify="center">
              <v-col cols="12" sm="6">
                <v-text-field
                    v-model="password2"
                    :rules="[passwordMatch, rules.required]"
                    :type="'password'"
                    label="Confirm password"
                    hint="Password must match"
                    counter
                    ></v-text-field>
              </v-col>
            </v-row>
            <v-row align="center" justify="center">
              <v-col cols="12" sm="6">
                <v-btn
                :disabled="!valid1"
                color="primary"
                class="d-flex justify-space-around mb-6"
                @click="continueTo2">
                Continue
                </v-btn>
              </v-col>
            </v-row>
          </v-container>
        </v-form>
      </v-stepper-content>

      <!-- Step 2 content -->
      <v-stepper-content step="2" >
        <v-form ref="form2" v-model="valid2" lazy-validation class="text-center">
          <v-container >
            <v-row align="center" justify="center">
              <v-col cols="12" sm="8" >
                <v-text-field
                    v-model="person.username"
                    :rules="[ rules.username , rules.required] "
                    label="Username:"
                    required
                    ></v-text-field>
              </v-col>
            </v-row>
          <v-row align="center" justify="center">
            <v-col cols="12" sm="8" >
              <v-text-field
                v-model="person.name"
                :rules="[ rules.required ] "
                label="Name:"
                required
                ></v-text-field>
            </v-col>
          </v-row>
          <v-row align="center" justify="center">
            <v-col cols="12" sm="8" >  
              <v-text-field
                v-model="person.surname"
                :rules="[ rules.required , rules.max ] "
                label="Surname:"
                required
                ></v-text-field>
            </v-col>
          </v-row>
          <v-row align="center" justify="center">
            <v-col cols="12" sm="8" >  
              <v-menu
                v-model="menu"
                :close-on-content-click="false"
                :nudge-right="40"
                transition="scale-transition"
                offset-y
                min-width="auto"
              >
                <template v-slot:activator="{ on, attrs }">
                  <v-text-field
                    v-model="person.birthDate"
                    label="Birth date"
                    prepend-icon="mdi-calendar"
                    readonly
                    v-bind="attrs"
                    v-on="on"
                  ></v-text-field>
                </template>
                <v-date-picker
                  v-model="person.birthDate"
                  @input="menu = false"
                ></v-date-picker>
              </v-menu>
            </v-col>
          </v-row>
          <v-row align="center" justify="center">
            <v-col cols="12" sm="8" >          
              <v-text-field
                v-model="person.telephone"
                label="Telephone:"
                ></v-text-field>
            </v-col>
          </v-row>
          <v-row align="center" justify="center">
            <v-col cols="12" sm="8" >          
              <div class="text-left">
                <label>Gender:</label>
              </div> 
              <v-radio-group
                v-model="person.gender"
                row
              >
                <v-radio
                  label="Male"
                  value="male"
                ></v-radio>
                <v-radio
                  label="Female"
                  value="female"
                ></v-radio>
              </v-radio-group> 
            </v-col>
          </v-row>

          <v-row align="center" justify="center">
              <v-col cols="12" sm="6" class="d-flex justify-space-around mb-6">
                <v-btn
                color="primary"
                :disabled="!valid2"
                @click="continueTo3">
                Continue
                </v-btn>
              <v-btn
              color="normal"
              class="d-flex justify-space-around mb-6"
              @click="e1=1;">
              Back
              </v-btn>
            </v-col>
          </v-row>
          </v-container>

        </v-form>
      </v-stepper-content>

      <!--Step3 content -->
      <v-stepper-content step="3">
        <v-form ref="form3" v-model="valid3" lazy-validation class="text-center">
          <v-container >
            <v-row align="center" justify="center">
              <v-col cols="12" sm="8" >
                <v-textarea
                  background-color="grey lighten-2"
                  color="cyan"
                  label="Biography"
                  v-model="person.biography"

                ></v-textarea>
              </v-col>
            </v-row>

            <v-row align="center" justify="center">
              <v-col cols="12" sm="8" >  
                <v-text-field
                  v-model="person.webSite"
                  label="Web site:"
                  ></v-text-field>
              </v-col>
            </v-row>

            <v-row align="center" justify="center">
              <v-col cols="12" sm="8" >
              <v-checkbox
                  v-model="isPrivate"
                  :label="`Private account`"
                ></v-checkbox>
              </v-col>
            </v-row>

          
            <v-row align="center" justify="center">
              <v-col cols="12" sm="8" >

                <v-combobox
                  v-model="person.interests"
                  :items="interests"
                  :rules="[ rules.required , rules.oneOrMoreElement ] "
                  chips
                  clearable
                  label="Your interests"
                  multiple
                  prepend-icon="mdi-filter-variant"
                  solo
                >
                  <template v-slot:selection="{ attrs, item, select, selected }">
                    <v-chip
                      v-bind="attrs"
                      :input-value="selected"
                      close
                      @click="select"
                      @click:close="remove(item)"
                    >
                      <strong>{{ item }}</strong>&nbsp;
                    </v-chip>
                  </template>
                </v-combobox>
              </v-col>
            </v-row>
          

          <v-row align="center" justify="center">
              <v-col cols="12" sm="6" class="d-flex justify-space-around mb-6">
                <v-btn
                color="primary"
                @click="register">
                Register
                </v-btn>
              <v-btn
              color="normal"
              class="d-flex justify-space-around mb-6"
              @click="e1=2">
              Back
              </v-btn>
            </v-col>
          </v-row>
          </v-container>

        </v-form>
      </v-stepper-content>
    </v-stepper-items>
  </v-stepper>
    </v-col>
  </v-row>
</v-container>
</template>

<script>
import axios from 'axios'
import * as comm from '../configuration/communication.js'
import * as validator from '../plugins/validator.js'

  export default {
    data() {return {
      alert: false,
      alertText : '',
      e1: 1,
      show: false,
      valid1: true,
      valid2: true,
      valid3: true,
      credentials: {
        email: '',
        password: '',
      },
      person: {
        username: '',
        name: '',
        surname: '',
        telephone: '',
        gender: '',
        birthDate: '',
        biography: '',
        webSite: '',
        interests: []
      } ,  
      interests: [],
      isPrivate: false,
      password2: '',
      rules: validator.rules,
      passwordMatch: () => (this.credentials.password === this.password2) || 'Password must match',
      menu: false,
    }},

    mounted(){
       if (this.isAvailable()){
          this.$router.push({name: 'NotFound'})
        }
    },
    created(){
      axios({
          method: "get",
          url: comm.protocol + "://" + comm.server +"/api/profile/interests",
        }).then(response => {
          if (response.status == 200) {
            this.interests = response.data.collection
          }
        })
        .catch(response => {
          console.log(response);
        });
    },

    methods: {
      isAvailable(){
        return comm.isUserLogged()
      },
      continueTo2() {
        if (this.$refs.form1.validate()){
            this.e1 = 2;
        }
      },
      continueTo3() {
        if (this.$refs.form2.validate()){
            this.e1 = 3;
        }
      },
      remove(item) {
        this.person.interests.splice(this.person.interests.indexOf(item), 1);
        this.person.interests = [...this.person.interests];
      },
      correctInterests() {
        if (this.person.interests.length == 0) {
          return false;
        }
        for (let item of this.person.interests){
          if(!this.interests.includes(item)){
            return false;
          }
        }
        return true;
      },
      register() {
        if (!this.correctInterests()) {
          alert('Enter valid interests!');
          return;
        }
        if (this.$refs.form3.validate()){
          let data = {
            username: this.person.username,
            password: this.credentials.password,
            name: this.person.name,
            surname: this.person.surname,
            email: this.credentials.email,
            telephone: this.person.telephone,
            gender: this.person.gender,
            birthDate: this.person.birthDate,
            isPrivate: this.isPrivate,
            biography: this.person.biography,
            webSite: this.person.webSite,
            interestedIn: this.person.interests
          }
          axios({
            method: "post",
            url: comm.protocol + "://" + comm.server +"/api/profile/",
            data: JSON.stringify(data),
          }).then((response) => {
            if (response.status == 200) {
              if(response.data.message == 'ok'){
                  alert('Check your email!');
                  this.alert = false;
              }
            if(response.data.message=="Invalid data."){
                this.alert = true;
                  if(response.data.errors.includes("Password")){
                    this.alertText = "Password is too weak. Please choose another password."
                  }
            }
            if(response.data.message == "Server error while registering."){
                    this.alert = true;
                    this.alertText = "Chosen e-mail already exists.Please choose another mail."
                }
            }
          })
        }
      }
  }
}
</script>