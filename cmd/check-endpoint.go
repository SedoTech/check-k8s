package main

import (
	"errors"
	"io"

	"github.com/benkeil/check-k8s/pkg/print"

	"github.com/benkeil/check-k8s/pkg/kube"
	"github.com/spf13/cobra"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type checkEndpointCmd struct {
	version   string
	out       io.Writer
	name      string
	namespace string
}

func newCheckEndpointCmd(out io.Writer) *cobra.Command {
	c := &checkEndpointCmd{out: out}

	cmd := &cobra.Command{
		Use:              "endpoint",
		Short:            "check if a k8s endpoint resource is healthy",
		TraverseChildren: true,
		SilenceUsage:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("endpoint name is required")
			}
			c.name = args[0]
			return c.run()
		},
	}

	cmd.PersistentFlags().StringVarP(&c.namespace, "namespace", "n", "", "the namespace where the deployment is")
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func (c *checkEndpointCmd) run() error {
	client, err := kube.GetKubeClient(settings.KubeContext)
	if err != nil {
		return err
	}
	resource, err := client.CoreV1().Endpoints(c.namespace).Get(c.name, meta_v1.GetOptions{})
	//labelSelector := labels.Set(map[string]string{"mylabel": "ourdaomain1"}).AsSelector()
	//resource, err := client.CoreV1().Endpoints(c.namespace).List(meta_v1.ListOptions{})
	if err != nil {
		return err
	}
	print.Fyaml(c.out, resource)
	return nil
}
