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
	//-> ../fixtures/format/unformatted/config.json
	//
	//@@ -11 +11,2 @@
	//-{
	//+    {
	//+      "slug": "one",
	//@@ -13 +13,0 @@
	//-      "slug": "one",
	//@@ -14,0 +15 @@
	//+      "auto_approve": true,
	//@@ -17 +17,0 @@
	//-      "auto_approve": true,
	//@@ -19,4 +19,4 @@
	//-            "Control-flow (conditionals)",
	//-            "Logic",
	//-            "Booleans",
	//-            "Integers"
	//+        "booleans",
	//+        "control_flow_conditionals",
	//+        "integers",
	//+        "logic"
	//@@ -25,0 +26 @@
	//+      "slug": "two",
	//@@ -27 +27,0 @@
	//-      "slug": "two",
	//@@ -32,5 +32,8 @@
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
	//  ../fixtures/format/unformatted/config.json
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
	//@@ -11 +11,2 @@
	//-{
	//+    {
	//+      "slug": "one",
	//@@ -13 +13,0 @@
	//-      "slug": "one",
	//@@ -14,0 +15 @@
	//+      "auto_approve": true,
	//@@ -17 +17,0 @@
	//-      "auto_approve": true,
	//@@ -19,4 +19,4 @@
	//-            "Control-flow (conditionals)",
	//-            "Logic",
	//-            "Booleans",
	//-            "Integers"
	//+        "booleans",
	//+        "control_flow_conditionals",
	//+        "integers",
	//+        "logic"
	//@@ -25,0 +26 @@
	//+      "slug": "two",
	//@@ -27 +27,0 @@
	//-      "slug": "two",
	//@@ -32,5 +32,8 @@
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
	//  ../fixtures/format/unformatted/config.json
	//../fixtures/format/unformatted/config/maintainers.json
}
