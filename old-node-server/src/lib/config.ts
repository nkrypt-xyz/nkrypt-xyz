type Config = {
  webServer: {
    http: {
      enabled: boolean;
      port: number;
    };
    https: {
      enabled: boolean;
      port: number;
      keyFilePath: string;
      certFilePath: string;
      caBundleFilePath: string;
    };
    contextPath: string;
  };
  database: {
    dir: string;
  };
  lockProvider: {
    dir: string;
  };
  blobStorage: {
    dir: string;
    maxFileSizeBytes: number;
  };
};

export { Config };
