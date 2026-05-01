# Test Fixtures

This directory contains mock JSON responses for the Tailscale v2 API, used by unit tests to avoid calling the real API.

## Purpose

Each `.json` file represents a typical response from a specific Tailscale API endpoint. They are loaded by test helpers (`loadFixture`) and served via `httptest.Server` instances to simulate API behavior.

## How to contribute

1. Pick an empty fixture file (currently all contain `{}`).
2. Fill it with a **valid JSON response** that matches the Tailscale API schema for the corresponding endpoint.
3. Run the related tests to verify the fixture is accepted.

API documentation: <https://tailscale.com/api>

## Important: use only fake data

**NEVER** put real data in these fixtures. This includes:

- Real node IDs, device IDs, or user IDs
- Real API tokens or auth keys
- Real IP addresses from your tailnet
- Real hostnames, email addresses, or domain names

Instead, use clearly fictional values. Examples:

| Field      | Fake value                |
|------------|---------------------------|
| nodeId     | `nTEST1234CNTRL`          |
| hostname   | `mock-server`             |
| IP address | `100.64.0.1`              |
| user email | `testuser@example.com`    |
| API key    | `tskey-auth-FAKE-secret`  |
| tailnet    | `test.example.com`        |

## File inventory

| Fixture file            | API endpoint                              |
|-------------------------|-------------------------------------------|
| `device_list.json`      | `GET /tailnet/{tailnet}/devices`           |
| `device_get.json`       | `GET /device/{deviceId}`                   |
| `device_routes.json`    | `GET /device/{deviceId}/routes`            |
| `acl_get.json`          | `GET /tailnet/{tailnet}/acl`               |
| `acl_validate.json`     | `POST /tailnet/{tailnet}/acl/validate`     |
| `key_list.json`         | `GET /tailnet/{tailnet}/keys`              |
| `key_create.json`       | `POST /tailnet/{tailnet}/keys`             |
| `key_get.json`          | `GET /tailnet/{tailnet}/keys/{keyId}`      |
| `dns_nameservers.json`  | `GET /tailnet/{tailnet}/dns/nameservers`   |
| `dns_preferences.json`  | `GET /tailnet/{tailnet}/dns/preferences`   |
| `dns_searchpaths.json`  | `GET /tailnet/{tailnet}/dns/searchpaths`   |
| `dns_split.json`        | `GET /tailnet/{tailnet}/dns/split-dns`     |
| `dns_config.json`       | `GET /tailnet/{tailnet}/dns`               |
| `user_list.json`        | `GET /tailnet/{tailnet}/users`             |
| `user_get.json`         | `GET /users/{userId}`                      |
| `settings_get.json`     | `GET /tailnet/{tailnet}/settings`          |
| `webhook_list.json`     | `GET /tailnet/{tailnet}/webhooks`          |
| `contact_get.json`      | `GET /tailnet/{tailnet}/contacts`          |
