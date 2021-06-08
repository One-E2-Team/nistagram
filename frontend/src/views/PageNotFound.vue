<template>
  <div class="home">
    <v-container>
        <v-row  align="center" justify="center" class="text-center">
        <v-col class="mb-4">
            <h1 class="display-2 font-weight-bold mb-3">
            Can't find page 
            </h1>
        </v-col>
        </v-row>

        <v-row justify="center" align="center" class="text-center">

            <v-col cols="12" sm="4" md="3" >
                <v-img
                    :src="require('../assets/page-not-found.jpg')"
                    class="my-5"
                    contain
                    
                    width="400"
                    style=" border: 2px solid rgba(0,0,0, 0.2); box-shadow: 20px 10px 100px; border-radius: 10%;" />
            </v-col>
        </v-row> 
    </v-container>
    
  </div>
</template>