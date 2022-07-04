Feature: Get Books By Author
    In order to get all the books by an author
    As an API client
    I want to make a request to the getBooksByAuthor transaction
    And receive the appropriate books

    Scenario: Request an author with multiple books
        Given there is a running "" test network
        And there are 3 books with prefix "book" by author "Jack"
        When I make a "GET" request to "/api/query/getBooksByAuthor" on port 880 with:
            """
            {
                "authorName": "Jack"
            }
            """
        Then the response code should be 200
        And the "result" field should have size 3

    Scenario: Request an author with no books
        Given there is a running "" test network
        When I make a "GET" request to "/api/query/getBooksByAuthor" on port 880 with:
            """
            {
                "authorName": "Mary"
            }
            """
        Then the response code should be 200
        And the "result" field should have size 0

    Scenario: Request an author with 2 books while there are other authors with more books
        Given there is a running "" test network
        Given there are 1 books with prefix "fantasy" by author "Missy"
        Given there are 2 books with prefix "cook" by author "John"
        When I make a "GET" request to "/api/query/getBooksByAuthor" on port 880 with:
            """
            {
                "authorName": "John"
            }
            """
        Then the response code should be 200
        And the "result" field should have size 2