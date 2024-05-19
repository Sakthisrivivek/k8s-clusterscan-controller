type ClusterScanSpec struct {
	Schedule   string               `json:"schedule,omitempty"`
	JobTemplate batchv1.JobTemplateSpec `json:"jobTemplate,omitempty"`
  }
  
  type ClusterScanStatus struct {
	LastRun string `json:"lastRun,omitempty"`
	Status  string `json:"status,omitempty"`
  }
  
  type ClusterScan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec   ClusterScanSpec   `json:"spec,omitempty"`
	Status ClusterScanStatus `json:"status,omitempty"`
  }
  
  
  type ClusterScanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterScan `json:"items"`
  }
  