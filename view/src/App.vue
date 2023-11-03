<template>
  <header>
    <div class="wrapper">
        <HelloWorld msg="You did it!"/>
    </div>
  </header>

  <main>
      <TheWelcome/>
      <h1>API ENV: {{ env.applicationFullName }} ({{ env.version }})</h1>
  </main>
</template>

<script setup>
import HelloWorld from './components/HelloWorld.vue'
import TheWelcome from './components/TheWelcome.vue'
</script>

<script>
export default {
    data: () => ({
        env: {
            applicationFullName: "",
            applicationShortName: "",
            version: "",
        }
    }),
    methods: {
        initEnv() {
            fetch("/api/v1/env")
                .then(response => response.json())
                .then(data => this.env = data)
        }
    },
    created() {
        this.initEnv()
  }
}
</script>
