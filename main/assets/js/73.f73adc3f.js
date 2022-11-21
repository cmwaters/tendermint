(window.webpackJsonp=window.webpackJsonp||[]).push([[73],{675:function(e,n,t){"use strict";t.r(n);var a=t(1),o=Object(a.a)({},(function(){var e=this,n=e.$createElement,t=e._self._c||n;return t("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[t("h1",{attrs:{id:"using-abci-cli"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#using-abci-cli"}},[e._v("#")]),e._v(" Using ABCI-CLI")]),e._v(" "),t("p",[e._v("To facilitate testing and debugging of ABCI servers and simple apps, we\nbuilt a CLI, the "),t("code",[e._v("abci-cli")]),e._v(", for sending ABCI messages from the command\nline.")]),e._v(" "),t("h2",{attrs:{id:"install"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#install"}},[e._v("#")]),e._v(" Install")]),e._v(" "),t("p",[e._v("Make sure you "),t("a",{attrs:{href:"https://golang.org/doc/install",target:"_blank",rel:"noopener noreferrer"}},[e._v("have Go installed"),t("OutboundLink")],1),e._v(".")]),e._v(" "),t("p",[e._v("Next, install the "),t("code",[e._v("abci-cli")]),e._v(" tool and example applications:")]),e._v(" "),t("tm-code-block",{staticClass:"codeblock",attrs:{language:"sh",base64:"Z2l0IGNsb25lIGh0dHBzOi8vZ2l0aHViLmNvbS90ZW5kZXJtaW50L3RlbmRlcm1pbnQuZ2l0CmNkIHRlbmRlcm1pbnQKbWFrZSBpbnN0YWxsX2FiY2kK"}}),e._v(" "),t("p",[e._v("Now run "),t("code",[e._v("abci-cli")]),e._v(" to see the list of commands:")]),e._v(" "),t("tm-code-block",{staticClass:"codeblock",attrs:{language:"sh",base64:"VXNhZ2U6CiAgYWJjaS1jbGkgW2NvbW1hbmRdCgpBdmFpbGFibGUgQ29tbWFuZHM6CiAgYmF0Y2ggICAgICAgICAgICBydW4gYSBiYXRjaCBvZiBhYmNpIGNvbW1hbmRzIGFnYWluc3QgYW4gYXBwbGljYXRpb24KICBjaGVja190eCAgICAgICAgIHZhbGlkYXRlIGEgdHJhbnNhY3Rpb24KICBjb21taXQgICAgICAgICAgIGNvbW1pdCB0aGUgYXBwbGljYXRpb24gc3RhdGUgYW5kIHJldHVybiB0aGUgTWVya2xlIHJvb3QgaGFzaAogIGNvbXBsZXRpb24gICAgICAgR2VuZXJhdGUgdGhlIGF1dG9jb21wbGV0aW9uIHNjcmlwdCBmb3IgdGhlIHNwZWNpZmllZCBzaGVsbAogIGNvbnNvbGUgICAgICAgICAgc3RhcnQgYW4gaW50ZXJhY3RpdmUgQUJDSSBjb25zb2xlIGZvciBtdWx0aXBsZSBjb21tYW5kcwogIGRlbGl2ZXJfdHggICAgICAgZGVsaXZlciBhIG5ldyB0cmFuc2FjdGlvbiB0byB0aGUgYXBwbGljYXRpb24KICBlY2hvICAgICAgICAgICAgIGhhdmUgdGhlIGFwcGxpY2F0aW9uIGVjaG8gYSBtZXNzYWdlCiAgaGVscCAgICAgICAgICAgICBIZWxwIGFib3V0IGFueSBjb21tYW5kCiAgaW5mbyAgICAgICAgICAgICBnZXQgc29tZSBpbmZvIGFib3V0IHRoZSBhcHBsaWNhdGlvbgogIGt2c3RvcmUgICAgICAgICAgQUJDSSBkZW1vIGV4YW1wbGUKICBwcmVwYXJlX3Byb3Bvc2FsIHByZXBhcmUgcHJvcG9zYWwKICBwcm9jZXNzX3Byb3Bvc2FsIHByb2Nlc3MgcHJvcG9zYWwKICBxdWVyeSAgICAgICAgICAgIHF1ZXJ5IHRoZSBhcHBsaWNhdGlvbiBzdGF0ZQogIHRlc3QgICAgICAgICAgICAgcnVuIGludGVncmF0aW9uIHRlc3RzCiAgdmVyc2lvbiAgICAgICAgICBwcmludCBBQkNJIGNvbnNvbGUgdmVyc2lvbgoKRmxhZ3M6CiAgICAgIC0tYWJjaSBzdHJpbmcgICAgICAgIGVpdGhlciBzb2NrZXQgb3IgZ3JwYyAoZGVmYXVsdCAmcXVvdDtzb2NrZXQmcXVvdDspCiAgICAgIC0tYWRkcmVzcyBzdHJpbmcgICAgIGFkZHJlc3Mgb2YgYXBwbGljYXRpb24gc29ja2V0IChkZWZhdWx0ICZxdW90O3RjcDovLzAuMC4wLjA6MjY2NTgmcXVvdDspCiAgLWgsIC0taGVscCAgICAgICAgICAgICAgIGhlbHAgZm9yIGFiY2ktY2xpCiAgICAgIC0tbG9nX2xldmVsIHN0cmluZyAgIHNldCB0aGUgbG9nZ2VyIGxldmVsIChkZWZhdWx0ICZxdW90O2RlYnVnJnF1b3Q7KQogIC12LCAtLXZlcmJvc2UgICAgICAgICAgICBwcmludCB0aGUgY29tbWFuZCBhbmQgcmVzdWx0cyBhcyBpZiBpdCB3ZXJlIGEgY29uc29sZSBzZXNzaW9uCgpVc2UgJnF1b3Q7YWJjaS1jbGkgW2NvbW1hbmRdIC0taGVscCZxdW90OyBmb3IgbW9yZSBpbmZvcm1hdGlvbiBhYm91dCBhIGNvbW1hbmQuCg=="}}),e._v(" "),t("h2",{attrs:{id:"kvstore-first-example"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#kvstore-first-example"}},[e._v("#")]),e._v(" KVStore - First Example")]),e._v(" "),t("p",[e._v("The "),t("code",[e._v("abci-cli")]),e._v(" tool lets us send ABCI messages to our application, to\nhelp build and debug them.")]),e._v(" "),t("p",[e._v("The most important messages are "),t("code",[e._v("deliver_tx")]),e._v(", "),t("code",[e._v("check_tx")]),e._v(", and "),t("code",[e._v("commit")]),e._v(",\nbut there are others for convenience, configuration, and information\npurposes.")]),e._v(" "),t("p",[e._v("We'll start a kvstore application, which was installed at the same time\nas "),t("code",[e._v("abci-cli")]),e._v(" above. The kvstore just stores transactions in a merkle\ntree. Its code can be found\n"),t("a",{attrs:{href:"https://github.com/tendermint/tendermint/blob/main/abci/cmd/abci-cli/abci-cli.go",target:"_blank",rel:"noopener noreferrer"}},[e._v("here"),t("OutboundLink")],1),e._v("\nand looks like the following:")]),e._v(" "),t("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"ZnVuYyBjbWRLVlN0b3JlKGNtZCAqY29icmEuQ29tbWFuZCwgYXJncyBbXXN0cmluZykgZXJyb3IgewoJbG9nZ2VyIDo9IGxvZy5OZXdUTUxvZ2dlcihsb2cuTmV3U3luY1dyaXRlcihvcy5TdGRvdXQpKQoKCS8vIENyZWF0ZSB0aGUgYXBwbGljYXRpb24gLSBpbiBtZW1vcnkgb3IgcGVyc2lzdGVkIHRvIGRpc2sKCXZhciBhcHAgdHlwZXMuQXBwbGljYXRpb24KCWlmIGZsYWdQZXJzaXN0ID09ICZxdW90OyZxdW90OyB7CgkJdmFyIGVyciBlcnJvcgoJCWZsYWdQZXJzaXN0LCBlcnIgPSBvcy5Na2RpclRlbXAoJnF1b3Q7JnF1b3Q7LCAmcXVvdDtwZXJzaXN0ZW50X2t2c3RvcmVfdG1wJnF1b3Q7KQoJCWlmIGVyciAhPSBuaWwgewoJCQlyZXR1cm4gZXJyCgkJfQoJfQoJYXBwID0ga3ZzdG9yZS5OZXdQZXJzaXN0ZW50S1ZTdG9yZUFwcGxpY2F0aW9uKGZsYWdQZXJzaXN0KQoJYXBwLigqa3ZzdG9yZS5QZXJzaXN0ZW50S1ZTdG9yZUFwcGxpY2F0aW9uKS5TZXRMb2dnZXIobG9nZ2VyLldpdGgoJnF1b3Q7bW9kdWxlJnF1b3Q7LCAmcXVvdDtrdnN0b3JlJnF1b3Q7KSkKCgkvLyBTdGFydCB0aGUgbGlzdGVuZXIKCXNydiwgZXJyIDo9IHNlcnZlci5OZXdTZXJ2ZXIoZmxhZ0FkZHJlc3MsIGZsYWdBYmNpLCBhcHApCglpZiBlcnIgIT0gbmlsIHsKCQlyZXR1cm4gZXJyCgl9CglzcnYuU2V0TG9nZ2VyKGxvZ2dlci5XaXRoKCZxdW90O21vZHVsZSZxdW90OywgJnF1b3Q7YWJjaS1zZXJ2ZXImcXVvdDspKQoJaWYgZXJyIDo9IHNydi5TdGFydCgpOyBlcnIgIT0gbmlsIHsKCQlyZXR1cm4gZXJyCgl9CgoJLy8gU3RvcCB1cG9uIHJlY2VpdmluZyBTSUdURVJNIG9yIENUUkwtQy4KCXRtb3MuVHJhcFNpZ25hbChsb2dnZXIsIGZ1bmMoKSB7CgkJLy8gQ2xlYW51cAoJCWlmIGVyciA6PSBzcnYuU3RvcCgpOyBlcnIgIT0gbmlsIHsKCQkJbG9nZ2VyLkVycm9yKCZxdW90O0Vycm9yIHdoaWxlIHN0b3BwaW5nIHNlcnZlciZxdW90OywgJnF1b3Q7ZXJyJnF1b3Q7LCBlcnIpCgkJfQoJfSkKCgkvLyBSdW4gZm9yZXZlci4KCXNlbGVjdCB7fQp9Cgo="}}),e._v(" "),t("p",[e._v("Start the application by running:")]),e._v(" "),t("tm-code-block",{staticClass:"codeblock",attrs:{language:"sh",base64:"YWJjaS1jbGkga3ZzdG9yZQo="}}),e._v(" "),t("p",[e._v("And in another terminal, run")]),e._v(" "),t("tm-code-block",{staticClass:"codeblock",attrs:{language:"sh",base64:"YWJjaS1jbGkgZWNobyBoZWxsbwphYmNpLWNsaSBpbmZvCg=="}}),e._v(" "),t("p",[e._v("You'll see something like:")]),e._v(" "),t("tm-code-block",{staticClass:"codeblock",attrs:{language:"sh",base64:"LSZndDsgZGF0YTogaGVsbG8KLSZndDsgZGF0YS5oZXg6IDY4NjU2QzZDNkYK"}}),e._v(" "),t("p",[e._v("and:")]),e._v(" "),t("tm-code-block",{staticClass:"codeblock",attrs:{language:"sh",base64:"LSZndDsgZGF0YTogeyZxdW90O3NpemUmcXVvdDs6MH0KLSZndDsgZGF0YS5oZXg6IDdCMjI3MzY5N0E2NTIyM0EzMDdECg=="}}),e._v(" "),t("p",[e._v("An ABCI application must provide two things:")]),e._v(" "),t("ul",[t("li",[e._v("a socket server")]),e._v(" "),t("li",[e._v("a handler for ABCI messages")])]),e._v(" "),t("p",[e._v("When we run the "),t("code",[e._v("abci-cli")]),e._v(" tool we open a new connection to the\napplication's socket server, send the given ABCI message, and wait for a\nresponse.")]),e._v(" "),t("p",[e._v("The server may be generic for a particular language, and we provide a\n"),t("a",{attrs:{href:"https://github.com/tendermint/tendermint/tree/main/abci/server",target:"_blank",rel:"noopener noreferrer"}},[e._v("reference implementation in\nGolang"),t("OutboundLink")],1),e._v(". See the\n"),t("a",{attrs:{href:"https://github.com/tendermint/awesome#ecosystem",target:"_blank",rel:"noopener noreferrer"}},[e._v("list of other ABCI implementations"),t("OutboundLink")],1),e._v(" for servers in\nother languages.")]),e._v(" "),t("p",[e._v("The handler is specific to the application, and may be arbitrary, so\nlong as it is deterministic and conforms to the ABCI interface\nspecification.")]),e._v(" "),t("p",[e._v("So when we run "),t("code",[e._v("abci-cli info")]),e._v(", we open a new connection to the ABCI\nserver, which calls the "),t("code",[e._v("Info()")]),e._v(" method on the application, which tells\nus the number of transactions in our Merkle tree.")]),e._v(" "),t("p",[e._v("Now, since every command opens a new connection, we provide the\n"),t("code",[e._v("abci-cli console")]),e._v(" and "),t("code",[e._v("abci-cli batch")]),e._v(" commands, to allow multiple ABCI\nmessages to be sent over a single connection.")]),e._v(" "),t("p",[e._v("Running "),t("code",[e._v("abci-cli console")]),e._v(" should drop you in an interactive console for\nspeaking ABCI messages to your application.")]),e._v(" "),t("p",[e._v("Try running these commands:")]),e._v(" "),t("tm-code-block",{staticClass:"codeblock",attrs:{language:"sh",base64:"Jmd0OyBlY2hvIGhlbGxvCi0mZ3Q7IGNvZGU6IE9LCi0mZ3Q7IGRhdGE6IGhlbGxvCi0mZ3Q7IGRhdGEuaGV4OiAweDY4NjU2QzZDNkYKCiZndDsgaW5mbyAKLSZndDsgY29kZTogT0sKLSZndDsgZGF0YTogeyZxdW90O3NpemUmcXVvdDs6MH0KLSZndDsgZGF0YS5oZXg6IDB4N0IyMjczNjk3QTY1MjIzQTMwN0QKCiZndDsgcHJlcGFyZV9wcm9wb3NhbCAmcXVvdDthYmMmcXVvdDsKLSZndDsgY29kZTogT0sKLSZndDsgbG9nOiBTdWNjZWVkZWQuIFR4OiBhYmMKCiZndDsgcHJvY2Vzc19wcm9wb3NhbCAmcXVvdDthYmMmcXVvdDsKLSZndDsgY29kZTogT0sKLSZndDsgc3RhdHVzOiBBQ0NFUFQKCiZndDsgY29tbWl0IAotJmd0OyBjb2RlOiBPSwotJmd0OyBkYXRhLmhleDogMHgwMDAwMDAwMDAwMDAwMDAwCgomZ3Q7IGRlbGl2ZXJfdHggJnF1b3Q7YWJjJnF1b3Q7Ci0mZ3Q7IGNvZGU6IE9LCgomZ3Q7IGluZm8gCi0mZ3Q7IGNvZGU6IE9LCi0mZ3Q7IGRhdGE6IHsmcXVvdDtzaXplJnF1b3Q7OjF9Ci0mZ3Q7IGRhdGEuaGV4OiAweDdCMjI3MzY5N0E2NTIyM0EzMTdECgomZ3Q7IGNvbW1pdCAKLSZndDsgY29kZTogT0sKLSZndDsgZGF0YS5oZXg6IDB4MDIwMDAwMDAwMDAwMDAwMAoKJmd0OyBxdWVyeSAmcXVvdDthYmMmcXVvdDsKLSZndDsgY29kZTogT0sKLSZndDsgbG9nOiBleGlzdHMKLSZndDsgaGVpZ2h0OiAyCi0mZ3Q7IGtleTogYWJjCi0mZ3Q7IGtleS5oZXg6IDYxNjI2MwotJmd0OyB2YWx1ZTogYWJjCi0mZ3Q7IHZhbHVlLmhleDogNjE2MjYzCgomZ3Q7IGRlbGl2ZXJfdHggJnF1b3Q7ZGVmPXh5eiZxdW90OwotJmd0OyBjb2RlOiBPSwoKJmd0OyBjb21taXQgCi0mZ3Q7IGNvZGU6IE9LCi0mZ3Q7IGRhdGEuaGV4OiAweDA0MDAwMDAwMDAwMDAwMDAKCiZndDsgcXVlcnkgJnF1b3Q7ZGVmJnF1b3Q7Ci0mZ3Q7IGNvZGU6IE9LCi0mZ3Q7IGxvZzogZXhpc3RzCi0mZ3Q7IGhlaWdodDogMwotJmd0OyBrZXk6IGRlZgotJmd0OyBrZXkuaGV4OiA2NDY1NjYKLSZndDsgdmFsdWU6IHh5egotJmd0OyB2YWx1ZS5oZXg6IDc4Nzk3QQoKJmd0OyBwcmVwYXJlX3Byb3Bvc2FsICZxdW90O3ByZXBhcmVkZWYmcXVvdDsKLSZndDsgY29kZTogT0sKLSZndDsgbG9nOiBTdWNjZWVkZWQuIFR4OiByZXBsYWNlZGVmCgomZ3Q7IHByb2Nlc3NfcHJvcG9zYWwgJnF1b3Q7cmVwbGFjZWRlZiZxdW90OwotJmd0OyBjb2RlOiBPSwotJmd0OyBzdGF0dXM6IEFDQ0VQVAoKJmd0OyBwcm9jZXNzX3Byb3Bvc2FsICZxdW90O3ByZXBhcmVkZWYmcXVvdDsKLSZndDsgY29kZTogT0sKLSZndDsgc3RhdHVzOiBSRUpFQ1QKCiZndDsgcHJlcGFyZV9wcm9wb3NhbCAKCiZndDsgcHJvY2Vzc19wcm9wb3NhbCAKLSZndDsgY29kZTogT0sKLSZndDsgc3RhdHVzOiBBQ0NFUFQKCiZndDsgY29tbWl0IAotJmd0OyBjb2RlOiBPSwotJmd0OyBkYXRhLmhleDogMHgwNDAwMDAwMDAwMDAwMDAwCg=="}}),e._v(" "),t("p",[e._v("Note that if we do "),t("code",[e._v('deliver_tx "abc"')]),e._v(" it will store "),t("code",[e._v("(abc, abc)")]),e._v(", but if\nwe do "),t("code",[e._v('deliver_tx "abc=efg"')]),e._v(" it will store "),t("code",[e._v("(abc, efg)")]),e._v(".")]),e._v(" "),t("p",[e._v("You could put the commands in a file and run\n"),t("code",[e._v("abci-cli --verbose batch < myfile")]),e._v(".")]),e._v(" "),t("p",[e._v("Note that the "),t("code",[e._v("abci-cli")]),e._v(" is designed strictly for testing and debugging. In a real\ndeployment, the role of sending messages is taken by Tendermint, which\nconnects to the app using three separate connections, each with its own\npattern of messages.")]),e._v(" "),t("p",[e._v("For examples of running an ABCI app with Tendermint, see the\n"),t("RouterLink",{attrs:{to:"/app-dev/getting-started.html"}},[e._v("getting started guide")]),e._v(".")],1),e._v(" "),t("h2",{attrs:{id:"bounties"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#bounties"}},[e._v("#")]),e._v(" Bounties")]),e._v(" "),t("p",[e._v("Want to write an app in your favorite language?! We'd be happy\nto add you to our "),t("a",{attrs:{href:"https://github.com/tendermint/awesome#ecosystem",target:"_blank",rel:"noopener noreferrer"}},[e._v("ecosystem"),t("OutboundLink")],1),e._v("!\nSee "),t("a",{attrs:{href:"https://github.com/interchainio/funding",target:"_blank",rel:"noopener noreferrer"}},[e._v("funding"),t("OutboundLink")],1),e._v(" opportunities from the\n"),t("a",{attrs:{href:"https://interchain.io",target:"_blank",rel:"noopener noreferrer"}},[e._v("Interchain Foundation"),t("OutboundLink")],1),e._v(" for implementations in new languages and more.")])],1)}),[],!1,null,null,null);n.default=o.exports}}]);