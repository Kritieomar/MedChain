// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract InteractionTracking {
    struct Interaction {
        string recordId;
        address patient;
        string action; // "Accepted" or "Rejected"
    }

    Interaction[] public interactions;

    event InteractionLogged(string recordId, address patient, string action);

    function logInteraction(string memory recordId, string memory action) public {
        interactions.push(Interaction(recordId, msg.sender, action));
        emit InteractionLogged(recordId, msg.sender, action);
    }

    function getInteractions() public view returns (Interaction[] memory) {
        return interactions;
    }
}
