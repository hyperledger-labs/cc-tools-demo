Feature: Create New Library
    In order to create a new library
    As an API client
    I want to make a request

    Background:
        Given there is a running "" test network

    Scenario: Invoke Create New Library transaction
        When I make a "POST" request to "/api/invoke/createNewLibrary" with:
            """
            {
                "name": "Maria's Library",
            }
            """
        Then the response code should be 200
        And the response should match json:
            """
            {
                "@key":         "library:3cab201f-9e2b-579d-b7b2-72297ed17f49",
                "@lastTouchBy": "org3MSP",
                "@lastTx":      "createNewLibrary",
                "@assetType":   "library",
                "name":         "Maria's Library",
            }
            """