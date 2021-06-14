<template>
  <div class="text-center">
    <v-dialog
      v-model="show"
      max-width="250px"
    >

      <v-card>
        <v-card-text>
            <v-btn
              label="Report"
              class="my-2"
              style="color:red"
              @click="report()"
              width="200"
            >Report</v-btn><br/>
            <v-btn
              label="Report"
              width="200"
              @click="unfollow()"
            >Unfollow</v-btn><br/>
            <v-btn
              label="Mute"
              class="my-2"
              width="200"
              @click="mute()"
            >Mute</v-btn>

        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
export default {
  props: ['visible'],
  methods:{
      report(){
          //TODO: send report axios
      },
      unfollow(){
          //TODO: send unfollow axios
      },
      mute(){
          //TODO: send mute axos
      }
  },
  computed: {
    show: {
      get () {
        return this.visible
      },
      set (value) {
        if (!value) {
          this.$emit('close')
        }
      }
    }
  }
}
</script>