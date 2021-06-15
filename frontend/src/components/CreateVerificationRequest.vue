<template>
  <div>
      <v-row justify="center">
    <v-col
      cols="12"
      sm="10"
      md="8"
      lg="6"
    >
      <v-form ref="form"  v-model="valid" lazy-validation>
      <v-card >
        <v-card-text>
          <v-text-field
            v-model="name"
            label="Name"
            :rules="[rules.required, rules.max255]"
            required
          ></v-text-field>

          <v-text-field
            v-model="surname"
            label="Surname"
            :rules="[rules.required, rules.max255]"
            required
          ></v-text-field>

          <v-autocomplete
            ref="selectedCategory"
            v-model="selectedCategory"
            :rules="[rules.required]"
            :items="categories"
            label="Category"
            placeholder="Select category..."
            required
          ></v-autocomplete>

          <v-file-input
            v-model="picture"
            accept="image/*"
            :rules="[rules.required]"
            chips
            label="Input picture.."
          ></v-file-input>
      </v-card-text>
      <v-card-actions>
          <v-btn text>
            Cancel
          </v-btn>
          <v-spacer></v-spacer>
          <v-slide-x-reverse-transition>
            <v-tooltip
              left
            >
              <template v-slot:activator="{ on, attrs }">
                <v-btn
                  icon
                  class="my-0"
                  v-bind="attrs"
                  @click="resetForm"
                  v-on="on"
                >
                  <v-icon>mdi-refresh</v-icon>
                </v-btn>
              </template>
              <span>Refresh form</span>
            </v-tooltip>
          </v-slide-x-reverse-transition>
          <v-btn
            color="primary"
            text
            @click="submit()"
          >
            Submit
          </v-btn>
        </v-card-actions>
    </v-card>
      </v-form>
    </v-col>
  </v-row>
  </div>
</template>

<script>
  import axios from 'axios'
  import * as comm from '../configuration/communication.js'
  import * as validator from '../plugins/validator.js'
  export default {

    name: 'CreateVerificationRequest',

    data() {return {
      valid: true,
      rules: validator.rules,
      name: '',
      surname: '',
      categories: [],
      selectedCategory : '',
      picture: null
    }},
    mounted(){
     // if( !comm.isUserLogged() )
       // this.$router.push({name: 'NotFound'});
       
     axios({
          method: "get",
          url: comm.protocol + "://" + comm.server + "/api/profile/categories"
        }).then(response => {
          console.log(response);
          this.categories = response.data.collection;
        })
    },
    methods: {
      resetForm () {
        this.valid= true
        this.name= ''
        this.surname= ''
        this.selectedCategory = null
        this.picture = null
      },
      submit () {
        if(this.$refs.form.validate() !== true )
          return;
        
        let dto = {"name" : this.name, "surname" : this.surname, "category" : this.selectedCategory}
        let json = JSON.stringify(dto);
        const data = new FormData();
        data.append("picture", this.picture)
        data.append("data", json);
        axios.defaults.headers.common['Authorization'] = 'Bearer ' + comm.getJWTToken().token;
        axios({
          method: "post",
          url: comm.protocol + "://" + comm.server + "/api/profile/verification-request",
          data: data,
          config: { headers: {...data.headers}}
        }).then(response => {
          console.log(response);
          delete axios.defaults.headers.common["Authorization"];
          alert("Post is successfully created!")
        })
        .catch(response => {
          delete axios.defaults.headers.common["Authorization"];
          console.log(response);
        });
    
      }
    }
  }
</script>