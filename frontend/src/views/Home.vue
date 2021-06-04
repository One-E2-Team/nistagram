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

export default {
  name: 'Home',
  components: {
    Welcome
  }
}
</script>
