import { boot } from "quasar/wrappers";
import draggable from "vuedraggable";

// "async" is optional;
// more info on params: https://v2.quasar.dev/quasar-cli/boot-files
export default boot(({ app }) => {
  // eslint-disable-next-line vue/multi-word-component-names
  app.component("draggable", draggable);
});
