<template>
  <v-form>
    <v-container fluid>
      <v-row>
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
      <v-row >
        <v-col
          cols="12"
          sm="4"
        >
          <v-text-field
            v-model="password2"
            :rules="[rules.passwordMatch]"
            :type="'password'"
            label="Repeat password "
          ></v-text-field>
        </v-col>
      </v-row>
    </v-container>
  </v-form>
</template>

<script>
  export default {
    data () {
      return {
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
  }
</script>