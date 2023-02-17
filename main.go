package main

import (
	"os"

	aw "github.com/deanishe/awgo"
)

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

func run() {
	prompts, err := GetPrompts()
	if err != nil {
		wf.FatalError(err)
	}

	for _, prompt := range prompts {
		wf.NewItem(prompt.Act).
			Subtitle(prompt.Prompt).
			Var("prompt", prompt.Prompt).
			Icon(&aw.Icon{Value: "./icon.png"}).
			Valid(true)
	}

	wf.Filter(os.Args[1])

	wf.WarnEmpty("No Chat-GPT Prompts.", "Try different prompts.")

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
