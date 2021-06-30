<template>
    <post-type-tabs-slider :posts="posts"/>
</template>

<script>
import PostTypeTabsSlider from '../components/PostTypeTabsSlider.vue'
import axios from 'axios'
import * as comm from '../configuration/communication.js'
export default {
  components: { PostTypeTabsSlider },
    data(){
        return{
            posts: [],
        }
    },
    created(){
        axios({
            method: "get",
            url: comm.protocol + "://" + comm.server +"/my-posts",
            headers: comm.getHeader(),
        }).then((response) => {
            if(response.status == 200){
                this.posts = response.data.collection;
            }
        });
    }
}
</script>