Feature: Get Number Of Books From Library
    In order to create the number of books from library
    As an API client
    I want to make a request

    Scenario: Query Get Number Of Books From Library that Exists
        Given there is a running "" test network
        And I make a "POST" request to "/api/invoke/createAsset" on port 880 with:
            """
            {
                "asset": [
                    {
                        "@assetType": "book",
                        "title":      "Meu Nome é Maria",
                        "author":     "Maria Viana"
                    }
                ]
            }
            """
        And I make a "POST" request to "/api/invoke/createAsset" on port 880 with:
            """
           {
                "asset": [{
                    "@assetType": "library",
                    "name": "Maria's Library",
                    "books": [
                        {
                            "@assetType": "book",
                            "@key": "book:a36a2920-c405-51c3-b584-dcd758338cb5"
                        }
                    ]
                }]
	        }
            """
        When I make a "GET" request to "/api/query/getNumberOfBooksFromLibrary" on port 880 with:
            """
            {
                "library": {
                    "@key": "library:3cab201f-9e2b-579d-b7b2-72297ed17f49",
                    "@assetType": "library"
		        }
            }
            """
        Then the response code should be 200
        And the response should have:
            """
            {
                "numberOfBooks": 1.0
            }
            """

    Scenario: Query Get Number Of Books From Library that Does Not Exists
        Given there is a running "" test network
        When I make a "GET" request to "/api/query/getNumberOfBooksFromLibrary" on port 880 with:
            """
            {
                "library": {
                    "@key": "library:5c5b201f-9e4c-579d-b7b2-72297ed17f78",
                    "@assetType": "library"
		        }
            }
            """
        Then the response code should be 400