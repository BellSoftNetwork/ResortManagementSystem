import { defineStore } from "pinia"
import { getServerConfig } from "src/api/v1/main"

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
      if (this.status.isLoaded && !force) return Promise.resolve()

      this.status.isLoading = true
      this.status.isLoaded = false

      return getServerConfig()
        .then((response) => {
          this.config = response.value
          this.status.isLoaded = true

          return response.value
        })
        .catch((error) => {
          this.config.isAvailableRegistration = false

          return Promise.reject(error)
        })
        .finally(() => {
          this.status.isLoading = false
        });
    },
  },
});
