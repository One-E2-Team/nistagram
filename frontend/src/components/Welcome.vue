<template>
  <v-container>
    <v-row class="text-center">
      <v-col class="mb-4">
        <h1 class="display-2 font-weight-bold mb-3">
          Welcome to Nistagram
        </h1>

        <p class="subheading font-weight-regular">
          For a lot of nudity, fake profiles and other things you can’t do without 
          <br>please join us
        </p>
      </v-col>

      <v-col
        class="mb-5"
        cols="12"
      >
        <v-row justify="center">
         <v-btn
                color="success"
                class="mr-4"
                elevation="8"
                large
                @click="login"
                >
                Log in
                </v-btn>
        </v-row>
      </v-col>

      <v-col
        class="mb-5"
        cols="12"
      >
        <v-btn
                color="primary"
                class="mr-4"
                elevation="8"
                large
                @click="register"
                >
                Register
                </v-btn>
         <v-img
          :src="require('../assets/minions.jpg')"
          class="my-3"
          contain
          height="200"
        />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
  export default {
    name: 'HelloWorld',

    methods: {
      login(){
        this.$router.push({ name: 'Login' })
      },
      register(){
        this.$router.push({ name: 'Register' })
      }
    }
  }
</script>
