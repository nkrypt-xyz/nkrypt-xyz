import { dialogService } from "./dialog-service";

export interface ProcessedError {
  message: string;
  code?: string;
  details?: any;
}

export const errorService = {
  processError(error: any): ProcessedError {
    if (typeof error === "string") {
      return { message: error };
    }

    if (error.response) {
      const status = error.response.status;
      const data = error.response.data;

      // Handle authentication errors specially
      if (status === 401 || status === 403) {
        // But still check if server provided a more specific message
        if (data?.error?.message) {
          return { 
            message: data.error.message, 
            code: data.error.code || "AUTH_ERROR",
            details: data.error.details 
          };
        }
        return { message: "Authentication failed. Please login again.", code: "AUTH_ERROR" };
      }

      // Handle 404 errors
      if (status === 404) {
        if (data?.error?.message) {
          return { 
            message: data.error.message, 
            code: data.error.code || "NOT_FOUND",
            details: data.error.details 
          };
        }
        return { message: "Resource not found.", code: "NOT_FOUND" };
      }

      // Check for nkrypt backend error format: { hasError: true, error: { code, message, details } }
      if (data?.error?.message) {
        return { 
          message: data.error.message, 
          code: data.error.code,
          details: data.error.details 
        };
      }

      // Fallback to old format (if any API still uses it)
      if (data?.message) {
        return { message: data.message, code: data.code };
      }

      return { message: `Server error: ${status}`, code: "SERVER_ERROR" };
    }

    if (error.message) {
      return { message: error.message };
    }

    return { message: "An unexpected error occurred.", code: "UNKNOWN_ERROR" };
  },

  async handleUnexpectedError(error: any): Promise<void> {
    try {
      const processed = this.processError(error);
      console.error("Error:", processed);
      await dialogService.alert("Error", processed.message);
    } catch (innerError) {
      console.error("Error in error handler:", innerError);
    }
  },

  async wrapAsyncHandler<T>(fn: () => Promise<T>): Promise<T | null> {
    try {
      return await fn();
    } catch (error) {
      await this.handleUnexpectedError(error);
      return null;
    }
  },
};
