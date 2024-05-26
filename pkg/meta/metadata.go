package meta

import "image"

type raw struct {
	comment *string
	image   image.Image
}

type Metadata struct {
	Comment        *Comment `json:"Comment,omitempty"`
	Description    string   `json:"Description"`
	GenerationTime *string  `json:"Generation time,omitempty"`
	Software       string   `json:"Software"`
	Source         string   `json:"Source"`

	raw *raw
}

type Comment struct {
	Prompt                                string   `json:"prompt"`
	Steps                                 int64    `json:"steps"`
	Height                                int64    `json:"height"`
	Width                                 int64    `json:"width"`
	Scale                                 float64  `json:"scale"`
	UncondScale                           float64  `json:"uncond_scale"`
	CFGRescale                            *float64 `json:"cfg_rescale,omitempty"`
	Seed                                  int64    `json:"seed"`
	NSamples                              int64    `json:"n_samples"`
	HideDebugOverlay                      *bool    `json:"hide_debug_overlay,omitempty"`
	NoiseSchedule                         *string  `json:"noise_schedule,omitempty"`
	LegacyV3Extend                        *bool    `json:"legacy_v3_extend,omitempty"`
	ReferenceInformationExtractedMultiple []any    `json:"reference_information_extracted_multiple"`
	ReferenceStrengthMultiple             []any    `json:"reference_strength_multiple"`
	Sampler                               string   `json:"sampler"`
	ControlnetStrength                    float64  `json:"controlnet_strength"`
	ControlnetModel                       *string  `json:"controlnet_model"`
	DynamicThresholding                   bool     `json:"dynamic_thresholding"`
	DynamicThresholdingPercentile         float64  `json:"dynamic_thresholding_percentile"`
	DynamicThresholdingMimicScale         float64  `json:"dynamic_thresholding_mimic_scale"`
	Sm                                    bool     `json:"sm"`
	SmDyn                                 bool     `json:"sm_dyn"`
	SkipCFGBelowSigma                     float64  `json:"skip_cfg_below_sigma"`
	LoraUnetWeights                       any      `json:"lora_unet_weights"`
	LoraClipWeights                       any      `json:"lora_clip_weights"`
	Strength                              *float64 `json:"strength,omitempty"`
	Noise                                 *float64 `json:"noise,omitempty"`
	ExtraNoiseSeed                        *int64   `json:"extra_noise_seed,omitempty"`
	Legacy                                *bool    `json:"legacy,omitempty"`
	Uc                                    string   `json:"uc"`
	RequestType                           string   `json:"request_type"`
	SignedHash                            *string  `json:"signed_hash,omitempty"`
}
