package main

import (
	"errors"
	"io"

	"github.com/benkeil/check-k8s/pkg/print"

	"github.com/benkeil/check-k8s/pkg/kube"
	"github.com/spf13/cobra"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type checkServiceCmd struct {
	version   string
	out       io.Writer
	name      string
	namespace string
}

func newCheckServiceCmd(out io.Writer) *cobra.Command {
	c := &checkServiceCmd{out: out}

	cmd := &cobra.Command{
		Use:              "service",
		Short:            "check if a k8s service resource is healthy",
		TraverseChildren: true,
		SilenceUsage:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("service name is required")
			}
			c.name = args[0]
			return c.run()
		},
	}

	cmd.PersistentFlags().StringVarP(&c.namespace, "namespace", "n", "", "the namespace where the deployment is")
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func (c *checkServiceCmd) run() error {
	client, err := kube.GetKubeClient(settings.KubeContext)
	if err != nil {
		return err
	}
	resource, err := client.CoreV1().Services(c.namespace).Get(c.name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}
	print.Fyaml(c.out, resource)
	return nil
}
