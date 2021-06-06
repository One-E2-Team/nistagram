<template>
<v-container>
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
        <v-form ref="form1" v-model="valid" lazy-validation>
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
              <v-col cols="12" sm="6">
                <v-btn
                :disabled="!valid"
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
        <v-form ref="form2" v-model="valid" lazy-validation class="text-center">
          <v-container >
            <v-row align="center" justify="center">
              <v-col cols="12" sm="8" >
                <v-text-field
                  v-model="person.username"
                  :rules="[ rules.required , rules.name] "
                  label="Username:"
                  required
                  ></v-text-field>
              </v-col>
            </v-row>
          <v-row align="center" justify="center">
            <v-col cols="12" sm="8" >
              <v-text-field
                v-model="person.name"
                :rules="[ rules.required , rules.name] "
                label="Name:"
                required
                ></v-text-field>
            </v-col>
          </v-row>
          <v-row align="center" justify="center">
            <v-col cols="12" sm="8" >  
              <v-text-field
                v-model="person.surname"
                :rules="[ rules.required , rules.name] "
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
                :rules="[ rules.required , rules.name] "
                label="Telephone:"
                required
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
                @click="continueTo3">
                Continue
                </v-btn>
              <v-btn
              color="normal"
              class="d-flex justify-space-around mb-6"
              @click="e1=1">
              Back
              </v-btn>
            </v-col>
          </v-row>
          </v-container>

        </v-form>
      </v-stepper-content>

      <!--Step3 content -->
      <v-stepper-content step="3">
        <v-form ref="form3" v-model="valid" lazy-validation class="text-center">
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
              @click="e1=1">
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

  export default {
    data: () => ({
      e1: 1,
      show: false,
      valid: true,
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

      rules: {
          required: value => !!value || 'Required.',
          min: v => v.length >= 8 || 'Min 8 characters',
          email: v => /.+@.+\..+/.test(v) || 'E-mail must be valid',
          name: v => (v && v.length <= 10) || 'Name must be less than 10 characters',
          emailMatch: () => (`The email and password you entered don't match`),
        },

        menu: false
    }),

    created(){
      axios({
          method: "get",
          url: "http://" + comm.server +"/api/profile/interests",
        }).then(response => {
          if (response.status == 200) {
            this.interests = response.data.collections
          }
        })
        .catch(response => {
          console.log(response);
        });
    },

    methods: {
      continueTo2 () {
        if (this.$refs.form1.validate()){
            this.e1 = 2 
        }
        this.e1 = 2
      },
      continueTo3(){
        if (this.$refs.form2.validate()){
            this.e1 = 3
        }
        this.e1 = 3
      },
      remove (item) {
        this.person.interests.splice(this.person.interests.indexOf(item), 1)
        this.person.interests = [...this.person.interests]
      },
      register (){
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
          interests: this.person.interests
        }
        axios({
          method: "post",
          url: "http://" + comm.server +"/api/profile/",
          data: JSON.stringify(data),
        }).then(function (response) {
          if (response.status == 200) {
            alert('Check your email!');
            //TODO: redirect on home page
          }
        })
        .catch(function (response) {
          console.log(response);
        });
      }
      
    },
  }
</script>