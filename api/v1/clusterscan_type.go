type ClusterScanSpec struct {
    Schedule string json:"schedule,omitempty"
    JobTemplate batchv1.JobSpec json:"jobTemplate"
}

type ClusterScanStatus struct {
    LastRunTime *metav1.Time json:"lastRunTime,omitempty"
    Conditions []metav1.Condition json:"conditions,omitempty"
}

