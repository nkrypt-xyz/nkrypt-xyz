/* eslint-disable no-undef */

import { validators } from "../../dist/validators.js";
import { GlobalPermission } from "../../dist/constant/global-permission.js";

import Joi from "joi";

export const directorySchema = Joi.object().optional().keys({
  _id: Joi.string().required(),
  name: Joi.string().min(4).max(32).required(),
  bucketId: Joi.string().min(1).max(64).required(),
  parentDirectoryId: Joi.string().min(1).max(64).allow(null).required(),
  encryptedMetaData: Joi.string().min(1).max(2048).allow(null).required(),
  metaData: Joi.object().required(),
  createdAt: validators.dateRequired,
  updatedAt: validators.dateRequired,
  createdByUserIdentifier: Joi.string().required(),
});

export const fileSchema = Joi.object()
  .keys({
    _id: Joi.string().required(),
    name: Joi.string().min(4).max(32).required(),
    bucketId: Joi.string().min(1).max(64).required(),
    parentDirectoryId: Joi.string().min(1).max(64).allow(null).required(),
    encryptedMetaData: Joi.string().min(1).max(2048).allow(null).required(),
    metaData: Joi.object().required(),
    sizeAfterEncryptionBytes: Joi.number().required(),
    createdAt: validators.dateRequired,
    updatedAt: validators.dateRequired,
    contentUpdatedAt: validators.dateRequired,
    createdByUserIdentifier: Joi.string().required(),
  })
  .optional();

export const bucketListSchema = Joi.array()
  .required()
  .items(
    Joi.object().keys({
      _id: Joi.string().required(),
      name: Joi.string().required(),
      cryptSpec: Joi.string().required(),
      cryptData: Joi.string().required(),
      metaData: Joi.object().required(),
      bucketAuthorizations: Joi.array()
        .required()
        .items(
          Joi.object().required().keys({
            userId: Joi.string().required(),
            notes: Joi.string().required(),
            permissions: validators.allBucketPermissions,
          })
        ),
      rootDirectoryId: Joi.string().required(),
      createdByUserIdentifier: Joi.string().required(),
      createdAt: validators.dateRequired,
      updatedAt: validators.dateRequired
    })
  );

export const userAssertion = {
  hasError: validators.hasErrorFalsy,
  apiKey: validators.apiKey,
  user: Joi.object().required().keys({
    _id: validators.id,
    userName: validators.userName,
    displayName: validators.displayName,
    globalPermissions: (() => {
      let keys = {};
      Object.keys(GlobalPermission).forEach(permission => {
        keys[permission] = Joi.boolean().required()
      });
      return Joi.object().required(keys);
    })(),
  }),
  session: Joi.object().required().keys({
    _id: validators.id,
  }),
};

export const errorOfCode = (code) => {
  return Joi.object()
    .required()
    .keys({
      code: Joi.string().required().valid(code),
      message: Joi.string().required(),
      details: Joi.object().required(),
    });
};

export const userListSchema = Joi.array()
  .items(
    Joi.object().keys({
      _id: validators.id,
      userName: validators.userName,
      displayName: validators.displayName,
      isBanned: Joi.boolean().required()
    })
  )
  .required();

export const userListWithPermissionsSchema = Joi.array()
  .items(
    Joi.object().keys({
      _id: validators.id,
      userName: validators.userName,
      displayName: validators.displayName,
      isBanned: Joi.boolean().required(),
      globalPermissions: validators.allGlobalPermissions
    })
  )
  .required();

export const sessionListSchema = Joi.array()
  .items(
    Joi.object().keys({
      isCurrentSession: Joi.boolean().required(),
      hasExpired: Joi.boolean().required(),
      expireReason: Joi.string().allow(null).required(),
      createdAt: Joi.number().required(),
      expiredAt: Joi.number().required().allow(null)
    })
  )
  .required();