# Test Data Directory

This directory contains test data for aiview CLI automated testing.

## Structure

```
testdata/
├── expected_outputs/     # Expected output samples for validation
│   ├── bilibili_hot.json
│   ├── douyin_hot.json
│   └── xiaohongshu_hot.json
├── mock_responses/       # Mock API responses for unit tests
│   ├── bilibili_api.json
│   ├── douyin_api.json
│   └── xiaohongshu_api.json
└── README.md
```

## Usage

### expected_outputs/
Contains expected JSON output samples from CLI commands. Used for validation in E2E tests.

### mock_responses/
Contains mock API responses from various platforms. Used in unit tests to simulate API calls without network access.

## Adding New Test Data

1. **Expected outputs**: Run a command with `--json` flag and save the output
2. **Mock responses**: Create JSON files matching the API response structure

## Notes

- All test data should be anonymized (no real user data)
- Keep file sizes small (< 100KB)
- Use realistic but fake data
