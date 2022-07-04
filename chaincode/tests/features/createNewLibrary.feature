Feature: Create New Library
    In order to create a new library
    As an API client
    I want to make a request with the name of the desired library

    Scenario: Create a new library
        Given there is a running "" test network from scratch
        When I make a "POST" request to "/api/invoke/createNewLibrary" on port 880 with:
            """
            {
                "name": "Elizabeth's Library"
            }
            """
        Then the response code should be 200
        And the response should have:
            """
            {
                "@key":         "library:9cf6726a-a327-568a-baf1-5881393073bf",
                "@lastTouchBy": "orgMSP",
                "@lastTx":      "createNewLibrary",
                "@assetType":   "library",
                "name":         "Elizabeth's Library"
            }
            """

    Scenario: Try to create a new library with a name that already exists
        Given there is a running "" test network
        Given there is a library with name "John's Library"
        When I make a "POST" request to "/api/invoke/createNewLibrary" on port 880 with:
            """
            {
                "name": "John's Library"
            }
            """
        Then the response code should be 409
