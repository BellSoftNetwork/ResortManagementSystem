// Utilities
import { defineStore } from "pinia"
import { api } from "boot/axios"

export const useAuthStore = defineStore("auth", {
  state: () => ({
    status: {
      isFirstRequest: true,
      isLoggedIn: false,
    },
    user: {
      email: String,
      name: String,
      role: String,
      profileImageUrl: String,
      createdAt: String,
    },
  }),

  getters: {
    isFirstRequest: (state) => state.status.isFirstRequest,
    isLoggedIn: (state) => state.status.isLoggedIn,
    isNormalRole: (state) => ["NORMAL", "ADMIN", "SUPER_ADMIN"].includes(state.user.role),
    isAdminRole: (state) => ["ADMIN", "SUPER_ADMIN"].includes(state.user.role),
    isSuperAdminRole: (state) => ["SUPER_ADMIN"].includes(state.user.role),
  },

  actions: {
    loadAccountInfo() {
      return api.post("/api/v1/whoami").then((response) => {
        const responseValue = response.data.value

        this.user.email = responseValue.email
        this.user.name = responseValue.name
        this.user.role = responseValue.role
        this.user.profileImageUrl = responseValue.profileImageUrl
        this.user.createdAt = responseValue.createdAt
        this.status.isLoggedIn = true
      }).catch(() => {
        this.status.isLoggedIn = false
        this.user.email = ""
        this.user.name = ""
        this.user.role = ""
        this.user.profileImageUrl = "https://gravatar.com/avatar"
        this.user.createdAt = ""
      }).finally(() => {
        this.status.isFirstRequest = false
      })
    },

    login(email, password) {
      this.status.isLoggedIn = false
      this.user.email = ""
      this.user.name = ""
      this.user.role = ""
      this.user.profileImageUrl = "https://gravatar.com/avatar"
      this.user.createdAt = ""

      const account = {
        email: email,
        password: password,
      }

      return api.post("/api/v1/auth/login", account).then((response) => {
        const responseValue = response.data.value

        this.user.email = responseValue.email
        this.user.name = responseValue.name
        this.user.role = responseValue.role
        this.user.profileImageUrl = responseValue.profileImageUrl
        this.user.createdAt = responseValue.createdAt
        this.status.isLoggedIn = true
      })
    },

    logout() {
      return api.post("/api/v1/auth/logout").then(() => {
        this.status.isLoggedIn = false
        this.user.email = ""
        this.user.name = ""
        this.user.role = ""
        this.user.profileImageUrl = "https://gravatar.com/avatar"
        this.user.createdAt = ""
      })
    },
  },
})
