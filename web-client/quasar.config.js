/* eslint-env node */

const { configure } = require("quasar/wrappers");

module.exports = configure(function (/* ctx */) {
  return {
    eslint: {
      exclude: ["src-pwa/**/*", "node_modules/nkrypt-xyz-core-nodejs/**/*"],
      warnings: true,
      errors: true,
    },

    boot: ["dark-mode", "credentials"],

    css: ["app.scss", "dark-theme.scss"],

    extras: [
      "roboto-font",
      "material-icons",
    ],

    build: {
      target: {
        browser: ["es2019", "edge88", "firefox78", "chrome87", "safari13.1"],
        node: "node16",
      },
      sourcemap: true,
      vueRouterMode: "hash",
      extendViteConf(viteConf) {
        Object.assign(viteConf.resolve.alias, {
          src: "/src",
          components: "/src/components",
          layouts: "/src/layouts",
          pages: "/src/pages",
          assets: "/src/assets",
          boot: "/src/boot",
          stores: "/src/stores",
          services: "/src/services",
          utils: "/src/utils",
          lib: "/src/lib",
          integration: "/src/integration",
          constants: "/src/constants",
          models: "/src/models",
        });
      },
    },

    devServer: {
      open: true,
      port: 9042,
    },

    framework: {
      config: {
        dark: false,
      },
      plugins: ["Notify", "Dialog"],
    },

    animations: [],
  };
});
