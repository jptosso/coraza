package engine

import (
	"github.com/corazawaf/coraza/v3/testing/profile"
)

var _ = profile.RegisterProfile(profile.Profile{
	Meta: profile.ProfileMeta{
		Author:      "sts",
		Description: "Test if disruptive actions trigger an interruption",
		Enabled:     true,
		Name:        "disruptive_actions.yaml",
	},
	Tests: []profile.ProfileTest{
		{
			Title: "disruptive_actions",
			Stages: []profile.ProfileStage{
				// Phase 1
				{
					Input: profile.ProfileStageInput{
						URI: "/redirect1",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{1},
						Interruption: &profile.ExpectedInterruption{
							Status: 302,
							Data:   "https://www.example.com",
							RuleID: 1,
							Action: "redirect",
						},
					},
				},
				{
					Input: profile.ProfileStageInput{
						URI: "/deny1",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{2},
						Interruption: &profile.ExpectedInterruption{
							Status: 500,
							Data:   "",
							RuleID: 2,
							Action: "deny",
						},
					},
				},
				{
					Input: profile.ProfileStageInput{
						URI: "/drop1",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{3},
						Interruption: &profile.ExpectedInterruption{
							Status: 0,
							Data:   "",
							RuleID: 3,
							Action: "drop",
						},
					},
				},
				// Phase 2
				{
					Input: profile.ProfileStageInput{
						URI: "/redirect2",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{21},
						Interruption: &profile.ExpectedInterruption{
							Status: 302,
							Data:   "https://www.example.com",
							RuleID: 21,
							Action: "redirect",
						},
					},
				},
				{
					Input: profile.ProfileStageInput{
						URI: "/deny2",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{22},
						Interruption: &profile.ExpectedInterruption{
							Status: 500,
							Data:   "",
							RuleID: 22,
							Action: "deny",
						},
					},
				},
				{
					Input: profile.ProfileStageInput{
						URI: "/drop2",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{23},
						Interruption: &profile.ExpectedInterruption{
							Status: 0,
							Data:   "",
							RuleID: 23,
							Action: "drop",
						},
					},
				},
				// Phase 3
				{
					Input: profile.ProfileStageInput{
						URI: "/redirect3",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{31},
						Interruption: &profile.ExpectedInterruption{
							Status: 302,
							Data:   "https://www.example.com",
							RuleID: 31,
							Action: "redirect",
						},
					},
				},
				{
					Input: profile.ProfileStageInput{
						URI: "/deny3",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{32},
						Interruption: &profile.ExpectedInterruption{
							Status: 500,
							Data:   "",
							RuleID: 32,
							Action: "deny",
						},
					},
				},
				{
					Input: profile.ProfileStageInput{
						URI: "/drop3",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{33},
						Interruption: &profile.ExpectedInterruption{
							Status: 0,
							Data:   "",
							RuleID: 33,
							Action: "drop",
						},
					},
				},
				// Phase 4
				{
					Input: profile.ProfileStageInput{
						URI: "/redirect4",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{41},
						Interruption: &profile.ExpectedInterruption{
							Status: 302,
							Data:   "https://www.example.com",
							RuleID: 41,
							Action: "redirect",
						},
					},
				},
				{
					Input: profile.ProfileStageInput{
						URI: "/deny4",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{42},
						Interruption: &profile.ExpectedInterruption{
							Status: 500,
							Data:   "",
							RuleID: 42,
							Action: "deny",
						},
					},
				},
				{
					Input: profile.ProfileStageInput{
						URI: "/drop4",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{43},
						Interruption: &profile.ExpectedInterruption{
							Status: 0,
							Data:   "",
							RuleID: 43,
							Action: "drop",
						},
					},
				},
				// Phase 5
				{
					Input: profile.ProfileStageInput{
						URI: "/redirect5",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{51},
						Interruption: &profile.ExpectedInterruption{
							Status: 302,
							Data:   "https://www.example.com",
							RuleID: 51,
							Action: "redirect",
						},
					},
				},
				{
					Input: profile.ProfileStageInput{
						URI: "/deny5",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{52},
						Interruption: &profile.ExpectedInterruption{
							Status: 500,
							Data:   "",
							RuleID: 52,
							Action: "deny",
						},
					},
				},
				{
					Input: profile.ProfileStageInput{
						URI: "/drop5",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{53},
						Interruption: &profile.ExpectedInterruption{
							Status: 0,
							Data:   "",
							RuleID: 53,
							Action: "drop",
						},
					},
				},
				{
					Input: profile.ProfileStageInput{
						URI: "/default/block",
					},
					Output: profile.ExpectedOutput{
						TriggeredRules: []int{103},
						LogContains:    "WOOOP_BLOCKED_BY_CORAZA_TEST",
						Interruption: &profile.ExpectedInterruption{
							Status: 501,
							Data:   "",
							RuleID: 103,
							Action: "deny",
						},
					},
				},
			},
		},
	},
	Rules: `
SecRule REQUEST_URI "/redirect1$" "phase:1,id:1,log,status:302,redirect:https://www.example.com"
SecRule REQUEST_URI "/deny1$" "phase:1,id:2,log,status:500,deny"
SecRule REQUEST_URI "/drop1$" "phase:1,id:3,log,drop"

SecRule REQUEST_URI "/redirect2$" "phase:2,id:21,log,status:302,redirect:https://www.example.com"
SecRule REQUEST_URI "/deny2$" "phase:2,id:22,log,status:500,deny"
SecRule REQUEST_URI "/drop2$" "phase:2,id:23,log,drop"

SecRule REQUEST_URI "/redirect3$" "phase:3,id:31,log,status:302,redirect:https://www.example.com"
SecRule REQUEST_URI "/deny3$" "phase:3,id:32,log,status:500,deny"
SecRule REQUEST_URI "/drop3$" "phase:3,id:33,log,drop"

SecRule REQUEST_URI "/redirect4$" "phase:4,id:41,log,status:302,redirect:https://www.example.com"
SecRule REQUEST_URI "/deny4$" "phase:4,id:42,log,status:500,deny"
SecRule REQUEST_URI "/drop4$" "phase:4,id:43,log,drop"

SecRule REQUEST_URI "/redirect5$" "phase:5,id:51,log,status:302,redirect:https://www.example.com"
SecRule REQUEST_URI "/deny5$" "phase:5,id:52,log,status:500,deny"
SecRule REQUEST_URI "/drop5$" "phase:5,id:53,log,drop"

SecDefaultAction phase:2,deny,status:501,log,logdata:'WOOOP_BLOCKED_BY_CORAZA_TEST'
SecRule REQUEST_URI "/default/block" "id:103,block"
`,
})
