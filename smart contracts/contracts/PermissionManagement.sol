// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract PermissionManagement {
    struct Permission {
        string recordId;
        address grantedTo; // Doctor's wallet address
        bool isGranted;
    }

    mapping(string => Permission[]) public permissions; // Map recordId to a list of permissions

    event PermissionGranted(string recordId, address grantedTo);
    event PermissionRevoked(string recordId, address grantedTo);

    function grantPermission(string memory recordId, address doctor) public {
        permissions[recordId].push(Permission(recordId, doctor, true));
        emit PermissionGranted(recordId, doctor);
    }

    function revokePermission(string memory recordId, address doctor) public {
        permissions[recordId].push(Permission(recordId, doctor, false));
        emit PermissionRevoked(recordId, doctor);
    }

    function getPermissions(string memory recordId) public view returns (Permission[] memory) {
        return permissions[recordId];
    }
}
