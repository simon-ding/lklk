<template>
  <div id="app">
    <img src="./assets/logo.png">
      <input1></input1>
     <div>
      <el-button @click="startHacking">Start</el-button>
    </div>
  </div>
</template>

<script>
import input1 from "./components/input"
export default {
  methods: {
    startHacking () {
      this.$notify({
        title: 'It works!',
        type: 'success',
        message: 'We\'ve laid the ground work for you. It\'s time for you to build something epic!',
        duration: 5000
      })
    }
  },
  components: {
    'input1': input1
  },
  data() {
    return {
      input: ''
    }
  }

}
</script>

<style>
#app {
  font-family: Helvetica, sans-serif;
  text-align: center;
}
</style>
