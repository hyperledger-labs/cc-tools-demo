Feature: Get Books By Author
    In order to get all the books by an author
    As an API client
    I want to make a request to the getBooksByAuthor transaction
    And receive the appropriate books

    Background:
        Given there is a running "" test network

    Scenario: Request an author with multiple books
        Given there are 3 books with prefix "book" by author "Maria da Silva"
        When I make a "GET" request to "/api/query/getBooksByAuthor" on port 980 with:
            """
            {
                "authorName": "Maria da Silva"
            }
            """
        Then the response code should be 200
        # And the response should have size 3
        And the response should match json:
            """
            {
                "metadata": null,
                "result": [
                    {
                        "@assetType":   "book",
                        "@key":         "library:3cab201f-9e2b-579d-b7b2-72297ed17f49",
                        "@lastTouchBy": "org2MSP",
                        "@lastTx":      "createAsset",
                        "author":       "Maria da Silva",
                        "title":        "book1"
                    },
                    {
                        "@assetType":   "book",
                        "@key":         "library:3cab201f-9e2b-579d-b7b2-72297ed17f49",
                        "@lastTouchBy": "org2MSP",
                        "@lastTx":      "createAsset",
                        "author":       "Maria da Silva",
                        "title":        "book2"
                    },
                    {
                        "@assetType":   "book",
                        "@key":         "library:3cab201f-9e2b-579d-b7b2-72297ed17f49",
                        "@lastTouchBy": "org2MSP",
                        "@lastTx":      "createAsset",
                        "author":       "Maria da Silva",
                        "title":        "book3"
                    }
                ]
            }
            """