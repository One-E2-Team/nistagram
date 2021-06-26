<template>
    <div>
        <template>
          <div class="text-center" v-if="showCart" @close="showCart=false">
            <v-dialog
              v-model="showCart"
              width="600px"
            >
              <v-card class="mx-auto my-12" width="600px" >
                  <v-card-title class="justify-center">Your shopping cart</v-card-title>
              </v-card>
            </v-dialog>
          </div>
          <v-container>
            <v-row>
              <v-col v-for="p in products" :key="p.id" cols="12" sm="4">
                <v-card class="mx-auto my-12" max-width="374" >
                  <template slot="progress">
                    <v-progress-linear color="deep-purple" height="10" indeterminate ></v-progress-linear>
                  </template>

                  <v-img height="250" src="https://cdn.vuetifyjs.com/images/cards/cooking.png" ></v-img>

                  <v-card-title>{{p.name}}</v-card-title>

                  <v-card-text>
                    <div class="my-4 text-subtitle-1">
                      {{p.pricePerItem}} RSD
                    </div>

                    <div>{{p.description}}</div>
                  </v-card-text>

                  <v-divider class="mx-4"></v-divider>

                  <v-card-title>Choose amount:</v-card-title>

                  <v-card-text>
                    <v-chip-group column >
                      <v-chip>
                        <input v-model="productAmount[p.id]" type="number" min="1" value="1" style="width: 50px">
                      </v-chip>
                      <v-chip>
                        <v-card-actions>
                          <v-btn color="deep-purple lighten-2" text @click="addToCart(p)" >
                            Add to cart
                          </v-btn>
                        </v-card-actions>
                      </v-chip>
                    </v-chip-group>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-container>
        </template>
    </div>
</template>

<script>
import axios from 'axios'
import * as comm from '../configuration/communication.js'
import { bus } from '../main'
  export default {
    name: "HomePage",
    data() {return {
      products: [],
      productAmount:{},
      cart:[],
      showCart : false
    }},
    created(){
      bus.$on('show-cart', (data) => {
      this.showCart = data;
    })
    },
    mounted(){
       axios({
                method: "get",
                url: comm.protocol +'://' + comm.server + '/product',
                headers: comm.getHeader(),
            }).then(response => {
              if(response.status==200){
                this.products = response.data;
                console.log(this.products);
              }
            }).catch((response) => {
              console.log(response.data)
            });
    },
    methods: {
     addToCart(p){
       p.amount = this.productAmount[p.id];
       this.cart.push(p);
     }
    },
  }
</script>

