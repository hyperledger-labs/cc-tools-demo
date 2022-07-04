Feature: Update Book Tentant
    In order to update book tentant
    As an API client
    I want to make a request

    Scenario: Update Book With A Existing Tentant 
        # The first 3 statements will be used by all scenarios on this feature
        Given there is a running "" test network
        And I make a "POST" request to "/api/invoke/createAsset" on port 880 with:
            """
            {
                "asset": [{
                        "@assetType": "book",
                        "title":      "Meu Nome é Maria",
                        "author":     "Maria Viana"
                    }]
            }
            """
        And I make a "POST" request to "/api/invoke/createAsset" on port 880 with:
            """
           {
                "asset": [{
                    "@assetType": "person",
                    "name": "Maria",
                    "id": "31820792048"
                }]
	        }
            """
        When I make a "PUT" request to "/api/invoke/updateBookTenant" on port 880 with:
            """
            {
                "book": {
                    "@assetType": "book",
                    "@key": "book:a36a2920-c405-51c3-b584-dcd758338cb5"
		        },
                "tenant": {
                    "@assetType": "person",
                    "@key": "person:47061146-c642-51a1-844a-bf0b17cb5e19"
                }
            }
            """
        Then the response code should be 200
        And the response should have:
            """
            {
                "@key": "book:a36a2920-c405-51c3-b584-dcd758338cb5",
                "@lastTouchBy": "orgMSP",
                "@lastTx": "updateBookTenant",
                "currentTenant": {
			        "@assetType": "person",
			        "@key": "person:47061146-c642-51a1-844a-bf0b17cb5e19"
		        }
            }
            """

    Scenario: Update Book With A Not Existing Tentant
        Given there is a running "" test network
        When I make a "PUT" request to "/api/invoke/updateBookTenant" on port 880 with:
            """
            {
                "book": {
                    "@assetType": "book",
                    "@key": "book:a36a2920-c405-51c3-b584-dcd758338cb5"
		        },
                "tenant": {
                    "@assetType": "person",
                    "@key": "person:56891146-c6866-51a1-844a-bf0b17cb5e19"
                }
            }
            """
        Then the response code should be 404