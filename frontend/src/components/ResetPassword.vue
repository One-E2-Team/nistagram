<template>
  <v-form
    ref="form"
    v-model="valid"
    lazy-validation>
    <v-container fluid>
      <v-row align="center"
      justify="center">
        <v-col
          cols="12"
          sm="4"
        >
          <v-text-field
            v-model ="password1"
            :append-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
            :rules="[rules.required, rules.min]"
            :type="show ? 'text' : 'password'"
            label="Not visible"
            hint="At least 8 characters"
            class="input-group--focused"
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
          <v-text-field
            v-model="password2"
            :rules="[rules.required, rules.passwordMatch]"
            :type="'password'"
            label="Repeat password "
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
                @click="resetPassword"
                >
                Confirm
                </v-btn>
            </v-col>
        </v-row>
    </v-container>
  </v-form>
</template>

<script>
  export default {
    data () {
      return {
        valid: true,
        show: false,
        password1: '',
        password2: '',
        rules: {
          required: value => !!value || 'Required.',
          min: v => v.length >= 8 || 'Min 8 characters',
          passwordMatch: () => (this.password1 === this.password2) || 'Password must match'
        },
      }
    },
    methods:{
      resetPassword(){
        if (this.$refs.form.validate()){
            //TODO: send axios
            console.log("validation pass!")
        }
      }
    }
  }
</script>