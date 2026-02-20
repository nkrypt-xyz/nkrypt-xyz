export async function callPostJsonApi<T = any>(serverUrl: string, apiKey: string, endpoint: string, data: any): Promise<T> {
  const url = `${serverUrl}${endpoint}`;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (apiKey) {
    headers["Authorization"] = `Bearer ${apiKey}`;
  }

  const response = await fetch(url, {
    method: "POST",
    headers,
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => null);
    
    // Create a structured error object to pass to error service
    const errorObj: any = new Error();
    errorObj.response = {
      status: response.status,
      data: errorData,
    };
    
    throw errorObj;
  }

  return await response.json();
}
