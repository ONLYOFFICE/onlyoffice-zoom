package registry

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/go-micro/plugins/v4/registry/kubernetes"
	"github.com/go-micro/plugins/v4/registry/mdns"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/registry/cache"
)

// NewRegistry looks up envs and configures respective registries based on those variables. Defaults to memory
func NewRegistry(opts ...Option) registry.Registry {
	options := NewOptions(opts...)

	var r registry.Registry
	switch options.RegistryType {
	case Kubernetes:
		r = kubernetes.NewRegistry(
			registry.Addrs(options.Addresses...),
			registry.Secure(options.Secure),
		)
	case Consul:
		r = consul.NewRegistry(
			registry.Addrs(options.Addresses...),
			registry.Secure(options.Secure),
		)
	case Etcd:
		r = etcd.NewRegistry(
			registry.Addrs(options.Addresses...),
			registry.Secure(options.Secure),
		)
	case MDNS:
		r = mdns.NewRegistry(
			registry.Addrs(options.Addresses...),
			registry.Secure(options.Secure),
		)
	default:
		r = mdns.NewRegistry(
			registry.Addrs(options.Addresses...),
			registry.Secure(options.Secure),
		)
	}

	return cache.New(r, cache.WithTTL(options.CacheTTL))
}
