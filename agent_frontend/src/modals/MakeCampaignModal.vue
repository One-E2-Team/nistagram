<template>
  <v-row justify="space-around">
    <v-col cols="auto">
      <v-dialog transition="dialog-top-transition" max-width="600">
        <template v-slot:activator="{ on, attrs }">
            <v-btn v-bind="attrs" v-on="on" @click="createDialog()"> Make campaign </v-btn>
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-toolbar color="primary" dark>Campaign</v-toolbar>
            <v-card-text>
              <v-form ref="form" v-model="valid" lazy-validation class="text-center">
                <v-row justify="center">
                  <v-col cols="12" sm="4">
                   <v-menu v-model="dateMenuStart"
                        :close-on-content-click="false"
                        :nudge-right="40"
                        transition="scale-transition"
                        offset-y
                        min-width="auto"
                    >
                      <template v-slot:activator="{ on, attrs }">
                      <v-text-field v-model="start"
                          label="Start"
                          prepend-icon="mdi-calendar"
                          v-bind="attrs"
                          v-on="on"
                      ></v-text-field>
                      </template>
                      <v-date-picker v-model="start" @input="dateMenuStart = false"/>
                    </v-menu>
                  </v-col>
                  <v-col cols="12" sm="4">
                   <v-menu v-model="dateMenuEnd"
                        :close-on-content-click="false"
                        :nudge-right="40"
                        transition="scale-transition"
                        offset-y
                        min-width="auto"
                    >
                      <template v-slot:activator="{ on, attrs }">
                      <v-text-field v-model="end"
                          label="End"
                          prepend-icon="mdi-calendar"
                          v-bind="attrs"
                          v-on="on"
                      ></v-text-field>
                      </template>
                      <v-date-picker v-model="end" @input="dateMenuEnd = false"/>
                    </v-menu>
                  </v-col>
                </v-row>
                <v-row justify="center">
                  <v-col cols="12" sm="6">
                    <v-menu ref="menu" v-model="timeMenu"
                      :close-on-content-click="false"
                      :nudge-right="40"
                      :return-value.sync="selectedTime"
                      transition="scale-transition"
                      offset-y
                    >
                      <template v-slot:activator="{ on, attrs }">
                        <v-text-field
                          v-model="selectedTime"
                          label="Chose time"
                          prepend-icon="mdi-clock-time-four-outline"
                          readonly
                          v-bind="attrs"
                          v-on="on"
                        ></v-text-field>
                      </template>
                      <v-time-picker
                        v-if="timeMenu"
                        v-model="selectedTime"
                        full-width
                        @click:minute="$refs.menu.save(selectedTime)"
                      ></v-time-picker>
                    </v-menu>
                  </v-col>
                  <v-col cols="12" sm="2">
                    <v-btn @click="addTime()" color="success">Add</v-btn>
                  </v-col>
                </v-row>
                <v-row justify="space-around">
                  <v-col cols="12" sm="8" md="8">
                    <v-sheet elevation="17"  height="50" >
                      <v-chip-group mandatory active-class="primary--text">
                        <v-chip v-for="time in timestamps" :key="time" close @click:close="removeTime(time)" clearable>
                          {{ time }}
                        </v-chip>
                      </v-chip-group>
                    </v-sheet>
                  </v-col>
                </v-row>
                 <v-row justify="center">
                    <v-col cols="12" sm="8" >
                      <v-combobox v-model="interests" :items="allInterests" chips
                        clearable label="Target group" multiple solo >
                        <template v-slot:selection="{ attrs, item, select, selected }">
                          <v-chip
                            v-bind="attrs"
                            :input-value="selected"
                            close
                            @click="select"
                            @click:close="removeInterest(item)"
                          >
                            <strong>{{ item }}</strong>&nbsp;
                          </v-chip>
                        </template>
                      </v-combobox>
                    </v-col>
                  </v-row>
                  <v-row justify="center">
                    <v-col cols="12" sm="8" >
                      <v-combobox v-model="influensers" :items="allFollowers" chips
                        clearable label="Influensers" multiple solo >
                        <template v-slot:selection="{ attrs, item, select, selected }">
                          <v-chip
                            v-bind="attrs"
                            :input-value="selected"
                            close
                            @click="select"
                            @click:close="removeInfluenser(item)"
                          >
                            <strong>{{ item }}</strong>&nbsp;
                          </v-chip>
                        </template>
                      </v-combobox>
                    </v-col>
                  </v-row>
              </v-form>
            </v-card-text>
            <v-card-actions class="justify-end">
                <v-btn text @click="confirm()">Confirm</v-btn>
                <v-btn text @click="dialog.value = false">Close</v-btn>
            </v-card-actions>
          </v-card>
        </template>
      </v-dialog>
    </v-col>
  </v-row>
</template>

<script>
import * as validator from '../plugins/validator.js'
export default {
    props:['postId'],
    data() {
      return {
        token: '',
        rules:validator.rules,
        valid: true,

        timeMenu: false,
        dateMenuStart: false,
        dateMenuEnd: false,
        start: '',
        end: '',
        timestamps: [],
        selectedTime: '',
        interests: [],
        allInterests: [],
        influensers: [],
        allFollowers: [],
    }},
    methods: {
      createDialog(){
        //TODO: ucitaj sve interese i smesti ih u allInterests 
        //TODO: ucitaj sve moguce influensere i smesti ih u allFollowers
      },
       confirm(){
           if(!this.$refs.form.validate()){
               return
           }
            let startDate = new Date(this.start)
            let endDate = new Date(this.end)
            let data = {
              postId : this.postId,
              start : startDate.toISOString(),
              end: endDate.toISOString(),
              interests : this.interests,
              timestamps : this.getAllTimestampsAsDate(),
              influencerProfileIds : this.getAllInfluencerProfileIds()
            }
            console.log(data)
            //TODO: posalji zahtev za kreiranje kampanje

       },
       addTime(){
         if (!this.isTimeExists(this.selectedTime))
            this.timestamps.push(this.selectedTime)
       },
       removeTime (item) {
          this.timestamps.splice(this.timestamps.indexOf(item), 1)
          this.timestamps = [...this.timestamps]
        },
        removeInterest(item) {
          this.interests.splice(this.interests.indexOf(item), 1)
          this.interests = [...this.interests]
        },
        removeInfluenser(item) {
          this.influensers.splice(this.influensers.indexOf(item), 1)
          this.influensers = [...this.influensers]
        },
      isTimeExists(time){
        for (let t of this.timestamps){
          if (time == t){
            return true 
          }
        }
        return false
      },
      getAllTimestampsAsDate(){
        let ret = [];
        for(let t of this.timestamps){
          let date = new Date();
          let parts = t.split(":")
          date.setHours(parts[0])
          date.setMinutes(parts[1])
          ret.push(date.toISOString());
        }
        return ret;
      },
      getAllInfluencerProfileIds(){
        //TODO: izvuci id-eve iz profila influensera
        return []
      }
      
    },
}
</script>