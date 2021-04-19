Feature: Console

    Scenario: Captured Output
        When I write to console:
        """
        Hello World!

        I'm John
        """

        Then console output is:
        """
        Hello World!

        I'm John
        """
