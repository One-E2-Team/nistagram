<template>
    <v-container>
    <v-row justify="center">
        <v-col
        cols="12"
        sm="10"
        md="8"
        lg="6"
        >
        <v-form ref="form" lazy-validation>
        <v-card >
            <v-card-text>
            <v-text-field 
                v-model="name"
                label="Name"
            ></v-text-field>

            <v-text-field 
                v-model="description"
                label="Description"
            ></v-text-field>

            <v-text-field
                v-model="price"
                :counter="255"
                label="Price"
            ></v-text-field>

            <v-text-field
                v-model="quantity"
                :counter="255"
                label="Quantity"
            ></v-text-field>

            <v-file-input
                v-model="file"
                
                label="Input picture.."
            ></v-file-input>
        </v-card-text>
        <v-card-actions>
            <v-btn text @click="resetForm()">
                Cancel
            </v-btn>
            <v-spacer></v-spacer>
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
    </v-container>
</template>

<script>
  import axios from 'axios'
  import * as comm from '../configuration/communication.js'
  export default {

    name: 'NewProduct',

    data() {return {
      name: "",
      description: "",
      price: "",
      quantity: "",
      file: {}
    }},
    mounted(){
      if( !comm.hasRole("AGENT") )
        this.$router.push({name: 'NotFound'})
    },
    methods: {
      resetForm () {
        this.valid= true;
        this.name = "";
        this.description = "",
        this.price = "",
        this.quantity = "",
        this.file = {}
      },
      submit () {
        let dto = { 
            "name" : this.name, "quantity" : parseInt(this.quantity), 
            "pricePerItem" : parseFloat(this.price) 
        };
        let json = JSON.stringify(dto);
        let data = new FormData();
        data.append("data", json);
        data.append("file", this.file);
        axios.defaults.headers.common['Authorization'] = 'Bearer ' + comm.getJWTToken().token;
        axios({
          method: "post",
          url: comm.protocol + "://" + comm.server + "/product",
          data: data,
          config: { headers: {...data.headers}}
        }).then(response => {
          console.log(response);
          delete axios.defaults.headers.common["Authorization"];
          alert("Product is successfully created!")
        })
        .catch(response => {
          delete axios.defaults.headers.common["Authorization"];
          console.log(response);
        });
    
      },
      
    }
  }
</script>