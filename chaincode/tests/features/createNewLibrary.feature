Feature: Create New Library
    In order to create a new library
    As an API client
    I want to make a request with the name of the desired library

    Scenario: Create a new library
        # The following statement will be used by all scenarios on this feature
        Given there is a running "" test network
        When I make a "POST" request to "/api/invoke/createNewLibrary" on port 880 with:
            """
            {
                "name": "Maria's Library"
            }
            """
        Then the response code should be 200
        And the response should have:
            """
            {
                "@key":         "library:3cab201f-9e2b-579d-b7b2-72297ed17f49",
                "@lastTouchBy": "orgMSP",
                "@lastTx":      "createNewLibrary",
                "@assetType":   "library",
                "name":         "Maria's Library"
            }
            """

    Scenario: Try to create a new library with a name that already exists
        Given there is a library with name "John's Library"
        When I make a "POST" request to "/api/invoke/createNewLibrary" on port 880 with:
            """
            {
                "name": "John's Library"
            }
            """
        Then the response code should be 409
