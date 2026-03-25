# Terminal_config_API

# Mocking and testing
We can use mocks to simulate our functions in an isolated, controlled way. It helps to separate your code, and what you are testing, into individual components, allowing you to identify what is and isn't working. We can use them to easily simulate edge cases and errors, such as 403 Unauthorized, which would be difficult to replicate in another format.
Useful resources:
https://www.freecodecamp.org/news/unit-testing-in-go-a-beginners-guide/
https://dev.to/neelp03/essential-unit-testing-for-go-apis-build-code-with-confidence-ne3
``` 
    go install github.com/vektra/mockery/v2@latest
    go mod tidy
```

