
<template>
  <v-row justify="center">
    <v-col
      cols="12"
      sm="10"
      md="8"
      lg="6"
    >
      <v-card ref="form">
        <v-card-text>
          <v-text-field
            v-model="description"
            :counter="255"
            :rules="descriptionRules"
            label="Description"
            required
          ></v-text-field>

          <v-autocomplete
            ref="selectedPostType"
            v-model="selectedPostType"
            :rules="[() => !!selectedPostType || 'Post type is required']"
            :items="postTypes"
            label="Post type"
            placeholder="Select post type..."
            required
          ></v-autocomplete>

          <v-select
            v-model="selectedLocation"
            :items="locations"
            label="Location"
          ></v-select>

          <v-checkbox
            v-model="isHighlighted"
            label="Is highlighted?"
          ></v-checkbox>

          <v-checkbox
            v-model="isCampaign"
            label="Is campaign?"
          ></v-checkbox>

          <v-checkbox
            v-model="isCloseFriendsOnly"
            label="Is close friends only?"
          ></v-checkbox>

          <v-file-input
            v-model="files"
            multiple
            chips
            label="Input pictures.."
          ></v-file-input>
      </v-card-text>
      <v-card-actions>
          <v-btn text>
            Cancel
          </v-btn>
          <v-spacer></v-spacer>
          <v-slide-x-reverse-transition>
            <v-tooltip
              v-if="formHasErrors"
              left
            >
              <template v-slot:activator="{ on, attrs }">
                <v-btn
                  icon
                  class="my-0"
                  v-bind="attrs"
                  @click="resetForm"
                  v-on="on"
                >
                  <v-icon>mdi-refresh</v-icon>
                </v-btn>
              </template>
              <span>Refresh form</span>
            </v-tooltip>
          </v-slide-x-reverse-transition>
          <v-btn
            color="primary"
            text
            @click="submit"
          >
            Submit
          </v-btn>
        </v-card-actions>
    </v-card>
    </v-col>
  </v-row>
</template>

<script>
  import axios from 'axios'
  import * as comm from '../configuration/communication.js'
  export default {

    name: 'CreatePost',

    data: () => ({
      valid: true,
      description: '',
      descriptionRules: [
        v => !!v || 'Description is required',
        v => (v && v.length <= 255) || 'description must be less than 255 characters',
      ],
      selectedPostType: null,
      selectedLocation: null,
      postTypes: [
        'Post',
        'Story'
      ],
      locations: [
        'Novi Sad',
        'Belgrade'
      ],
      isHighlighted: false,
      isCampaign: false,
      isCloseFriendsOnly: false,
      files : null,
    }),

    methods: {
      resetForm () {
        Object.keys(this.form).forEach(f => {
          this.$refs[f].reset()
        })
      },
      submit () {
        let dto = {"description" : this.description, "isHighlighted" : this.isHighlighted, "isCampaign" : this.isCampaign,
        "isCloseFriendsOnly": this.isCloseFriendsOnly, "location" : this.selectedLocation, 
        "hashTags" : [], "taggedUsers" : [], "postType" : this.selectedPostType}
        let json = JSON.stringify(dto);
        console.log(json);
        const data = new FormData();
         data.append("files", this.files);
         data.append("data", json);
        axios({
          method: "post",
          url: "http://" + comm.server + "/api/post",
          data: data,
          config: { headers: {'Content-Type': 'multipart/form-data' }}
        }).then(function (response) {
          console.log(response);
        })
        .catch(function (response) {
          console.log(response);
        });
      }
    }
  }
</script>

<style>

</style>