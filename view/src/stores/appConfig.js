import { defineStore } from "pinia"
import { api } from "boot/axios"

export const useAppConfigStore = defineStore("appConfig", {
  state: () => ({
    status: {
      isLoading: false,
      isLoaded: false,
    },
    config: {
      isAvailableRegistration: false,
    },
  }),

  actions: {
    loadAppConfig(force = false) {
      if (this.status.isLoaded && !force)
        return Promise.resolve()

      this.status.isLoading = true
      this.status.isLoaded = false

      return api.get("/api/v1/config").then((response) => {
        this.config.isAvailableRegistration = response.data.value.isAvailableRegistration
        this.status.isLoaded = true

        return response.data.value
      }).catch((error) => {
        this.config.isAvailableRegistration = false

        return Promise.reject(error)
      }).finally(() => {
        this.status.isLoading = false
      })
    },
  },
})
