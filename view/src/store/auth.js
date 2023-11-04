// Utilities
import { defineStore } from "pinia"
import axios from "@/modules/axios-wrapper"

export const useAuthStore = defineStore("auth", {
  state: () => ({
    status: {
      isLoggedIn: false,
    },
    user: {
      email: String,
    },
  }),

  getters: {
    isLoggedIn: (state) => state.status.isLoggedIn,
  },

  actions: {
    loadAccountInfo() {
      return axios.post("/api/v1/whoami").then((response) => {
        this.user.email = response.data.email
        this.status.isLoggedIn = true
      }).catch((error) => {
        this.user.email = ""
        this.status.isLoggedIn = false
      })
    },

    login(email, password) {
      this.status.isLoggedIn = false
      this.user.email = ""

      const account = {
        email: email,
        password: password,
      }

      return axios.post("/api/v1/auth/login", account).then((response) => {
        this.user.email = response.data.email
        this.status.isLoggedIn = true
      })
    },

    logout() {
      return axios.post("/api/v1/auth/logout").then((response) => {
        this.status.isLoggedIn = false
        this.user.email = ""
      })
    },
  },
})
