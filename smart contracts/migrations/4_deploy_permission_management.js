
const PermissionManagement = artifacts.require("PermissionManagement");

module.exports = function (deployer) {

  deployer.deploy(PermissionManagement);
};
