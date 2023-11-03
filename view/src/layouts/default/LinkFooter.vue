<template>
  <v-footer
    color="grey lighten-2"
    padless
  >
    <v-row
      justify="center"
      no-gutters
    >
      <v-btn
        v-for="link in footerLinks"
        :key="link.url"
        :href="link.url"
        target="_blank"
        variant="text"
        rounded
        class="my-2"
      >
        {{ link.title }}
      </v-btn>
      <v-col
        class="lighten-2 py-4 text-center"
        cols="12"
      >
        {{ new Date().getFullYear() }} — <strong>RMS (Resort Management System)</strong>
        <br />
        Powered by Bell Soft Network
      </v-col>
    </v-row>
  </v-footer>
</template>

<script>
export default {
  data: () => ({
    footerLinks: [
      { title: "BSN GitLab", url: "https://gitlab.bellsoft.net" },
    ],
  }),
}
</script>
