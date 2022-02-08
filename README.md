# Ocean CD

A command-line interface to manage your [Ocean CD](https://spot.io/) resources.

## What is Ocean CD?
Ocean CD provides DevOps and Infrastructure teams with out of the box processes to reimplement 
and share complex and mission critical pieces of CD across different environments, such as 
progressive delivery and verification of the software deployments. Service owners are able to 
promote service changes to production without code or re-inventing deployment strategies. 
Developers are able to commit with confidence, now that the deployment phases are managed and 
visible.

## Why Ocean CD?
Ocean CD provides a central place to observe the deployment (e.g., state, progress, and resources). 
These visibility features allow quick identification of issues during and after the deployment process 
and ensure reliability at scale.

You will set up external verifications, the method Ocean CD uses to integrate your unique test 
outcomes as part of an orchestrated release process, and define webhook API notifications to 
communicate between Ocean CD and your external DevOps related tools.

Ocean CD performs automated Rollback, a mission critical feature not covered out of the box by 
Kubernetes. You will be able to define the failure policy which determines the type of rollback that 
will be employed

## Documentation
To learn more about Argo Rollouts go to the [complete documentation](https://docs.spot.io/ocean-cd/ocean-cd-overview).

## Installation

For macOS users, the easiest way to install `oceancd` is to use [Homebrew](https://brew.sh/):

```
$ brew install spotinst/tap/oceancd
```

Otherwise, please download the latest release from the [Releases](https://github.com/spotinst/spot-oceancd-cli/releases) page.

## Getting Started

Before using `oceancd`, you need to configure your Spot credentials. 
The quickest way to get started is to run the `oceancd configure` command:

```
$ spotctl configure
```

## Documentation

If you're new to Ocean CD and want to get started, please checkout our [Getting Started](https://docs.spot.io/ocean-cd/getting-started/) guide, available on the [Spot Documentation](https://help.spot.io/) website.

## Getting Help

We use GitHub issues for tracking bugs and feature requests. Please use these community resources for getting help:

- Join our Spot community on [Slack](http://slack.spot.io/).
- Open an [issue](https://github.com/spotinst/spot-oceancd-cli/issues/new).

## Community

- [Slack](http://slack.spot.io/)
- [Twitter](https://twitter.com/spot_hq/)

## License

Code is licensed under the [Apache License 2.0](LICENSE).