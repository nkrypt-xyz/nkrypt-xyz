export const validators = {
  required: (val: any) => !!val || "This field is required",

  minLength: (min: number) => (val: string) => (val && val.length >= min) || `Must be at least ${min} characters`,

  maxLength: (max: number) => (val: string) => !val || val.length <= max || `Must be at most ${max} characters`,

  email: (val: string) => !val || /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(val) || "Invalid email address",

  url: (val: string) => !val || /^https?:\/\/.+/.test(val) || "Invalid URL",

  username: [(val: string) => !!val || "Username is required", (val: string) => val.length >= 4 || "Username must be at least 4 characters"],

  password: [(val: string) => !!val || "Password is required", (val: string) => val.length >= 8 || "Password must be at least 8 characters"],

  serverUrl: [(val: string) => !!val || "Server URL is required", (val: string) => val.length >= 4 || "Invalid server URL"],
};
