# Code Conventions

## Error Handling

- Document all possible app-errors that a method can return in its comments:
  - Use the format: "App Errors: - ErrTypeName"
  - List each error on a new line
  - This helps callers handle specific errors and enables better IDE support

## Business Rules & Validation

### Single Record Operations

- For single record operations:
  1. Fetch the record from DB
  2. Perform validation in the Domain layer
  3. Update the entire entity in DB if valid

### Bulk Operations

- For bulk operations, consider implementing business rules directly in DB:
  - Pros:
    - Avoids loading large datasets into memory
    - Better performance by eliminating round trips
    - Atomic validation and updates
  - Use Cases:
    - Mass updates/deletes
    - Batch processing
    - Data migrations

## Project Structure

### Domain Layer

- Keep domain models free of infrastructure concerns (like SQL tags)
- Use separate DTO models for API responses in the api/dtos package
- Domain interfaces (like Repository) should be defined in the domain layer

### Workflows

Workflows are to be used for cross-domain write operations. However, use them only when absolutely necessary; each domain should handle its own operations independently whenever possible

- Define workflows in the `workflows` package to coordinate between different domains
- Use domain handlers to encapsulate business logic within each domain
- Ensure transactional integrity
  - Use one of the domain's repository to initiate a transaction
  - Then create command_handlers with that transaction
  - Use the new transactional handlers to perform all the operations
  - Commit or rollback accordingly
  - Document each workflow with clear steps and involved domains for better maintainability
- Example:
  - User registration involving authentication and profile management domains

### Error Management

- Use `AppError` for wrapping all domain errors
- Include location context using `NewAppError(err, "location")`

## Repository Pattern

- Implement `RepositoryTx[TRepo]` interface for transactional support
- Use the WithTx pattern for transaction handling
- Keep SQL-specific code in the infrastructure layer

## Logging

- Extract logger from request, if available otherwise from the `App.Logger`
- Include relevant context fields using the chaining method `WithMetadata`
  - Top layers should add the context fields
  - Lower layers can add the context fields, if they have generated new fields

## Configuration

- Use strongly-typed configuration structs
- Group related settings into specific config sections (ServerConfig, DBConfig, etc.)

## File Naming

- Use one word with no underscores for package names
  - Use same convention for the files named same as package/folder name
- Use snake_case for all other golang files

## Type Safety

- Use strongly-typed enums where possible (e.g., Provider type)
- Leverage interfaces for dependency injection and unit testing
