/* eslint-disable */

// Forces TS to apply `@quasar/app-vite` augmentations of `quasar` package
// Removing this would break `quasar/wrappers` imports as those typings are declared by `@quasar/app-vite`
/// <reference types="@quasar/app-vite" />

// Removes the need for `QComponent` type in TypeScript files
export * from "quasar";
