# kubeconfig-ca-fetch

Print a new kubeconfig on stdout, optionally print some messages about timing and
failures to the stderr.

# tl;dr

This is probably not what you want. But if you know who you are, and you need
this, then see also: [Howto kubeconfig][Kubelogin prereqs], for team members
in the `kingdon-ci:weave-gitops` GitHub OIDC unit.

Welcome to use this tool for learning, or adapt the code for other purposes.

# About

This is a command you should only need to run once. It connects to many
clusters and saves the certificate authority data that they present in a
Kubeconfig. A list of kubeconfigs does not dynamically update; it's hardcoded.

If the list of clusters changes, then you would need to edit the source code,
update the list, rebuild, and run it again to produce a new Kubeconfig. The
default timeout is 2s so this command always returns fast.

The application `./cmd/kubeconfig-ca-fetch` has zero knowledge about the remote
Kubernetes API services, and (generally speaking, insecurely) trusts that all
the kube APIs are valid in order to generate a kubeconfig for all the clusters
that our OIDC access token will authorize.

Use one of the supported methods to [download kubelogin][Kubelogin prereqs].
The exec mode configuration will [run kubelogin][] to get an id-token.

[Kubelogin prereqs]: https://howto.howard.moomboo.space/#prerequisites
[run kubelogin]: https://howto.howard.moomboo.space/#tldr-run-kubelogin

### Safety

It should be safe to use this in environments that you have good reason to know
your connection to the Kubernetes servers has not been MITM'ed (they are local)
and you know your DNS records that direct to the Kubers are surely trustworthy.
It is safer to source certificate authority data through a verified provenance.

Safety is relative. My goal was to enable an unpermissioned client (Chromebook)
to access my Kubernetes clusters using only an OIDC token, and this does that.
Is it safer than copying Kubeconfigs around as an operational best-practice?

IMHO yes it is more secure, but YMMV.

Now my home lab's kubeconfigs access tokens expire and refreshing according to
my GitHub account requires a 2-factor authorization. 📈

## Instructions

```bash
kubeconfig-ca-fetch > kube.config
```

## Build

To build this, you can run:

```
go install github.com/kingdon-ci/kubeconfig-ca-fetch/cmd/kubeconfig-ca-fetch@latest
```

or if you cloned this repo locally, try:

```bash
make tidy && make build
```
