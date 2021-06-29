<template>
    <post-type-tabs-slider :posts="posts"/>
</template>

<script>
import PostTypeTabsSlider from '../components/PostTypeTabsSlider.vue'
export default {
  components: { PostTypeTabsSlider },
    data(){
        return{
            posts:[]
        }
    },
    created(){
        //TODO: get my campaigns
    }
}
</script>