# rebind

## Name

*rebind* - rebinds a domain from one IP address to another IP address to facilitate testing [DNS Rebinding vulnerabilities](https://en.wikipedia.org/wiki/DNS_rebinding).

## Description

This is a [CoreDNS](https://github.com/coredns/coredns/) plugin. It rebinds domains from one IP address to another IP address. Use this plugin to learn more about DNS rebinding attacks or to test proof of concepts as a security researcher in a responsible manner.

This plugin is inspired by [nccgroup/singularity](https://github.com/nccgroup/singularity) and [brannondorsey/whonow](https://github.com/brannondorsey/whonow).


## Syntax

~~~ corefile
rebind example.com {
  first_ip 1.2.3.4
  second_ip 0.0.0.0
  strategy first_then_second
}
~~~

- **first_ip** is the first IP address. This is usually an IP address that you own
- **second_ip** is the second IP address to rebind to. This is usually the target IP address of the vulnerable server
- **strategy** is one of the following:
  - first_then_second: responds with the `first_ip` and then responds with the `second_ip` address for all subsequent requests
  - random: responds with a random selection of `first_ip` and `second_ip`
  - round_robin: responds in a round robin fashion of `first_ip` and then `second_ip`

## Examples

In this configuration, a DNS request to `rebind.example.com` will receive a response of `1.2.3.4`. All future DNS requests will respond with `0.0.0.0`.

~~~ corefile
example.com {
  rebind rebind.example.com {
    first_ip 1.2.3.4
    second_ip 0.0.0.0
  }
}
~~~

## Compilation

This package will always be compiled as part of CoreDNS and not in a standalone way. It will require you to use `go get` or as a dependency on [plugin.cfg](https://github.com/coredns/coredns/blob/master/plugin.cfg).

The [manual](https://coredns.io/manual/toc/#what-is-coredns) will have more information about how to configure and extend the server with external plugins.

A simple way to consume this plugin, is by adding the following on [plugin.cfg](https://github.com/coredns/coredns/blob/master/plugin.cfg), and recompile it as [detailed on coredns.io](https://coredns.io/2017/07/25/compile-time-enabling-or-disabling-plugins/#build-with-compile-time-configuration-file).

~~~
rebind:github.com/ivantsepp/rebind
~~~

Put this early in the plugin list, so that *rebind* is executed before any of the other plugins.

After this you can compile coredns by:

``` sh
go generate
go build
```

Or you can instead use make:

``` sh
make
```
