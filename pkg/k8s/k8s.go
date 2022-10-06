package k8s

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	admissionv1beta1 "k8s.io/api/admissionregistration/v1beta1"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8s struct {
	clientset kubernetes.Interface
}

type AdmissionRegistrationVersion string

const (
	admissionRegistrationV1      AdmissionRegistrationVersion = "v1"
	admissionRegistrationV1beta1 AdmissionRegistrationVersion = "v1beta1"
)

func New(kubeconfig string) *K8s {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.WithError(err).Fatal("error building kubernetes config")
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.WithError(err).Fatal("error creating kubernetes client")
	}

	return &K8s{clientset: c}
}

// PatchWebhookConfigurations will patch validatingWebhook and mutatingWebhook clientConfig configurations with
// the provided ca data. If failurePolicy is provided, patch all webhooks with this value.
func (k8s *K8s) PatchWebhookConfigurations(
	ctx context.Context,
	configurationNames string,
	ca []byte,
	failurePolicy string,
	patchMutating bool,
	patchValidating bool,
	version AdmissionRegistrationVersion,
) error {
	log.Infof(
		"patching webhook configurations '%s' mutating=%t, validating=%t, failurePolicy=%s",
		configurationNames, patchMutating, patchValidating, failurePolicy,
	)

	if patchValidating {
		if err := k8s.patchValidatingWebhookConfiguration(ctx, version, configurationNames, ca, failurePolicy); err != nil {
			return err
		}
	} else {
		log.Debug("validating hook patching not required")
	}

	if patchMutating {
		if err := k8s.patchMutatingWebhookConfiguration(ctx, version, configurationNames, ca, failurePolicy); err != nil {
			return err
		}
	} else {
		log.Debug("mutating hook patching not required")
	}

	log.Info("Patched hook(s)")

	return nil
}

func (k8s *K8s) patchValidatingWebhookConfiguration(
	ctx context.Context,
	version AdmissionRegistrationVersion,
	configurationNames string,
	ca []byte,
	failurePolicy string,
) error {
	switch version {
	case admissionRegistrationV1beta1:
		failurePolicyV1beta1 := admissionv1beta1.FailurePolicyType(failurePolicy)
		return k8s.patchValidatingWebhookConfigurationV1beta1(ctx, configurationNames, ca, &failurePolicyV1beta1)
	case admissionRegistrationV1:
		failurePolicyV1 := admissionv1.FailurePolicyType(failurePolicy)
		return k8s.patchValidatingWebhookConfigurationV1(ctx, configurationNames, ca, &failurePolicyV1)
	default:
		return errors.Errorf("invalid admissionregistration.k8s.io version: %s", version)
	}
}

func (k8s *K8s) patchValidatingWebhookConfigurationV1beta1(
	ctx context.Context,
	configurationNames string,
	ca []byte,
	failurePolicy *admissionv1beta1.FailurePolicyType,
) error {
	valHook, err := k8s.clientset.
		AdmissionregistrationV1beta1().
		ValidatingWebhookConfigurations().
		Get(ctx, configurationNames, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed getting admissionregistration.k8s.io/v1beta1 validating webhook")
	}

	for i := range valHook.Webhooks {
		h := &valHook.Webhooks[i]
		h.ClientConfig.CABundle = ca
		if *failurePolicy != "" {
			h.FailurePolicy = failurePolicy
		}
	}

	if _, err = k8s.clientset.AdmissionregistrationV1beta1().
		ValidatingWebhookConfigurations().
		Update(ctx, valHook, metav1.UpdateOptions{}); err != nil {
		return errors.Wrap(err, "failed patching admissionregistration.k8s.io/v1beta1 validating webhook")
	}
	log.Debug("patched admissionregistration.k8s.io/v1beta1 validating hook")

	return nil
}

func (k8s *K8s) patchValidatingWebhookConfigurationV1(
	ctx context.Context,
	configurationNames string,
	ca []byte,
	failurePolicy *admissionv1.FailurePolicyType,
) error {
	valHook, err := k8s.clientset.
		AdmissionregistrationV1().
		ValidatingWebhookConfigurations().
		Get(ctx, configurationNames, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed getting admissionregistration.k8s.io/v1 validating webhook")
	}

	for i := range valHook.Webhooks {
		h := &valHook.Webhooks[i]
		h.ClientConfig.CABundle = ca
		if *failurePolicy != "" {
			h.FailurePolicy = failurePolicy
		}
	}

	if _, err = k8s.clientset.AdmissionregistrationV1().
		ValidatingWebhookConfigurations().
		Update(ctx, valHook, metav1.UpdateOptions{}); err != nil {
		return errors.Wrap(err, "failed patching admissionregistration.k8s.io/v1 validating webhook")
	}
	log.Debug("patched admissionregistration.k8s.io/v1 validating hook")

	return nil
}

func (k8s *K8s) patchMutatingWebhookConfiguration(
	ctx context.Context,
	version AdmissionRegistrationVersion,
	configurationNames string,
	ca []byte,
	failurePolicy string,
) error {
	switch version {
	case admissionRegistrationV1beta1:
		failurePolicyV1beta1 := admissionv1beta1.FailurePolicyType(failurePolicy)
		return k8s.patchMutatingWebhookConfigurationV1beta1(ctx, configurationNames, ca, &failurePolicyV1beta1)
	case admissionRegistrationV1:
		failurePolicyV1 := admissionv1.FailurePolicyType(failurePolicy)
		return k8s.patchMutatingWebhookConfigurationV1(ctx, configurationNames, ca, &failurePolicyV1)
	default:
		return errors.Errorf("invalid admissionregistration.k8s.io version: %s", version)
	}
}

func (k8s *K8s) patchMutatingWebhookConfigurationV1beta1(
	ctx context.Context,
	configurationNames string,
	ca []byte,
	failurePolicy *admissionv1beta1.FailurePolicyType,
) error {
	mutHook, err := k8s.clientset.
		AdmissionregistrationV1beta1().
		MutatingWebhookConfigurations().
		Get(ctx, configurationNames, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed getting admissionregistration.k8s.io/v1beta1 mutating webhook")
	}

	for i := range mutHook.Webhooks {
		h := &mutHook.Webhooks[i]
		h.ClientConfig.CABundle = ca
		if *failurePolicy != "" {
			h.FailurePolicy = failurePolicy
		}
	}

	if _, err = k8s.clientset.AdmissionregistrationV1beta1().
		MutatingWebhookConfigurations().
		Update(ctx, mutHook, metav1.UpdateOptions{}); err != nil {
		return errors.Wrap(err, "failed patching admissionregistration.k8s.io/v1beta1 mutating webhook")
	}
	log.Debug("patched admissionregistration.k8s.io/v1beta1 mutating hook")

	return nil
}

func (k8s *K8s) patchMutatingWebhookConfigurationV1(
	ctx context.Context,
	configurationNames string,
	ca []byte,
	failurePolicy *admissionv1.FailurePolicyType,
) error {
	mutHook, err := k8s.clientset.
		AdmissionregistrationV1().
		MutatingWebhookConfigurations().
		Get(ctx, configurationNames, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed getting admissionregistration.k8s.io/v1 mutating webhook")
	}

	for i := range mutHook.Webhooks {
		h := &mutHook.Webhooks[i]
		h.ClientConfig.CABundle = ca
		if *failurePolicy != "" {
			h.FailurePolicy = failurePolicy
		}
	}

	if _, err = k8s.clientset.AdmissionregistrationV1().
		MutatingWebhookConfigurations().
		Update(ctx, mutHook, metav1.UpdateOptions{}); err != nil {
		return errors.Wrap(err, "failed patching admissionregistration.k8s.io/v1 mutating webhook")
	}
	log.Debug("patched admissionregistration.k8s.io/v1 mutating hook")

	return nil
}

// GetCaFromSecret will check for the presence of a secret. If it exists, will return the content of the
// "ca" from the secret, otherwise will return nil.
func (k8s *K8s) GetCaFromSecret(secretName string, namespace string, caName string) ([]byte, error) {
	log.Debugf("getting secret '%s' in namespace '%s'", secretName, namespace)
	secret, err := k8s.clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.WithField("err", err).Infof("secret %s/%s does not exist", namespace, secretName)
			return nil, nil
		}
		log.WithField("err", err).Fatal("error getting secret")
	}

	data := secret.Data[caName]
	if data == nil {
		return nil, errors.Errorf("secret %s/%s does not contain '%s' key", namespace, secretName, caName)
	}
	log.Debug("got secret")
	return data, nil
}

// SaveCertsToSecret saves the provided ca, cert and key into a secret in the specified namespace.
func (k8s *K8s) SaveCertsToSecret(ctx context.Context, secretName, namespace, caName, certName, keyName string, ca, cert, key []byte) error {
	log.Debugf("saving to secret '%s' in namespace '%s'", secretName, namespace)
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		Data: map[string][]byte{caName: ca, certName: cert, keyName: key},
	}

	log.Debug("saving secret")
	_, err := k8s.clientset.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed creating secret %s/%s", namespace, secretName)
	}
	log.Debug("saved secret")

	return nil
}
