// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MedicalRecords {
    struct Record {
        string recordId;
        string ipfsHash;
        string owner; // Patient ID
        address addedBy; // Doctor's wallet address
    }

    mapping(string => Record) public records; // Mapping of recordId to Record

    event RecordAdded(string recordId, string ipfsHash, string owner, address addedBy);

    function addRecord(string memory recordId, string memory ipfsHash, string memory owner) public {
        require(bytes(records[recordId].recordId).length == 0, "Record already exists");

        records[recordId] = Record(recordId, ipfsHash, owner, msg.sender);
        emit RecordAdded(recordId, ipfsHash, owner, msg.sender);
    }

    function getRecord(string memory recordId) public view returns (Record memory) {
        require(bytes(records[recordId].recordId).length > 0, "Record does not exist");
        return records[recordId];
    }
}
