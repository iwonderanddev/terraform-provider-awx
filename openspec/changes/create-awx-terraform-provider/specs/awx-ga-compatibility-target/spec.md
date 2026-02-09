## ADDED Requirements

### Requirement: GA compatibility scope
The provider SHALL declare support for AWX 24.6.1 API v2 as the GA compatibility target and SHALL not claim backward compatibility for older AWX versions.

#### Scenario: Compatibility statement publication
- **WHEN** GA documentation is generated
- **THEN** the compatibility section explicitly states support for AWX 24.6.1 API v2 only

### Requirement: Compatibility validation against target version
The provider SHALL include validation workflows that exercise core resource and data source behavior against AWX 24.6.1.

#### Scenario: Target-version validation run
- **WHEN** a compatibility validation run is executed against AWX 24.6.1
- **THEN** core lifecycle and import behaviors pass for the supported scope
