// Utilities
import { defineStore } from "pinia"
import { User } from "src/schema/user"
import { postLogin, postLogout } from "src/api/v1/auth"
import { postMy } from "src/api/v1/main"

interface State {
  status: {
    isFirstRequest: boolean;
  };
  user: User | null;
}

export const useAuthStore = defineStore("auth", {
  state: (): State => ({
    status: {
      isFirstRequest: true,
    },

    user: null,
  }),

  getters: {
    isFirstRequest: (state) => state.status.isFirstRequest,
    isLoggedIn: (state) => state.user !== null,
    isNormalRole: (state) =>
      ["NORMAL", "ADMIN", "SUPER_ADMIN"].includes(state.user?.role),
    isAdminRole: (state) => ["ADMIN", "SUPER_ADMIN"].includes(state.user?.role),
    isSuperAdminRole: (state) => ["SUPER_ADMIN"].includes(state.user?.role),
  },

  actions: {
    loadAccountInfo() {
      return postMy()
        .then((response) => {
          this.user = response.value
        })
        .catch(() => {
          this.user = null
        })
        .finally(() => {
          this.status.isFirstRequest = false
        });
    },

    login(email: string, password: string) {
      this.user = null

      const account = {
        email: email,
        password: password,
      };

      return postLogin(account).then((response) => {
        this.user = response.value
      });
    },

    logout() {
      return postLogout().then(() => {
        this.user = null
      });
    },
  },
});
