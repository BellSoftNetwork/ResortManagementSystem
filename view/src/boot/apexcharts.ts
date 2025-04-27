import { boot } from "quasar/wrappers";
import VueApexCharts from "vue3-apexcharts";

export default boot(({ app }) => {
  // eslint-disable-next-line vue/multi-word-component-names
  app.component("apexchart", VueApexCharts);
});
