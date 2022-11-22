package cmd

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/SquareFactory/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1"
	"github.com/SquareFactory/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1/cluster"
	"github.com/creasty/defaults"
	"github.com/k0sproject/dig"
	"github.com/k0sproject/rig"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

// DefaultK0sYaml is pretty much what "k0s default-config" outputs
var DefaultK0sYaml = []byte(`apiVersion: k0s.k0sproject.io/v1beta1
kind: Cluster
metadata:
  name: k0s
spec:
  api:
    port: 6443
    k0sApiPort: 9443
  storage:
    type: etcd
  network:
    podCIDR: 10.244.0.0/16
    serviceCIDR: 10.96.0.0/12
    provider: calico
    calico:
      mode: 'vxlan'
      overlay: Always
      mtu: 1450
      wireguard: false
    kubeProxy:
      disabled: false
      mode: iptables
  podSecurityPolicy:
    defaultPolicy: 00-k0s-privileged
  telemetry:
    enabled: false
  installConfig:
    users:
      etcdUser: etcd
      kineUser: kube-apiserver
      konnectivityUser: konnectivity-server
      kubeAPIserverUser: kube-apiserver
      kubeSchedulerUser: kube-scheduler
  konnectivity:
    agentPort: 8132
    adminPort: 8133

  extensions:
  helm:
    repositories:
      - name: traefik
        url: https://helm.traefik.io/traefik
      - name: bitnami
        url: https://charts.bitnami.com/bitnami
      - name: jetstack
        url: https://charts.jetstack.io
      - name: csi-driver-nfs
        url: https://raw.githubusercontent.com/kubernetes-csi/csi-driver-nfs/master/charts
      - name: cloudve
        url: https://github.com/CloudVE/helm-charts/raw/master
    charts:
      - name: metallb
        chartname: bitnami/metallb
        version: '4.1.2'
        namespace: metallb

      - name: traefik
        chartname: traefik/traefik
        version: '10.24.2'
        namespace: traefik
        values: |
          deployment:
            kind: DaemonSet

          ingressClass:
            enabled: true
            isDefaultClass: true

          service:
            enabled: true
            annotations:
              metallb.universe.tf/address-pool: main-pool
              metallb.universe.tf/allow-shared-ip: traefik-lb-key
            spec:
              externalTrafficPolicy: Cluster
              loadBalancerIP: 192.168.1.100

          providers:
            kubernetesCRD:
              enabled: true
              allowCrossNamespace: true
              allowExternalNameServices: true
              namespaces: []
            kubernetesIngress:
              enabled: true
              allowExternalNameServices: true
              namespaces: []
              ingressClass: traefik
              publishedService:
                enabled: true

          globalArguments:
            - '--global.checknewversion'
            - '--api.dashboard=true'

          additionalArguments:
            - '--entryPoints.websecure.proxyProtocol.insecure'
            - '--entryPoints.websecure.forwardedHeaders.insecure'

          ingressRoute:
            dashboard:
              enabled: false

          ports:
            traefik:
              port: 9000
              expose: false
              exposedPort: 9000
              protocol: TCP
            dns-tcp:
              port: 8053
              expose: true
              exposedPort: 53
              protocol: TCP
            dns-udp:
              port: 8054
              expose: true
              exposedPort: 53
              protocol: UDP
            web:
              port: 80
              expose: true
              exposedPort: 80
              protocol: TCP
            websecure:
              port: 443
              expose: true
              exposedPort: 443
              protocol: TCP
              # You MUST open port 443 UDP!
              # HTTP3 upgrades the connection from TCP to UDP.
              http3: true
              tls:
                enabled: true
            metrics:
              port: 9100
              expose: false
              exposedPort: 9100
              protocol: TCP

          experimental:
            http3:
              enabled: true

          securityContext:
            capabilities:
              drop: [ALL]
              add: [NET_BIND_SERVICE]
            readOnlyRootFilesystem: true
            runAsGroup: 0
            runAsNonRoot: false
            runAsUser: 0

          podSecurityContext:
            fsGroup: 65532

      - name: cert-manager
        chartname: jetstack/cert-manager
        version: 'v1.9.1'
        namespace: cert-manager
        values: |
          installCRDs: true

      - name: csi-driver-nfs
        chartname: csi-driver-nfs/csi-driver-nfs
        version: 'v4.1.0'
        namespace: csi-driver-nfs
        values: |
          driver:
            mountPermissions: 0775
          kubeletDir: /var/lib/k0s/kubelet
`)

var defaultHosts = cluster.Hosts{
	&cluster.Host{
		Connection: rig.Connection{
			SSH: &rig.SSH{
				Address: "10.0.0.1",
				User:    "root",
				Port:    22,
				KeyPath: "~/.ssh/id_ed25519",
			},
		},
		Role:             "controller+worker",
		NoTaints:         true,
		PrivateInterface: "eno1",
		PrivateAddress:   "10.0.0.1",
		InstallFlags: cluster.Flags{
			"--debug",
			"--labels=\"topology.kubernetes.io/region=ch-sion,topology.kubernetes.io/zone=ch-sion-1\"",
			"--disable-components coredns",
		},
	},
	&cluster.Host{
		Connection: rig.Connection{
			SSH: &rig.SSH{
				Address: "10.0.0.2",
			},
		},
		Role:             "worker",
		PrivateInterface: "eno1",
		PrivateAddress:   "10.0.0.2",
		InstallFlags: cluster.Flags{
			"--debug",
			"--labels=\"topology.kubernetes.io/region=ch-sion,topology.kubernetes.io/zone=ch-sion-1\"",
		},
	},
}

func hostFromAddress(addr, role, user, keypath string) *cluster.Host {
	port := 22

	if idx := strings.Index(addr, "@"); idx > 0 {
		user = addr[:idx]
		addr = addr[idx+1:]
	}

	if idx := strings.Index(addr, ":"); idx > 0 {
		pstr := addr[idx+1:]
		if p, err := strconv.Atoi(pstr); err == nil {
			port = p
		}
		addr = addr[:idx]
	}

	host := &cluster.Host{
		Connection: rig.Connection{
			SSH: &rig.SSH{
				Address: addr,
				Port:    port,
			},
		},
	}
	if role != "" {
		host.Role = role
	} else {
		host.Role = "worker"
	}
	if user != "" {
		host.SSH.User = user
	}
	if keypath != "" {
		host.SSH.KeyPath = &keypath
	}

	_ = defaults.Set(host)

	return host
}

func buildHosts(addresses []string, ccount int, user, keypath string) cluster.Hosts {
	var hosts cluster.Hosts
	role := "controller"
	for _, a := range addresses {
		// strip trailing comments
		if idx := strings.Index(a, "#"); idx > 0 {
			a = a[:idx]
		}
		a = strings.TrimSpace(a)
		if a == "" || strings.HasPrefix(a, "#") {
			// skip empty and comment lines
			continue
		}

		if len(hosts) >= ccount {
			role = "worker"
		}

		hosts = append(hosts, hostFromAddress(a, role, user, keypath))
	}

	if len(hosts) == 0 {
		return defaultHosts
	}

	return hosts
}

var initCommand = &cli.Command{
	Name:        "init",
	Usage:       "Create a configuration template",
	Description: "Outputs a new cfctl configuration. When a list of addresses are provided, hosts are generated into the configuration. The list of addresses can also be provided via stdin.",
	ArgsUsage:   "[[user@]address[:port] ...]",
	Before:      actions(initLogging),
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "k0s",
			Usage: "Include a skeleton k0s config section",
		},
		&cli.StringFlag{
			Name:    "cluster-name",
			Usage:   "Cluster name",
			Aliases: []string{"n"},
			Value:   "k0s-cluster",
		},
		&cli.IntFlag{
			Name:    "controller-count",
			Usage:   "The number of controllers to create when addresses are given",
			Aliases: []string{"C"},
			Value:   1,
		},
		&cli.StringFlag{
			Name:    "user",
			Usage:   "Host user when addresses given",
			Aliases: []string{"u"},
		},
		&cli.StringFlag{
			Name:    "key-path",
			Usage:   "Host key path when addresses given",
			Aliases: []string{"i"},
		},
	},
	Action: func(ctx *cli.Context) error {
		var addresses []string

		// Read addresses from stdin
		stat, err := os.Stdin.Stat()
		if err == nil {
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				rd := bufio.NewReader(os.Stdin)
				for {
					row, _, err := rd.ReadLine()
					if err != nil {
						break
					}
					addresses = append(addresses, string(row))
				}
				if err != nil {
					return err
				}

			}
		}

		// Read addresses from args
		addresses = append(addresses, ctx.Args().Slice()...)

		cfg := v1beta1.Cluster{
			APIVersion: v1beta1.APIVersion,
			Kind:       "Cluster",
			Metadata:   &v1beta1.ClusterMetadata{Name: ctx.String("cluster-name")},
			Spec: &cluster.Spec{
				Hosts: buildHosts(addresses, ctx.Int("controller-count"), ctx.String("user"), ctx.String("key-path")),
				K0s:   &cluster.K0s{},
			},
		}

		if err := defaults.Set(&cfg); err != nil {
			return err
		}

		if ctx.Bool("k0s") {
			cfg.Spec.K0s.Config = dig.Mapping{}
			if err := yaml.Unmarshal(DefaultK0sYaml, &cfg.Spec.K0s.Config); err != nil {
				return err
			}
		}

		encoder := yaml.NewEncoder(os.Stdout)
		return encoder.Encode(&cfg)
	},
}
