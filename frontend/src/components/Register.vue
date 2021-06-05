<template>
  <v-stepper v-model="e1">
    <v-stepper-header>
      <v-stepper-step :complete="e1 > 1" step="1">
        Credentials
      </v-stepper-step>

      <v-divider></v-divider>

      <v-stepper-step :complete="e1 > 2" step="2">
        Personal data
      </v-stepper-step>

      <v-divider></v-divider>

      <v-stepper-step step="3">
        Name of step 3
      </v-stepper-step>
    </v-stepper-header>

    <v-stepper-items class="text-center">
      <!-- Step 1  component-->
      <v-stepper-content step="1">
        <v-form ref="form1" v-model="valid" lazy-validation>
          <v-container >
            <v-row align="center" justify="center">
              <v-col cols="12" sm="4" >
                <v-text-field
                    v-model="credentials.email"
                    :rules="[ rules.email , rules.required] "
                    label="Mail:"
                    required
                    ></v-text-field>
              </v-col>
            </v-row>
            <v-row align="center" justify="center">
              <v-col cols="12" sm="4">
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
              <v-col cols="12" sm="4">
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
              <v-col cols="12" sm="6" >
                <v-text-field
                  v-model="person.username"
                  :rules="[ rules.required , rules.name] "
                  label="Username:"
                  required
                  ></v-text-field>
              </v-col>
            </v-row>
          <v-row align="center" justify="center">
            <v-col cols="12" sm="6" >
              <v-text-field
                v-model="person.name"
                :rules="[ rules.required , rules.name] "
                label="Name:"
                required
                ></v-text-field>
            </v-col>
          </v-row>
          <v-row align="center" justify="center">
            <v-col cols="12" sm="6" >  
              <v-text-field
                v-model="person.surname"
                :rules="[ rules.required , rules.name] "
                label="Surname:"
                required
                ></v-text-field>
            </v-col>
          </v-row>
          <v-row align="center" justify="center">
            <v-col cols="12" sm="6" >  
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
            <v-col cols="12" sm="6" >          
              <v-text-field
                v-model="person.telephone"
                :rules="[ rules.required , rules.name] "
                label="Telephone:"
                required
                ></v-text-field>
            </v-col>
          </v-row>
          <v-row align="center" justify="center">
            <v-col cols="12" sm="6" >          
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
              <v-col cols="12" sm="4" class="d-flex justify-space-around mb-6">
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
        <v-card
          class="mb-12"
          color="grey lighten-1"
          height="200px"
        ></v-card>

        <v-btn
          color="primary"
          @click="e1 = 1"
        >
          Continue
        </v-btn>

        <v-btn text>
          Cancel
        </v-btn>
      </v-stepper-content>
    </v-stepper-items>
  </v-stepper>
</template>

<script>
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
        birthDate: ''
      } ,     

      rules: {
          required: value => !!value || 'Required.',
          min: v => v.length >= 8 || 'Min 8 characters',
          email: v => /.+@.+\..+/.test(v) || 'E-mail must be valid',
          name: v => (v && v.length <= 10) || 'Name must be less than 10 characters',
          emailMatch: () => (`The email and password you entered don't match`),
        },

        menu: false
    }),

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
      }
      
    },
  }
</script>