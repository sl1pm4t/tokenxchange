package credentials

import (
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"time"
)

type ExecCredential struct {
	metav1.TypeMeta

	Status *ExecCredentialStatus `json:"status"`
}

type ExecCredentialStatus struct {
	ClientCertificateData string       `json:"clientCertificateData"`
	ExpirationTimestamp   *metav1.Time `json:"expirationTimestamp"`
	Token                 string       `json:"token"`
}

func FormatExecCredential(token string, expiration time.Time) string {
	expirationTimestamp := metav1.NewTime(expiration)
	ec := ExecCredential{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ExecCredential",
			APIVersion: "client.authentication.k8s.io/v1beta1",
		},
		Status: &ExecCredentialStatus{
			ExpirationTimestamp: &expirationTimestamp,
			Token:               token,
		},
	}

	return MarshalExecCredential(ec)
}

func MarshalExecCredential(ec ExecCredential) string {
	b, err := json.MarshalIndent(&ec, "", "  ")
	if err != nil {
		fmt.Printf("could not marshal ExecCredentials: %s", err)
		os.Exit(1)
	}
	return string(b)
}
