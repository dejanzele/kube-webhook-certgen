package cmd

import (
	"context"
	"os"

	"github.com/jet/kube-webhook-certgen/pkg/k8s"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var patch = &cobra.Command{
	Use:    "patch",
	Short:  "Patch a validatingwebhookconfiguration and mutatingwebhookconfiguration 'webhook-name' by using the ca from 'secret-name' in 'namespace'",
	Long:   "Patch a validatingwebhookconfiguration and mutatingwebhookconfiguration 'webhook-name' by using the ca from 'secret-name' in 'namespace'",
	PreRun: prePatchCommand,
	RunE:   patchCommand,
}

func prePatchCommand(cmd *cobra.Command, args []string) {
	configureLogging(cmd, args)
	if cfg.patchMutating == false && cfg.patchValidating == false {
		log.Fatal("patch-validating=false, patch-mutating=false. You must patch at least one kind of webhook, otherwise this command is a no-op")
		os.Exit(1)
	}
	switch cfg.patchFailurePolicy {
	case "":
		break
	case "Ignore":
	case "Fail":
		failurePolicy = cfg.patchFailurePolicy
	default:
		log.Fatalf("patch-failure-policy %s is not valid", cfg.patchFailurePolicy)
		os.Exit(1)
	}
}

func patchCommand(_ *cobra.Command, _ []string) error {
	k := k8s.New(cfg.kubeconfig)
	ca, err := k.GetCaFromSecret(cfg.secretName, cfg.namespace, cfg.caName)
	if err != nil {
		return err
	}

	if ca == nil {
		log.Fatalf("no secret with '%s' in '%s'", cfg.secretName, cfg.namespace)
	}

	return k.PatchWebhookConfigurations(
		context.Background(),
		cfg.webhookName,
		ca,
		failurePolicy,
		cfg.patchMutating,
		cfg.patchValidating,
		k8s.AdmissionRegistrationVersion(cfg.admissionRegistrationVersion),
	)
}

func init() {
	rootCmd.AddCommand(patch)
	patch.Flags().StringVar(&cfg.secretName, "secret-name", "", "Name of the secret where certificate information will be read from")
	patch.Flags().StringVar(&cfg.namespace, "namespace", "", "Namespace of the secret where certificate information will be read from")
	patch.Flags().StringVar(&cfg.webhookName, "webhook-name", "", "Name of validatingwebhookconfiguration and mutatingwebhookconfiguration that will be updated")
	patch.Flags().BoolVar(&cfg.patchValidating, "patch-validating", true, "If true, patch validatingwebhookconfiguration")
	patch.Flags().BoolVar(&cfg.patchMutating, "patch-mutating", true, "If true, patch mutatingwebhookconfiguration")
	patch.Flags().StringVar(&cfg.patchFailurePolicy, "patch-failure-policy", "", "If set, patch the webhooks with this failure policy. Valid options are Ignore or Fail")
	create.Flags().StringVar(&cfg.admissionRegistrationVersion, "admission-registration-version", "v1", "admissionregistration.k8s.io api version")
	_ = patch.MarkFlagRequired("secret-name")
	_ = patch.MarkFlagRequired("namespace")
	_ = patch.MarkFlagRequired("webhook-name")
}
