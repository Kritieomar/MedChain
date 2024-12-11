// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract UserManagement {
    struct User {
        string userId;
        string userType; // "Doctor" or "Patient"
        address wallet;
    }

    mapping(string => User) public users; // Mapping of userId to User

    event UserCreated(string userId, string userType, address wallet);

    function createUser(string memory userId, string memory userType) public {
        require(bytes(users[userId].userId).length == 0, "User already exists");
        require(
            keccak256(bytes(userType)) == keccak256(bytes("Doctor")) ||
            keccak256(bytes(userType)) == keccak256(bytes("Patient")),
            "Invalid user type"
        );

        users[userId] = User(userId, userType, msg.sender);
        emit UserCreated(userId, userType, msg.sender);
    }

    function getUser(string memory userId) public view returns (User memory) {
        require(bytes(users[userId].userId).length > 0, "User does not exist");
        return users[userId];
    }
}
