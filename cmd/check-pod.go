package main

import (
	"errors"
	"io"

	"git.i.sedorz.net/infrastructure/icinga/check-k8s/pkg/print"

	"git.i.sedorz.net/infrastructure/icinga/check-k8s/pkg/kube"
	"github.com/spf13/cobra"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type checkPodCmd struct {
	version   string
	out       io.Writer
	name      string
	namespace string
}

func newCheckPodCmd(out io.Writer) *cobra.Command {
	c := &checkPodCmd{out: out}

	cmd := &cobra.Command{
		Use:              "pod",
		Short:            "check if a k8s pod resource is healthy",
		TraverseChildren: true,
		SilenceUsage:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("pod name is required")
			}
			c.name = args[0]
			return c.run()
		},
	}

	cmd.PersistentFlags().StringVarP(&c.namespace, "namespace", "n", "", "the namespace where the deployment is")
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func (c *checkPodCmd) run() error {
	client, err := kube.GetKubeClient(settings.KubeContext)
	if err != nil {
		return err
	}
	resource, err := client.CoreV1().Pods(c.namespace).List(meta_v1.ListOptions{})
	if err != nil {
		return err
	}
	print.Fyaml(c.out, resource)
	return nil
}
