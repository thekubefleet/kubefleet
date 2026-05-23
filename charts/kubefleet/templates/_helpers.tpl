{{- define "kubefleet.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "kubefleet.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s" (include "kubefleet.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "kubefleet.namespace" -}}
{{- default .Release.Namespace .Values.namespaceOverride -}}
{{- end -}}

{{- define "kubefleet.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "kubefleet.labels" -}}
helm.sh/chart: {{ include "kubefleet.chart" . }}
app.kubernetes.io/name: {{ include "kubefleet.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "kubefleet.agentName" -}}
{{- printf "%s-agent" (include "kubefleet.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "kubefleet.dashboardName" -}}
{{- printf "%s-dashboard" (include "kubefleet.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "kubefleet.agentServiceAccountName" -}}
{{- if .Values.agent.serviceAccount.create -}}
{{- default (include "kubefleet.agentName" .) .Values.agent.serviceAccount.name -}}
{{- else -}}
{{- required "agent.serviceAccount.name is required when agent.serviceAccount.create is false" .Values.agent.serviceAccount.name -}}
{{- end -}}
{{- end -}}
