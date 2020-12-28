package steps

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
	"math/big"
	"time"
)

type SetupEncryption struct {
	environment                    models.Environment
	executor                       *executor.Executor
	nodeRepository                 components.NodeRepository
	instanceRepository             components.InstanceRepository
	certificateAuthorityRepository components.CertificateAuthorityRepository
}

func (s *SetupEncryption) Do(nodeId uuid.UUID) error {
	node, err := s.nodeRepository.FindById(nodeId)
	if err != nil {
		return err
	}
	instance, err := s.instanceRepository.FindById(node.InstanceId)
	if err != nil {
		return err
	}
	rootCA, err := s.certificateAuthorityRepository.FindByProjectId(instance.ProjectId)
	if err != nil {
		return err
	}
	bundle, err := tls.X509KeyPair([]byte(rootCA.Crt), []byte(rootCA.Key))
	if err != nil {
		return err
	}
	rootCertificate, err := x509.ParseCertificate(bundle.Certificate[0])
	if err != nil {
		return err
	}
	serverCertificateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	now := time.Now()
	serverCertificateRaw, err := x509.CreateCertificate(
		rand.Reader,
		&x509.Certificate{
			SerialNumber: big.NewInt(now.Unix()),
			IPAddresses:  node.Metadata.IPs(),
			DNSNames:     node.Names(s.environment.CloudflareDomain),
			Subject: pkix.Name{
				Organization: []string{"ScaleChamp"},
			},
			NotBefore:    now,
			NotAfter:     now.AddDate(10, 0, 0),
			SubjectKeyId: []byte{1, 2, 3, 4, 6},
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:     x509.KeyUsageDigitalSignature,
		},
		rootCertificate,
		serverCertificateKey.Public(),
		bundle.PrivateKey,
	)
	if err != nil {
		return err
	}
	template := &x509.Certificate{
		IPAddresses:  node.Metadata.IPs(),
		DNSNames:     node.Names(s.environment.CloudflareDomain),
		SerialNumber: big.NewInt(now.Unix()),
		Subject: pkix.Name{
			Organization: []string{"ScaleChamp"},
		},
		NotBefore:             now,
		NotAfter:              now.AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	leafCA, err := x509.CreateCertificate(
		rand.Reader,
		template,
		template,
		key.Public(),
		key,
	)
	if err != nil {
		return err
	}
	if err := s.executor.Wait(node); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	if err := s.executor.PutBytes(node, s.pem("CERTIFICATE", serverCertificateRaw), "/tmp/cert.pem"); err != nil {
		return err
	}
	if err := s.executor.PutBytes(node, s.pem("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(serverCertificateKey)), "/tmp/key.pem"); err != nil {
		return err
	}
	if err := s.executor.PutBytes(node, s.pem("CERTIFICATE", leafCA), "/tmp/ca.pem"); err != nil {
		return err
	}
	caBundle := append(s.pem("CERTIFICATE", leafCA), rootCA.Crt...)
	if err := s.executor.PutBytes(node, caBundle, "/tmp/ca_bundle.pem"); err != nil {
		return err
	}
	cmds := []string{
		"sudo mv /tmp/cert.pem /etc/cert.pem",
		"sudo mv /tmp/key.pem /etc/key.pem",
		"sudo mv /tmp/ca.pem /etc/ca.pem",
		"sudo mv /tmp/ca_bundle.pem /etc/ca_bundle.pem",
	}
	if err := s.executor.MultiRun(node, cmds); err != nil {
		return err
	}
	if err := s.executor.Run(node, "sudo chmod 0600 /etc/cert.pem /etc/key.pem /etc/ca.pem /etc/ca_bundle.pem"); err != nil {
		return err
	}

	switch instance.Kind {
	case "pg":
		return s.executor.Run(node, "sudo chown postgres:postgres /etc/cert.pem /etc/key.pem /etc/ca.pem /etc/ca_bundle.pem")
	case "redis":
		return s.executor.Run(node, "sudo chown redis:redis /etc/cert.pem /etc/key.pem /etc/ca.pem /etc/ca_bundle.pem")
	case "keydb", "keydb-pro":
		return s.executor.Run(node, "sudo chown keydb:keydb /etc/cert.pem /etc/key.pem /etc/ca.pem /etc/ca_bundle.pem")
	case "mysql":
		if err := s.executor.Run(node, "sudo cp /etc/cert.pem /etc/key.pem /etc/ca.pem /etc/ca_bundle.pem /etc/mysql/."); err != nil {
			return err
		}
		return s.executor.Run(node, "sudo chown mysql:mysql /etc/mysql/cert.pem /etc/mysql/key.pem /etc/mysql/ca.pem /etc/mysql/ca_bundle.pem")
	default:
		return errors.New("not supported kind")
	}
}

func (s *SetupEncryption) pem(t string, b []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: t, Bytes: b})
}

func (s *SetupEncryption) Undo(uuid.UUID) error {
	return nil
}

func NewSetupEncryption(
	e models.Environment,
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	certificateAuthorityRepository components.CertificateAuthorityRepository,
) *SetupEncryption {
	return &SetupEncryption{
		e,
		executor,
		nodeRepository,
		instanceRepository,
		certificateAuthorityRepository,
	}
}

// echo `curl https://api.ipify.org`;
// docker stop 62723d5ebff0
// iptables -t nat -D POSTROUTING 1;
// curl https://api.ipify.org;
