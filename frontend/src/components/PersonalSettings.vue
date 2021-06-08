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
            :rules="[ rules.required , rules.email] "
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
import * as validator from '../plugins/validator.js'
import axios from 'axios'
  import * as comm from '../configuration/communication.js'
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
            rules: validator.rules,
            menu: false
        }
    },
    created(){
        axios({
            method: 'get',
            url: "http://" + comm.server + "/api/profile/my-personal-data",
            headers: comm.getHeader(),
            }).then(response => {
                if(response.status == 200){
                    this.person = response.data
                }
            });
    },
    methods:{
        updateSettings(){
            axios({
            method: 'put',
            url: "http://" + comm.server + "/api/profile/my-personal-data",
            headers: comm.getHeader(),
            data: JSON.stringify(this.person)
            }).then(response => {
                if(response.status == 200){
                console.log(response.data)
                }
            });
        }
    }
}
</script>