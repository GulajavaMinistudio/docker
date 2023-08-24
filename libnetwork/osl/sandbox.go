// Package osl describes structures and interfaces which abstract os entities
package osl

import (
	"net"

	"github.com/docker/docker/libnetwork/types"
)

// SandboxType specify the time of the sandbox, this can be used to apply special configs
type SandboxType int

const (
	// SandboxTypeIngress indicates that the sandbox is for the ingress
	SandboxTypeIngress = iota
	// SandboxTypeLoadBalancer indicates that the sandbox is a load balancer
	SandboxTypeLoadBalancer = iota
)

type Iface struct {
	SrcName, DstPrefix string
}

// IfaceOption is a function option type to set interface options.
type IfaceOption func(i *Interface) error

// NeighOption is a function option type to set neighbor options.
type NeighOption func(nh *neigh)

// Sandbox represents a network sandbox, identified by a specific key.  It
// holds a list of Interfaces, routes etc, and more can be added dynamically.
type Sandbox interface {
	// Key returns the path where the network namespace is mounted.
	Key() string

	// AddInterface adds an existing Interface to this sandbox. The operation will rename
	// from the Interface SrcName to DstName as it moves, and reconfigure the
	// interface according to the specified settings. The caller is expected
	// to only provide a prefix for DstName. The AddInterface api will auto-generate
	// an appropriate suffix for the DstName to disambiguate.
	AddInterface(SrcName string, DstPrefix string, options ...IfaceOption) error

	// SetGateway sets the default IPv4 gateway for the sandbox.
	SetGateway(gw net.IP) error

	// SetGatewayIPv6 sets the default IPv6 gateway for the sandbox.
	SetGatewayIPv6(gw net.IP) error

	// UnsetGateway the previously set default IPv4 gateway in the sandbox.
	UnsetGateway() error

	// UnsetGatewayIPv6 unsets the previously set default IPv6 gateway in the sandbox.
	UnsetGatewayIPv6() error

	// GetLoopbackIfaceName returns the name of the loopback interface
	GetLoopbackIfaceName() string

	// AddAliasIP adds the passed IP address to the named interface
	AddAliasIP(ifName string, ip *net.IPNet) error

	// RemoveAliasIP removes the passed IP address from the named interface
	RemoveAliasIP(ifName string, ip *net.IPNet) error

	// DisableARPForVIP disables ARP replies and requests for VIP addresses
	// on a particular interface.
	DisableARPForVIP(ifName string) error

	// AddStaticRoute adds a static route to the sandbox.
	AddStaticRoute(*types.StaticRoute) error

	// RemoveStaticRoute removes a static route from the sandbox.
	RemoveStaticRoute(*types.StaticRoute) error

	// AddNeighbor adds a neighbor entry into the sandbox.
	AddNeighbor(dstIP net.IP, dstMac net.HardwareAddr, force bool, option ...NeighOption) error

	// DeleteNeighbor deletes neighbor entry from the sandbox.
	DeleteNeighbor(dstIP net.IP, dstMac net.HardwareAddr, osDelete bool) error

	// InvokeFunc invoke a function in the network namespace.
	InvokeFunc(func()) error

	// Destroy destroys the sandbox.
	Destroy() error

	// Restore restores the sandbox.
	Restore(ifsopt map[Iface][]IfaceOption, routes []*types.StaticRoute, gw net.IP, gw6 net.IP) error

	// ApplyOSTweaks applies operating system specific knobs on the sandbox.
	ApplyOSTweaks([]SandboxType)

	Info
}

// Info represents all possible information that
// the driver wants to place in the sandbox which includes
// interfaces, routes and gateway
type Info interface {
	// Interfaces returns the collection of Interface previously added with the AddInterface
	// method. Note that this doesn't include network interfaces added in any
	// other way (such as the default loopback interface which is automatically
	// created on creation of a sandbox).
	Interfaces() []*Interface

	// Gateway returns the IPv4 gateway for the sandbox.
	Gateway() net.IP

	// GatewayIPv6 returns the IPv6 gateway for the sandbox.
	GatewayIPv6() net.IP

	// StaticRoutes returns additional static routes for the sandbox. Note that
	// directly connected routes are stored on the particular interface they
	// refer to.
	StaticRoutes() []*types.StaticRoute
}
