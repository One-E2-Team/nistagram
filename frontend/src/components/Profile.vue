<template>
  <div></div>
</template>

<script>
import axios from 'axios'
import * as comm from '../configuration/communication.js'
export default {
name: "Profile",

data: () => ({
    profile: {}
}),
methods: {

},
created(){
    let username = comm.getUrlVars()['username'];
    axios({
        method: "get",
        url: 'http://' + comm.server + '/api/profile/get/' + username,
    }).then(response => {
        if(response.status==200){
            this.profile = response.data;
        }
    }).catch(reason => {
        console.log(reason);
        alert("User with username '" + username + "' does not exist!");
    });
}

}
</script>

<style>

</style>