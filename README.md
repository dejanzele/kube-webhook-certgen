[![Go Report Card](https://goreportcard.com/badge/github.com/dejanzele/kube-webhook-certgen)](https://goreportcard.com/report/github.com/dejanzele/kube-webhook-certgen)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/dejanzele/kube-webhook-certgen?sort=semver)](https://github.com/dejanzele/kube-webhook-certgen/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/dpejcev/kube-webhook-certgen?color=blue)](https://hub.docker.com/r/dpejcev/kube-webhook-certgen/tags)

# Kubernetes webhook certificate generator and patcher

This is repo is a fork from [jet/kube-webhook-certgen](https://github.com/jet/kube-webhook-certgen) as the original project is no longer maintained.
This repo will diverge from [jet/kube-webhook-certgen](https://github.com/jet/kube-webhook-certgen) and [kubernetes/ingress-nginx](https://github.com/kubernetes/ingress-nginx/tree/main/images/kube-webhook-certgen)
as it has a roadmap of its own.

## Overview
Generates a CA and leaf certificate with a long (100y) expiration, then patches [Kubernetes Admission Webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)
by setting the `caBundle` field with the generated CA. 
Can optionally patch the hooks `failurePolicy` setting - useful in cases where a single Helm chart needs to provision resources
and hooks at the same time as patching.

The utility works in two parts, optimized to work better with the Helm provisioning process that leverages pre-install and post-install hooks to execute this as a Kubernetes job.

## Security Considerations
This tool may not be adequate in all security environments. If a more complete solution is required, you may want to 
seek alternatives such as [jetstack/cert-manager](https://github.com/jetstack/cert-manager)

## Command line options
```
Use this to create a ca and signed certificates and patch admission webhooks to allow for quick
                   installation and configuration of validating and admission webhooks.

Usage:
  kube-webhook-certgen [flags]
  kube-webhook-certgen [command]

Available Commands:
  create      Generate a ca and server cert+key and store the results in a secret 'secret-name' in 'namespace'
  help        Help about any command
  patch       Patch a validatingwebhookconfiguration and mutatingwebhookconfiguration 'webhook-name' by using the ca from 'secret-name' in 'namespace'
  version     Prints the CLI version information

Flags:
  -h, --help                           help for kube-webhook-certgen
      --admission-registration-version admissionregistration.k8s.io api version
      --kubeconfig string              Path to kubeconfig file: e.g. ~/.kube/kind-config-kind
      --log-format string              Log format: text|json (default "text")
      --log-level string               Log level: panic|fatal|error|warn|info|debug|trace (default "info")
```

### Create
```
Generate a ca and server cert+key and store the results in a secret 'secret-name' in 'namespace'

Usage:
  kube-webhook-certgen create [flags]

Flags:
      --cert-name string     Name of cert file in the secret (default "cert")
  -h, --help                 help for create
      --host string          Comma-separated hostnames and IPs to generate a certificate for
      --key-name string      Name of key file in the secret (default "key")
      --namespace string     Namespace of the secret where certificate information will be written
      --secret-name string   Name of the secret where certificate information will be written

Global Flags:
      --kubeconfig string   Path to kubeconfig file: e.g. ~/.kube/kind-config-kind
      --log-format string   Log format: text|json (default "json")
      --log-level string    Log level: panic|fatal|error|warn|info|debug|trace (default "info")
```

### Patch
```
Patch a validatingwebhookconfiguration and mutatingwebhookconfiguration 'webhook-name' by using the ca from 'secret-name' in 'namespace'

Usage:
  kube-webhook-certgen patch [flags]

Flags:
  -h, --help                          help for patch
      --namespace string              Namespace of the secret where certificate information will be read from
      --patch-failure-policy string   If set, patch the webhooks with this failure policy. Valid options are Ignore or Fail
      --patch-mutating                If true, patch mutatingwebhookconfiguration (default true)
      --patch-validating              If true, patch validatingwebhookconfiguration (default true)
      --secret-name string            Name of the secret where certificate information will be read from
      --webhook-name string           Name of validatingwebhookconfiguration and mutatingwebhookconfiguration that will be updated

Global Flags:
      --kubeconfig string   Path to kubeconfig file: e.g. ~/.kube/kind-config-kind
      --log-format string   Log format: text|json (default "text")
      --log-level string    Log level: panic|fatal|error|warn|info|debug|trace (default "info")
```

## Recent changes
* Updated go version to v1.19
* Added support for `--admission-registration-version` flag which allows users to select which version of admissionregistration.k8s.io they want to use (v1 or v1beta1)
