<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog transition="dialog-top-transition" max-width="600">
        <template v-slot:activator="{ on, attrs }">
            <v-btn v-bind="attrs" v-on="on" color="primary"  @click="clearInputs()">Add</v-btn>
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-toolbar color="primary" dark>Select image</v-toolbar>
            <v-card-text>
                <v-file-input
                    v-model="file"
                    label="Input pictures.."
                ></v-file-input>

                <v-text-field v-model="link" label="Link"></v-text-field>

            </v-card-text>
            <v-card-actions class="justify-end">
                <v-btn text @click="confirm()" >Confirm</v-btn>
                <v-btn text @click="dialog.value = false" >Close</v-btn>
            </v-card-actions>
          </v-card>
        </template>
      </v-dialog>
    </v-col>
  </v-row>
</template>

<script>
export default {
    data() {return {
        file: {},
        link: ''
    }},
    methods: {
       confirm(){
           let data = {
               file : this.file,
               link : this.link
           }
           this.$emit('addPictureWithLink',data)
       },
       clearInputs(){
           this.file = {},
           this.link = ''
       }
    }
}
</script>

<style>

</style>