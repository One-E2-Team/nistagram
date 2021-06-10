<template>
  <div class="home">
    <Welcome />
    <v-img
          :src="require('../assets/minions.jpg')"
          class="my-3"
          contain
          height="200"
        />
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
  beforeCreate(){
    if (this.isAvailable()){
      this.$router.push({name: 'NotFound'})
    }
  },
  methods:{
    isAvailable(){
      return !comm.isUserLogged()
    }
  }
}
</script>
