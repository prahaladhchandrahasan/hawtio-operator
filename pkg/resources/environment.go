package resources

import (
	"path"
	"strings"

	corev1 "k8s.io/api/core/v1"

	hawtiov1alpha1 "github.com/hawtio/hawtio-operator/pkg/apis/hawtio/v1alpha1"
)

const (
	HawtioTypeEnvVar                = "HAWTIO_ONLINE_MODE"
	HawtioNamespaceEnvVar           = "HAWTIO_ONLINE_NAMESPACE"
	HawtioOAuthClientEnvVar         = "HAWTIO_OAUTH_CLIENT_ID"
	HawtioRbacEnvVar                = "HAWTIO_ONLINE_RBAC_ACL"
	HawtioGatewayEnvVar             = "HAWTIO_ONLINE_GATEWAY"
	OpenShiftClusterVersionEnvVar   = "OPENSHIFT_CLUSTER_VERSION"
	OpenShiftWebConsoleUrlEnvVar    = "OPENSHIFT_WEB_CONSOLE_URL"
	NginxClientBodyBufferSize       = "NGINX_CLIENT_BODY_BUFFER_SIZE"
	NginxProxyBuffers               = "NGINX_PROXY_BUFFERS"
	NginxSubrequestOutputBufferSize = "NGINX_SUBREQUEST_OUTPUT_BUFFER_SIZE"
)

func envVarsForHawtio(deploymentType hawtiov1alpha1.HawtioDeploymentType, name string) []corev1.EnvVar {
	oauthClientId := name
	if deploymentType == hawtiov1alpha1.ClusterHawtioDeploymentType {
		// Pin to a known name for cluster-wide OAuthClient
		oauthClientId = OAuthClientName
	}

	envVars := []corev1.EnvVar{
		{
			Name:  HawtioTypeEnvVar,
			Value: strings.ToLower(string(deploymentType)),
		},
		{
			Name:  HawtioOAuthClientEnvVar,
			Value: oauthClientId,
		},
	}

	if deploymentType == hawtiov1alpha1.NamespaceHawtioDeploymentType {
		envVars = append(envVars, corev1.EnvVar{
			Name: HawtioNamespaceEnvVar,
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "metadata.namespace",
				},
			},
		})
	}

	return envVars
}

func envVarsForOpenshift4(openShiftVersion string, openShiftConsoleURL string) []corev1.EnvVar {
	envVars := []corev1.EnvVar{
		{
			Name:  HawtioGatewayEnvVar,
			Value: "true",
		},
		{
			Name:  OpenShiftClusterVersionEnvVar,
			Value: openShiftVersion,
		},
		{
			Name:  OpenShiftWebConsoleUrlEnvVar,
			Value: openShiftConsoleURL,
		},
	}
	return envVars
}

func envVarForRBAC(rbacConfigMapName string) corev1.EnvVar {
	value := ""
	if rbacConfigMapName != "" {
		value = path.Join(rbacConfigMapVolumeMountPath, RBACConfigMapKey)
	}

	return corev1.EnvVar{
		Name:  HawtioRbacEnvVar,
		Value: value,
	}
}

func envVarsForNginx(nginx hawtiov1alpha1.HawtioNginx) []corev1.EnvVar {
	var envVars []corev1.EnvVar
	if nginx.ClientBodyBufferSize != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  NginxClientBodyBufferSize,
			Value: nginx.ClientBodyBufferSize,
		})
	}
	if nginx.ProxyBuffers != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  NginxProxyBuffers,
			Value: nginx.ProxyBuffers,
		})
	}
	if nginx.SubrequestOutputBufferSize != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  NginxSubrequestOutputBufferSize,
			Value: nginx.SubrequestOutputBufferSize,
		})
	}
	return envVars
}
