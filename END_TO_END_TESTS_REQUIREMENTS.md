# Production-Ready End-to-End (E2E) Testing Strategy for Rule Engine

## 1. Overview & Goals

This document outlines the requirements and implementation plan for creating a robust, production-ready End-to-End (E2E) testing suite for the Rule Engine.

The primary goals of this test suite are:
- **Verify Real-World Scenarios:** Ensure the entire system works as expected from an end-user's perspective, simulating complete workflows.
- **Build Confidence in Releases:** Create a reliable gate that prevents regressions and ensures new features are integrated correctly.
- **Ensure Component Integration:** Validate that all microservices and backing services (database, cache, external APIs) work together seamlessly.

## 2. Core Principles

- **Isolation:** Tests must be completely isolated from each other and from any external, non-deterministic dependencies.
- **Realism:** The test environment should mirror the production stack as closely as possible, using the same Docker containers and infrastructure components.
- **Automation:** The entire test lifecycle—spinning up infrastructure, running tests, and tearing down—must be fully automated.
- **Maintainability:** Test code should be clean, well-structured, and easy to understand and extend.

## 3. Technology Stack

To achieve the principles above, the E2E test suite will use the following technologies:

| Component                 | Technology                                                              | Purpose                                                                                             |
| ------------------------- | ----------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| **Test Runner**           | **Go (standard `testing` package)**                                     | To write and execute the tests.                                                                     |
| **Infrastructure Mgmt.**  | **Testcontainers for Go**                                               | To programmatically manage the lifecycle of all required services (Rule Engine, DB, etc.) in Docker. |
| **API Interaction**       | **Auto-generated Go Client** (`/clients/go`)                            | To interact with the Rule Engine API in a type-safe manner, exactly as a real client would.         |
| **External API Mocking**  | **Hoverfly (`hoverfly/hoverfly`)**                                      | To simulate external APIs (e.g., webhooks, external data sources) with precision and reliability.   |

## 4. Test Environment Setup with Testcontainers

The E2E tests will run in a completely ephemeral environment orchestrated by **Testcontainers**.

### 4.1. Managed Containers

The test setup will programmatically define and launch the following services for each test run:
1.  **PostgreSQL:** The primary database for storing rules, actions, and triggers.
2.  **Redis:** The caching layer and message queue.
3.  **Hoverfly:** The external API simulator. The Rule Engine will be configured to point to the Hoverfly container for all outbound HTTP requests.
4.  **Rule Engine Service:** The application itself, built from the current source code. It will be configured to connect to the PostgreSQL, Redis, and Hoverfly containers.

### 4.2. Network Configuration

All containers will be attached to a single Testcontainers network, allowing them to communicate with each other using container names as hostnames. The Rule Engine's configuration (e.g., via environment variables) will be dynamically updated with the connection strings for the other containers.

## 5. Test Structure

Tests should be located in the `tests/e2e/` directory (to be created).

### 5.1. File and Test Case Structure

- Use a `main_test.go` to set up the Testcontainers environment once for the entire test suite for efficiency.
- Organize tests by feature or workflow (e.g., `rule_workflow_test.go`, `action_execution_test.go`).
- Use table-driven tests to cover multiple scenarios with different inputs and expected outcomes concisely.

### 5.2. Test Case Logic (Arrange-Act-Assert)

Each test should follow the Arrange-Act-Assert pattern:

1.  **Arrange:**
    - **Load Mocks:** Configure Hoverfly with the specific API simulation required for the test case (e.g., a JSON file defining expected requests and responses).
    - **Instantiate API Client:** Create an instance of the auto-generated Go client from the `/clients/go` directory.
    - **Prepare System:** Use the API client to create any prerequisite resources (e.g., create a base rule that another rule will interact with).

2.  **Act:**
    - Perform the primary action of the test using the API client. This is typically a single, focused operation like creating a new rule, sending an event that should trigger a rule, or updating a configuration.

3.  **Assert:**
    - **API Verification:** Use the API client to fetch resources and assert that their state is correct.
    - **Mock Verification:** Use the Hoverfly API to verify that the expected outbound calls were made by the Rule Engine. For example, assert that a specific webhook endpoint was called exactly once with the correct payload.
    - **Database Verification (Optional):** In rare cases, connect directly to the test database to assert that a specific low-level state change occurred that is not verifiable through the API.

## 6. API Interaction

**Crucially, all interactions with the Rule Engine API MUST be performed through the auto-generated Go client located in `/clients/go`.** This ensures that the tests are a true representation of how a client integrates with the system and validates the API contract, including any potential breaking changes in the generated client itself.

## 7. Mocking with Hoverfly

Hoverfly is essential for isolating our tests from the internet and ensuring deterministic behavior.

- **Simulations:** For each distinct external API interaction, a `simulation.json` file will be created. These files define the expected request from the Rule Engine and the mock response Hoverfly should return.
- **Dynamic Loading:** Before each test (or group of tests), the relevant simulation file will be loaded into the Hoverfly container via its API.
- **Example Scenario (Webhook Action):**
    1.  A test needs to verify that a rule correctly triggers a webhook.
    2.  A `simulation-webhook-success.json` is created. It specifies that a `POST` request to `https://example.com/webhook` with a specific JSON body should be expected. Hoverfly is configured to return a `200 OK` response.
    3.  The test runs, triggering the rule.
    4.  The test asserts that the Rule Engine considers the action successful.
    5.  The test uses the Hoverfly verification API to confirm that the webhook endpoint was indeed called as expected.

## 8. Key Scenarios to Implement

The test suite should cover the following critical user journeys:

- **Full Happy Path:**
    - Create a Rule with a Trigger and an Action (e.g., call a webhook).
    - Send an event that matches the Trigger's condition.
    - Verify the Action was executed (i.e., Hoverfly received the webhook call).
    - Clean up the created resources.
- **Complex Rule Logic:**
    - Test rules with multiple `AND`/`OR` conditions.
    - Test rules with different data types and operators (`=`, `>`, `<`, `CONTAINS`, etc.).
- **Action Execution Failures:**
    - Configure Hoverfly to return a `500 Internal Server Error` for a webhook call.
    - Verify that the Rule Engine correctly marks the action execution as failed and performs any configured retry logic.
- **Invalid Input:**
    - Attempt to create rules/triggers/actions with invalid data (e.g., missing fields, incorrect data types).
    - Assert that the API returns the expected `4xx` error codes and descriptive error messages.
- **Rule Deletion and Updates:**
    - Create a rule, update it, and verify the new behavior.
    - Create a rule, delete it, and verify it no longer triggers.

By following these guidelines, we will build a powerful and reliable E2E testing suite that provides maximum confidence in the correctness and stability of the Rule Engine.
