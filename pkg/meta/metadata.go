package meta

type Metadata struct {
	Comment        *Comment `json:"Comment,omitempty"`
	Description    string   `json:"Description"`
	GenerationTime *string  `json:"Generation time,omitempty"`
	Software       string   `json:"Software"`
	Source         string   `json:"Source"`
}

type Comment struct {
	CFGRescale                            *float64 `json:"cfg_rescale,omitempty"`
	ControlnetModel                       *string  `json:"controlnet_model"`
	ControlnetStrength                    float64  `json:"controlnet_strength"`
	DynamicThresholding                   bool     `json:"dynamic_thresholding"`
	DynamicThresholdingMimicScale         float64  `json:"dynamic_thresholding_mimic_scale"`
	DynamicThresholdingPercentile         float64  `json:"dynamic_thresholding_percentile"`
	Height                                int64    `json:"height"`
	HideDebugOverlay                      *bool    `json:"hide_debug_overlay,omitempty"`
	LegacyV3Extend                        *bool    `json:"legacy_v3_extend,omitempty"`
	LoraClipWeights                       any      `json:"lora_clip_weights"`
	LoraUnetWeights                       any      `json:"lora_unet_weights"`
	NSamples                              int64    `json:"n_samples"`
	NoiseSchedule                         *string  `json:"noise_schedule,omitempty"`
	Prompt                                string   `json:"prompt"`
	ReferenceInformationExtractedMultiple []any    `json:"reference_information_extracted_multiple,omitempty"`
	ReferenceStrengthMultiple             []any    `json:"reference_strength_multiple,omitempty"`
	RequestType                           string   `json:"request_type"`
	Sampler                               string   `json:"sampler"`
	Scale                                 float64  `json:"scale"`
	Seed                                  int64    `json:"seed"`
	SignedHash                            string   `json:"signed_hash"`
	SkipCFGBelowSigma                     float64  `json:"skip_cfg_below_sigma"`
	Sm                                    bool     `json:"sm"`
	SmDyn                                 bool     `json:"sm_dyn"`
	Steps                                 int64    `json:"steps"`
	Uc                                    string   `json:"uc"`
	UncondScale                           float64  `json:"uncond_scale"`
	Width                                 int64    `json:"width"`
	ExtraNoiseSeed                        *int64   `json:"extra_noise_seed,omitempty"`
	Legacy                                *bool    `json:"legacy,omitempty"`
	Noise                                 *float64 `json:"noise,omitempty"`
	Strength                              *float64 `json:"strength,omitempty"`
}
