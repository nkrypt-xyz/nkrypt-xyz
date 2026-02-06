enum GlobalPermission {
  MANAGE_ALL_USER = "MANAGE_ALL_USER",
  CREATE_USER = "CREATE_USER",
  CREATE_BUCKET = "CREATE_BUCKET",
}

const getDefaultGlobalPermissionsForNewStandardUser = () => {
  return {
    MANAGE_ALL_USER: false,
    CREATE_USER: false,
    CREATE_BUCKET: true
  }
}

export { GlobalPermission, getDefaultGlobalPermissionsForNewStandardUser };
