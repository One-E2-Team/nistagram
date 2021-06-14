<template>
  <div class="text-center">
    <v-dialog
      v-model="dialog"
      width="500"
    >

      <v-card>
        <v-card-text>
            <v-btn
              label="Report"
              color="red"
              @click="report()"
            ></v-btn>

        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
export default {
    methods:{
        report(){
            alert('Report')
        }
    }
}
</script>