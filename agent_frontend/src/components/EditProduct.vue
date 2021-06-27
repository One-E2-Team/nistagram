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
            <v-text-field readonly
                v-model="id"
                label="Id"
            ></v-text-field>
            
            <v-text-field 
                v-model="name"
                label="Name"
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
  import { bus } from '../main'
  export default {

    name: 'EditProduct',

    data() {return {
      id : 0,  
      name: "",
      price: "",
      quantity: ""
    }},
    created(){
      bus.$on('product-data', (item) => {
      this.id = item.id;
      this.name = item.name;
      this.price = item.pricePerItem;
      this.quantity = item.quantity;
    })
    },
    mounted(){
      if( !comm.hasRole("AGENT") )
        this.$router.push({name: 'NotFound'})
    },
    methods: {
      resetForm () {
        this.valid= true;
        this.name = "";
        this.price = "";
        this.quantity = "";
      },
      submit () {
        let dto = { 
            "productId" : this.id, "name" : this.name, "quantity" : parseInt(this.quantity), 
            "pricePerItem" : parseFloat(this.price) 
        };
        let json = JSON.stringify(dto);
        axios({
          method: "put",
          url: comm.protocol + "://" + comm.server + "/product",
          data: json,
          headers: comm.getHeader()
        }).then(response => {
          console.log(response);
          delete axios.defaults.headers.common["Authorization"];
          alert("Product is successfully edited!")
        })
        .catch(response => {
          delete axios.defaults.headers.common["Authorization"];
          console.log(response);
        });
    
      },
      
    }
  }
</script>