package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/exercism/configlet/ui"
)

func ExampleFormatVerboseFlag() {
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

	fmtVerbose = true
	fmtTest = false

	runFmt("../fixtures/format/unformatted/", unformattedDir)

	// Output:
	//	-> ../fixtures/format/unformatted/config.json
	//
	//@@ -16 +16,2 @@
	//-{
	//+    {
	//+      "slug": "one",
	//@@ -18 +18,0 @@
	//-      "slug": "one",
	//@@ -19,0 +20 @@
	//+      "auto_approve": true,
	//@@ -22 +22,0 @@
	//-      "auto_approve": true,
	//@@ -24,4 +24,4 @@
	//-            "Control-flow (conditionals)",
	//-            "Logic",
	//-            "Booleans",
	//-            "Integers"
	//+        "booleans",
	//+        "control_flow_conditionals",
	//+        "integers",
	//+        "logic"
	//@@ -30,0 +31 @@
	//+      "slug": "two",
	//@@ -32 +32,0 @@
	//-      "slug": "two",
	//@@ -37,5 +37,8 @@
	//-        "Time",
	//-        "Mathematics",
	//-        "Text formatting",
	//-        "Equality"
	//-      ]}]}
	//+        "equality",
	//+        "mathematics",
	//+        "text_formatting",
	//+        "time"
	//+      ]
	//+    }
	//+  ]
	//+}
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
	/* ../fixtures/format/unformatted/config.json*/
	//../fixtures/format/unformatted/config/maintainers.json
}

func ExampleFormatTestFlag() {
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

	fmtVerbose = false
	fmtTest = true

	runFmt("../fixtures/format/unformatted/", unformattedDir)

	// Output:
	//-> ../fixtures/format/unformatted/config.json
	//
	//@@ -16 +16,2 @@
	//-{
	//+    {
	//+      "slug": "one",
	//@@ -18 +18,0 @@
	//-      "slug": "one",
	//@@ -19,0 +20 @@
	//+      "auto_approve": true,
	//@@ -22 +22,0 @@
	//-      "auto_approve": true,
	//@@ -24,4 +24,4 @@
	//-            "Control-flow (conditionals)",
	//-            "Logic",
	//-            "Booleans",
	//-            "Integers"
	//+        "booleans",
	//+        "control_flow_conditionals",
	//+        "integers",
	//+        "logic"
	//@@ -30,0 +31 @@
	//+      "slug": "two",
	//@@ -32 +32,0 @@
	//-      "slug": "two",
	//@@ -37,5 +37,8 @@
	//-        "Time",
	//-        "Mathematics",
	//-        "Text formatting",
	//-        "Equality"
	//-      ]}]}
	//+        "equality",
	//+        "mathematics",
	//+        "text_formatting",
	//+        "time"
	//+      ]
	//+    }
	//+  ]
	//+}
	//
	//-> ../fixtures/format/unformatted/config/maintainers.json
	//
	//@@ -2,2 +2,2 @@
	//-	"docs_url": "http://docs.example.com",
	//-	"maintainers": []
	//+  "docs_url": "http://docs.example.com",
	//+  "maintainers": []
	//
	//-> no changes were made to:
	/* ../fixtures/format/unformatted/config.json*/
	//../fixtures/format/unformatted/config/maintainers.json
}
