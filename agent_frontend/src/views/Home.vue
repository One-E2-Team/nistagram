<template>
  <div class="home">
    <Welcome />
  </div>
</template>

<script>
// @ is an alias to /src
import Welcome from '@/components/Welcome.vue'
import * as comm from '../configuration/communication.js'

export default {
  name: 'Home',
  components: {
    Welcome
  },
  mounted(){
    if (this.isAvailable()){
      this.$router.push({name: 'NotFound'})
    }
  },
  methods:{
    isAvailable(){
      return comm.isUserLogged()
    }
  }
}
</script>
