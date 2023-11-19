# scwtui

> **Warning**
> Not an official Scaleway project, use at your own risk.

`scwtui` is a terminal user interface for Scaleway powered by [bubbletea](https://github.com/charmbracelet/bubbletea).

It allows you to view all the resources in your Scaleway account and perform basic operations on them.

## Running

The easiest way to run `scwtui` is to use Docker.

```bash
docker run -it --rm -v "$HOME/.config/scw:/root/.config/scw" cyclimse/scwtui:0.1
```

You can also provide your Scaleway credentials using environment variables. See the [scaleway-cli](https://github.com/scaleway/scaleway-cli/blob/master/docs/commands/config.md) documentation for more information.

## Keybindings

| Keybinding      | Description                              |
|-----------------|------------------------------------------|
| `esc`, `ctrl+c` | Quit                                     |
| `\`             | Search                                   |
| `d`             | Describe selected resource               |
| `x`             | Delete selected resource                 |
| `l`             | View Cockpit logs for selected resource  |
| `t`             | View quick actions for selected resource |

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

This feature relies on the Cockpit Loki API. It will generate a token for each project to allow you to view the logs of the resources in the project. As such, you will need to provide the `ObservabilityFullAccess` permission set to your API token.

### Quick Actions

Quick actions are available for some resources. You can view the available actions by pressing `t` when a resource is selected. This will open a new window with the available actions.

## Supported Resources

| Resource             | List | Describe | Delete | Logs |      Actions       |
|----------------------|:----:|:--------:|:------:|:----:|:------------------:|
| Project              |  ✅   |    ✅     |   ✅    |  ❌   | `Activate Cockpit` |
| Cockpit              |  ✅   |    ✅     |   ✅    |  ❌   |   `Open Grafana`   |
| Serverless Function  |  ✅   |    ✅     |   ✅    |  ✅   |     (planned)      |
| Serverless Container |  ✅   |    ✅     |   ✅    |  ✅   |     (planned)      |
| Serverless Job       |  ✅   |    ✅     |   ✅    |  ✅   |  `Start`, `Stop`   |
| Registry Namespace   |  ✅   |    ✅     |   ✅    |  ❌   |                    |
| RDB Instance         |  ✅   |    ✅     |   ✅    |  ✅   |                    |
| Kapsule Cluster      |  ✅   |    ✅     |   ✅    |  ✅   |                    |
| Instance             |  ✅   |    ✅     |   ✅    |  ❌   |     (planned)      |

While it is possible to delete projects, it will require you to have deleted all the resources in the project first. In the future, this could be improved by deleting all the resources in the project first.

## Troubleshooting

### IAM Permissions

If you're using an IAM scoped token, you will need to provide the `ProjectReadyOnly` permission set to your API token. This is needed to retrieve the projects in your account, so that all resources can be listed.

In addition, you will need read permissions for all the products you want to use with `scwtui`. For instance, if you want to use this tool with all products, you can use the `AllProductsReadOnly` permission set.
