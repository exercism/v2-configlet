package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/exercism/configlet/ui"
)

func ExampleFormat() {
	oldOut := ui.Out
	oldErrOut := ui.ErrOut
	ui.Out = os.Stdout
	ui.ErrOut = os.Stderr
	defer func() {
		ui.Out = oldOut
		ui.ErrOut = oldErrOut
	}()

	unformattedDir, err := ioutil.TempDir("", "unformatted")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(unformattedDir)

	runFmt("../fixtures/format/unformatted/", unformattedDir, true)

	// Output:
	// -> ../fixtures/format/unformatted/config.json
	//
	// @@ -11 +11,2 @@
	// -{
	// +    {
	// +      "slug": "one",
	// @@ -13 +13,0 @@
	// -      "slug": "one",
	// @@ -18,4 +18,4 @@
	// -            "Control-flow (conditionals)",
	// -            "Logic",
	// -            "Booleans",
	// -            "Integers"
	// +        "booleans",
	// +        "control_flow_conditionals",
	// +        "integers",
	// +        "logic"
	// @@ -24,0 +25 @@
	// +      "slug": "two",
	// @@ -26 +26,0 @@
	// -      "slug": "two",
	// @@ -31,5 +31,8 @@
	// -        "Time",
	// -        "Mathematics",
	// -        "Text formatting",
	// -        "Equality"
	// -      ]}]}
	// +        "equality",
	// +        "mathematics",
	// +        "text_formatting",
	// +        "time"
	// +      ]
	// +    }
	// +  ]
	// +}
	//
	//-> ../fixtures/format/unformatted/config/maintainers.json
	//
	//@@ -2,2 +2,2 @@
	//-	"docs_url": "http://docs.example.com",
	//-	"maintainers": []
	//+  "docs_url": "http://docs.example.com",
	//+  "maintainers": []
	//
	//-> changes made to:
	//  ../fixtures/format/unformatted/config.json
	// ../fixtures/format/unformatted/config/maintainers.json
}
