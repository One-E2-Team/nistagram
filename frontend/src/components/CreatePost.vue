
<template>
  <v-row justify="center">
    <v-alert
    v-if="alert"
    :value="alertText"
    color="red"
    type="error"
    dismissible
    text
    v-model="alertText"
  >{{alertText}}</v-alert>
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
            label="Description"
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

          <v-text-field
            v-model="location"
            :counter="255"
            label="Location"
          ></v-text-field>

          <v-checkbox v-if="selectedPostType === 'Story'"
            v-model="isHighlighted"
            label="Is highlighted?"
          ></v-checkbox>

          <v-checkbox v-if="selectedPostType === 'Story'"
            v-model="isCloseFriendsOnly"
            label="Is close friends only?"
          ></v-checkbox>

          <v-text-field
            v-model="hashTags"
            :counter="255"
            label="Hash tags"
          ></v-text-field>

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

    data() {return {
      alert : false,
      alertText : '',
      valid: true,
      description: '',
      descriptionRules: [
        v => !!v || 'Description is required',
        v => (v && v.length <= 255) || 'description must be less than 255 characters',
      ],
      selectedPostType: null,
      location: '',
      hashTags: '',
      postTypes: [
        'Post',
        'Story'
      ],
      isHighlighted: false,
      isCloseFriendsOnly: false,
      files : [],
    }},
    mounted(){
      if( !comm.isUserLogged() )
        this.$router.push({name: 'NotFound'})
    },
    methods: {
      resetForm () {
        Object.keys(this.form).forEach(f => {
          this.$refs[f].reset()
        })
      },
      submit () {
        if(this.files.length == 0){
          this.alert = true;
           this.alertText = "Please choose at least one picture or video."
        } else if(this.postTypes != 'Post' && this.postTypes != 'Story' ){
          this.alert = true;
          this.alertText = "Post type must be selected"
        }
        else{
          this.alert = false;
          this.alertText = ""
          let dto = {"description" : this.description, "isHighlighted" : this.isHighlighted,
          "isCloseFriendsOnly": this.isCloseFriendsOnly, "location" : this.location, 
          "hashTags" : this.hashTags, "taggedUsers" : [], "postType" : this.selectedPostType}
          let json = JSON.stringify(dto);
          const data = new FormData();
          for(let i = 0;i < this.files.length;i++){
              data.append("file" + i, this.files[i], this.files[i].name);
            }
          data.append("data", json);
          axios.defaults.headers.common['Authorization'] = 'Bearer ' + comm.getJWTToken().token;
          axios({
            method: "post",
            url: comm.protocol + "://" + comm.server + "/api/post",
            data: data,
            config: { headers: {...data.headers}}
          }).then(response => {
            console.log(response);
            alert("Post is successfully created!")
          })
          .catch(response => {
            console.log(response);
          });
      }
      }
    }
  }
</script>

<style>

</style>