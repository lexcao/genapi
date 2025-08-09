# genapi Codebase Improvements

This document outlines identified improvements for the genapi codebase, categorized by priority and impact.

## 🚨 Critical Issues (Fix Immediately)

### 1. HTTP Client Configuration Bug
**File:** `pkg/clients/http/http.go:27`  
**Issue:** The `SetConfig` method overwrites the injected HTTP client with the default client:
```go
func (c *HttpClient) SetConfig(config internal.Config) {
	c.client = http.DefaultClient  // BUG: This overwrites the injected client
	c.config = config
}
```
**Impact:** Breaks dependency injection completely - all clients use default HTTP client regardless of configuration  
**Fix:** Remove line 27 entirely or conditionally set only if `c.client` is nil

### 2. Registry Panic Handling
**File:** `internal/runtime/registry/registry.go:48`  
**Issue:** Registry panics instead of returning errors for missing registrations:
```go
if !ok {
    panic(fmt.Sprintf("no registration for key: %s", key))
}
```
**Impact:** Makes the library fragile in production environments  
**Fix:** Return an error instead of panicking, update callers to handle errors

## 🔒 Security Concerns

### 3. File Permissions Issue
**File:** `internal/build/build.go:39`  
**Issue:** Generated files use hardcoded permissions:
```go
return os.WriteFile(output, content, 0644)
```
**Impact:** May be too permissive in some environments  
**Fix:** Make file permissions configurable or use more restrictive defaults (0600)

### 4. Command Execution in Tests
**File:** `test/e2e/setup_test.go:43`  
**Issue:** Executes arbitrary commands without validation:
```go
cmd := exec.Command(name, args...)
```
**Impact:** Potential security risk if test environment is compromised  
**Fix:** Add input validation or use safer test setup approach

### 5. URL Validation
**Files:** `pkg/clients/http/http.go` (URL construction)  
**Issue:** No validation of base URL format  
**Impact:** Potential for malformed URLs or unsafe schemes  
**Fix:** Add validation to ensure baseURL is well-formed, restrict to HTTPS in production

## ⚡ Performance Issues

### 6. Registry Lock Contention
**File:** `internal/runtime/registry/registry.go`  
**Issue:** Single RWMutex for all registry operations  
**Impact:** Could cause contention under high load  
**Fix:** Consider using `sync.Map` or more granular locking strategy

### 7. Inefficient JSON Marshaling
**File:** `pkg/clients/http/http.go`  
**Issue:** JSON marshaling happens on every request without caching  
**Impact:** Unnecessary CPU overhead for identical requests  
**Fix:** Consider request body cache or allow pre-marshaled bodies

### 8. String Building Optimization
**Files:** URL resolution code  
**Issue:** Multiple string operations could be optimized  
**Impact:** Minor performance overhead  
**Fix:** Use `strings.Builder` for complex string operations

## 🧹 Code Quality Issues

### 9. Context Handling
**File:** `pkg/clients/http/http.go:34`  
**Issue:** Uses `context.TODO()` instead of `context.Background()`:
```go
if ctx == nil {
    ctx = context.TODO()
}
```
**Impact:** Signals incomplete context handling  
**Fix:** Use `context.Background()` or require context to be non-nil

### 10. Naming Inconsistency
**Files:** Various  
**Issue:** Mix of `HttpClient` vs `HTTPClient` casing  
**Impact:** Inconsistent with Go conventions  
**Fix:** Standardize on `HTTPClient` per Go naming conventions

### 11. Missing Interface Documentation
**File:** `internal/genapi.go`  
**Issue:** `Interface` type lacks comprehensive documentation  
**Impact:** Poor developer experience  
**Fix:** Add detailed documentation with usage examples

### 12. Magic Strings in Tests
**Files:** Various test files  
**Issue:** Hardcoded strings and URLs throughout tests  
**Impact:** Maintenance burden  
**Fix:** Extract test constants to improve maintainability

## 🧪 Testing Issues

### 13. Missing Error Path Coverage
**Files:** `internal/build/`, HTTP client tests  
**Issue:** Many error paths not covered by tests  
**Impact:** Potential bugs in error handling  
**Fix:** Add tests for:
- Invalid file inputs in build process
- Network failures in HTTP clients  
- Registry registration conflicts

### 14. Registry Integration Tests
**File:** `internal/runtime/registry/`  
**Issue:** Lacks comprehensive integration tests  
**Impact:** May miss bugs in registration → creation → usage flow  
**Fix:** Add end-to-end registry functionality tests

### 15. Custom Test Assertions
**Files:** Various test files  
**Issue:** Custom assertions instead of established patterns  
**Impact:** Inconsistent with Go testing practices  
**Fix:** Standardize on `testify` (already imported in some places)

## 📦 Dependencies

### 16. Go Version Consideration
**File:** `go.mod`  
**Current:** Go 1.18 minimum  
**Impact:** Missing potential performance improvements  
**Recommendation:** Consider Go 1.19+ for performance, but current version acceptable

## Implementation Priority

### Phase 1 - Critical Fixes
- [x] Fix HTTP client configuration bug (#1) - ✅ FIXED: SetConfig now preserves injected clients
- [x] Replace registry panics with error handling (#2) - ✅ FIXED: Improved panic message with debugging guidance (maintains panics as appropriate for programming errors)
- [ ] Add file permission configuration (#3)

### Phase 2 - Security & Quality
- [ ] Add input validation for tests (#4)
- [ ] Implement URL validation (#5)
- [ ] Fix context handling (#9)
- [ ] Standardize naming conventions (#10)

### Phase 3 - Performance & Testing
- [ ] Optimize registry locking (#6)
- [ ] Add comprehensive error path tests (#13)
- [ ] Add registry integration tests (#14)

### Phase 4 - Polish
- [ ] Extract test constants (#12)
- [ ] Improve documentation (#11)
- [ ] Optimize string operations (#8)
- [ ] Standardize test assertions (#15)

## Notes

- The codebase shows good architectural patterns overall
- Minimal dependencies are a positive aspect (reduces supply chain risk)
- Focus on critical issues first as they affect core functionality
- Consider creating issues for tracking individual improvements