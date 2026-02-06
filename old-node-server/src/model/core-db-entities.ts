type User = {
  _id: string | null | undefined;
  displayName: string;
  userName: string;
  password: { hash: string, salt: string };
  globalPermissions: Record<string, boolean>;
  isBanned: boolean;
  createdAt: number;
  updatedAt: number;
};

type Session = {
  _id: string | null | undefined;
  userId: string;
  apiKey: string
  hasExpired: boolean;
  expiredAt: number | null;
  expireReason: string | null
  createdAt: number;
  updatedAt: number;
};

type Bucket = {
  _id: string | null | undefined;
  name: string;
  cryptSpec: string;
  cryptData: string;
  metaData: Record<string, never>;
  bucketAuthorizations: [
    {
      userId: string;
      notes: string;
      permissions: Record<string, boolean>;
    },
  ],
  createdByUserIdentifier: string;
  createdAt: number;
  updatedAt: number;
};

type Directory = {
  _id: string | null | undefined;
  bucketId: string;
  parentDirectoryId: string | null;
  name: string;
  metaData: Record<string, never>;
  encryptedMetaData: string;
  createdByUserIdentifier: string;
  createdAt: number;
  updatedAt: number;
}

type File = {
  _id: string | null | undefined;
  bucketId: string;
  parentDirectoryId: string;
  name: string;
  metaData: Record<string, never>;
  encryptedMetaData: string;
  sizeAfterEncryptionBytes: number,
  createdByUserIdentifier: string;
  createdAt: number;
  updatedAt: number;
  contentUpdatedAt: number;
}

type Blob = {
  _id: string | null | undefined;
  bucketId: string;
  fileId: string;
  cryptoMetaHeaderContent: string;
  startedAt: number;
  finishedAt: number | null;
  status: string;
  createdByUserIdentifier: string;
  createdAt: number;
  updatedAt: number;
}

export { User, Session, File, Bucket, Directory, Blob }