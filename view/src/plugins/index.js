/**
 * plugins/index.js
 *
 * Automatically included in `./src/main.js`
 */

// Plugins
import vuetify from "./vuetify"
import pinia from "@/store"
import router from "@/router"

// Font Awesome
import { library } from "@fortawesome/fontawesome-svg-core"
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome"
import { fas } from "@fortawesome/free-solid-svg-icons"
import { fab } from "@fortawesome/free-brands-svg-icons"

library.add(fas)
library.add(fab)

export function registerPlugins(app) {
  app.component("font-awesome-icon", FontAwesomeIcon)

  app.use(vuetify)
  app.use(pinia)
  app.use(router)
}
