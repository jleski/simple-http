package main

import (
	"context"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/networking/v1"
	rancher2 "github.com/pulumi/pulumi-rancher2/sdk/v3/go/rancher2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/ffuser"
	"log"
	"time"
)

func checkFeatureFlag(user ffuser.User, key string) func() bool {
	return func() bool {
		if ok, _ := ffclient.BoolVariation(key, user, false); ok {
			return true
		} else {
			return false
		}
	}
}

func main() {
	// Init ffclient with a file retriever.
	err := ffclient.Init(ffclient.Config{
		PollingInterval: 10 * time.Second,
		//		Logger:          log.New(os.Stdout, "FLAGS ", 0),
		Context: context.Background(),
		Retriever: &ffclient.FileRetriever{
			Path: "flags.yaml",
		},
	})

	// Check init errors.
	if err != nil {
		log.Fatal(err)
	}

	// defer closing ffclient
	defer ffclient.Close()

	// get functions for retrieving feature flag states
	flagDeployUsingPulumi := checkFeatureFlag(ffuser.User{}, "deployAppUsingPulumi")
	flagConfigureNamespaceUsingPulumi := checkFeatureFlag(ffuser.User{}, "configureNamespaceUsingPulumi")
	pulumi.Run(func(ctx *pulumi.Context) error {
		appNameSpace := "simple-http"
		projectId := "c-9dj55:p-lvx5w"
		if flagConfigureNamespaceUsingPulumi() {
			_, err := rancher2.NewNamespace(ctx, appNameSpace, &rancher2.NamespaceArgs{
				ContainerResourceLimit: &rancher2.NamespaceContainerResourceLimitArgs{
					LimitsCpu:      pulumi.String("20m"),
					LimitsMemory:   pulumi.String("20Mi"),
					RequestsCpu:    pulumi.String("1m"),
					RequestsMemory: pulumi.String("1Mi"),
				},
				Description: pulumi.String("simple-http demo namespace"),
				ProjectId:   pulumi.String(projectId), // format: cluster_id:project_id, check it from Rancher API
				ResourceQuota: &rancher2.NamespaceResourceQuotaArgs{
					Limit: &rancher2.NamespaceResourceQuotaLimitArgs{
						LimitsCpu:       pulumi.String("100m"),
						LimitsMemory:    pulumi.String("100Mi"),
						RequestsStorage: pulumi.String("1Gi"),
					},
				},
			})
			if err != nil {
				return err
			}

		}

		if flagDeployUsingPulumi() {

			//var rancherNameSpace pulumi.StringOutput
			// simple http service
			appName := "simple-http"
			numberOfReplicas := 3
			simpleHttpServicePort := 8080
			containerImage := "jledev.azurecr.io/simple-http:latest"
			appLabels := pulumi.StringMap{
				"app": pulumi.String(appName),
			}
			_, err := appsv1.NewDeployment(ctx, appName, &appsv1.DeploymentArgs{
				Metadata: metav1.ObjectMetaArgs{
					Name:      pulumi.String(appName),
					Namespace: pulumi.String(appNameSpace),
				},
				Spec: appsv1.DeploymentSpecArgs{
					Selector: &metav1.LabelSelectorArgs{
						MatchLabels: appLabels,
					},
					Replicas: pulumi.Int(numberOfReplicas),
					Template: &corev1.PodTemplateSpecArgs{
						Metadata: &metav1.ObjectMetaArgs{
							Labels: appLabels,
						},
						Spec: &corev1.PodSpecArgs{
							Containers: corev1.ContainerArray{
								corev1.ContainerArgs{
									Name:  pulumi.String(appName),
									Image: pulumi.String(containerImage),
								},
							},
							ImagePullSecrets: corev1.LocalObjectReferenceArray{
								corev1.LocalObjectReferenceArgs{
									Name: pulumi.String("jledev-azurecr-cred"),
								},
							},
						},
					},
				},
			})
			if err != nil {
				return err
			}
			frontendServiceType := "ClusterIP"

			frontend, err := corev1.NewService(ctx, appName, &corev1.ServiceArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Name:      pulumi.Sprintf("%s-service", appName),
					Labels:    appLabels,
					Namespace: pulumi.String(appNameSpace),
				},
				Spec: &corev1.ServiceSpecArgs{
					Type: pulumi.String(frontendServiceType),
					Ports: &corev1.ServicePortArray{
						corev1.ServicePortArgs{
							Port:       pulumi.Int(simpleHttpServicePort),
							TargetPort: pulumi.Int(simpleHttpServicePort),
							Protocol:   pulumi.String("TCP"),
						},
					},
					Selector: appLabels,
				},
			})

			if err != nil {
				return err
			}

			frontendIngress, err := v1.NewIngress(ctx, appName, &v1.IngressArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Name:      pulumi.Sprintf("%s-ingress", appName),
					Labels:    appLabels,
					Namespace: pulumi.String(appNameSpace),
				},
				Spec: &v1.IngressSpecArgs{
					Rules: v1.IngressRuleArray{
						v1.IngressRuleArgs{
							Host: pulumi.String("simple-http.apps-aks.alusta.cloud"),
							Http: v1.HTTPIngressRuleValueArgs{
								Paths: v1.HTTPIngressPathArray{
									v1.HTTPIngressPathArgs{
										Path:     pulumi.String("/"),
										PathType: pulumi.String("Prefix"),
										Backend: v1.IngressBackendArgs{
											Service: v1.IngressServiceBackendArgs{
												Name: pulumi.String("simple-http-service"),
												Port: v1.ServiceBackendPortArgs{
													Number: pulumi.Int(8080),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			})

			if err != nil {
				return err
			}

			// Export the public IP
			ctx.Export("frontendIp", frontend.Spec.ApplyT(func(spec *corev1.ServiceSpec) *string {
				return spec.ClusterIP
			}))
			ctx.Export("loadBalancerAddress", frontendIngress.Spec.ApplyT(func(spec *v1.IngressSpec) *string {
				return spec.Rules[0].Host
			}))

		}

		return nil
	})
}
