### SERVER

- Improve logging
- Direct Messages between users
- User login - bind connection to usrname
- List chat users - http (not enforced, could be ws if so desired)
  Nice to have:
- Multiple rooms - users can chat on different rooms
  - should we change username? how?s
- Load Chat History when connected \*\*\*
  - will have to preserve state - how? files, db?

### CLIENT

- Accept input form stdin
- Quit / Disconnect
- Config server connection params: -url, port etc
  - options:
    - envvars
    - file config
- NOTE: does
