<template>
  <v-form
    ref="form"
    v-model="valid"
    lazy-validation
  >
    <v-container >
        <v-row align="center"
      justify="center">
            <v-col
            cols="12"
            sm="4"
            >
            <v-text-field
                v-model="name"
                :rules="[ rules.email , rules.required] "
                label="Mail:"
                required
                ></v-text-field>
            </v-col>
        </v-row>
        <v-row align="center"
      justify="center">
            <v-col
            cols="12"
            sm="4"
            >
            <v-text-field
                v-model="password"
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
        <v-row align="center"
      justify="center">
            <v-col
            cols="12"
            sm="4"
            >
                <v-btn
                :disabled="!valid"
                color="success"
                class="mr-4"
                @click="login"
                >
                Log in
                </v-btn>
            </v-col>
        </v-row>
    </v-container>
  </v-form>
</template>

<script>
  export default {
    data: () => ({
      show: false,
      valid: true,
      name: '',      
      email: '',
      password: '',

      rules: {
          required: value => !!value || 'Required.',
          min: v => v.length >= 8 || 'Min 8 characters',
          email: v => /.+@.+\..+/.test(v) || 'E-mail must be valid',
          name: v => (v && v.length <= 10) || 'Name must be less than 10 characters',
          emailMatch: () => (`The email and password you entered don't match`),
        },
    }),

    methods: {
      login () {
        if (this.$refs.form.validate()){
            //TODO: send axios
            console.log("validation pass!")
        }
      },
      
    },
  }
</script>