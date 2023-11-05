# TODOs

## Features

### General

- [ ] Add more Scaleway resources

### Discovery

- [ ] Improve the retry on when discovering resources.
- [ ] Refresh the list of resources after a while. Currently, the list of resources is never refreshed from the API, only from the local store.

### Logging

- [ ] Add support for selecting a time range
- [ ] Add support for filtering logs. This can be done via the Loki API directly, so we could add a search bar to filter logs and rely on the Loki API to do the filtering.
