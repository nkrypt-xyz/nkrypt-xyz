import { useSessionStore } from "stores/session";
import { callPostJsonApi } from "utils/api-utils";

export async function callMetricsGetSummaryApi(data: Record<string, never>) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/metrics/get-summary", data);
}
