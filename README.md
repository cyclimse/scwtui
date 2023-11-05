# scwtui

> **Note**
> Not an official Scaleway project, use at your own risk.

`scwtui` is a terminal user interface for Scaleway powered by [bubbletea](https://github.com/charmbracelet/bubbletea). Inspired by [k9s](https://k9scli.io/).

It allows you to view all the resources in your Scaleway account and perform basic operations on them.

## Running

The easiest way to run `scwtui` is to use Docker.

```bash
docker run -it --rm -v "$HOME/.config/scw:/root/.config/scw" cyclimse/scwtui:0.1.0
```

## Keybindings

| Keybinding      | Description                             |
|-----------------|-----------------------------------------|
| `esc`, `ctrl+c` | Quit                                    |
| `\`             | Search                                  |
| `d`             | Describe selected resource              |
| `x`             | Delete selected resource                |
| `l`             | View Cockpit logs for selected resource |

## Features

### Search

The search is powered by [bleve](https://github.com/blevesearch/bleve).

You can use fields to search for specific resources. For example, to search for all resources in the `fr-par` region, you can use `region:fr-par`.

### Describe

You can view the details of a resource by pressing `d` when it is selected. This will show all the details of the resource as a JSON object.

### Delete

You can delete a resource by pressing `x` when it is selected. This will prompt you to confirm the deletion of the resource.

### Logs

> **Note**
> This feature is only available for resources integrated with Cockpit.

You can view the logs for a resource by pressing `l` when it is selected. This will open a new window with the logs for the resource.

In practice, this feature relies on the Cockpit Loki API. It will generate tokens for each project to allow you to view the logs for the project.

This requires your API key to have the `ObservabilityFullAccess` permission set to work.

## Supported Resources

| Resource             | List | Describe | Delete | Logs |
|----------------------|------|----------|--------|------|
| Serverless Function  | ✅    | ✅        | ✅      | ✅    |
| Serverless Container | ✅    | ✅        | ✅      | ✅    |
| Registry Namespace   | ✅    | ✅        | ✅      | ❌    |
| RDB Instance         | ✅    | ✅        | ✅      | ✅    |
