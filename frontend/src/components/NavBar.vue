<template>
    <div id="nav">
      <router-link to="/">Home</router-link> |
      <router-link to="/about">About</router-link>

      <v-menu
            bottom
            left
          >
            <template v-slot:activator="{ on, attrs }">
              <v-btn
                dark
                icon
                v-bind="attrs"
                v-on="on"
              >
                <v-icon color="blue">mdi-dots-vertical</v-icon>
              </v-btn>
            </template>

            <v-list>
              <v-list-item
                v-for="(item, i) in options"
                :key="i"
              >
                <v-list-item-title @click="relocate(item.name)">{{ item.title }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
    </div>
</template>

<script>

export default {
    name: "NavBar",

    data(){
        return{
            options: [
                {title: 'Personal settings',
                name: 'PersonalSettings'},
                {title: 'Profile settings',
                name: 'ProfileSettings'}],
        }
    },

    methods:{
        relocate(componentName){
            console.log(componentName)
            this.$router.push({name:componentName})
        }
    }
}
</script>