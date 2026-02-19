var GlobalPermission;
(function (GlobalPermission) {
    GlobalPermission["MANAGE_ALL_USER"] = "MANAGE_ALL_USER";
    GlobalPermission["CREATE_USER"] = "CREATE_USER";
    GlobalPermission["CREATE_BUCKET"] = "CREATE_BUCKET";
})(GlobalPermission || (GlobalPermission = {}));
const getDefaultGlobalPermissionsForNewStandardUser = () => {
    return {
        MANAGE_ALL_USER: false,
        CREATE_USER: false,
        CREATE_BUCKET: true
    };
};
export { GlobalPermission, getDefaultGlobalPermissionsForNewStandardUser };
