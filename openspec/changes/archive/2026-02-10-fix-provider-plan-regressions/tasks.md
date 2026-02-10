## 1. Reproduce and Scope

- [x] 1.1 Reproduce provider failures with `terraform plan` in `/Users/damien/git/mockshop-in-cloud-2/terraform/infrastructure/awx/dev`
- [x] 1.2 Identify root causes for write-only field type mismatch and pagination `next` parsing

## 2. Provider Runtime Fixes

- [x] 2.1 Update object data source write-only state assignment to use typed nulls
- [x] 2.2 Update object resource write-only snapshot/state handling to preserve typed values
- [x] 2.3 Keep write-only semantics (do not repopulate values from API reads)

## 3. Client Pagination Fixes

- [x] 3.1 Update `ListAll` pagination follow logic to resolve `next` URLs using reference semantics
- [x] 3.2 Support absolute, relative, and query-only `next` URL forms without path encoding errors

## 4. Regression Coverage

- [x] 4.1 Add provider tests for typed write-only preservation and typed-null fallback
- [x] 4.2 Add client test for query-only pagination link handling
- [x] 4.3 Run `go test ./internal/provider ./internal/client`
- [x] 4.4 Run `make test`

## 5. End-to-End Verification

- [x] 5.1 Build provider binary (`make build`)
- [x] 5.2 Re-run `terraform plan` in `/Users/damien/git/mockshop-in-cloud-2/terraform/infrastructure/awx/dev` and verify provider errors are resolved
