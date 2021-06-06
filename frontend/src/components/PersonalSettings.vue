<template>
<v-container>
  <v-row align="center" justify="center">
    <v-col cols="12" sm="6" >
      <h1 class="display-2 font-weight-bold mb-3">
         Personal data
      </h1>
      <v-form
        ref="form"
      >
        <v-text-field
            v-model="person.username"
            :rules="[ rules.required , rules.name] "
            label="Username:"
            required
            ></v-text-field>
            
        <v-text-field
            v-model="person.email"
            :rules="[ rules.required , rules.name] "
            label="Email:"
            required
            ></v-text-field>

        <v-text-field
            v-model="person.name"
            :rules="[ rules.required , rules.name] "
            label="Name:"
            required
            ></v-text-field>   
        <v-text-field
            v-model="person.surname"
            :rules="[ rules.required , rules.name] "
            label="Surname:"
            required
            ></v-text-field>

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

        <v-text-field
            v-model="person.telephone"
            :rules="[ rules.required , rules.name] "
            label="Telephone:"
            required
            ></v-text-field>
        
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
        <v-btn
          color="success"
          class="mr-4"
          @click="updateSettings"
        >
          Confirm
        </v-btn>

      </v-form>
    </v-col>
  </v-row>
</v-container>
</template>

<script>
export default {
    data(){
        return{
            person: {
                username: '',
                name: '',
                surname: '',
                email: '',
                telephone: '',
                gender: '',
                birthDate: '',
                biography: '',
                webSite: '',
            } , 
            rules: {
                required: value => !!value || 'Required.',
                min: v => v.length >= 8 || 'Min 8 characters',
                email: v => /.+@.+\..+/.test(v) || 'E-mail must be valid',
                name: v => (v && v.length <= 10) || 'Name must be less than 10 characters',
                emailMatch: () => (`The email and password you entered don't match`),
                }, 
            menu: false
        }
    },
    methods:{
        updateSettings(){

        }
    }
}
</script>