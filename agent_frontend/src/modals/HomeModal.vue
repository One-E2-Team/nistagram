<template>
  <div class="text-center">
    <v-dialog
      v-model="show"
      max-width="250px"
    >

      <v-card>
        <v-card-text>
            <v-btn
              label="Log out"
              class="my-2"
              style="color:red"
              @click="logOut()"
              width="200"
            >Log out</v-btn><br/>
        </v-card-text>
      </v-card>
    </v-dialog>
    </div>
</template>

<script>
import * as comm from '../configuration/communication.js'
export default {
  props: ['visible'],
  data(){
    return {
    }
  },
  created(){
  },
  methods: {
      logOut(){
        comm.logOut();
        this.$router.push({name: 'Home'});
      }
  },
  computed: {
    show: {
      get () {
        return this.visible;
      },
      set (value) {
        if (!value) {
          this.$emit('close');
        }
      }
    }
  }
}
</script>