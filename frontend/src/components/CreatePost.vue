
<template>
  <v-row justify="center">
    <v-col
      cols="12"
      sm="10"
      md="8"
      lg="6"
    >
      <v-form ref="form"  v-model="valid" lazy-validation>
      <v-card >
        <v-card-text>
          <v-text-field @keyup="e => findTag(e)"
            v-model="description"
            label="Description"
            :rules="[rules.max255]"
          ></v-text-field>

          <v-list rounded>
              <v-list-item-group
                color="primary"
              >
                <v-list-item @click="setTag(item)"
                  v-for="(item, i) in searchedTaggedUsers"
                  :key="i"
                >
                  <v-list-item-icon>
                     <v-icon>mdi-account</v-icon>
                  </v-list-item-icon>
                  <v-list-item-content>
                    <v-list-item-title v-text="item"></v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </v-list-item-group>
            </v-list>

          <v-autocomplete
            ref="selectedPostType"
            v-model="selectedPostType"
            :rules="[rules.required, validPostType]"
            :items="postTypes"
            label="Post type"
            placeholder="Select post type..."
            required
          ></v-autocomplete>

          <v-text-field
            v-model="location"
            :counter="255"
            :rules="[rules.max255]"
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
            :rules="[rules.max255]"
            label="Hash tags"
          ></v-text-field>

          <v-file-input
            v-model="files"
            multiple
            :rules="[rules.oneOrMoreElement]"
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
            @click="submit()"
          >
            Submit
          </v-btn>
        </v-card-actions>
    </v-card>
      </v-form>
    </v-col>
  </v-row>
</template>

<script>
  import axios from 'axios'
  import * as comm from '../configuration/communication.js'
  import * as validator from '../plugins/validator.js'
  export default {

    name: 'CreatePost',

    data() {return {
      valid: true,
      description: '',
      searchedTaggedUsers : [],
      cursorStart: -1,
      cursorEnd: -1,
      rules: validator.rules,
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
      validPostType: (v) => (v == 'Post' || v == 'Story') || 'Post type must be selected'
    }},
    mounted(){
      if( !comm.isUserLogged() )
        this.$router.push({name: 'NotFound'})
    },
    methods: {
      resetForm () {
        this.valid= true
        this.description= ''
        this.selectedPostType= null
        this.location= ''
        this.hashTags= ''
        
        this.isHighlighted= false
        this.isCloseFriendsOnly= false
        this.files = []
      },
      submit () {
        if(this.$refs.form.validate() !== true )
          return
        
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
          delete axios.defaults.headers.common["Authorization"];
          alert("Post is successfully created!")
        })
        .catch(response => {
          delete axios.defaults.headers.common["Authorization"];
          console.log(response);
        });
    
      },
      findTag(e){
        let end = e.target.selectionStart -1;
        if (this.description[end] == '@')
          return;
        for(let i = end; i >= 0; i--){
          if(this.description[i] == ' ')
            break;
          if(this.description[i] == '@'){
              this.cursorStart = i + 1;
              this.cursorEnd = end;
              this.searchUsername(this.description.slice(i+1, end + 1));
              return;
          }
        }
        this.searchedTaggedUsers = [];
      },
      searchUsername(username){
        axios({
          method: "get",
          url: comm.protocol + '://' + comm.server + '/api/profile/search-for-tag/' + username,
          headers: comm.getHeader()
        }).then(response => {
          if(response.status==200){
            this.searchedTaggedUsers = response.data.collection;
          }
        })
      },
      setTag(item){
        this.description = this.description.slice(0, this.cursorStart - 1) +
                           item + this.description.slice(this.cursorEnd + 1, this.description.length - 1);
        this.searchedTaggedUsers = [];
      }
    }
  }
</script>