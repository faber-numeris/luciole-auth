# luciole-auth
AuthN and AuthZ components used on the luciole ecosystem

## TODO Management

This project uses automated TODO tracking. When you push code with TODO comments to the main branch, a GitHub Action will automatically create issues from those TODOs.

### Adding TODOs

Use the following format for TODO comments:

```go
// TODO: Description of what needs to be done
```

or

```go
// TODO(username): Description with optional assignee
```

### How it works

1. When code is pushed to `main` branch, the TODO-to-Issue workflow runs
2. The workflow scans all files for TODO comments
3. For each TODO found, it creates a GitHub issue (if one doesn't already exist)
4. Issues are automatically labeled with `todo` and `automated`
5. When a TODO is removed from code, the corresponding issue is automatically closed

### Viewing TODOs

- See `TODOS.md` for a comprehensive list of current TODOs and their context
- Check the GitHub Issues tab for individual TODO tracking
