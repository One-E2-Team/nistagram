<template>
  <v-form>
    <v-container>
      <v-row align="center" justify="center">
        <v-col cols="12" sm="9" md="5" >
          <v-text-field
            v-model="searchParams"
            label="Search .."
          ></v-text-field>
        </v-col>

        <v-col cols="12" sm="6" md="2">
         <v-card flat class="py-0">
            <v-card-text>    
                <v-btn-toggle
                mandatory
                >
                <v-btn @click="searchType='accounts'">
                    <v-icon>mdi-account</v-icon>
                </v-btn>
                <v-btn @click="searchType='locations'">
                    <v-icon>mdi-map-marker</v-icon>
                </v-btn>
                <v-btn  @click="searchType='hashtags'">
                    <v-icon>mdi-pound</v-icon>
                </v-btn>
                </v-btn-toggle>
            </v-card-text>
        </v-card>
        </v-col>

        <v-col cols="12" sm="6" md="1" >
            <v-btn  @click="search()">
                <v-icon >mdi-magnify</v-icon> Search
            </v-btn>
        </v-col>
      </v-row>
    </v-container>
  </v-form>
</template>

<script>
export default {
    name: "Search",
    data(){
        return{
            searchParams : '',
            searchType: 'accounts' // possible values (accounts|locations|hashtags)
    }},
    methods:{
        search(){
            if (this.searchType == 'locations'){
                this.$emit('searched-result', this.searchLocation())
            }else if (this.searchType == 'hashtags'){
                this.$emit('searched-result', this.searchHashTags())
            }
        },
        searchLocation(){
        let ret = [];
        this.allPosts.forEach((post) => {
                if((post.location.toLowerCase()).includes(this.location.toLowerCase())){
                  ret.push(post);
                }
            });
        return ret;
      },
      searchHashTags(){
        let ret = [];
        this.allPosts.forEach((post) => {
                if((post.hashTags.toLowerCase()).includes(this.hashTags.toLowerCase())){
                  ret.push(post);
                }
            });
        return ret;
      }
    }
}
</script>

<style>

</style>