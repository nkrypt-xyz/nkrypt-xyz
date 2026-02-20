export type User = {
  userName: string;
  displayName: string;
  userId: string;
  globalPermissions: Record<string, boolean>;
};

export type Session = {
  serverUrl: string;
  apiKey: string;
};

export type Settings = {
  uploadMechanism: string;
  downloadMechanism: string;
  plainTextEditorNoRestrictions: boolean;
  darkMode: boolean | null;
};

export type Bucket = {
  _id: string;
  name: string;
  cryptSpec: string;
  cryptData: string;
  rootDirectoryId: string;
  metaData: Record<string, any>;
  encryptedMetaData: string;
  bucketAuthorizations: BucketAuthorization[];
};

export type BucketAuthorization = {
  userId: string;
  userName: string;
  isInherited: boolean;
};

export type Directory = {
  _id: string;
  name: string;
  bucketId: string;
  parentDirectoryId: string | null;
  metaData: Record<string, any>;
  encryptedMetaData: string;
};

export type File = {
  _id: string;
  name: string;
  bucketId: string;
  parentDirectoryId: string;
  blobId: string;
  metaData: Record<string, any>;
  encryptedMetaData: string;
};

export type EntityStackItem = {
  bucket: Bucket;
  directory: Directory;
};

export type ClipboardData = {
  action: string;
  isDirectory: boolean;
  entity: Directory | File;
};
